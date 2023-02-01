use crate::example_signature::ExampleContext;

#[cfg(target_arch = "wasm32")]
pub type Context = ExampleContext;

#[cfg(target_arch = "wasm32")]
pub fn new() -> Context {
    Context {
        data: "".to_string(),
    }
}