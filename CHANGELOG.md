# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changes 

- Added `runtime.NewSignature` type to signify a factory function for creating a new `signature.Signature` type.
- Updating the `runtime.New` function to accept a `signature.New` factory function (`runtime.NewSignature`) instead of a `signature.Signature` type. 
- Updated `runtime_test.go` to use the new `runtime.New` function signature.
- Updated the `runtime.Instance` function to make the `next` argument optional. 

## [v0.1.1] - 2022-11-28

### Changes

- Updating https://github.com/loopholelabs/scale-signature to `v0.1.1`
- Adding missing license headers

## [v0.1.0] - 2022-11-25

### Features

- Initial release of the Scale Runtime library.

[unreleased]: https://github.com/loopholelabs/scale/compare/v0.1.1...HEAD
[v0.1.1]: https://github.com/loopholelabs/scale/compare/v0.1.1
[v0.1.0]: https://github.com/loopholelabs/scale/compare/v0.1.0
