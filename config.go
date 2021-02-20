//go:generate mapstructure-to-hcl2 -type Config

package main

import (
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/template/interpolate"
)

// Config is a post-processor's configuration
// PostProcessor generates it using Packer's configuration in `Configure()` method
type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	AccessConfig        `mapstructure:",squash"`

	Identifier   string   `mapstructure:"identifier"`
	KeepReleases int      `mapstructure:"keep_releases"`
	KeepDays     int      `mapstructure:"keep_days"`
	Regions      []string `mapstructure:"regions"`
	DryRun       bool     `mapstructure:"dry_run"`

	ctx interpolate.Context
}
