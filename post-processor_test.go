package main

import (
	"bytes"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/packer/packer"
)

//go:generate mockgen -source vendor/github.com/aws/aws-sdk-go/service/ec2/ec2iface/interface.go -destination mock.go -package main

func testUI() *packer.BasicUi {
	return &packer.BasicUi{
		Reader: new(bytes.Buffer),
		Writer: new(bytes.Buffer),
	}
}

func TestPostProcessor_ImplementsPostProcessor(t *testing.T) {
	var _ packer.PostProcessor = new(PostProcessor)
}

func TestPostProcessor_Configure_validConfig(t *testing.T) {
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
		"identifier":    "",
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

func TestPostProcessor_PostProcess_emptyImages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ec2mock := NewMockEC2API(ctrl)

	ec2mock.EXPECT().DescribeImages(&ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:Amazon_AMI_Management_Identifier"),
				Values: []*string{
					aws.String("packer-example"),
				},
			},
		},
	}).Return(&ec2.DescribeImagesOutput{
		Images: []*ec2.Image{},
	}, nil)

	p := PostProcessor{
		testMode: true,
		ec2conn:  ec2mock,
	}
	p.config.Identifier = "packer-example"
	p.config.Regions = []string{"us-east-1"}
	artifact := &packer.MockArtifact{}
	_, keep, err := p.PostProcess(testUI(), artifact)

	if !keep {
		t.Fatal("should keep")
	}

	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPostProcessor_PostProcess_fewImages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ec2mock := NewMockEC2API(ctrl)

	ec2mock.EXPECT().DescribeImages(&ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:Amazon_AMI_Management_Identifier"),
				Values: []*string{
					aws.String("packer-example"),
				},
			},
		},
	}).Return(&ec2.DescribeImagesOutput{
		Images: []*ec2.Image{&ec2.Image{
			CreationDate: aws.String("2016-08-01T15:04:05.000Z"),
		}, &ec2.Image{
			CreationDate: aws.String("2016-08-04T15:04:05.000Z"),
		}},
	}, nil)

	p := PostProcessor{
		testMode: true,
		ec2conn:  ec2mock,
	}
	p.config.Identifier = "packer-example"
	p.config.KeepReleases = 3
	p.config.Regions = []string{"us-east-1"}
	artifact := &packer.MockArtifact{}
	_, keep, err := p.PostProcess(testUI(), artifact)

	if !keep {
		t.Fatal("should keep")
	}

	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPostProcessor_PostProcess_manyImages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ec2mock := NewMockEC2API(ctrl)

	ec2mock.EXPECT().DescribeImages(&ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:Amazon_AMI_Management_Identifier"),
				Values: []*string{
					aws.String("packer-example"),
				},
			},
		},
	}).Return(&ec2.DescribeImagesOutput{
		Images: []*ec2.Image{&ec2.Image{
			ImageId:      aws.String("ami-12345a"),
			CreationDate: aws.String("2016-08-01T15:04:05.000Z"),
		}, &ec2.Image{
			ImageId:      aws.String("ami-12345b"),
			CreationDate: aws.String("2016-08-04T15:04:05.000Z"),
		}, &ec2.Image{
			ImageId:      aws.String("ami-12345c"),
			CreationDate: aws.String("2016-07-29T15:04:05.000Z"),
			BlockDeviceMappings: []*ec2.BlockDeviceMapping{&ec2.BlockDeviceMapping{
				Ebs: &ec2.EbsBlockDevice{
					SnapshotId: aws.String("snap-12345a"),
				},
			}, &ec2.BlockDeviceMapping{
				Ebs: &ec2.EbsBlockDevice{
					SnapshotId: aws.String("snap-12345b"),
				},
			}},
		}},
	}, nil)

	ec2mock.EXPECT().DeregisterImage(&ec2.DeregisterImageInput{
		ImageId: aws.String("ami-12345c"),
		DryRun:  aws.Bool(false),
	}).Return(&ec2.DeregisterImageOutput{}, nil)
	ec2mock.EXPECT().DeleteSnapshot(&ec2.DeleteSnapshotInput{
		SnapshotId: aws.String("snap-12345a"),
		DryRun:     aws.Bool(false),
	}).Return(&ec2.DeleteSnapshotOutput{}, nil)
	ec2mock.EXPECT().DeleteSnapshot(&ec2.DeleteSnapshotInput{
		SnapshotId: aws.String("snap-12345b"),
		DryRun:     aws.Bool(false),
	}).Return(&ec2.DeleteSnapshotOutput{}, nil)

	p := PostProcessor{
		testMode: true,
		ec2conn:  ec2mock,
	}
	p.config.Identifier = "packer-example"
	p.config.KeepReleases = 2
	p.config.Regions = []string{"us-east-1"}
	artifact := &packer.MockArtifact{}
	_, keep, err := p.PostProcess(testUI(), artifact)

	if !keep {
		t.Fatal("should keep")
	}

	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPostProcessor_PostProcess_ephemeralDevise(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ec2mock := NewMockEC2API(ctrl)

	ec2mock.EXPECT().DescribeImages(&ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:Amazon_AMI_Management_Identifier"),
				Values: []*string{
					aws.String("packer-example"),
				},
			},
		},
	}).Return(&ec2.DescribeImagesOutput{
		Images: []*ec2.Image{&ec2.Image{
			ImageId:      aws.String("ami-12345a"),
			CreationDate: aws.String("2016-08-20T12:19:56.000Z"),
			BlockDeviceMappings: []*ec2.BlockDeviceMapping{&ec2.BlockDeviceMapping{
				Ebs: &ec2.EbsBlockDevice{
					SnapshotId: aws.String("snap-12345a"),
				},
			}, &ec2.BlockDeviceMapping{
				Ebs: nil,
			}, &ec2.BlockDeviceMapping{
				Ebs: nil,
			}},
		}},
	}, nil)

	ec2mock.EXPECT().DeregisterImage(&ec2.DeregisterImageInput{
		ImageId: aws.String("ami-12345a"),
		DryRun:  aws.Bool(false),
	}).Return(&ec2.DeregisterImageOutput{}, nil)
	ec2mock.EXPECT().DeleteSnapshot(&ec2.DeleteSnapshotInput{
		SnapshotId: aws.String("snap-12345a"),
		DryRun:     aws.Bool(false),
	}).Return(&ec2.DeleteSnapshotOutput{}, nil)

	p := PostProcessor{
		testMode: true,
		ec2conn:  ec2mock,
	}
	p.config.Identifier = "packer-example"
	p.config.KeepReleases = 0
	p.config.Regions = []string{"us-east-1"}
	artifact := &packer.MockArtifact{}
	_, keep, err := p.PostProcess(testUI(), artifact)

	if !keep {
		t.Fatal("should keep")
	}

	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPostProcessor_PostProcess_ephemeralDevise_withDryRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ec2mock := NewMockEC2API(ctrl)

	ec2mock.EXPECT().DescribeImages(&ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:Amazon_AMI_Management_Identifier"),
				Values: []*string{
					aws.String("packer-example"),
				},
			},
		},
	}).Return(&ec2.DescribeImagesOutput{
		Images: []*ec2.Image{&ec2.Image{
			ImageId:      aws.String("ami-12345a"),
			CreationDate: aws.String("2016-08-20T12:19:56.000Z"),
			BlockDeviceMappings: []*ec2.BlockDeviceMapping{&ec2.BlockDeviceMapping{
				Ebs: &ec2.EbsBlockDevice{
					SnapshotId: aws.String("snap-12345a"),
				},
			}, &ec2.BlockDeviceMapping{
				Ebs: nil,
			}, &ec2.BlockDeviceMapping{
				Ebs: nil,
			}},
		}},
	}, nil)

	ec2mock.EXPECT().DeregisterImage(&ec2.DeregisterImageInput{
		ImageId: aws.String("ami-12345a"),
		DryRun:  aws.Bool(true),
	}).Return(nil, awserr.New(
		"DryRunOperation",
		"Request would have succeeded, but DryRun flag is set.",
		errors.New("Request would have succeeded, but DryRun flag is set."),
	))
	ec2mock.EXPECT().DeleteSnapshot(&ec2.DeleteSnapshotInput{
		SnapshotId: aws.String("snap-12345a"),
		DryRun:     aws.Bool(true),
	}).Return(nil, awserr.New(
		"DryRunOperation",
		"Request would have succeeded, but DryRun flag is set.",
		errors.New("Request would have succeeded, but DryRun flag is set."),
	))

	p := PostProcessor{
		testMode: true,
		ec2conn:  ec2mock,
	}
	p.config.Identifier = "packer-example"
	p.config.KeepReleases = 0
	p.config.Regions = []string{"us-east-1"}
	p.config.DryRun = true
	artifact := &packer.MockArtifact{}
	_, keep, err := p.PostProcess(testUI(), artifact)

	if !keep {
		t.Fatal("should keep")
	}

	if err != nil {
		t.Fatalf("err: %s", err)
	}
}
