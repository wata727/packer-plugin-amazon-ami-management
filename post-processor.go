package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
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
	KeepDays     int      `mapstructure:"keep_days"`
	Regions      []string `mapstructure:"regions"`
	DryRun       bool     `mapstructure:"dry_run"`

	ctx interpolate.Context
}

// PostProcessor is the core of this library
// Packer performs `PostProcess()` method of this processor
type PostProcessor struct {
	testMode bool
	cleaner  AbstractCleaner
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
func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	log.Println("Running the post-processor")

	for _, region := range p.config.Regions {
		ui.Message(p.uiMessage(fmt.Sprintf("Processing in %s", region)))

		if !p.testMode {
			sess, err := p.config.AccessConfig.Session()
			if err != nil {
				return nil, true, err
			}
			p.cleaner = NewCleaner(
				ec2.New(sess.Copy(&aws.Config{Region: aws.String(region)})),
				p.config,
			)
		}

		images, err := p.cleaner.RetrieveCandidateImages()
		if err != nil {
			return nil, true, err
		}
		log.Println("Deleting old images...")
		for _, image := range images {
			ui.Message(p.uiMessage(fmt.Sprintf("Deleting image: %s", *image.ImageId)))
			err := p.cleaner.DeleteImage(image)
			if err != nil {
				return nil, true, err
			}
		}
	}

	return artifact, true, nil
}

func (p *PostProcessor) uiMessage(message string) string {
	if p.config.DryRun {
		return "[DryRun] " + message
	}
	return message
}
