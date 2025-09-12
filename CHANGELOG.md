## 1.30.5 (Sep 11, 2025)

IMPROVEMENTS:

* Added cronScheduleValidator function using robfig/cron/v3 for Xray reports

* Added GetBoolEnvVar function to support the XRAY_SKIP_VERSION_CHECK environment variable

## 1.30.4 (Aug 16, 2025)

BUG FIXES:

* Fix `ValidateCatalogHealth` to check catalog health

## 1.30.3 (Aug 08, 2025)

IMPROVEMENTS:

* Added func `ValidateCatalogHealth` to check catalog health

## 1.30.2 (July 29, 2025)

IMPROVEMENTS:

* Added func `ValidateXrayConfig` to check xray version supported

## 1.30.1 (June 18, 2025)

IMPROVEMENTS:

* Added `GetAccessVersion` function to retrieve the instance version from Access.

## 1.30.0 (June 2, 2025)

BUG FIXES:

* Fixed a validation failure for cron expression for value "0 0 2 ? * SAT"

## 1.29.0 (February 24, 2025)

IMPROVEMENTS:

* Add `PreApply` stage to config plan checks func

## 1.28.0 (December 20, 2024)

IMPROVEMENTS:

* Add `JFrogProvider` and associates functions for Terraform provider interface

## 1.27.1 (November 22, 2024)

IMPROVEMENTS:

* Upgrade Resty to 2.16.2

## 1.27.0 (November 15, 2024)

IMPROVEMENTS:

* Add `JFrogResource` and `JFrogError` structs

## 1.26.0 (October 15, 2024)

IMPROVEMENTS:

* Add support for multiple TFC Workload Identity Token env var

## 1.25.5 (June 18, 2024)

IMPROVEMENTS:

* Add `RegexNotMatches` Plugin Framework validator

## 1.25.4 (June 4, 2024)

IMPROVEMENTS:

* Move Identifiable() out of `/util/sdk` package path to `/util`

## 1.25.3 (May 24, 2024)

IMPROVEMENTS:

* Add RepoKey validator
* Check resty log for header existence before redacting it

## 1.25.2 (May 10, 2024)

IMPROVEMENTS:

* Add project_key validator for Framework
* Update CheckPlan func to improve text output and filtering

## 1.25.1 (May 2, 2024)

IMPROVEMENTS:

* Switch SendUsage func arg to Resty Request to make gorountine-safe

## 1.25.0 (Apr 30, 2024)

IMPROVEMENTS:

* Add plan checking func for acceptance test

## 1.24.0 (Apr 12, 2024)

IMPROVEMENTS:

* Fix typo in provider metadata type name
* Improve error checking on applyTelemetry()

## 1.23.0 (Apr 11, 2024)

IMPROVEMENTS:

* Add support for exchanging OIDC ID token for access token using JFrog OIDC configuration

