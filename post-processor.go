package main

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	awscommon "github.com/hashicorp/packer/builder/amazon/common"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

// Config is a post-processor's configuration
// PostProcessor generates it using Packer's configuration in `Configure()` method
type Config struct {
	common.PackerConfig    `mapstructure:",squash"`
	awscommon.AccessConfig `mapstructure:",squash"`

	Identifier   string   `mapstructure:"identifier"`
	KeepReleases int      `mapstructure:"keep_releases"`
	Regions      []string `mapstructure:"regions"`

	ctx interpolate.Context
}

// PostProcessor is the core of this library
// Packer performs `PostProcess()` method of this processor
type PostProcessor struct {
	testMode bool
	ec2conn  ec2iface.EC2API
	config   Config
}

// Configure generates post-processor's configuration
func (p *PostProcessor) Configure(raws ...interface{}) error {
	p.config.ctx.Funcs = awscommon.TemplateFuncs
	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return err
	}

	if p.config.Identifier == "" {
		return errors.New("empty `identifier` is not allowed. Please make sure that it is set correctly")
	}
	if p.config.KeepReleases < 1 {
		return errors.New("`keep_releases` must be greater than 1. Please make sure that it is set correctly")
	}
	if len(p.config.Regions) == 0 {
		return errors.New("empty `regions` is not allowed. Please make sure that it is set correctly")
	}

	return nil
}

// PostProcess deletes old AMI and snapshot so as to maintain the number of AMIs expected
func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	log.Println("Running the post-processor")

	for _, region := range p.config.Regions {
		ui.Message(fmt.Sprintf("Processing in %s", region))

		if !p.testMode {
			sess, err := p.config.AccessConfig.Session()
			if err != nil {
				return nil, true, err
			}
			p.ec2conn = ec2.New(sess.Copy(&aws.Config{Region: aws.String(region)}))
		}

		if err := p.manageAMIs(ui); err != nil {
			return nil, true, err
		}
	}

	return artifact, true, nil
}

func (p *PostProcessor) manageAMIs(ui packer.Ui) error {
	log.Println("Describing images")
	output, err := p.ec2conn.DescribeImages(&ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Amazon_AMI_Management_Identifier"),
				Values: []*string{
					aws.String(p.config.Identifier),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	// Sort in descending order by creation date
	sort.Slice(output.Images, func(i, j int) bool {
		iTime, _ := time.Parse("2006-01-02T15:04:05.000Z", *output.Images[i].CreationDate)
		jTime, _ := time.Parse("2006-01-02T15:04:05.000Z", *output.Images[j].CreationDate)
		return iTime.After(jTime)
	})

	log.Println("Deleting old images...")
	for i, image := range output.Images {
		if i < p.config.KeepReleases {
			continue
		}
		ui.Message(fmt.Sprintf("Deleting image: %s", *image.ImageId))
		log.Printf("Deleting AMI (%s)", *image.ImageId)
		if _, err := p.ec2conn.DeregisterImage(&ec2.DeregisterImageInput{
			ImageId: image.ImageId,
		}); err != nil {
			return err
		}

		// DeregisterImage method only performs to AMI
		// Because it retains snapshot. Following operation is deleting snapshots.
		log.Printf("Deleting snapshot related to AMI (%s)", *image.ImageId)
		for _, device := range image.BlockDeviceMappings {
			// skip delete if use ephemeral devise
			if device.Ebs == nil {
				continue
			}
			log.Printf("Deleting snapshot (%s) related to AMI (%s)", *device.Ebs.SnapshotId, *image.ImageId)
			if _, err := p.ec2conn.DeleteSnapshot(&ec2.DeleteSnapshotInput{
				SnapshotId: device.Ebs.SnapshotId,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}
