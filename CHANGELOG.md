# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.3.4] - 2023-02-17

### Changes

- Bumping `scale-signature` version to `v0.2.9`
- Bumping `scale-signature-http` version to `v0.3.4`
- Bumping `scalefile` version to `v0.1.7`
- Updating `storage` libraries for both TS and Go to use the new `scalefile` library

## [v0.3.3] - 2023-02-17

### Changes

- Changed the implementation of `List` in `go/storage` so that it returns an `[]Entry` which contains the `Organization` and the `Hash` of the scale function
- Changed the implementation of `Get` in `go/storage` and `ts/storage` so they both Return an `Entry`

## [v0.3.2] - 2023-02-16

### Changes

- Fixing bug in `go/storage` where the `List` function was not appending the storage BasePath properly
- Fixing bug in `go/storage` where the `Get` function would return an error if the required scale function was not found (it now returns nil)
- Updating the `golang` compile template's `main.go` file to not have a dependency on the `scale` runtime itself

## [v0.3.1] - 2023-02-15

### Changes

- Adding API Client for both `Golang` and `Typescript`
- Adding Registry Functionality for both `Golang` and `Typescript`
- Adding `scalefile` support for `Typescript`
- Bumping `scale-signature` version to `v0.2.7`
- Bumping `scale-signature-http` version to `v0.3.2`

## [v0.3.0] - 2023-02-14

### Changes

- Bumping `scale-signature` version to `v0.2.2`
- Bumping `scale-signature-http` version to `v0.3.0`
- Bumping `scalefile` version to `v0.1.5`

## [v0.2.2] - 2023-02-03

### Changes

- Removing `wee_alloc` from all dependencies and from the rust generator
- Bumping `scale-signature` version to `v0.2.1`
- Bumping `scale-signature-http` version to `v0.2.4`

## [v0.2.1] - 2023-02-02

### Changes

- Bumping `scale-signature-http` version to `v0.2.3`

## [v0.2.0] - 2023-02-01

### Changes

- Adding Typescript support for the Scale Runtime
- Adding Rust support for the Scale Guest
- Adding support for custom signatures in any Scale Runtime
- Adding test cases to guarantee cross-language compatibility (Go, Rust, Typescript)

## [v0.1.4] - 2023-01-12

### Changes 

- Added `runtime.NewSignature` type to signify a factory function for creating a new `signature.Signature` type.
- Updating the `runtime.New` function to accept a `signature.New` factory function (`runtime.NewSignature`) instead of a `signature.Signature` type. 
- Updated `runtime_test.go` to use the new `runtime.New` function signature.
- Updated the `runtime.Instance` function to make the `next` argument optional. 
- `instance.RuntimeContext` is now a private function since we don't expect developers to use it directly.

## [v0.1.1] - 2022-11-28

### Changes

- Updating https://github.com/loopholelabs/scale-signature to `v0.1.1`
- Adding missing license headers

## [v0.1.0] - 2022-11-25

### Features

- Initial release of the Scale Runtime library.

[unreleased]: https://github.com/loopholelabs/scale/compare/v0.3.4...HEAD
[v0.3.4]: https://github.com/loopholelabs/scale/compare/v0.3.4
[v0.3.3]: https://github.com/loopholelabs/scale/compare/v0.3.3
[v0.3.2]: https://github.com/loopholelabs/scale/compare/v0.3.2
[v0.3.1]: https://github.com/loopholelabs/scale/compare/v0.3.1
[v0.3.0]: https://github.com/loopholelabs/scale/compare/v0.3.0
[v0.2.2]: https://github.com/loopholelabs/scale/compare/v0.2.2
[v0.2.1]: https://github.com/loopholelabs/scale/compare/v0.2.1
[v0.2.0]: https://github.com/loopholelabs/scale/compare/v0.2.0
[v0.1.4]: https://github.com/loopholelabs/scale/compare/v0.1.4
[v0.1.1]: https://github.com/loopholelabs/scale/compare/v0.1.1
[v0.1.0]: https://github.com/loopholelabs/scale/compare/v0.1.0
