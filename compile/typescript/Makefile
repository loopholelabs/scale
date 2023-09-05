.PHONY: builder interpreter clean
.DEFAULT_GOAL := builder

current_dir := $(shell pwd)
install:
	cargo install --path builder

builder: interpreter
		cd builder && cargo build --release && cd -

interpreter:
		cd interpreter && cargo build --release --target=wasm32-wasi && cd -

interpreter-embedded-source:
		cd interpreter && cargo build --release --target=wasm32-wasi --features embedded_source && cd -

interpreter-runtime-source::
		cd interpreter && cargo build --release --target=wasm32-wasi --features runtime_source && cd -

clean: clean-cargo

clean-cargo:
		cargo clean
