package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
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
	resolvedAliases map[string]string
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
		resolvedAliases: map[string]string{},
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
	filters := c.genTagsFilter()
	output, err := c.ec2conn.DescribeImages(&ec2.DescribeImagesInput{
		Filters: filters,
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

			imageId := *ltv.LaunchTemplateData.ImageId

			if c.config.ResolveAliases {
				resolvedImageId, err := c.resolveImageAlias(imageId, lt.LaunchTemplateId, ltv.VersionNumber)

				if err != nil {
					return err
				}

				imageId = resolvedImageId
			}

			c.used[imageId] = &Used{
				ID:   fmt.Sprintf("%s (%d)", *ltv.LaunchTemplateName, *ltv.VersionNumber),
				Type: "launch template",
			}
		}
	}
	return nil
}

// If the passed imageAlias is an alias, It will return the resolved imageId.
// Otherwise, it will return the original imageId.
// See https://docs.aws.amazon.com/autoscaling/ec2/userguide/using-systems-manager-parameters.html to understand this use case.
func (c *Cleaner) resolveImageAlias(imageAlias string, launchTemplateId *string, launchTemplateVersion *int64) (string, error) {
	// If we need to resolve aliases. We need to perform an additional DescribeLaunchTemplateVersions
	// for each version we want to resolve. Because when calling the DescribeLaunchTemplateVersions operation: Resource aliasing (resolveAlias)
	// is only supported when doing single version describe
	imageIdIsAliased := strings.HasPrefix(imageAlias, "resolve:ssm:")
	resolvedImageId, isAliasResolved := c.resolvedAliases[imageAlias]
	imageId := imageAlias

	if imageIdIsAliased && isAliasResolved {
		// If already resolved, Just use the resolved value.
		imageId = resolvedImageId
	} else if imageIdIsAliased && !isAliasResolved {
		// If have not been already resolved, resolve the alias.
		aliasedVersions, err := c.ec2conn.DescribeLaunchTemplateVersions(&ec2.DescribeLaunchTemplateVersionsInput{
			LaunchTemplateId: launchTemplateId,
			Versions: []*string{
				aws.String(strconv.FormatInt(*launchTemplateVersion, 10)),
			},
			ResolveAlias: aws.Bool(true),
		})

		if err != nil {
			return "", err
		}

		if len(aliasedVersions.LaunchTemplateVersions) == 0 {
			return imageAlias, nil
		}

		// Save the resolved value.
		ltv := aliasedVersions.LaunchTemplateVersions[0]
		imageId = *ltv.LaunchTemplateData.ImageId

		// We track if a given alias was already resolved to avoid unnecesary additional DescribeLaunchTemplateVersions requests.
		c.resolvedAliases[imageAlias] = imageId
	}

	return imageId, nil
}

func (c *Cleaner) genTagsFilter() []*ec2.Filter {
	var (
		filters []*ec2.Filter
	)
	if c.config.Identifier != "" {
		filters = append(filters, &ec2.Filter{
			Name: aws.String("tag:Amazon_AMI_Management_Identifier"),
			Values: []*string{
				aws.String(c.config.Identifier),
			},
		})
	} else {
		for k, v := range c.config.Tags {
			filters = append(filters, &ec2.Filter{
				Name: aws.String(fmt.Sprintf("tag:%s", k)),
				Values: []*string{
					aws.String(v),
				},
			})
		}
	}
	return filters
}
