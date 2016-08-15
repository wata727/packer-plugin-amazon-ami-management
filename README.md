# packer-post-processor-amazon-ami-management
[![Build Status](https://travis-ci.org/wata727/packer-post-processor-amazon-ami-management.svg?branch=master)](https://travis-ci.org/wata727/packer-post-processor-amazon-ami-management)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

Packer post-processor plugin for AMI generation management

## Description
The Packer Amazon AMI Management post-processor assists your AMI generation management.  
You can set `identifer` and `keep_releases`, it works on the basis of AMI tags.  
It delete old AMIs and EBS snapshots related to AMI.

## Installation
Packer supports plugin system. Please read document the following:  
https://www.packer.io/docs/extend/plugins.html

You can download binary built for your architecture from [latest releases](https://github.com/wata727/packer-post-processor-amazon-ami-management/releases/latest).

## Usage
The following example `template.json`:

```
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
    "keep_releases": "3"
  }]
}
```

### Configuration

Type: `amazon-ami-management`

Required:
  - `identifier` (string) 
    - The Identifier of AMIs. This plugin looks `Amazon_AMI_Management_Identifier` tag. If `identifier` matches tag value, these AMI becomes to management target.
  - `keep_releases` (interger)
    - The number of generations.
  - `access_key` (string)
    - The access key used to communicate with AWS. If you can use environment values or [shared credentials](https://blogs.aws.amazon.com/security/post/Tx3D6U6WSFGOK2H/A-New-and-Standardized-Way-to-Manage-Credentials-in-the-AWS-SDKs), not required this parameter.
  - `secret_key` (string)
    - The secret key used to communicate with AWS. If you can use environment values or [shared credentials](https://blogs.aws.amazon.com/security/post/Tx3D6U6WSFGOK2H/A-New-and-Standardized-Way-to-Manage-Credentials-in-the-AWS-SDKs), not required this parameter.
  - `region` (string)
    - The name of the region, such as `us-east-1` in which to manage AMIs.  If you can use environment values, not required this parameter.

## Developing Plugin

If you wish to build this plugin on your environment, you can use GNU Make build system.  
But this Makefile depends on [Go](https://golang.org/). At First, you should install Go.  
And we use [godep](https://github.com/tools/godep) for dependency management. Please looks the [reference](https://godoc.org/github.com/tools/godep)

### Run Test
```
make test
go get github.com/tools/godep
godep restore
go get ./...
go test ./...
?       github.com/wata727/packer-post-processor-amazon-ami-management  [no test files]
?       github.com/wata727/packer-post-processor-amazon-ami-management/awsmock  [no test files]
ok      github.com/wata727/packer-post-processor-amazon-ami-management/plugin   0.029s
```
Running unit tests in developing plugin. You can use `awsmock` package.

### Installation
```
make install
go get github.com/tools/godep
godep restore
go get ./...
go test ./...
?       github.com/wata727/packer-post-processor-amazon-ami-management  [no test files]
?       github.com/wata727/packer-post-processor-amazon-ami-management/awsmock  [no test files]
ok      github.com/wata727/packer-post-processor-amazon-ami-management/plugin   0.023s
go build ./
mkdir -p ~/.packer.d/plugins
install ./packer-post-processor-amazon-ami-management ~/.packer.d/plugins/
```
Run tests, Build and Move to plugin directory.

### Release
```
make release
go get github.com/tools/godep
godep restore
go get ./...
go test ./...
?       github.com/wata727/packer-post-processor-amazon-ami-management  [no test files]
?       github.com/wata727/packer-post-processor-amazon-ami-management/awsmock  [no test files]
ok      github.com/wata727/packer-post-processor-amazon-ami-management/plugin   0.020s
...
go get github.com/mitchellh/gox
gox --output 'dist/{{.OS}}_{{.Arch}}/{{.Dir}}'
...
```
Run tests, Build for each architecture and Archive binaries.