#![cfg(target_arch = "wasm32")]

use crate::bad_signature::BadContext;

pub type Context = BadContext;

pub fn new() -> Context {
    Context {
        data: 0,
    }
}