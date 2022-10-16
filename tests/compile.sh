#!/bin/bash
#	Copyright 2022 Loophole Labs
#
#	Licensed under the Apache License, Version 2.0 (the "License");
#	you may not use this file except in compliance with the License.
#	You may obtain a copy of the License at
#
#		   http://www.apache.org/licenses/LICENSE-2.0
#
#	Unless required by applicable law or agreed to in writing, software
#	distributed under the License is distributed on an "AS IS" BASIS,
#	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#	See the License for the specific language governing permissions and
#	limitations under the License.
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
