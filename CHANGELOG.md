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
