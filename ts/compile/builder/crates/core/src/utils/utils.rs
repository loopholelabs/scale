#![allow(dead_code)]
#![allow(unused_imports)]

extern crate quickjs_wasm_sys;

use quickjs_wasm_sys::{
  ext_js_exception, ext_js_null, ext_js_undefined, size_t as JS_size_t, JSCFunctionData,
  JSContext, JSValue, JS_Eval, JS_FreeCString, JS_GetGlobalObject, JS_NewArray, JS_NewBigInt64,
  JS_NewBool_Ext, JS_NewCFunctionData, JS_NewContext, JS_NewFloat64_Ext, JS_NewInt32_Ext,
  JS_NewInt64_Ext, JS_NewObject, JS_NewRuntime, JS_NewStringLen, JS_NewUint32_Ext,
  JS_ToCStringLen2, JS_EVAL_TYPE_GLOBAL, JS_GetPropertyStr, JS_GetPropertyUint32, 
  JS_DefinePropertyValueUint32, JS_PROP_C_W_E,
};
use std::os::raw::{c_char, c_int, c_void};

use std::ffi::CString;
use std::io::{Cursor, Write};

use lazy_static::lazy_static;
use std::sync::Mutex;

//lazy_static! {
  pub static mut READ_BUFFER: Vec<u8> = Vec::new();
  pub static mut RETURN_BUFFER: Vec<u8> = Vec::new();
  pub static mut NEXT_READ_BUFFER: Vec<u8> = Vec::new();
//}

pub fn set_buffer(v: Vec<u8>) -> (u32, u32) {
  unsafe {
    RETURN_BUFFER = v;
    let ptr = RETURN_BUFFER.as_ptr() as u32;
    let len = RETURN_BUFFER.len() as u32;
    return (ptr, len);
  }
}

pub unsafe fn resize_buffer(size: u32) -> *const u8 {
  let existing_cap = READ_BUFFER.capacity() as u32;
  if (existing_cap == size) {   // Quick shortcut
    // return ptr
    let ptr = READ_BUFFER.as_ptr();
    return ptr
  }
  READ_BUFFER.reserve_exact((size - existing_cap) as usize);
  READ_BUFFER.resize(size as usize, 0);
  let ptr = READ_BUFFER.as_ptr();
  return ptr
}

pub fn set_next_buffer(v: Vec<u8>) -> (u32, u32) {
  unsafe {
    NEXT_READ_BUFFER = v;
    let ptr = NEXT_READ_BUFFER.as_ptr() as u32;
    let len = NEXT_READ_BUFFER.len() as u32;
    return (ptr, len);
  }
}

pub fn pack_uint32(ptr: u32, len: u32) -> u64 {
    return ((ptr as u64) << 32) | len as u64;
}

pub fn unpack_uint32(packed: u64) -> (u32, u32) {
	return ((packed >> 32) as u32, packed as u32)
}

// Convert a vec to a JSValue array of bytes
pub fn vec_to_js(context: *mut JSContext, v: &Vec<u8>) -> JSValue {
  unsafe {
    let array = JS_NewArray(context);
    let mut index: u32 = 0;
    for val in v.iter() {
      let jval = JS_NewInt32_Ext(context, i32::from(*val));
      JS_DefinePropertyValueUint32(
          context,
          array,
          index,
          jval,
          JS_PROP_C_W_E as i32,
      );
      index += 1;
    }
    return array;
  }
}

pub fn js_to_vec(context: *mut JSContext, v: JSValue) -> Vec<u8> {
  let mut ve = Vec::new();

  // The input (jsval2) is expected to be an array of bytes.
  let cstring_key = CString::new("length").unwrap();
  let len = unsafe { JS_GetPropertyStr(context, v, cstring_key.as_ptr()) } as u32;

  for i in 0..len {
    let v = unsafe { JS_GetPropertyUint32(context, v, i) } as u8;
//    let nval:&[u8] = &[v];
    ve.push(v);
  }

  ve.shrink_to_fit();
  return ve;
}
