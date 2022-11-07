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
use context::{RunContext, PTR_LEN};
use generated::{Context};
use scale::scale;
use utils::pack_uint32;
use std::mem;
use std::mem::{MaybeUninit};

#[cfg_attr(all(target_arch = "wasm32"), export_name = "run")]
#[no_mangle]
pub unsafe extern "C" fn run() -> u64 {
    //  host calls resize first, which sets PTR_LEN
    let ptr = PTR_LEN.lock().unwrap().0;
    let len = PTR_LEN.lock().unwrap().1;
    let mut vec = Vec::from_raw_parts(ptr as *mut u8, len as usize, len as usize);
    let mut constructed = Cursor::new(&mut vec);

    let context: Context = RunContext::new();
    let cont = scale(context.from_read_buffer(&mut constructed));
    let ptr_len = cont.to_write_buffer();
    return pack_uint32(ptr_len.0, ptr_len.1);
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "resize")]
#[no_mangle]
pub unsafe extern "C" fn resize(size: u32) -> *const MaybeUninit<u8> {
    let vec: Vec<MaybeUninit<u8>> = Vec::with_capacity(size as usize);

    let ptr = vec.as_ptr();
    PTR_LEN.lock().unwrap().0 = ptr as u32;
    PTR_LEN.lock().unwrap().1 = size;
    mem::forget(vec);  // prevents deallocation in Rust
                       // vec still exists in mem, but
                       // rust doesn't have any concept of it

   return ptr
}
