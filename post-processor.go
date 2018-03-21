package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	awscommon "github.com/hashicorp/packer/builder/amazon/common"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig    `mapstructure:",squash"`
	awscommon.AccessConfig `mapstructure:",squash"`

	Identifier   string `mapstructure:"identifier"`
	KeepReleases int    `mapstructure:"keep_releases"`
	AccessKey    string `mapstructure:"access_key"`
	SecretKey    string `mapstructure:"secret_key"`
	Region       string `mapstructure:"region"`

	ctx interpolate.Context
}

type PostProcessor struct {
	ec2conn ec2iface.EC2API
	config  Config
}

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

	return nil
}

func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	log.Println("Running the post-processor")

	ec2conn := p.ec2conn
	if ec2conn == nil {
		// If processor doesn't set connection, make the real connection
		config := aws.NewConfig().WithRegion(p.config.Region).WithMaxRetries(11)
		sess := session.New(config)
		creds := credentials.NewChainCredentials([]credentials.Provider{
			&credentials.StaticProvider{Value: credentials.Value{
				AccessKeyID:     p.config.AccessKey,
				SecretAccessKey: p.config.SecretKey,
			}},
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{Filename: "", Profile: ""},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(sess),
			},
		})

		log.Println("Creating a session")
		ec2Session := session.New(config.WithCredentials(creds))
		ec2conn = ec2.New(ec2Session)
	}

	log.Println("Describing images")
	output, err := ec2conn.DescribeImages(&ec2.DescribeImagesInput{
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
		return nil, true, err
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
		if _, err := ec2conn.DeregisterImage(&ec2.DeregisterImageInput{
			ImageId: image.ImageId,
		}); err != nil {
			return nil, true, err
		}

		// DeregisterImage method only perform to AMI
		// Because it retain snapshots. Following operation is deleting snapshots.
		log.Printf("Deleting snapshot related to AMI (%s)", *image.ImageId)
		for _, device := range image.BlockDeviceMappings {
			// skip delete if use ephemeral devise
			if device.Ebs == nil {
				continue
			}
			log.Printf("Deleting snapshot (%s) related to AMI (%s)", *device.Ebs.SnapshotId, *image.ImageId)
			if _, err := ec2conn.DeleteSnapshot(&ec2.DeleteSnapshotInput{
				SnapshotId: device.Ebs.SnapshotId,
			}); err != nil {
				return nil, true, err
			}
		}
	}

	return artifact, true, nil
}
