#![cfg(target_arch = "wasm32")]

use crate::context::Context;
use crate::example_signature::{Decode, Encode, ExampleContext};
use scale_signature::{Context as ContextTrait, GuestContext as GuestContextTrait};
use std::io::{Cursor, Error, ErrorKind};

pub static mut READ_BUFFER: Vec<u8> = Vec::new();
pub static mut WRITE_BUFFER: Vec<u8> = Vec::new();

pub type GuestContext = Context;

impl ContextTrait for Context {
    fn guest_context(&mut self) -> &mut dyn GuestContextTrait {
        self
    }
}

impl GuestContextTrait for GuestContext {
    unsafe fn to_write_buffer(&mut self) -> (u32, u32) {
        let mut cursor = Cursor::new(Vec::new());
        cursor = match ExampleContext::encode(self.clone(), &mut cursor) {
            Ok(_) => cursor,
            Err(err) => return self.error_write_buffer(&err.to_string()),
        };

        let vec = cursor.into_inner();

        let existing_cap = WRITE_BUFFER.capacity() as usize;
        WRITE_BUFFER.reserve_exact((vec.len() - existing_cap) as usize);

        WRITE_BUFFER.clone_from(&vec);
        WRITE_BUFFER.shrink_to_fit();

        return (WRITE_BUFFER.as_ptr() as u32, WRITE_BUFFER.len() as u32);
    }

    unsafe fn error_write_buffer(&mut self, error: &str) -> (u32, u32) {
        let mut cursor = Cursor::new(Vec::new());
        Encode::internal_error(self.clone(), &mut cursor, error);

        let vec = cursor.into_inner();

        let existing_cap = WRITE_BUFFER.capacity() as usize;
        WRITE_BUFFER.reserve_exact((vec.len() - existing_cap) as usize);

        WRITE_BUFFER.clone_from(&vec);
        WRITE_BUFFER.shrink_to_fit();

        return (WRITE_BUFFER.as_ptr() as u32, WRITE_BUFFER.len() as u32);
    }

    unsafe fn from_read_buffer(&mut self) -> Option<Error> {
        let mut cursor = Cursor::new(&mut READ_BUFFER);
        let result = ExampleContext::decode(&mut cursor);
        return match result {
            Ok(context) => {
                *self = context.unwrap();
                None
            }
            Err(_) => Some(Error::new(ErrorKind::InvalidInput, "decoding error")),
        };
    }
}

impl Context {
    pub fn next(&mut self) -> Result<&mut Self, Error> {
        unsafe {
            let (ptr, len) = self.to_write_buffer();
            _next(ptr, len);
            return match self.from_read_buffer() {
                Some(err) => Err(err),
                None => Ok(self),
            };
        }
    }
}

pub unsafe fn resize(size: u32) -> *const u8 {
    let existing_cap = READ_BUFFER.capacity() as usize;
    READ_BUFFER.reserve_exact(size as usize - existing_cap);
    return READ_BUFFER.as_ptr();
}

#[link(wasm_import_module = "env")]
extern "C" {
    #[link_name = "next"]
    fn _next(ptr: u32, size: u32);
}
