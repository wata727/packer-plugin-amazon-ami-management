package main

import (
	"bytes"
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

func testUI() *packer.BasicUi {
	return &packer.BasicUi{
		Reader: new(bytes.Buffer),
		Writer: new(bytes.Buffer),
	}
}

func TestPostProcessor_ImplementsPostProcessor(t *testing.T) {
	var _ packer.PostProcessor = new(PostProcessor)
}

func TestConfigCases(t *testing.T) {
	var (
		defaultRegions    = []string{"us-east-1"}
		defaultIndentifer = "packer-example"
		defaultTags       = map[string]string{
			"Amazon_AMI_Management_Identifier": "packer-example",
		}

		configTestCases = []struct {
			Name           string
			ExptectedError string
			Config         map[string]interface{}
		}{
			{
				Name:           "Missing Regions",
				ExptectedError: "empty `regions` is not allowed. Please make sure that it is set correctly",
				Config: map[string]interface{}{
					"identifier":    defaultIndentifer,
					"keep_releases": 3,
				},
			},
			{
				Name:           "Missing Regions",
				ExptectedError: "empty `regions` is not allowed. Please make sure that it is set correctly",
				Config: map[string]interface{}{
					"tags":          defaultTags,
					"keep_releases": 3,
				},
			},
			{
				Name:           "Invalid KeepReleases",
				ExptectedError: "`keep_releases` must be greater than 1. Please make sure that it is set correctly",
				Config: map[string]interface{}{
					"identifier":    defaultIndentifer,
					"keep_releases": -1,
				},
			},
			{
				Name:           "Invalid KeepReleases",
				ExptectedError: "`keep_releases` must be greater than 1. Please make sure that it is set correctly",
				Config: map[string]interface{}{
					"tags":          defaultTags,
					"keep_releases": -1,
				},
			},
			{
				Name:           "Invalid KeepDays",
				ExptectedError: "`keep_days` must be greater than 1. Please make sure that it is set correctly",
				Config: map[string]interface{}{
					"regions":    defaultRegions,
					"identifier": defaultIndentifer,
					"keep_days":  -1,
				},
			},
			{
				Name:           "Invalid KeepDays",
				ExptectedError: "`keep_days` must be greater than 1. Please make sure that it is set correctly",
				Config: map[string]interface{}{
					"regions":   defaultRegions,
					"tags":      defaultTags,
					"keep_days": -1,
				},
			},
			{
				Name:           "Set KeepReleases and KeepDays",
				ExptectedError: "`keep_releases` and `keep_days` cannot be set as the same time",
				Config: map[string]interface{}{
					"regions":       defaultRegions,
					"identifier":    defaultIndentifer,
					"keep_releases": 3,
					"keep_days":     10,
				},
			},
			{
				Name:           "Set KeepReleases and KeepDays",
				ExptectedError: "`keep_releases` and `keep_days` cannot be set as the same time",
				Config: map[string]interface{}{
					"regions":       defaultRegions,
					"tags":          defaultTags,
					"keep_releases": 3,
					"keep_days":     10,
				},
			},
			{
				Name:           "Neither KeepReleases nor KeepDays is set",
				ExptectedError: "`keep_releases` or `keep_days` must be greater than 1. Please make sure that it is set correctly",
				Config: map[string]interface{}{
					"regions":    defaultRegions,
					"identifier": defaultIndentifer,
				},
			},
			{
				Name:           "validate config with tags",
				ExptectedError: "`keep_releases` or `keep_days` must be greater than 1. Please make sure that it is set correctly",
				Config: map[string]interface{}{
					"regions": defaultRegions,
					"tags":    defaultTags,
				},
			},
			{
				Name:           "Empty Identifier and Tags",
				ExptectedError: "`identifier` or `tags` must be defined. Please make sure that it is set correctly",
				Config: map[string]interface{}{
					"regions":       defaultRegions,
					"keep_releases": 3,
				},
			},
			{
				Name: "Configure valid config with KeepDays",
				Config: map[string]interface{}{
					"regions":    defaultRegions,
					"identifier": defaultIndentifer,
					"keep_days":  10,
				},
			},
			{
				Name: "Configure valid config with KeepDays",
				Config: map[string]interface{}{
					"regions":   defaultRegions,
					"tags":      defaultTags,
					"keep_days": 10,
				},
			},
			{
				Name: "Configure valid config with KeepReleases",
				Config: map[string]interface{}{
					"regions":       defaultRegions,
					"identifier":    defaultIndentifer,
					"keep_releases": 3,
				},
			},
			{
				Name: "Configure valid config with KeepReleases",
				Config: map[string]interface{}{
					"regions":       defaultRegions,
					"tags":          defaultTags,
					"keep_releases": 3,
				},
			},
			{
				Name: "Configure valid config with ResolveAliases set to true",
				Config: map[string]interface{}{
					"regions":         defaultRegions,
					"tags":            defaultTags,
					"keep_releases":   3,
					"resolve_aliases": true,
				},
			},
			{
				Name: "Configure valid config with ResolveAliases set to false",
				Config: map[string]interface{}{
					"regions":         defaultRegions,
					"tags":            defaultTags,
					"keep_releases":   3,
					"resolve_aliases": false,
				},
			},
			{
				Name: "Configure valid config with ResolveAliases set to true",
				Config: map[string]interface{}{
					"regions":         defaultRegions,
					"tags":            defaultTags,
					"keep_days":       10,
					"resolve_aliases": true,
				},
			},
			{
				Name: "Configure valid config with ResolveAliases set to false",
				Config: map[string]interface{}{
					"regions":         defaultRegions,
					"tags":            defaultTags,
					"keep_days":       10,
					"resolve_aliases": false,
				},
			},
		}
	)

	for _, c := range configTestCases {
		t.Run(c.Name, func(t *testing.T) {
			p := new(PostProcessor)
			err := p.Configure(c.Config)
			if c.ExptectedError != "" {
				if err == nil {
					t.Fatalf("case: %s should cause validation errors", c.Name)
				}
				if err.Error() != c.ExptectedError {
					t.Fatalf("case: %s unexpected error occurred: %s", err, c.Name)
				}
			} else {
				if err != nil {
					t.Fatalf("case: %s unexpected error occurred: %s", err, c.Name)
				}
			}
		})
	}
}

func TestPostProcessor_PostProcess(t *testing.T) {
	var (
		defaultRegions = []string{"us-east-1"}
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cleanermock := NewMockCleanable(ctrl)

	cleanermock.EXPECT().RetrieveCandidateImages().Return(
		[]*ec2.Image{
			{ImageId: aws.String("ami-12345a")},
			{ImageId: aws.String("ami-12345b")},
		},
		nil,
	)
	cleanermock.EXPECT().IsUsed(&ec2.Image{ImageId: aws.String("ami-12345a")}).Return(nil)
	cleanermock.EXPECT().DeleteImage(&ec2.Image{ImageId: aws.String("ami-12345a")}).Return(nil)
	cleanermock.EXPECT().IsUsed(&ec2.Image{ImageId: aws.String("ami-12345b")}).Return(&Used{
		ID:   "i-12345678",
		Type: "instance",
	})

	p := PostProcessor{
		testMode: true,
		cleaner:  cleanermock,
		config: Config{
			Regions: defaultRegions,
		},
	}

	_, keep, forceOverride, err := p.PostProcess(context.Background(), testUI(), &packer.MockArtifact{})
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
	if !keep {
		t.Fatal("should keep")
	}
	if forceOverride {
		t.Fatal("should not override")
	}
}