PR: [#59](https://github.com/jfrog/terraform-provider-shared/pull/59)

## 1.22.4 (Apr 4, 2024)

IMPROVEMENTS:

* Enable Resty's debug logging when `TF_LOG` is set to `DEBUG` or `TRACE`.

Issue [#16](https://github.com/jfrog/terraform-provider-shared/issues/16)
PR: [#58](https://github.com/jfrog/terraform-provider-shared/pull/58)

## 1.22.3 (Apr 4, 2024)

BUG FIXES:

* Consolidate CheckArtifactoryLicense funcs into single version for both SDKv2 and Framework

PR: [#57](https://github.com/jfrog/terraform-provider-shared/pull/57)

## 1.22.2 (Apr 4, 2024)

BUG FIXES:

* Remove special error code handling in response middleware
* Fix error handling in check license and version func

PR: [#56](https://github.com/jfrog/terraform-provider-shared/pull/56)

## 1.22.1 (Mar 20, 2024)

IMPROVEMENTS:

* Improve `RepoKey` validation to match Artifactory web UI.

PR: [#55](https://github.com/jfrog/terraform-provider-shared/pull/55)

## 1.22.0 (Mar 11, 2024)

IMPROVEMENTS:

* Add `RepoKey` validation

PR: [#53](https://github.com/jfrog/terraform-provider-shared/pull/53)

## 1.21.3 (Jan 3, 2024)

IMPROVEMENTS:

* Upgrade Resty to 2.11.1

PR: [#51](https://github.com/jfrog/terraform-provider-shared/pull/51)

## 1.21.2 (Dec 21, 2023)

IMPROVEMENTS:

* Upgrade Resty to 2.9.1
* Compile using Go 1.21

PR: [#50](https://github.com/jfrog/terraform-provider-shared/pull/50)

## 1.21.1 (Dec 5, 2023)

BUG FIXES:

* Fix incorrect template package for `ExecuteTemplate` function

PR [48](https://github.com/jfrog/terraform-provider-shared/pull/48)

## 1.21.0 (Dec 4, 2023)

IMPROVEMENTS:

* Add `ProductId` to `ProviderMetadata` struct to support Framework resource
* Add send usage funcs for each resource method for Framework resource
* Move utility function `ExecuteTemplate` that are not specific to SDKv2 or Framework to `util.go` so they can be used without pulling in either library.

PR [47](https://github.com/jfrog/terraform-provider-shared/pull/47)

## 1.20.4 (Nov 22, 2023)

IMPROVEMENTS:

* Update Terraform packages.

PR [45](https://github.com/jfrog/terraform-provider-shared/pull/45)

## 1.20.3 (Nov 20, 2023)

IMPROVEMENTS:

* Move utility functions that are not specific to SDKv2 or Framework to `util.go` so they can be used without pulling in either library.

PR [44](https://github.com/jfrog/terraform-provider-shared/pull/44)

## 1.20.2 (Oct 30, 2023)

BUG FIXES:

* Bump google.golang.org/grpc from 1.53.0 to 1.56.3

PR [43](https://github.com/jfrog/terraform-provider-shared/pull/43)

## 1.20.1 (Oct 12, 2023)

BUG FIXES:

* Bump golang.org/x/net from 0.8.0 to 0.17.0

PR [42](https://github.com/jfrog/terraform-provider-shared/pull/42)

## 1.20.0 (Oct 4, 2023)

NEW FEATURES:

* added new string validator to check if string is valid URL.
* re-organize framework string validator files.
* added unit tests for framework string validators.

PR [#41](https://github.com/jfrog/terraform-provider-shared/pull/41)

## 1.19.0 (Sep 28, 2023)

NEW FEATURES:

* added new string validator to check if string is valid cron expression.

PR [#40](https://github.com/jfrog/terraform-provider-shared/pull/40)

## 1.18.0 (Aug 18, 2023)

NEW FEATURES:

* added `CheckXrayVersion` function to get instance version from Xray.

PR [#39](https://github.com/jfrog/terraform-provider-shared/pull/39)

## 1.17.0 (May 30, 2023)

NEW FEATURES:

* added new set validator to check if string is in a pre-defined list of strings.
* added new object validator to check if an attribute is set with value (not unknown or null) when the object is set.

## 1.16.1 (May 15, 2023)

BUG FIXES:

* fixed verification for fw.ValidateBool, values are not compared against each other, but against the value, passed in the function caller.

## 1.16.0 (May 2, 2023)

NEW FEATURES:

* added validator to check if two boolean attributes are set to the same value in Plugin Framework resources.

## 1.15.0 (May 2, 2023)

NEW FEATURES: 

* Added `utilfw.go` file with functions used in the resources, migrated to Terraform Plugin Framework. Moved from SDK to Plugin package.

PR [#33](https://github.com/jfrog/terraform-provider-shared/pull/33)

## 1.14.0 (March 29, 2023)

NEW FEATURES:

* Added `CheckArtifactoryVersion` function to get instance version from Artifactory.

PR [#32](https://github.com/jfrog/terraform-provider-shared/pull/32)

## 1.13.0 (March 28, 2023)

NEW FEATURES:

* Increase the allowed project key length to 32 characters, since Artifactory 7.56.2 expands the maximum length for project key to 32.

PR [#31](https://github.com/jfrog/terraform-provider-shared/pull/31)

## 1.12.0 (March 27, 2023)

NEW FEATURES:

* Added `ProviderMetadata` struct to support passing more enriched metadata from Terraform Provider to resources.
* Added `CheckVersion` function to verify if a version is same or later than the supported version.

Issue [#705](https://github.com/jfrog/terraform-provider-artifactory/issues/705)
PR [#30](https://github.com/jfrog/terraform-provider-shared/pull/30)

## 1.11.1 (March 20, 2023)

BUG FIXES:

* Add nil checking for `CastToStringArr` to avoid panic. PR [#29](https://github.com/jfrog/terraform-provider-shared/pull/29)

## 1.11.0 (February 24, 2023)

BUG FIXES:

* Update regex for `ProjectKey` validator to allow 2-20 characters since 7.49.3. PR [#28](https://github.com/jfrog/terraform-provider-shared/pull/28)

# 1.10.0 (January 23, 2023)

NEW FEATURES:

* added `CheckImportState` function to verify the import state in import tests. PR [#26](https://github.com/jfrog/terraform-provider-shared/pull/26)

## 1.9.0 (January 6, 2023)

NEW FEATURES:

* added validation for the cron expression length to only allow 6-7 parts expressions. PR [#25](https://github.com/jfrog/terraform-provider-shared/pull/25)

## 1.8.0 (December 21, 2022)

BUG FIXES:

* Update regex for `ProjectKey` validator to allow 2-10 characters. PR [#23](https://github.com/jfrog/terraform-provider-shared/pull/23)

## 1.7.0 (August 9, 2022)

BUG FIXES:

* Update regex for `ProjectKey` validator

## 1.6.0 (July 28, 2022)

BUG FIXES:

* Revert changes to size limit for `ProjectKey` validator

## 1.5.0 (July 27, 2022)

BUG FIXES:

* Update string size limit for `ProjectKey` validator

## 1.4.0 (July 1, 2022)

REFACTOR:

* Fix client user agent string was hardcoded to Artifactory

## 1.3.0 (June 14, 2022)

REFACTOR:

* MergeSchema is now MergeMap and has been genericized
* 1 more util function added

## 1.2.0 (June 14, 2022)

REFACTOR:

* revert changes to how field properties are fetched.

## 1.0.1 (June 14, 2022)

REFACTOR:

* `unpack` Moved `unpack` code into this module even though it's still WIP.
* `Identifiable` has been moved into this module as it is shared
* `GetOk` - removed deprecated usage

## 1.0.0 (June 14, 2022)

REFACTOR:

* `util` package stripped down. Predicates moved to `predicate` packages and packers to `packer` package.
* Some remainder of sharable code from artifactory code was moved in.
* `NoPassword` predicate now no longer also includes `NoClass` - done for distinction and clarity

## 0.7.0 (May 11, 2022)

BUG FIXES:

* `util.CheckArtifactoryLicense`: Fix only checking for `Enterprise` license type. Now support a list of license types to check.

## 0.6.0 (May 11, 2022)

NEW FEATURES:

* Add `StringIsNotURL()` to `validator` package

## 0.5.0 (May 4, 2022)

NEW FEATURES:

* Add `util.AddTelemetry()` and `util.CheckArtifactoryLicense()`

## 0.4.0 (May 4, 2022)

NEW FEATURES:

* Add `util.ApplyTelemetry()` and `util.SendUsage()`

## 0.3.0 (May 2, 2022)

NEW FEATURES:

* Add `StringInSlice()` and `IntAtLeast()` to `validator` package

## 0.2.0 (Apr 29, 2022)

NEW FEATURES:

* Add `test.ExecuteTemplate()` and `util.UniversalPack()`

## 0.1.0 (Apr 28, 2022)

NEW FEATURES:

* New packages for all JFrog Terraform providers.
