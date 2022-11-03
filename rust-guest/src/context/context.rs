#![allow(dead_code)]
#![allow(unused_imports)]
#![allow(unused_variables)]

use lazy_static::lazy_static;
use std::sync::Mutex;

use std::io::Cursor;
use std::mem;
use std::collections::HashMap;
use super::generated::{Encode, Decode, Context, Request, Response};

lazy_static! {
    pub static ref PTR_LEN: Mutex<(u32, u32)> = Mutex::new((0, 0));
}

pub trait RunContext {
    fn new() -> Self;
    fn from_read_buffer(self, read_buff_global: &mut Cursor<&mut Vec<u8>>) -> Self;
    fn to_write_buffer(self) -> (u32, u32);
    fn next(self) -> Self;
    fn request(&mut self) -> &mut Request;
    fn response(&mut self) -> &mut Response;
    fn default() -> Self;
}

impl RunContext for Context {
    fn new()  -> Context {
          Self::default()
    }

    fn from_read_buffer(self, read_buff_global: &mut Cursor<&mut Vec<u8>>) -> Self {
          Decode::decode(read_buff_global).unwrap().unwrap()
    }

    fn to_write_buffer(self) -> (u32, u32) {
        let mut cursor = Cursor::new(Vec::new());
        let _ = Encode::encode(self, &mut cursor);

        let mut vec = cursor.into_inner();
        vec.shrink_to_fit();
        let ptr = vec.as_ptr() as u32;
        let len = vec.len() as u32;
        mem::forget(vec);  // prevents deallocation in Rust
                           // vec still exists in mem, but
                           // rust doesn't have any concept of it

        return (ptr, len)
    }

    fn request(&mut self) -> &mut Request {
        &mut self.request
    }

    fn response(&mut self) -> &mut Response {
        &mut self.response
    }

    fn next(self) -> Self {
         unsafe {
           let ptr_len = self.to_write_buffer();

           //  calls resize from host side, which sets PTR_LEN
           _next(ptr_len.0, ptr_len.1);

           let ptr = PTR_LEN.lock().unwrap().0;
           let len = PTR_LEN.lock().unwrap().1;

           let mut vec = Vec::from_raw_parts(ptr as *mut u8, len as usize, len as usize);
           let mut constructed = Cursor::new(&mut vec);

           let empty_context: Context = Self::default();

           let from_buf = empty_context.from_read_buffer(&mut constructed);
           return from_buf;
         }
    }

    fn default() -> Self {
            Context {
                    request: Request {
                        headers: HashMap::new(),
                        method: "".to_string(),
                        content_length: 0,
                        protocol: "".to_string(),
                        i_p: "".to_string(),
                        body: Vec::new()
                    },
                    response: Response {
                        headers: HashMap::new(),
                        status_code: 0,
                        body: Vec::new()
                    },
           }
    }
}

#[link(wasm_import_module = "env")]
extern "C" {
    #[link_name = "next"]
    fn _next(ptr: u32, size: u32);
}
