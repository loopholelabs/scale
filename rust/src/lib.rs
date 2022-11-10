#![allow(dead_code)]
#![allow(unused_imports)]

#[path = "context/context.rs"]
mod context;

#[path = "context/request.rs"]
mod request;

#[path = "context/response.rs"]
mod response;

#[path = "compile/scale.rs"]
mod scale;

#[path = "utils/utils.rs"]
mod utils;

#[path = "generated/generated.rs"]
mod generated;

use lazy_static::lazy_static;
use std::sync::Mutex;

use std::io::Cursor;
use context::{RunContext, PTR, LEN, READ_BUFFER};
use generated::{Context};
use scale::scale;
use utils::pack_uint32;
use std::mem;
use std::mem::{MaybeUninit};
extern crate wee_alloc;

#[global_allocator]
pub static ALLOC: wee_alloc::WeeAlloc = wee_alloc::WeeAlloc::INIT;

#[cfg_attr(all(target_arch = "wasm32"), export_name = "run")]
#[no_mangle]
pub unsafe extern "C" fn run() -> u64 {
    //  Host calls resize first, which sets PTR and LEN.
    //  This unsafe pointer/len reconstruction gets around the os-level mutex restrictions from
    //  rust's Mutex, which are required for the read buffer global with static_mut
    let ptr = PTR.lock().unwrap().clone();
    let len = LEN.lock().unwrap().clone();
    let mut vec = Vec::from_raw_parts(ptr as *mut u8, len as usize, len as usize);

    let mut constructed = Cursor::new(&mut vec);
    let context: Context = RunContext::new();
    let cont = scale(context.from_read_buffer(&mut constructed));
    let ptr_len = cont.to_write_buffer();
    return pack_uint32(ptr_len.0, ptr_len.1);
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "resize")]
#[no_mangle]
pub unsafe extern "C" fn resize(size: u32) -> *const u8 {
   let existing_cap = READ_BUFFER.lock().unwrap().capacity() as u32;
   READ_BUFFER.lock().unwrap().reserve_exact((size - existing_cap) as usize);
   let ptr = READ_BUFFER.lock().unwrap().as_ptr();

   *PTR.lock().unwrap() = ptr as u32;
   *LEN.lock().unwrap() = size;

   return ptr
}
