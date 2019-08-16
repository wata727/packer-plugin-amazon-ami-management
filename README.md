# packer-post-processor-amazon-ami-management
[![Build Status](https://travis-ci.org/wata727/packer-post-processor-amazon-ami-management.svg?branch=master)](https://travis-ci.org/wata727/packer-post-processor-amazon-ami-management)
[![GitHub release](https://img.shields.io/github/release/wata727/packer-post-processor-amazon-ami-management.svg)](https://github.com/wata727/packer-post-processor-amazon-ami-management/releases/latest)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fwata727%2Fpacker-post-processor-amazon-ami-management.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fwata727%2Fpacker-post-processor-amazon-ami-management?ref=badge_shield)

Packer post-processor plugin for AMI management

## Description
This post-processor cleanups old AMIs and EBS snapshots using `amazon-ebs` builder's access configuration after baking a new AMI.

## Installation
Packer supports plugin system. Please read the following documentation:

https://www.packer.io/docs/extend/plugins.html

You can download binary built for your architecture from [latest releases](https://github.com/wata727/packer-post-processor-amazon-ami-management/releases/latest).

For example, to install v0.7.0 for 64bit OSX

For Linux based OS, you can use the install_linux.sh to automate the installation process

```sh
mkdir -p ~/.packer.d/plugins
wget https://github.com/wata727/packer-post-processor-amazon-ami-management/releases/download/v0.6.2/packer-post-processor-amazon-ami-management_0.7.0_darwin_amd64.zip -P /tmp/
cd ~/.packer.d/plugins
unzip -j /tmp/packer-post-processor-amazon-ami-management_0.7.0_darwin_amd64.zip -d ~/.packer.d/plugins
```

## Usage
The following example is a template to keep only the latest 3 AMIs.

```json
{
  "builders": [{
    "type": "amazon-ebs",
    "region": "us-east-1",
    "source_ami": "ami-6869aa05",
    "instance_type": "t2.micro",
    "ssh_username": "ec2-user",
    "ssh_pty": "true",
    "ami_name": "packer-example {{timestamp}}",
    "tags": {
        "Amazon_AMI_Management_Identifier": "packer-example"
    }
  }],
  "provisioners":[{
    "type": "shell",
    "inline": [
      "echo 'running...'"
    ]
  }],
  "post-processors":[{
    "type": "amazon-ami-management",
    "regions": ["us-east-1"],
    "identifier": "packer-example",
    "keep_releases": "3"
  }]
}
```

### Configuration

Type: `amazon-ami-management`

Required:
  - `identifier` (string) - An identifier of AMIs. This plugin looks `Amazon_AMI_Management_Identifier` tag. If `identifier` matches tag value, these AMI becomes to management target.
  - `keep_releases` (integer) - The number of AMIs. This value is invalid when `keep_days` is set.
  - `keep_days` (integer) - The number of days to keep AMIs. For example, if you specify `10`, AMIs created before 10 days will be deleted. This value is invalid when `keep_releases` is set.
  - `regions` (array of strings) - A list of regions, such as `us-east-1` in which to manage AMIs. **NOTE:** Before v0.3.0, this parameter was `region`. Since 0.4.0, `region` is not used.

Optional:
  - `dry_run` (boolean) - If `true`, the post-processor doesn't actually delete AMIs.

The following attibutes are also available. These are optional and used in the same way as AWS Builder:

- `access_key`
- `secret_key`
- `profile`
- `token`
- `mfa_code`
- `custom_endpoint_ec2`
- `skip_region_validation`
- `skip_metadata_api_check`

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
        "ec2:TerminateInstances"
      ],
      "Resource" : "*"
  }]
}
```

## Developing Plugin

If you wish to build this plugin on your environment, you can use GNU Make build system.
But this Makefile depends on [Go](https://golang.org/) 1.12 or more. At First, you should install Go.

```
$ GO111MODULE=on make build
```


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fwata727%2Fpacker-post-processor-amazon-ami-management.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fwata727%2Fpacker-post-processor-amazon-ami-management?ref=badge_large)