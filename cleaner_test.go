package main

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
)

//go:generate mockgen -source vendor/github.com/aws/aws-sdk-go/service/ec2/ec2iface/interface.go -destination ec2iface_mock.go -package main

func TestCleaner_RetrieveCandidateImages_KeepReleases(t *testing.T) {
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

	cleaner := &Cleaner{
		ec2conn: ec2mock,
		config: Config{
			Identifier:   "packer-example",
			KeepReleases: 2,
		},
		now: time.Now().UTC(),
	}

	images, err := cleaner.RetrieveCandidateImages()
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
	if len(images) != 1 {
		t.Fatalf("Unexpected image count: %d", len(images))
	}
	if *images[0].ImageId != "ami-12345c" {
		t.Fatalf("Unexpected image: %s", *images[0].ImageId)
	}
}

func TestCleaner_RetrieveCandidateImages_KeepDays(t *testing.T) {
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

	cleaner := &Cleaner{
		ec2conn: ec2mock,
		config: Config{
			Identifier: "packer-example",
			KeepDays:   10,
		},
		now: time.Date(2016, time.August, 11, 11, 0, 0, 0, time.UTC),
	}

	images, err := cleaner.RetrieveCandidateImages()
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
	if len(images) != 1 {
		t.Fatalf("Unexpected image count: %d", len(images))
	}
	if *images[0].ImageId != "ami-12345c" {
		t.Fatalf("Unexpected image: %s", *images[0].ImageId)
	}
}

func TestCleaner_DeleteImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ec2mock := NewMockEC2API(ctrl)

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

	cleaner := &Cleaner{
		ec2conn: ec2mock,
	}

	err := cleaner.DeleteImage(&ec2.Image{
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
	})
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
}

func TestCleaner_DeleteImage_EphemeralDevise(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ec2mock := NewMockEC2API(ctrl)

	ec2mock.EXPECT().DeregisterImage(&ec2.DeregisterImageInput{
		ImageId: aws.String("ami-12345c"),
		DryRun:  aws.Bool(false),
	}).Return(&ec2.DeregisterImageOutput{}, nil)
	ec2mock.EXPECT().DeleteSnapshot(&ec2.DeleteSnapshotInput{
		SnapshotId: aws.String("snap-12345a"),
		DryRun:     aws.Bool(false),
	}).Return(&ec2.DeleteSnapshotOutput{}, nil)

	cleaner := &Cleaner{
		ec2conn: ec2mock,
	}

	err := cleaner.DeleteImage(&ec2.Image{
		ImageId:      aws.String("ami-12345c"),
		CreationDate: aws.String("2016-07-29T15:04:05.000Z"),
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{&ec2.BlockDeviceMapping{
			Ebs: &ec2.EbsBlockDevice{
				SnapshotId: aws.String("snap-12345a"),
			},
		}, &ec2.BlockDeviceMapping{
			Ebs: nil,
		}, &ec2.BlockDeviceMapping{
			Ebs: nil,
		}},
	})
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
}

func TestCleaner_DeleteImage_DryRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ec2mock := NewMockEC2API(ctrl)

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

	cleaner := &Cleaner{
		ec2conn: ec2mock,
		config: Config{
			DryRun: true,
		},
	}

	err := cleaner.DeleteImage(&ec2.Image{
		ImageId:      aws.String("ami-12345a"),
		CreationDate: aws.String("2016-07-29T15:04:05.000Z"),
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{&ec2.BlockDeviceMapping{
			Ebs: &ec2.EbsBlockDevice{
				SnapshotId: aws.String("snap-12345a"),
			},
		}},
	})
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
}
