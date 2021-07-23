// Code generated by "mapstructure-to-hcl2 -type Config,AssumeRoleConfig"; DO NOT EDIT.
package main

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// FlatAssumeRoleConfig is an auto-generated flat version of AssumeRoleConfig.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatAssumeRoleConfig struct {
	AssumeRoleARN               *string           `mapstructure:"role_arn" required:"false" cty:"role_arn"`
	AssumeRoleDurationSeconds   *int              `mapstructure:"duration_seconds" required:"false" cty:"duration_seconds"`
	AssumeRoleExternalID        *string           `mapstructure:"external_id" required:"false" cty:"external_id"`
	AssumeRolePolicy            *string           `mapstructure:"policy" required:"false" cty:"policy"`
	AssumeRolePolicyARNs        []string          `mapstructure:"policy_arns" required:"false" cty:"policy_arns"`
	AssumeRoleSessionName       *string           `mapstructure:"session_name" required:"false" cty:"session_name"`
	AssumeRoleTags              map[string]string `mapstructure:"tags" required:"false" cty:"tags"`
	AssumeRoleTransitiveTagKeys []string          `mapstructure:"transitive_tag_keys" required:"false" cty:"transitive_tag_keys"`
}

// FlatMapstructure returns a new FlatAssumeRoleConfig.
// FlatAssumeRoleConfig is an auto-generated flat version of AssumeRoleConfig.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*AssumeRoleConfig) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatAssumeRoleConfig)
}

// HCL2Spec returns the hcl spec of a AssumeRoleConfig.
// This spec is used by HCL to read the fields of AssumeRoleConfig.
// The decoded values from this spec will then be applied to a FlatAssumeRoleConfig.
func (*FlatAssumeRoleConfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"role_arn":            &hcldec.AttrSpec{Name: "role_arn", Type: cty.String, Required: false},
		"duration_seconds":    &hcldec.AttrSpec{Name: "duration_seconds", Type: cty.Number, Required: false},
		"external_id":         &hcldec.AttrSpec{Name: "external_id", Type: cty.String, Required: false},
		"policy":              &hcldec.AttrSpec{Name: "policy", Type: cty.String, Required: false},
		"policy_arns":         &hcldec.AttrSpec{Name: "policy_arns", Type: cty.List(cty.String), Required: false},
		"session_name":        &hcldec.AttrSpec{Name: "session_name", Type: cty.String, Required: false},
		"tags":                &hcldec.BlockAttrsSpec{TypeName: "tags", ElementType: cty.String, Required: false},
		"transitive_tag_keys": &hcldec.AttrSpec{Name: "transitive_tag_keys", Type: cty.List(cty.String), Required: false},
	}
	return s
}

// FlatConfig is an auto-generated flat version of Config.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatConfig struct {
	PackerBuildName      *string               `mapstructure:"packer_build_name" cty:"packer_build_name"`
	PackerBuilderType    *string               `mapstructure:"packer_builder_type" cty:"packer_builder_type"`
	PackerCoreVersion    *string               `mapstructure:"packer_core_version" cty:"packer_core_version"`
	PackerDebug          *bool                 `mapstructure:"packer_debug" cty:"packer_debug"`
	PackerForce          *bool                 `mapstructure:"packer_force" cty:"packer_force"`
	PackerOnError        *string               `mapstructure:"packer_on_error" cty:"packer_on_error"`
	PackerUserVars       map[string]string     `mapstructure:"packer_user_variables" cty:"packer_user_variables"`
	PackerSensitiveVars  []string              `mapstructure:"packer_sensitive_variables" cty:"packer_sensitive_variables"`
	AccessKey            *string               `mapstructure:"access_key" cty:"access_key"`
	AssumeRole           *FlatAssumeRoleConfig `mapstructure:"assume_role" required:"false" cty:"assume_role"`
	CustomEndpointEc2    *string               `mapstructure:"custom_endpoint_ec2" cty:"custom_endpoint_ec2"`
	MFACode              *string               `mapstructure:"mfa_code" cty:"mfa_code"`
	ProfileName          *string               `mapstructure:"profile" cty:"profile"`
	SecretKey            *string               `mapstructure:"secret_key" cty:"secret_key"`
	SkipMetadataAPICheck *bool                 `mapstructure:"skip_metadata_api_check" cty:"skip_metadata_api_check"`
	Token                *string               `mapstructure:"token" cty:"token"`
	SkipValidation       *bool                 `mapstructure:"skip_region_validation" cty:"skip_region_validation"`
	Identifier           *string               `mapstructure:"identifier" cty:"identifier"`
	KeepReleases         *int                  `mapstructure:"keep_releases" cty:"keep_releases"`
	KeepDays             *int                  `mapstructure:"keep_days" cty:"keep_days"`
	Regions              []string              `mapstructure:"regions" cty:"regions"`
	DryRun               *bool                 `mapstructure:"dry_run" cty:"dry_run"`
}

