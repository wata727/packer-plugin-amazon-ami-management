package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

//go:generate mockgen -source cleaner.go -destination cleaner_mock.go -package main

// Cleanable is an interface of Cleaner
type Cleanable interface {
	RetrieveCandidateImages() ([]*ec2.Image, error)
	DeleteImage(*ec2.Image) error
	IsUsed(*ec2.Image) *Used
}

// Cleaner is a wrapper of aws-sdk client
type Cleaner struct {
	ec2conn         ec2iface.EC2API
	autoscalingconn autoscalingiface.AutoScalingAPI
	config          Config
	now             time.Time
	used            map[string]*Used
}

// Used is metadata about the details of the image being used
type Used struct {
	Type string
	ID   string
}

// NewCleaner returns a new cleaner
func NewCleaner(sess *session.Session, config Config) (*Cleaner, error) {
	cleaner := &Cleaner{
		ec2conn:         ec2.New(sess),
		autoscalingconn: autoscaling.New(sess),
		config:          config,
		now:             time.Now().UTC(),
		used:            map[string]*Used{},
	}

	err := cleaner.setInstanceUsed()
	if err != nil {
		return nil, err
	}

	err = cleaner.setLaunchConfigurationUsed()
	if err != nil {
		return nil, err
	}

	err = cleaner.setLaunchTemplateUsed()
	if err != nil {
		return nil, err
	}

	return cleaner, nil
}

// RetrieveCandidateImages returns the images of candidate to be deleted.
// These images are sorted in descending order by creation date.
// Please be aware that these are candidates. Not all images are deleted due to output to the Packer UI.
func (c *Cleaner) RetrieveCandidateImages() ([]*ec2.Image, error) {
	log.Println("Describing images")
	output, err := c.ec2conn.DescribeImages(&ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Amazon_AMI_Management_Identifier"),
				Values: []*string{
					aws.String(c.config.Identifier),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Sort in descending order by creation date
	sort.Slice(output.Images, func(i, j int) bool {
		iTime, _ := time.Parse("2006-01-02T15:04:05.000Z", *output.Images[i].CreationDate)
		jTime, _ := time.Parse("2006-01-02T15:04:05.000Z", *output.Images[j].CreationDate)
		return iTime.After(jTime)
	})

	images := []*ec2.Image{}
	for i, image := range output.Images {
		if c.config.KeepReleases != 0 && i < c.config.KeepReleases {
			continue
		}

		if c.config.KeepDays != 0 {
			creationDate, err := time.ParseInLocation("2006-01-02T15:04:05.000Z", *image.CreationDate, time.UTC)
			if err != nil {
				return []*ec2.Image{}, err
			}
			if creationDate.Add(time.Duration(c.config.KeepDays) * 24 * time.Hour).After(c.now) {
				continue
			}
		}

		images = append(images, image)
	}

	return images, nil
}

// DeleteImage deletes a passed image and related snapshots
func (c *Cleaner) DeleteImage(image *ec2.Image) error {
	log.Printf("Deleting AMI (%s)", *image.ImageId)
	if _, err := c.ec2conn.DeregisterImage(&ec2.DeregisterImageInput{
		ImageId: image.ImageId,
		DryRun:  aws.Bool(c.config.DryRun),
	}); err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "DryRunOperation" {
			// noop
		} else {
			return err
		}
	}

	// DeregisterImage method only performs to AMI
	// Because it retains snapshot. Following operation is deleting snapshots.
	log.Printf("Deleting snapshot related to AMI (%s)", *image.ImageId)
	for _, device := range image.BlockDeviceMappings {
		// skip delete if use ephemeral devise
		if device.Ebs == nil || device.Ebs.SnapshotId == nil {
			continue
		}
		log.Printf("Deleting snapshot (%s) related to AMI (%s)", *device.Ebs.SnapshotId, *image.ImageId)
		if _, err := c.ec2conn.DeleteSnapshot(&ec2.DeleteSnapshotInput{
			SnapshotId: device.Ebs.SnapshotId,
			DryRun:     aws.Bool(c.config.DryRun),
		}); err != nil {
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "DryRunOperation" {
				// noop
			} else {
				return err
			}
		}
	}

	return nil
}

// IsUsed checks whether a passed image is used.
// If used, it returns Used instead of nil.
func (c *Cleaner) IsUsed(image *ec2.Image) *Used {
	return c.used[*image.ImageId]
}

func (c *Cleaner) setInstanceUsed() error {
	ret, err := c.ec2conn.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("pending"),
					aws.String("running"),
					aws.String("shutting-down"),
					aws.String("stopping"),
					aws.String("stopped"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	for _, reservation := range ret.Reservations {
		for _, instance := range reservation.Instances {
			if instance.ImageId == nil {
				continue
			}

			c.used[*instance.ImageId] = &Used{
				ID:   *instance.InstanceId,
				Type: "instance",
			}
		}
	}
	return nil
}

func (c *Cleaner) setLaunchConfigurationUsed() error {
	ret, err := c.autoscalingconn.DescribeLaunchConfigurations(&autoscaling.DescribeLaunchConfigurationsInput{})
	if err != nil {
		return err
	}
	for _, lc := range ret.LaunchConfigurations {
		if lc.ImageId == nil {
			continue
		}

		c.used[*lc.ImageId] = &Used{
			ID:   *lc.LaunchConfigurationName,
			Type: "launch configuration",
		}
	}
	return nil
}

func (c *Cleaner) setLaunchTemplateUsed() error {
	ret, err := c.ec2conn.DescribeLaunchTemplates(&ec2.DescribeLaunchTemplatesInput{})
	if err != nil {
		return err
	}

	for _, lt := range ret.LaunchTemplates {
		versions, err := c.ec2conn.DescribeLaunchTemplateVersions(&ec2.DescribeLaunchTemplateVersionsInput{
			LaunchTemplateId: lt.LaunchTemplateId,
		})
		if err != nil {
			return err
		}

		for _, ltv := range versions.LaunchTemplateVersions {
			if ltv.LaunchTemplateData == nil || ltv.LaunchTemplateData.ImageId == nil {
				continue
			}

			c.used[*ltv.LaunchTemplateData.ImageId] = &Used{
				ID:   fmt.Sprintf("%s (%d)", *ltv.LaunchTemplateName, *ltv.VersionNumber),
				Type: "launch template",
			}
		}
	}
	return nil
}
