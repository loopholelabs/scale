// Code generated by scale-signature , DO NOT EDIT.
// output: signature

pub mod types;
use crate::types::{Encode, Decode};

use std::io::Cursor;
use polyglot_rs::{Encoder};

static HASH: &'static str = "3a592aa345d412faa2e6285ee048ca2ab5aa64b0caa2f9ca67b2c1e0792101e5";

static mut READ_BUFFER: Vec<u8> = Vec::new();
static mut WRITE_BUFFER: Vec<u8> = Vec::new();

// write serializes the signature into the global WRITE_BUFFER and returns the pointer to the buffer and its size
//
// Users should not use this method.
pub unsafe fn write(ctx: Option<&mut types::ModelWithAllFieldTypes>) -> (u32, u32) {
    let mut cursor = Cursor::new(Vec::new());
    match ctx {
        Some(ctx) => {
            cursor = match types::ModelWithAllFieldTypes::encode(Some(ctx), &mut cursor) {
                Ok(_) => cursor,
                Err(err) => return error(err),
            };
        }
        None => {
            cursor = match types::ModelWithAllFieldTypes::encode(None, &mut cursor) {
                Ok(_) => cursor,
                Err(err) => return error(err),
            };
        }
    }
    let vec = cursor.into_inner();

    WRITE_BUFFER.resize(vec.len() as usize, 0);
    WRITE_BUFFER.copy_from_slice(&vec);

    return (WRITE_BUFFER.as_ptr() as u32, WRITE_BUFFER.len() as u32);
}

// read deserializes signature from the global READ_BUFFER
//
// Users should not use this method.
pub unsafe fn read() -> Result<Option<types::ModelWithAllFieldTypes>, Box<dyn std::error::Error>> {
    let mut cursor = Cursor::new(&mut READ_BUFFER);
    types::ModelWithAllFieldTypes::decode(&mut cursor)
}

// error serializes an error into the global WRITE_BUFFER and returns a pointer to the buffer and its size
//
// Users should not use this method.
pub unsafe fn error(error: Box<dyn std::error::Error>) -> (u32, u32) {
    let mut cursor = Cursor::new(Vec::new());
    return match cursor.encode_error(error) {
        Ok(_) => {
            let vec = cursor.into_inner();

            WRITE_BUFFER.resize(vec.len() as usize, 0);
            WRITE_BUFFER.copy_from_slice(&vec);

            (WRITE_BUFFER.as_ptr() as u32, WRITE_BUFFER.len() as u32)
        }
        Err(_) => {
            (0, 0)
        }
    };
}

// resize resizes the global READ_BUFFER to the given size and returns the pointer to the buffer
//
// Users should not use this method.
pub unsafe fn resize(size: u32) -> *const u8 {
    READ_BUFFER.resize(size as usize, 0);
    return READ_BUFFER.as_ptr();
}

// hash returns the hash of the Scale Signature
//
// Users should not use this method.
pub unsafe fn hash() -> (u32, u32) {
    let mut cursor = Cursor::new(Vec::new());
    return match cursor.encode_string(&String::from(HASH)) {
        Ok(_) => {
            let vec = cursor.into_inner();

            WRITE_BUFFER.resize(vec.len() as usize, 0);
            WRITE_BUFFER.copy_from_slice(&vec);

            (WRITE_BUFFER.as_ptr() as u32, WRITE_BUFFER.len() as u32)
        }
        Err(_) => {
            (0, 0)
        }
    };
}

// next calls the next function in the Scale Function Chain
pub fn next(ctx: Option<&mut types::ModelWithAllFieldTypes>) -> Result<Option<types::ModelWithAllFieldTypes>, Box<dyn std::error::Error>> {
    unsafe {
        let (ptr, len) = write(ctx);
        _next(ptr, len);
        read()
    }
}

#[link(wasm_import_module = "env")]
extern "C" {
    #[link_name = "next"]
    fn _next(ptr: u32, size: u32);
}