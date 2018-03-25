# packer-post-processor-amazon-ami-management
[![Build Status](https://travis-ci.org/wata727/packer-post-processor-amazon-ami-management.svg?branch=master)](https://travis-ci.org/wata727/packer-post-processor-amazon-ami-management)
[![GitHub release](https://img.shields.io/github/release/wata727/packer-post-processor-amazon-ami-management.svg)](https://github.com/wata727/packer-post-processor-amazon-ami-management/releases/latest)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

Packer post-processor plugin for AMI management

## Description
This post-processor assists your AMI management. It deletes old AMI and EBS snapshot after bake new AMI.

## Installation
Packer supports plugin system. Please read the following documentation:

https://www.packer.io/docs/extend/plugins.html

You can download binary built for your architecture from [latest releases](https://github.com/wata727/packer-post-processor-amazon-ami-management/releases/latest).

## Usage
The following example `template.json`:

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
    "region": "us-east-1",
    "identifier": "packer-example",
    "keep_releases": "3",
    "ami_regions": ["ap-northeast-1"]
  }]
}
```

### Configuration

Type: `amazon-ami-management`

Required:
  - `identifier` (string)
    - An identifier of AMIs. This plugin looks `Amazon_AMI_Management_Identifier` tag. If `identifier` matches tag value, these AMI becomes to management target.
  - `keep_releases` (interger)
    - The number of AMIs.
  - `access_key` (string)
    - The access key used in AWS. If you can use environment values or [shared credentials](https://blogs.aws.amazon.com/security/post/Tx3D6U6WSFGOK2H/A-New-and-Standardized-Way-to-Manage-Credentials-in-the-AWS-SDKs), not required this parameter.
  - `secret_key` (string)
    - The secret key used in AWS. If you can use environment values or [shared credentials](https://blogs.aws.amazon.com/security/post/Tx3D6U6WSFGOK2H/A-New-and-Standardized-Way-to-Manage-Credentials-in-the-AWS-SDKs), not required this parameter.
  - `region` (string)
    - The name of the region, such as `us-east-1` in which to manage AMIs. If you can use environment values, not required this parameter.

Optional:
  - `ami_regions` (array of strings)
    - A list of regions where copied AMI is located. It is useful when specifying `ami_regions` in builder config.

## Developing Plugin

If you wish to build this plugin on your environment, you can use GNU Make build system.
But this Makefile depends on [Go](https://golang.org/) 1.9 or more. At First, you should install Go.

```
$ make build
```
