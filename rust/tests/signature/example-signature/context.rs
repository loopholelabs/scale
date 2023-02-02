#![cfg(target_arch = "wasm32")]

use crate::example_signature::ExampleContext;

pub type Context = ExampleContext;

pub fn new() -> Context {
    Context {
        data: "".to_string(),
    }
}