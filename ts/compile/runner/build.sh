#!/bin/bash

# Build the typescript targets
npm run build

# Now convert the js to wasm modules...


for TARGET in module_middleware module_endpoint module_error
do
#TARGET=module_middleware

# Find the nomodule version of the compiled js
JS=`cat dist/${TARGET}/index.html | tr ">" "\n" | grep nomodule | awk -F'"' '{print $2}'`

# Use the builder to build a wasm module
BUILDER=../builder/target/release/jsbuilder

${BUILDER} dist/${TARGET}${JS} -o wasm/${TARGET}_opt.wasm

SRC=`realpath dist/${TARGET}${JS}`

cat ${SRC} | gzip > ${SRC}.gz

# Build a non-optimized version
#cp dist/${TARGET}${JS} ../host/crates/core/src/index.js
cd ../builder
rm -rf crates/core/target_jssource
JS_SOURCE=${SRC} make jssource
cp crates/core/target_jssource/wasm32-wasi/release/jsbuilder_core.wasm ../runner/wasm/${TARGET}.wasm
rm -rf crates/core/target_jssource
JS_SOURCE=${SRC}.gz make jssource
cp crates/core/target_jssource/wasm32-wasi/release/jsbuilder_core.wasm ../runner/wasm/${TARGET}_gz.wasm

cd ../runner
done
