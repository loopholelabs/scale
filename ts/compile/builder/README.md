# JS Builder

JSBuilder is a tool for building a .wasm file from some js for use with scale functions.

## Usage

There are three options when it comes to building.

### Build to a wasm with embedded js

```
JS_SOURCE=/path/to/index.js make jssource
cp crates/core/target_jssource/wasm32-wasi/release/jsbuilder_core.wasm index.wasm
```

### With size optimization (Compress js)

```
gzip /path/to/index.js
JS_SOURCE=/path/to/index.js.gz make jssource
cp crates/core/target_jssource/wasm32-wasi/release/jsbuilder_core.wasm index.wasm
```

### With speed optimization for fast cold start

```
make

Should be an executable in `target/release/jsbuilder`

target/release/jsbuilder -o something.wasm index.js
```
