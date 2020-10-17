## v0.9.0 (2020-10-17)

### Bug Fixes

- Fix panic by nil pointer dereference ([#133](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/133))

### Others

- Fix download URL ([#82](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/82))
- Bump github.com/hashicorp/hcl/v2 from 2.3.0 to 2.7.0 ([#93](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/93) [#130](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/130))
- Bump github.com/golang/mock from 1.4.0 to 1.4.4 ([#107](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/107))
- Bump github.com/zclconf/go-cty from 1.2.1 to 1.6.1 ([#115](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/115))
- Bump github.com/hashicorp/packer from 1.5.4 to 1.6.4 ([#121](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/121))
- Bump github.com/aws/aws-sdk-go from 1.29.8 to 1.35.9 ([#125](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/125))
- Switch into GitHub's Dependabot ([#126](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/126))
- Update actions/checkout requirement to v2.3.3 ([#127](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/127))
- Bump actions/setup-go from v1 to v2.1.3 ([#128](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/128))
- Bump actions/cache from v1 to v2.1.2 ([#129](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/129))
- Upgrade Go 1.15 ([#131](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/131))
- Make mock from interface ([#132](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/132))

## v0.8.0 (2020-02-22)

### Enhacements

- Add support for Packer 1.5.4 ([#59](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/59))

### Others

- Stop GitHub Actions ([#48](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/48))
- Fix wrong download URL ([#57](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/57))
- Remove FOSSA ([#63](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/63))
- Retry GitHub Actions ([#50](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/50))
- Bump github.com/golang/mock from 1.3.1 to 1.4.0 ([#61](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/61))
- Bump github.com/aws/aws-sdk-go from 1.24.1 to 1.29.8 ([#62](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/62))
- Run tests on GitHub Actions ([#64](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/64))

## v0.7.0 (2019-08-16)

### Enhancements

- Bump packer and others dependencies ([#42](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/42))

### Others

- Migrate main.workflow to new yaml syntaxes ([#43](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/43))
- Update README for AWS configuration attributes ([#47](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/47))

## v0.6.2 (2019-04-11)

Re-release due to [#151](https://github.com/wata727/packer-post-processor-amazon-ami-management/issues/39). There is no change from v0.6.1. 

### Others

- Fix installation scripts ([#38](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/38))

## v0.6.1 (2019-03-22)

There is no change from v0.6.0. Only changes related to the project, such as documentation.

### Others

- making for easier copy and paste ([#33](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/33))
- Add license scan report and status ([#34](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/34))
- automate the Installation for linux environment ([#35](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/35))
- Add NOTICE.md ([#36](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/36))
- Create releases with GoReleaser and GitHub Action ([#37](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/37))

## v0.6.0 (2019-03-04)

### Changes

- Remove only unused images ([#32](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/32))
  - Previously, the post-processor removes images even if the image is already used elsewhere.
- Change the minimal set permissions necessary ([#32](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/32))
  - Additional API operations are required to check whether if an image is in use. See [here](https://github.com/wata727/packer-post-processor-amazon-ami-management/tree/v0.6.0#iam-task-or-instance-role) for updated permissions.

### Enhancements

- Add `dry_run` option ([#26](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/26))
- Add `keep_days` option ([#31](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/31))

### Others

- Update README.md ([#23](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/23))
- readme: added osx installation instructions ([#28](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/28))
- Go 1.12 ([#29](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/29))

## v0.5.0 (2018-06-16)

### Changes

- Validate post-processor's configuration ([#22](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/22))

### Others

- Add note about `region` parameter ([#20](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/20))
- CI against Go 1.10 ([#21](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/21))

## v0.4.0 (2018-04-14)

### Changes

- Inherit from Packer access config ([#18](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/18))

## v0.3.0 (2018-03-29)

### Enhancements

- Support multiple regions ([#14](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/14))

### Others

- Improve plugin architecture ([#9](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/9))
- Fix lint issues ([#10](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/10))
- Revise docs ([#11](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/11))

## v0.2.0 (2017-05-27)

### Changes

- Use own credentials ([#5](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/5))

## v0.1.1 (2016-08-21)

### Bug Fixes

- Fix panic when use ephemeral devises ([#1](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/1))

### Others

- Fix SideCI issues ([#2](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/2))

## v0.1.0 (2016-08-15)

First release
