//go:generate packer-sdc mapstructure-to-hcl2 -type Config,AssumeRoleConfig

package main

import (
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

// Config is a post-processor's configuration
// PostProcessor generates it using Packer's configuration in `Configure()` method
type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	AccessConfig        `mapstructure:",squash"`

	KeepReleases int               `mapstructure:"keep_releases"`
	KeepDays     int               `mapstructure:"keep_days"`
	Regions      []string          `mapstructure:"regions"`
	DryRun       bool              `mapstructure:"dry_run"`
	Tags         map[string]string `mapstructure:"tags"`

	ctx interpolate.Context
}
