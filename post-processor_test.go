package main

import (
	"bytes"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/packer/packer"
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

func TestPostProcessor_Configure_validConfigWithKeepReleases(t *testing.T) {
	p := new(PostProcessor)
	err := p.Configure(map[string]interface{}{
		"regions":       []string{"us-east-1"},
		"identifier":    "packer-example",
		"keep_releases": 3,
	})

	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPostProcessor_Configure_validConfigWithKeepDays(t *testing.T) {
	p := new(PostProcessor)
	err := p.Configure(map[string]interface{}{
		"regions":    []string{"us-east-1"},
		"identifier": "packer-example",
		"keep_days":  10,
	})

	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPostProcessor_Configure_missingRegions(t *testing.T) {
	p := new(PostProcessor)
	err := p.Configure(map[string]interface{}{
		"region":        "us-east-1",
		"identifier":    "packer-example",
		"keep_releases": 3,
	})

	if err == nil {
		t.Fatal("should cause validation errors")
	}
	if err.Error() != "empty `regions` is not allowed. Please make sure that it is set correctly" {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
}

func TestPostProcessor_Configure_emptyIdentifier(t *testing.T) {
	p := new(PostProcessor)
	err := p.Configure(map[string]interface{}{
		"regions":       []string{"us-east-1"},
		"keep_releases": 3,
	})

	if err == nil {
		t.Fatal("should cause validation errors")
	}
	if err.Error() != "empty `identifier` is not allowed. Please make sure that it is set correctly" {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
}

func TestPostProcessor_Configure_invalidKeepReleases(t *testing.T) {
	p := new(PostProcessor)
	err := p.Configure(map[string]interface{}{
		"regions":       []string{"us-east-1"},
		"identifier":    "packer-example",
		"keep_releases": -1,
	})

	if err == nil {
		t.Fatal("should cause validation errors")
	}
	if err.Error() != "`keep_releases` must be greater than 1. Please make sure that it is set correctly" {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
}

func TestPostProcessor_Configure_invalidKeepDays(t *testing.T) {
	p := new(PostProcessor)
	err := p.Configure(map[string]interface{}{
		"regions":    []string{"us-east-1"},
		"identifier": "packer-example",
		"keep_days":  -1,
	})

	if err == nil {
		t.Fatal("should cause validation errors")
	}
	if err.Error() != "`keep_days` must be greater than 1. Please make sure that it is set correctly" {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
}

func TestPostProcessor_Configure_setKeepReleasesAndKeepDays(t *testing.T) {
	p := new(PostProcessor)
	err := p.Configure(map[string]interface{}{
		"regions":       []string{"us-east-1"},
		"identifier":    "packer-example",
		"keep_releases": 3,
		"keep_days":     10,
	})

	if err == nil {
		t.Fatal("should cause validation errors")
	}
	if err.Error() != "`keep_releases` and `keep_days` cannot be set as the same time" {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
}

func TestPostProcessor_Configure_NeitherKeepReleasesNorKeepDaysIsSet(t *testing.T) {
	p := new(PostProcessor)
	err := p.Configure(map[string]interface{}{
		"regions":    []string{"us-east-1"},
		"identifier": "packer-example",
	})

	if err == nil {
		t.Fatal("should cause validation errors")
	}
	if err.Error() != "`keep_releases` or `keep_days` must be greater than 1. Please make sure that it is set correctly" {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
}

func TestPostProcessor_PostProcess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cleanermock := NewMockAbstractCleaner(ctrl)

	cleanermock.EXPECT().RetrieveCandidateImages().Return(
		[]*ec2.Image{
			&ec2.Image{ImageId: aws.String("ami-12345a")},
			&ec2.Image{ImageId: aws.String("ami-12345b")},
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
			Regions: []string{"us-east-1"},
		},
	}

	_, keep, err := p.PostProcess(testUI(), &packer.MockArtifact{})
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
	if !keep {
		t.Fatal("should keep")
	}
}
