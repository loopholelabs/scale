#!/bin/bash
set -eo pipefail

for compiled in modules/*.wasm ; do
  rm -rf "$compiled"
done