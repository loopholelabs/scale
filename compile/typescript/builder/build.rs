/*
    Copyright 2022 Loophole Labs

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

           http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

use std::env;
use std::fs;
use std::path::PathBuf;

fn main() {
    if let Ok("cargo-clippy") = env::var("CARGO_CFG_FEATURE").as_ref().map(String::as_str) {
        stub_engine_for_clippy();
    } else {
        copy_engine_binary();
    }
}

// When using clippy, we need to write a stubbed engine.wasm file to ensure compilation succeeds. This
// skips building the actual engine.wasm binary that would be injected into the CLI binary.
fn stub_engine_for_clippy() {
    let engine_path = PathBuf::from(env::var("OUT_DIR").unwrap()).join("interpreter.wasm");
    if !engine_path.exists() {
        fs::write(engine_path, []).expect("failed to write empty interpreter.wasm stub");
        println!("cargo:warning=using stubbed interpreter.wasm for static analysis purposes...");
    }
}

// Copy the engine binary build from the `core` crate
fn copy_engine_binary() {
    let input_path = env::var("JS_INTERPRETER_PATH")
        .unwrap_or("interpreter/target/wasm32-wasi/release/js_interpreter.wasm".into());
    let output_path = format!("{}/interpreter.wasm", env::var("OUT_DIR").unwrap());
    println!("cargo:rerun-if-changed={:?}", input_path);
    println!("cargo:rerun-if-changed=build.rs");
    println!(
        "cargo:warning=using js_interpreter.wasm from {:?}...",
        input_path
    );
    fs::copy(&input_path, output_path).expect("failed to copy js_interpreter.wasm");
}
