package amimanagement

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

//go:generate mockgen -destination ec2iface_mock.go -package amimanagement github.com/aws/aws-sdk-go/service/ec2/ec2iface EC2API

// PostProcessor is the core of this library
// Packer performs `PostProcess()` method of this processor
type PostProcessor struct {
	testMode bool
	cleaner  AbstractCleaner
	config   Config
}

// ConfigSpec returns HCL object spec
func (p *PostProcessor) ConfigSpec() hcldec.ObjectSpec {
	return p.config.FlatMapstructure().HCL2Spec()
}

// Configure generates post-processor's configuration
func (p *PostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		PluginType:         "packer.post-processor.amazon-ami-management",
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
	if p.config.KeepReleases != 0 && p.config.KeepDays != 0 {
		return errors.New("`keep_releases` and `keep_days` cannot be set as the same time")
	}
	if p.config.KeepReleases == 0 && p.config.KeepDays == 0 {
		return errors.New("`keep_releases` or `keep_days` must be greater than 1. Please make sure that it is set correctly")
	}
	if p.config.KeepReleases < 1 && p.config.KeepDays == 0 {
		return errors.New("`keep_releases` must be greater than 1. Please make sure that it is set correctly")
	}
	if p.config.KeepDays < 1 && p.config.KeepReleases == 0 {
		return errors.New("`keep_days` must be greater than 1. Please make sure that it is set correctly")
	}
	if len(p.config.Regions) == 0 {
		return errors.New("empty `regions` is not allowed. Please make sure that it is set correctly")
	}

	return nil
}

// PostProcess deletes old AMI and snapshot so as to maintain the number of AMIs expected
func (p *PostProcessor) PostProcess(ctx context.Context, ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, bool, error) {
	log.Println("Running the post-processor")

	for _, region := range p.config.Regions {
		ui.Message(p.uiMessage(fmt.Sprintf("Processing in %s", region)))

		if !p.testMode {
			sess, err := p.config.AccessConfig.Session()
			if err != nil {
				return nil, true, false, err
			}
			p.cleaner, err = NewCleaner(
				sess.Copy(&aws.Config{Region: aws.String(region)}),
				p.config,
			)
			if err != nil {
				return nil, true, false, err
			}
		}

		images, err := p.cleaner.RetrieveCandidateImages()
		if err != nil {
			return nil, true, false, err
		}
		log.Println("Deleting old images...")
		for _, image := range images {
			ui.Message(p.uiMessage(fmt.Sprintf("Deleting image: %s", *image.ImageId)))
			used := p.cleaner.IsUsed(image)
			if used != nil {
				ui.Message(fmt.Sprintf("[WARN] %s is used in %s: %s. Skipped...", *image.ImageId, used.Type, used.ID))
			} else {
				err := p.cleaner.DeleteImage(image)
				if err != nil {
					return nil, true, false, err
				}
			}
		}
	}

	return artifact, true, false, nil
}

func (p *PostProcessor) uiMessage(message string) string {
	if p.config.DryRun {
		return "[DryRun] " + message
	}
	return message
}
