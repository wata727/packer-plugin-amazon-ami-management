## v1.2.0 (2021-07-23)

### Enhancements

- [#207](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/207): Add support for `assume_role` access config

### Chores

- [#192](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/192): Bump actions/cache from 2.1.5 to 2.1.6
- [#196](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/196): Bump github.com/hashicorp/packer-plugin-sdk from 0.2.0 to 0.2.3
- [#197](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/197): Bump github.com/golang/mock from 1.5.0 to 1.6.0
- [#203](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/203): Bump github.com/zclconf/go-cty from 1.8.2 to 1.9.0
- [#206](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/206): Bump github.com/aws/aws-sdk-go from 1.38.25 to 1.40.5

## v1.1.2 (2021-04-25)

- [#184](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/184): build: Fix installation command for Go 1.16 style

## v1.1.1 (2021-04-25)

- [#183](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/183): build: Fix Go version in release action

## v1.1.0 (2021-04-25)

### Changes

- [#174](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/174): Upgrade to Go 1.16
  - darwin/arm64 is now available. See also https://golang.org/doc/go1.16

### Chores

- [#173](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/173): Remove support for the installation script
- [#175](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/175): Bump actions/cache from v2.1.4 to v2.1.5
- [#176](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/176): Bump github.com/hashicorp/packer-plugin-sdk from 0.1.0 to 0.2.0
- [#178](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/178): Bump github.com/hashicorp/hcl/v2 from 2.9.1 to 2.10.0
- [#179](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/179): Bump github.com/golang/mock from 1.4.4 to 1.5.0
- [#180](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/180): Bump github.com/hashicorp/aws-sdk-go-base from 0.6.0 to 0.7.1
- [#181](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/181): Bump github.com/aws/aws-sdk-go from 1.38.0 to 1.38.25
- [#182](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/182): Bump github.com/zclconf/go-cty from 1.8.1 to 1.8.2

## v1.0.0 (2021-02-22)

This release contains some major changes for Packer v1.7 support. If you want to use Packer < v1.7, please use v0.x versions.

### Breaking Changes

- [#165](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/165): Remove the AWS builder dependency
  - `clean_resource_name` function support was removed from the post-processor attributes.
  - Undocumented AWS access config attributes were removed.
    - `assume_role`
    - `shared_credentials_file`
    - `decode_authorization_messages`
    - `insecure_skip_tls_verify`
    - `max_retries`
    - `region`
    - `skip_credential_validation`
    - `vault_aws_engine`
    - `aws_polling`
- [#166](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/166): Change license: MIT -> MPL 2.0
  - Changed to meet licensing requirements due to porting code from Packer core.
- [#167](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/167): Upgrade the plugin to be compatible with Packer v1.7
  - Drop support for Packer < v1.7
- [#169](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/169) [#171](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/171): Make the plugin multi component plugin
  - Support automatic installation by `packer init`.
  - Rename to `packer-plugin-amazon-ami-management` from `packer-post-processor-amazon-ami-management`.
  - Drop pre-built binary support for netbsd/openbsd.
  - Add pre-built binary support for arm64.

### Chores

- [#139](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/139): Bump actions/checkout from v2.3.3 to v2.3.4
- [#152](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/152): Bump github.com/zclconf/go-cty from 1.6.1 to 1.7.1
- [#156](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/156): Bump github.com/hashicorp/hcl/v2 from 2.7.0 to 2.8.2
- [#160](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/160): Bump actions/cache from v2.1.2 to v2.1.4
- [#170](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/170): Bump github.com/aws/aws-sdk-go from v1.36.5 to 1.37.15
- [#172](https://github.com/wata727/packer-plugin-amazon-ami-management/pull/172): Small refactoring

## v0.9.0 (2020-10-17)

### Changes

- Upgrade Go 1.15 ([#131](https://github.com/wata727/packer-post-processor-amazon-ami-management/pull/131))
  - darwin/386 build will no longer available from the release. See also https://golang.org/doc/go1.15#darwin

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
