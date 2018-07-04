# Changelog
This file contains a record of all non-trivial changes to kurly

## [Unreleased]

## [1.2.2] 20180704

## Added
* Gitlab CI/CD integration to build and archive the binary and associated files

## Fixed
* Uploading files and Sending request bodies
* Printing the headers of the Redirecting response
* Rebuilding the URL when it has no trailing slash
* Github references to Gitlab

## Removed
* Travis CI support
* snap package support
* rpm build spec file
* Dockerfile

## [1.2.1] 20180312

### Added
* Improved verbosity
* TLS Verbosity
* Support for insecure HTTPS
* Added man page
* Behind-the-scenes refactor for future maintenance
* Handle multiple URLs
* Snap installation via desktop UI

## [1.1.0] 20171229

### Added
* Resume transfer from offset
* Cookie and cookie jar support
* Send HTTP multipart post data
* Installation for Linux via snap package
* Installation for Arch Linux via Arch User Repos

## [1.0.0] 20170501

Initial public release
