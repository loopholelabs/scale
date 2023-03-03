#!/usr/bin/env bash

set -euo pipefail

PATH_TO_SDK="wasi-sdk"
if [[ ! -d $PATH_TO_SDK ]]; then
    TMPGZ=$(mktemp)
    VERSION_MAJOR="12"
    VERSION_MINOR="0"
    if [[ "$(uname -s)" == "Darwin" ]]; then
        curl --fail --location --silent https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-${VERSION_MAJOR}/wasi-sdk-${VERSION_MAJOR}.${VERSION_MINOR}-macos.tar.gz --output $TMPGZ
    else
        curl --fail --location --silent https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-${VERSION_MAJOR}/wasi-sdk-${VERSION_MAJOR}.${VERSION_MINOR}-linux.tar.gz --output $TMPGZ
    fi
    mkdir $PATH_TO_SDK
    tar xf $TMPGZ -C $PATH_TO_SDK --strip-components=1
fi

# NB This doesn't set externally, so we'd need to build here or something...
echo "You should set QUICKJS_WASM_SYS_WASI_SDK_PATH to " `pwd`/${PATH_TO_SDK}