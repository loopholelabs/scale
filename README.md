<br/>
<div align="center">
  <a href="https://scale.sh">
    <img src="docs/logo/dark.svg" alt="Logo" height="90">
  </a>
  <h3 align="center">
    A framework for building high-performance plugin systems into any application, all powered by WebAssembly.
  </h3>

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-brightgreen.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Discord](https://dcbadge.vercel.app/api/server/JYmFhtdPeu?style=flat)](https://loopholelabs.io/discord)
</div>

With [Scale Functions](https://scale.sh) your users can write fully typed plugins in any language they choose, and your application can easily and safely
run those plugins with the Scale Runtime, which provides state-of-the-art sandboxing, low startup times, and extremely high performance.

Currently, guest plugins can be written in [Golang](https://golang.org), [Rust](https://www.rust-lang.org/), and [Typescript](https://www.typescriptlang.org/), with the Runtime supporting [Golang](https://golang.org) and [Typescript](https://www.typescriptlang.org/) host applications.

## Getting Started

First, install the [CLI](https://scale.sh/docs/getting-started/quick-start#install-the-scale-cli).

Create a new function, passing `<name>:<tag>` to the `new` command:

```bash
scale new hello:1.0
```

The following files will be generated:

```yaml
version: v1alpha
name: hello
tag: 1.0
signature: http@v0.3.4
language: go
dependencies:
- name: github.com/loopholelabs/scale-signature
version: v0.2.9
- name: github.com/loopholelabs/scale-signature-http
version: v0.3.4
source: scale.go
```

```Go
//go:build tinygo || js || wasm
package scale

import (
    signature "github.com/loopholelabs/scale-signature-http"
)

func Scale(ctx *signature.Context) (*signature.Context, error) {
    ctx.Response().SetBody("Hello, World!")
    return ctx.Next()
}
```

```Go
module scale

go 1.18

require github.com/loopholelabs/scale-signature v0.2.9
require github.com/loopholelabs/scale-signature-http v0.3.4
```

For more information on these files, see the full [Quick Start Guide](https://scale.sh/docs/getting-started/quick-start#create-a-new-function).

Build the above function:

```bash
scale function build
```
And run:

```bash
scale run local/hello:1.0
```

This will start a local HTTP server on port `8080` and will run the function whenever you make a request to it.

Et Voilà! Your first Scale Function! 🎉

----

Functions be [chained together](https://scale.sh/docs/languages/golang/overview#guest-support), [embedded in other language's apps](https://scale.sh/docs/languages/javascript-typescript/overview#embedding-scale-functions), and used independently. For more information, as well as usage with other supported
language, including Rust and TypeScript/JavaScript, see the [documentation](https://scale.sh/docs).

## Documentation

Full instructions and documentation for Scale is available at [https://scale.sh/docs](https://scale.sh/docs).

## Contributing

Bug reports and pull requests are welcome on GitHub at [https://github.com/loopholelabs/scale][gitrepo]. For more
contribution information check
out [the contribution guide](https://github.com/loopholelabs/scale/blob/main/CONTRIBUTING.md).

## License

The Scale project is available as open source under the terms of
the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).

## Code of Conduct

Everyone interacting in the Scale project’s codebases, issue trackers, chat rooms and mailing lists is expected to follow the [CNCF Code of Conduct](https://github.com/cncf/foundation/blob/main/code-of-conduct.md).

## Project Managed By:

[![https://loopholelabs.io][loopholelabs]](https://loopholelabs.io)

[gitrepo]: https://github.com/loopholelabs/scale
[loopholelabs]: https://cdn.loopholelabs.io/loopholelabs/LoopholeLabsLogo.svg
[loophomepage]: https://loopholelabs.io

