# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.3.16] - 2023-04-15

### Changes

- Updating the `Golang` and `Typescript` API Clients 

### Features

- Adding `Typescript` Client Support to the Scale Runtime @jimmyaxod

### Dependencies

- Bumping `wazero` to `v1.0.1`

## [v0.3.15] - 2023-03-19

### Dependencies

- Bumping `scale-signature-http` to `v0.3.8`

## [v0.3.14] - 2023-03-12

### Fixes

- Fix a series of bugs in the JS Compiler that were introduced in `v0.3.13`

### Dependencies

- Bumping `scalefile` to `v0.1.9`
- Bumping `scale-signature` to `v0.2.11`
- Bumping `scale-signature-http` to `v0.3.7`

## [v0.3.13] - 2023-03-12

### Changes

- Updating the `typescript` and `golang` API clients from the latest scale API
- Removing the `webpack` loader (moved to the `scalefile` package)

### Dependencies

- Bumping `scalefile` to `v0.1.8`
- Bumping `scale-signature` to `v0.2.10`
- Bumping `scale-signature-http` to `v0.3.6`

## [v0.3.12] - 2023-02-28

### Features

- Updating the `Typescript` runtime to function properly with `NextJS`
- Updating the `Typescript` registry to function properly in Browser environments
- Adding a `webpack` loader for `Typescript` runtimes that need to import scale functions directly
- Allowing the `Typescript` scale runtime to be instantiated with Promises of `ScaleFunc` objects

### Changes

- Changing the `registry.New` function for both `Typescript` and `Golang` to `registry.Download`

### Dependencies

- Bumping `golang.org/x/net` to `v0.7.0`
- Bumping `golang.org/x/text` to `v0.7.0`

## [v0.3.11] - 2023-02-20

### Features

- Updated API Client to expose the new `DeleteFunction` endpoint from the Scale API

## [v0.3.10] - 2023-02-19

### Fixes

- Fixing bug in `Go` and `TS` Registry implementations where the computed hashes would get encoded in `base64` instead of `hex` leading to incorrect hashes

## [v0.3.9] - 2023-02-19

### Changes

- Bumping `auth` version to `v0.2.26`

### Features

- Added `WithStorage` option to `go/registry` to allow for pre-configured Storage Clients to be used
- Updated API Client to expose the new `UserInfo` endpoint from the Scale API

## [v0.3.8] - 2023-02-19

### Features

- Added `WithClient` option to `go/registry` to allow for pre-configured OpenAPI Clients to be used

## [v0.3.7] - 2023-02-19

### Fixes

- Fixing bug in `Go` Runtime where passing in `nil` as the `Next` function would cause a panic

## [v0.3.6] - 2023-02-19

### Features

- In the `Go` runtime module actions are now cancellable via the given context

### Changes

- Removing `parcel.js` and using the `typescript` compiler directly to build typescript libraries
- Renaming `@loopholelabs/scale-ts` library to `@loopholelabs/scale`
- Bumping `wazero` to `1.0.0-pre.9`

### Fixes

- Fixing bugs in the `DisabledWASI` Polyfill implementation where the proper error codes would not be returned (`fd_write`, `fd_read`, `environ...`, `args...`)
- Fixing bug in the `DisabledWASI` Polyfill implementation where the proper clock time would not get returned
- Making sure the `client`, `registry`, and `storage` typescript packages get exported and packaged properly
- Making sure modules that return an error get thrown away properly instead of being recycled

## [v0.3.5] - 2023-02-17

### Fixes

- Fixing a bug in `go/storage` and `ts/storage` where the entire function path would be used to parse the filename instead of just the base file name

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

[unreleased]: https://github.com/loopholelabs/scale/compare/v0.3.16...HEAD
[v0.3.16]: https://github.com/loopholelabs/scale/compare/v0.3.16
[v0.3.15]: https://github.com/loopholelabs/scale/compare/v0.3.15
[v0.3.14]: https://github.com/loopholelabs/scale/compare/v0.3.14
[v0.3.13]: https://github.com/loopholelabs/scale/compare/v0.3.13
[v0.3.12]: https://github.com/loopholelabs/scale/compare/v0.3.12
[v0.3.11]: https://github.com/loopholelabs/scale/compare/v0.3.11
[v0.3.10]: https://github.com/loopholelabs/scale/compare/v0.3.10
[v0.3.9]: https://github.com/loopholelabs/scale/compare/v0.3.9
[v0.3.8]: https://github.com/loopholelabs/scale/compare/v0.3.8
[v0.3.7]: https://github.com/loopholelabs/scale/compare/v0.3.7
[v0.3.6]: https://github.com/loopholelabs/scale/compare/v0.3.6
[v0.3.5]: https://github.com/loopholelabs/scale/compare/v0.3.5
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
