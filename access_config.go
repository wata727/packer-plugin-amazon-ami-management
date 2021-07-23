// @see https://github.com/hashicorp/packer-plugin-amazon/blob/v1.0.0/builder/common/access_config.go

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsCredentials "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awsbase "github.com/hashicorp/aws-sdk-go-base"
	"github.com/hashicorp/go-cleanhttp"
)

// AssumeRoleConfig lets users set configuration options for assuming a special
// role when executing this plugin.
type AssumeRoleConfig struct {
	AssumeRoleARN               string            `mapstructure:"role_arn" required:"false"`
	AssumeRoleDurationSeconds   int               `mapstructure:"duration_seconds" required:"false"`
	AssumeRoleExternalID        string            `mapstructure:"external_id" required:"false"`
	AssumeRolePolicy            string            `mapstructure:"policy" required:"false"`
	AssumeRolePolicyARNs        []string          `mapstructure:"policy_arns" required:"false"`
	AssumeRoleSessionName       string            `mapstructure:"session_name" required:"false"`
	AssumeRoleTags              map[string]string `mapstructure:"tags" required:"false"`
	AssumeRoleTransitiveTagKeys []string          `mapstructure:"transitive_tag_keys" required:"false"`
}

// AccessConfig is for common configuration related to AWS access
type AccessConfig struct {
	AccessKey            string           `mapstructure:"access_key"`
	AssumeRole           AssumeRoleConfig `mapstructure:"assume_role" required:"false"`
	CustomEndpointEc2    string           `mapstructure:"custom_endpoint_ec2"`
	MFACode              string           `mapstructure:"mfa_code"`
	ProfileName          string           `mapstructure:"profile"`
	SecretKey            string           `mapstructure:"secret_key"`
	SkipMetadataAPICheck bool             `mapstructure:"skip_metadata_api_check"`
	Token                string           `mapstructure:"token"`

	// SkipValidation is not used, but it is still a valid option to keep backward compatibility.
	SkipValidation bool `mapstructure:"skip_region_validation"`

	session *session.Session
}

// Session returns a valid session.Session object for access to AWS services, or
// an error if the authentication and region couldn't be resolved
func (c *AccessConfig) Session() (*session.Session, error) {
	if c.session != nil {
		return c.session, nil
	}

	// Create new AWS config
	config := aws.NewConfig().WithCredentialsChainVerboseErrors(true)

	if c.CustomEndpointEc2 != "" {
		config = config.WithEndpoint(c.CustomEndpointEc2)
	}

	config = config.WithHTTPClient(cleanhttp.DefaultClient())
	transport := config.HTTPClient.Transport.(*http.Transport)
	transport.Proxy = http.ProxyFromEnvironment

	// Figure out which possible credential providers are valid; test that we
	// can get credentials via the selected providers, and set the providers in
	// the config.
	creds, err := c.GetCredentials(config)
	if err != nil {
		return nil, err
	}
	config.WithCredentials(creds)

	// Create session options based on our AWS config
	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            *config,
	}

	if c.ProfileName != "" {
		opts.Profile = c.ProfileName
	}

	if c.MFACode != "" {
		opts.AssumeRoleTokenProvider = func() (string, error) {
			return c.MFACode, nil
		}
	}

	sess, err := session.NewSessionWithOptions(opts)
	if err != nil {
		return nil, err
	}
	log.Printf("Found region %s", *sess.Config.Region)
	c.session = sess

	cp, err := c.session.Config.Credentials.Get()

	if IsAWSErr(err, "NoCredentialProviders", "") {
		return nil, c.NewNoValidCredentialSourcesError(err)
	}

	if err != nil {
		return nil, fmt.Errorf("Error loading credentials for AWS Provider: %s", err)
	}

	log.Printf("[INFO] AWS Auth provider used: %q", cp.ProviderName)

	return c.session, nil
}

// GetCredentials gets credentials from the environment, shared credentials,
// the session (which may include a credential process), or ECS/EC2 metadata
// endpoints. GetCredentials also validates the credentials and the ability to
// assume a role or will return an error if unsuccessful.
func (c *AccessConfig) GetCredentials(config *aws.Config) (*awsCredentials.Credentials, error) {
	// Reload values into the config used by the Packer-Terraform shared SDK
	awsbaseConfig := &awsbase.Config{
		AccessKey:                   c.AccessKey,
		AssumeRoleARN:               c.AssumeRole.AssumeRoleARN,
		AssumeRoleDurationSeconds:   c.AssumeRole.AssumeRoleDurationSeconds,
		AssumeRoleExternalID:        c.AssumeRole.AssumeRoleExternalID,
		AssumeRolePolicy:            c.AssumeRole.AssumeRolePolicy,
		AssumeRolePolicyARNs:        c.AssumeRole.AssumeRolePolicyARNs,
		AssumeRoleSessionName:       c.AssumeRole.AssumeRoleSessionName,
		AssumeRoleTags:              c.AssumeRole.AssumeRoleTags,
		AssumeRoleTransitiveTagKeys: c.AssumeRole.AssumeRoleTransitiveTagKeys,
		DebugLogging:                false,
		Profile:                     c.ProfileName,
		SecretKey:                   c.SecretKey,
		SkipMetadataApiCheck:        c.SkipMetadataAPICheck,
		Token:                       c.Token,
	}

	return awsbase.GetCredentials(awsbaseConfig)
}

// IsAWSErr returns true if the error matches all these conditions:
//  * err is of type awserr.Error
//  * Error.Code() matches code
//  * Error.Message() contains message
func IsAWSErr(err error, code string, message string) bool {
	if err, ok := err.(awserr.Error); ok {
		return err.Code() == code && strings.Contains(err.Message(), message)
	}
	return false
}

// NewNoValidCredentialSourcesError returns user-friendly errors for authentication failed.
func (c *AccessConfig) NewNoValidCredentialSourcesError(err error) error {
	return fmt.Errorf("No valid credential sources found for amazon-ami-management post processor. "+
		"Please see https://github.com/wata727/packer-plugin-amazon-ami-management "+
		"for more information on providing credentials for the post processor. "+
		"Error: %w", err)
}