// FlatMapstructure returns a new FlatConfig.
// FlatConfig is an auto-generated flat version of Config.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*Config) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatConfig)
}

// HCL2Spec returns the hcl spec of a Config.
// This spec is used by HCL to read the fields of Config.
// The decoded values from this spec will then be applied to a FlatConfig.
func (*FlatConfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"packer_build_name":          &hcldec.AttrSpec{Name: "packer_build_name", Type: cty.String, Required: false},
		"packer_builder_type":        &hcldec.AttrSpec{Name: "packer_builder_type", Type: cty.String, Required: false},
		"packer_core_version":        &hcldec.AttrSpec{Name: "packer_core_version", Type: cty.String, Required: false},
		"packer_debug":               &hcldec.AttrSpec{Name: "packer_debug", Type: cty.Bool, Required: false},
		"packer_force":               &hcldec.AttrSpec{Name: "packer_force", Type: cty.Bool, Required: false},
		"packer_on_error":            &hcldec.AttrSpec{Name: "packer_on_error", Type: cty.String, Required: false},
		"packer_user_variables":      &hcldec.BlockAttrsSpec{TypeName: "packer_user_variables", ElementType: cty.String, Required: false},
		"packer_sensitive_variables": &hcldec.AttrSpec{Name: "packer_sensitive_variables", Type: cty.List(cty.String), Required: false},
		"access_key":                 &hcldec.AttrSpec{Name: "access_key", Type: cty.String, Required: false},
		"assume_role":                &hcldec.BlockSpec{TypeName: "assume_role", Nested: hcldec.ObjectSpec((*FlatAssumeRoleConfig)(nil).HCL2Spec())},
		"custom_endpoint_ec2":        &hcldec.AttrSpec{Name: "custom_endpoint_ec2", Type: cty.String, Required: false},
		"mfa_code":                   &hcldec.AttrSpec{Name: "mfa_code", Type: cty.String, Required: false},
		"profile":                    &hcldec.AttrSpec{Name: "profile", Type: cty.String, Required: false},
		"secret_key":                 &hcldec.AttrSpec{Name: "secret_key", Type: cty.String, Required: false},
		"skip_metadata_api_check":    &hcldec.AttrSpec{Name: "skip_metadata_api_check", Type: cty.Bool, Required: false},
		"token":                      &hcldec.AttrSpec{Name: "token", Type: cty.String, Required: false},
		"skip_region_validation":     &hcldec.AttrSpec{Name: "skip_region_validation", Type: cty.Bool, Required: false},
		"identifier":                 &hcldec.AttrSpec{Name: "identifier", Type: cty.String, Required: false},
		"keep_releases":              &hcldec.AttrSpec{Name: "keep_releases", Type: cty.Number, Required: false},
		"keep_days":                  &hcldec.AttrSpec{Name: "keep_days", Type: cty.Number, Required: false},
		"regions":                    &hcldec.AttrSpec{Name: "regions", Type: cty.List(cty.String), Required: false},
		"dry_run":                    &hcldec.AttrSpec{Name: "dry_run", Type: cty.Bool, Required: false},
	}
	return s
}
