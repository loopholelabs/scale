#!/bin/bash
set -eo pipefail

# This script is used to compile the wasm modules in the tests directory
# and then run the tests.

# The wasm modules are compiled with the following command:
# tinygo build -o <module name>.wasm -scheduler=none --no-debug -wasm-abi=c -target wasm ./

if [ -d "scale-go-compile" ]
then
    echo "Repo already cloned, deleting and cloning again"
    rm -rf scale-go-compile
fi

# scale-go-compile is the template directory to compile scale functions in go
git clone https://github.com/loopholelabs/scale-go-compile

for compiled in modules/*.wasm ; do
  rm -rf "$compiled"
done

for module in modules/*/*.go ; do
    directory=$(dirname "$module")
    test=$(basename "$directory")
    echo "Compiling $test"
    cp "$module" scale-go-compile/scale/scale.go
    cd scale-go-compile
    tinygo build -o "../modules/$test.wasm" -scheduler=none --no-debug -wasm-abi=c -target wasm ./
    cd ..
done

rm -rf scale-go-compile
