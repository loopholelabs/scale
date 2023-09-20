.PHONY: builder interpreter clean
.DEFAULT_GOAL := builder

current_dir := $(shell pwd)
install:
	cargo install --path builder

builder: interpreter
		cargo build --release

builder-x86_64-unknown-linux-gnu:
		cargo build --release --target=x86_64-unknown-linux-gnu && cp target/x86_64-unknown-linux-gnu/release/js_builder js_builder-x86_64-unknown-linux-gnu

builder-aarch64-unknown-linux-gnu:
		cargo build --release --target=aarch64-unknown-linux-gnu && cp target/aarch64-unknown-linux-gnu/release/js_builder js_builder-aarch64-unknown-linux-gnu

builder-x86_64-apple-darwin:
		cargo build --release --target=x86_64-apple-darwin && cp target/x86_64-apple-darwin/release/js_builder js_builder-x86_64-apple-darwin

builder-aarch64-apple-darwin:
		cargo build --release --target=aarch64-apple-darwin && cp target/aarch64-apple-darwin/release/js_builder js_builder-aarch64-apple-darwin

builder-x86_64-pc-windows-msvc:
		cargo build --release --target=x86_64-pc-windows-msvc && cp target/x86_64-pc-windows-msvc/release/js_builder.exe js_builder-x86_64-pc-windows-msvc.exe

interpreter:
		cd interpreter && cargo build --release --target=wasm32-wasi && cd -

interpreter-embedded-source:
		cd interpreter && cargo build --release --target=wasm32-wasi --features embedded_source && cd -

interpreter-runtime-source::
		cd interpreter && cargo build --release --target=wasm32-wasi --features runtime_source && cd -

clean: clean-cargo

clean-cargo:
		cargo clean
