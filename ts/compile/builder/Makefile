.PHONY: cli core clean
.DEFAULT_GOAL := cli

current_dir := $(shell pwd)
export QUICKJS_WASM_SYS_WASI_SDK_PATH := $(current_dir)/wasi-sdk

download-wasi-sdk:
	./install-wasi-sdk.sh

install:
	cargo install --path crates/cli

cli: core
		cd crates/cli && cargo build --release && cd -

core:
		cd crates/core \
				&& RUSTFLAGS='-C target-feature=-bulk-memory' cargo build --release --target=wasm32-wasi \
				&& cd -

jssource:
		cd crates/core \
				&& RUSTFLAGS='-C target-feature=-bulk-memory' cargo build --release --target=wasm32-wasi --no-default-features --features js_source --target-dir target_jssource \
				&& cd -

clean: clean-wasi-sdk clean-cargo

clean-cargo:
		cargo clean

clean-wasi-sdk:
		rm -r wasi-sdk 2> /dev/null || true
