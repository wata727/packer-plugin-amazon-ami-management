# packer-plugin-amazon-ami-management

[![Build Status](https://github.com/wata727/packer-plugin-amazon-ami-management/workflows/build/badge.svg?branch=master)](https://github.com/wata727/packer-plugin-amazon-ami-management/actions)
[![GitHub release](https://img.shields.io/github/release/wata727/packer-plugin-amazon-ami-management.svg)](https://github.com/wata727/packer-plugin-amazon-ami-management/releases/latest)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](LICENSE)

Packer post-processor plugin for Amazon AMI management

## Description

This post-processor cleanups outdated AMIs and EBS snapshots after baking a new AMI.

You can configure the quantity or number of days you want to keep them, and the post-processor will delete all unused AMIs (not used in an instance, launch template, or launch configuration) according to your configuration.

## Installation

Packer >= v1.7.0 supports third-party plugin installation by `init` command. You can install the plugin automatically after adding the `required_plugin` block.

```hcl
packer {
  required_plugins {
    amazon-ami-management = {
      version = ">= 1.0.0"
      source = "github.com/wata727/amazon-ami-management"
    }
  }
}
```

See the [Packer documentation](https://www.packer.io/docs/plugins#installing-plugins) for more details.

## Usage

The following examples is a templates to keep only the latest 3 AMIs.

### An example with defined option `identifier`

```hcl
source "amazon-ebs" "example" {
  region        = "us-east-1"
  source_ami    = "ami-6869aa05"
  instance_type = "t2.micro"
  ssh_username  = "ec2-user"
  ssh_pty       = true
  ami_name      = "packer-example ${formatdate("YYYYMMDDhhmmss", timestamp())}"
  tags = {
    Amazon_AMI_Management_Identifier = "packer-example"
  }
}

build {
  sources = ["source.amazon-ebs.example"]

  provisioner "shell" {
    inline = ["echo 'running...'"]
  }

  post-processor "amazon-ami-management" {
    regions       = ["us-east-1"]
    identifier    = "packer-example"
    keep_releases = 3
  }
}
```

### An example with defined option `tags`

```hcl
locals {
  tags = {
    version    = 1.23
    department = "dev"
  }
}

source "amazon-ebs" "example" {
  region        = "us-east-1"
  source_ami    = "ami-6869aa05"
  instance_type = "t2.micro"
  ssh_username  = "ec2-user"
  ssh_pty       = true
  ami_name      = "packer-example ${formatdate("YYYYMMDDhhmmss", timestamp())}"
  tags          = local.tags
}

build {
  sources = ["source.amazon-ebs.example"]

  provisioner "shell" {
    inline = ["echo 'running...'"]
  }

  post-processor "amazon-ami-management" {
    regions       = ["us-east-1"]
    keep_releases = 3
    tags          = local.tags
  }
}
```

### Configuration

Type: `amazon-ami-management`

Required:

- `identifier` (string) - An identifier of AMIs. This plugin identifies AMIs as managed if the value matches the `Amazon_AMI_Management_Identifier` tag.
- `tags` (map of strings) - The tags to indetify AMI. It can be used when a single `identifier` tag is not sufficient. If `identifier` is set, this parameter is ignored.
- `keep_releases` (integer) - The number of AMIs. This value is invalid when `keep_days` is set.
- `keep_days` (integer) - The number of days to keep AMIs. For example, if you specify `10`, AMIs created before 10 days will be deleted. This value is invalid when `keep_releases` is set.
- `regions` (array of strings) - A list of regions, such as `us-east-1` in which to manage AMIs.

Optional:

- `resolve_aliases` (boolean) - If `true`, the post-processor resolves the AWS Systems Manager parameter when the launch template uses it to specify the AMI ID. See [AWS documentation](https://docs.aws.amazon.com/autoscaling/ec2/userguide/using-systems-manager-parameters.html). **Important**: If you set this to `true`, you must add `ssm:GetParameters` permission to the IAM Role.
- `dry_run` (boolean) - If `true`, the post-processor doesn't actually delete AMIs.

The following attibutes are also available. These are optional and used in the same way as AWS Builder:

- `access_key`
- `secret_key`
- `assume_role`
- `custom_endpoint_ec2`
- `mfa_code`
- `profile`
- `skip_metadata_api_check`
- `token`

### IAM Task or Instance Role

The post-processor requires additional permissions to work. Below is the difference from [the minimum permissions required by Packer](https://www.packer.io/docs/builders/amazon.html#iam-task-or-instance-role).

```diff
{
  "Version": "2012-10-17",
  "Statement": [{
      "Effect": "Allow",
      "Action" : [
+       "autoscaling:DescribeLaunchConfigurations",
        "ec2:AttachVolume",
        "ec2:AuthorizeSecurityGroupIngress",
        "ec2:CopyImage",
        "ec2:CreateImage",
        "ec2:CreateKeypair",
        "ec2:CreateSecurityGroup",
        "ec2:CreateSnapshot",
        "ec2:CreateTags",
        "ec2:CreateVolume",
        "ec2:DeleteKeyPair",
        "ec2:DeleteSecurityGroup",
        "ec2:DeleteSnapshot",
        "ec2:DeleteVolume",
        "ec2:DeregisterImage",
        "ec2:DescribeImageAttribute",
        "ec2:DescribeImages",
        "ec2:DescribeInstances",
        "ec2:DescribeInstanceStatus",
+       "ec2:DescribeLaunchTemplates",
+       "ec2:DescribeLaunchTemplateVersions",
        "ec2:DescribeRegions",
        "ec2:DescribeSecurityGroups",
        "ec2:DescribeSnapshots",
        "ec2:DescribeSubnets",
        "ec2:DescribeTags",
        "ec2:DescribeVolumes",
        "ec2:DetachVolume",
        "ec2:GetPasswordData",
        "ec2:ModifyImageAttribute",
        "ec2:ModifyInstanceAttribute",
        "ec2:ModifySnapshotAttribute",
        "ec2:RegisterImage",
        "ec2:RunInstances",
        "ec2:StopInstances",
        "ec2:TerminateInstances",
+       "ssm:GetParameters" // If "resolve_aliases" is enabled
      ],
      "Resource" : "*"
  }]
}
```

## Developing Plugin

To use the plugin built locally with Packer, you can use `make install`.

```
$ make install
```

This command runs `go build` to generate the plugin binary and then installs the plugin with `packer plugins install`. This requires that you have Go v1.23+ and Packer v1.7+ installed.
