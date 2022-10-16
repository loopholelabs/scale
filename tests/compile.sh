#!/bin/bash
set -eo pipefail

COMPILE="../compile"
TESTS=$(pwd)

# This script is used to compile the wasm modules in the tests directory
# and then run the tests.

# The wasm modules are compiled with the following command:
# tinygo build -o <module name>.wasm -scheduler=none --no-debug -wasm-abi=c -target wasm ./

for compiled in modules/*.wasm ; do
  rm -rf "$compiled"
done

for module in modules/*/*.go ; do
    directory=$(dirname "$module")
    test=$(basename "$directory")
    echo "Compiling $test from $module"
    mv "$COMPILE/scale/scale.go" "$COMPILE/scale/scale.bak"
    cp "$module" $COMPILE/scale/scale.go
    cd $COMPILE
    tinygo build -o "$TESTS/modules/$test.wasm" -scheduler=none --no-debug -wasm-abi=c -target wasm ./
    mv "$COMPILE/scale/scale.bak" "$COMPILE/scale/scale.go"
    cd $TESTS
done
