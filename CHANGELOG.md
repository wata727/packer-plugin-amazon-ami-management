## v0.6.1 (2019-03-22)

There is no change from v0.6.0. Only changes related to the project, such as documentation.

## Others

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
