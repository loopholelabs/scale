#![allow(dead_code)]
#![allow(unused_imports)]

#[macro_use]
extern crate lazy_static;

#[path = "utils/utils.rs"]
mod utils;


extern crate quickjs_wasm_sys;
extern crate once_cell;

use quickjs_wasm_sys::{
  ext_js_exception, ext_js_null, ext_js_undefined, size_t as JS_size_t, JSCFunctionData,
  JSContext, JSValue, JS_Eval, JS_FreeCString, JS_GetGlobalObject, JS_NewArray, JS_NewBigInt64,
  JS_Call, JS_NewBool_Ext, JS_NewCFunctionData, JS_NewContext, JS_NewFloat64_Ext, JS_NewInt32_Ext,
  JS_NewInt64_Ext, JS_NewObject, JS_NewRuntime, JS_NewStringLen, JS_NewUint32_Ext,
  JS_ToCStringLen2, JS_EVAL_TYPE_GLOBAL, JS_GetPropertyStr, JS_GetPropertyUint32, 
  JS_DefinePropertyValueStr, JS_DefinePropertyValueUint32, JS_PROP_C_W_E,
  JS_TAG_BIG_INT, JS_TAG_BOOL, JS_TAG_EXCEPTION, JS_TAG_INT, JS_TAG_NULL,
  JS_TAG_OBJECT, JS_TAG_STRING, JS_TAG_UNDEFINED,
};
use std::os::raw::{c_char, c_int, c_void};

use utils::{pack_uint32, unpack_uint32, vec_to_js, js_to_vec, set_buffer, resize_buffer, set_next_buffer};
use utils::{READ_BUFFER, RETURN_BUFFER, NEXT_READ_BUFFER};

use std::io::{self, Cursor, Read, Write};
use std::ffi::CString;

extern crate wee_alloc;

#[cfg(not(test))]
#[global_allocator]
static ALLOC: wee_alloc::WeeAlloc = wee_alloc::WeeAlloc::INIT;


use once_cell::sync::OnceCell;
static mut JS_CONTEXT: OnceCell<*mut JSContext> = OnceCell::new();

static mut ENTRY_EXPORTS: OnceCell<JSValue> = OnceCell::new();

static mut ENTRY_MAIN: OnceCell<JSValue> = OnceCell::new();
static mut ENTRY_RUN: OnceCell<JSValue> = OnceCell::new();
static mut ENTRY_RESIZE: OnceCell<JSValue> = OnceCell::new();

static SCRIPT_NAME: &str = "script.js";

// The function env.next exported by the host
#[link(wasm_import_module = "env")]
extern "C" {
    #[link_name = "next"]
    fn _next(ptr: u32, size: u32) -> u64;
}

// Wrap the exported next function so it can be called from js
fn nextwrap(context: *mut JSContext, jsval1: JSValue, int1: c_int, jsval2: *mut JSValue, int2: c_int) -> JSValue {
  unsafe {
    let vec = js_to_vec(context, *jsval2);
    let (ptr, len) = set_next_buffer(vec);

    let packed = _next(ptr, len);
    let (ptr, len) = unpack_uint32(packed);    
    let rvec = Vec::from_raw_parts(ptr as *mut u8, len as usize, len as usize);
    return vec_to_js(context, &rvec);
  }
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "resize")]
#[no_mangle]
pub unsafe extern "C" fn resize(size: u32) -> *const u8 {
  return resize_buffer(size);
}


#[export_name = "wizer.initialize"]
pub extern "C" fn init() {
    unsafe {
        let runtime = JS_NewRuntime();
        if runtime.is_null() {
          panic!("Couldn't create JavaScript runtime");
        }
        let context = JS_NewContext(runtime);
        if context.is_null() {
          panic!("Couldn't create JavaScript context");
        }

        let mut contents = String::new();
        io::stdin().read_to_string(&mut contents).unwrap();

        let len = contents.len() - 1;
        let input = CString::new(contents).unwrap();
        let script_name = CString::new(SCRIPT_NAME).unwrap();
        JS_Eval(context, input.as_ptr(), len as _, script_name.as_ptr(), JS_EVAL_TYPE_GLOBAL as i32);

        let global = JS_GetGlobalObject(context);
        let exports_key = CString::new("Exports").unwrap();
        let exports = JS_GetPropertyStr(context, global, exports_key.as_ptr());

        let main_key = CString::new("main").unwrap();
        let main_fn = JS_GetPropertyStr(context, exports, main_key.as_ptr());
        ENTRY_MAIN.set(main_fn).unwrap();

        let run_key = CString::new("run").unwrap();
        let run_fn = JS_GetPropertyStr(context, exports, run_key.as_ptr());
        ENTRY_RUN.set(run_fn).unwrap();

        // Setup console.log and console.error to pipe through to io::stderr
        let log_cb = console_log_to(io::stderr());
        let error_cb = console_log_to(io::stderr());

        let console = JS_NewObject(context);
        set_callback(context, console, "log", log_cb);
        set_callback(context, console, "error", error_cb);

        let console_name = CString::new("console").unwrap();
        JS_DefinePropertyValueStr(context, global, console_name.as_ptr(), console, JS_PROP_C_W_E as i32);

        // Setup a function called next() in the global_object
        set_callback(context, global, "scale_fn_next", &nextwrap);

        ENTRY_EXPORTS.set(exports).unwrap();

        JS_CONTEXT.set(context).unwrap();
    }
}

fn set_callback<F>(context: *mut JSContext, global: JSValue, fn_name: impl Into<Vec<u8>>, f: F)
where
  F: FnMut(*mut JSContext, JSValue, c_int, *mut JSValue, c_int) -> JSValue,
{
  unsafe {
    let trampoline = build_trampoline(&f);
    let data = &f as *const _ as *mut c_void as *mut JSValue;
    let cb = JS_NewCFunctionData(context, trampoline, 0, 1, 1, data);

    let name_fn = CString::new(fn_name).unwrap();

    JS_DefinePropertyValueStr(context, global, name_fn.as_ptr(), cb, JS_PROP_C_W_E as i32);
  }
}

fn build_trampoline<F>(_f: &F) -> JSCFunctionData
where
    F: FnMut(*mut JSContext, JSValue, c_int, *mut JSValue, c_int) -> JSValue,
{
    // We build a trampoline to jump between c <-> rust and allow closing over a specific context.
    // For more info around how this works, see https://adventures.michaelfbryan.com/posts/rust-closures-in-ffi/.
    unsafe extern "C" fn trampoline<F>(
        ctx: *mut JSContext,
        this: JSValue,
        argc: c_int,
        argv: *mut JSValue,
        magic: c_int,
        data: *mut JSValue,
    ) -> JSValue
    where
        F: FnMut(*mut JSContext, JSValue, c_int, *mut JSValue, c_int) -> JSValue,
    {
        let closure_ptr = data;
        let closure: &mut F = &mut *(closure_ptr as *mut F);
        (*closure)(ctx, this, argc, argv, magic)
    }

    Some(trampoline::<F>)
}

fn console_log_to<T>(
  mut stream: T,
) -> impl FnMut(*mut JSContext, JSValue, c_int, *mut JSValue, c_int) -> JSValue
where
  T: Write,
{
  move |ctx: *mut JSContext, _this: JSValue, argc: c_int, argv: *mut JSValue, _magic: c_int| {
      let mut len: JS_size_t = 0;
      for i in 0..argc {
          if i != 0 {
              write!(stream, " ").unwrap();
          }

          let str_ptr = unsafe { JS_ToCStringLen2(ctx, &mut len, *argv.offset(i as isize), 0) };
          if str_ptr.is_null() {
              return unsafe { ext_js_exception };
          }

          let str_ptr = str_ptr as *const u8;
          let str_len = len as usize;
          let buffer = unsafe { std::slice::from_raw_parts(str_ptr, str_len) };

          stream.write_all(buffer).unwrap();
          unsafe { JS_FreeCString(ctx, str_ptr as *const i8) };
      }

      writeln!(stream,).unwrap();
      unsafe { ext_js_undefined }
  }
}

fn main() {
    unsafe {
        let context = JS_CONTEXT.get().unwrap();
        let exports = ENTRY_EXPORTS.get().unwrap();
        let main = ENTRY_MAIN.get().unwrap();

        let args: Vec<JSValue> = Vec::new();
        let ret = JS_Call(*context, *main, *exports, args.len() as i32, args.as_slice().as_ptr() as *mut JSValue);
    }
}


#[export_name = "run"]
fn run() -> u64 {
  unsafe {
    let context = JS_CONTEXT.get().unwrap();
    let exports = ENTRY_EXPORTS.get().unwrap();
    let runfn = ENTRY_RUN.get().unwrap();

    let input_vals = vec_to_js(*context, &READ_BUFFER);
    let mut args: Vec<JSValue> = Vec::new();
    args.push(input_vals);
    let ret = JS_Call(*context, *runfn, *exports, args.len() as i32, args.as_slice().as_ptr() as *mut JSValue);

    let ret_tag = (ret >> 32) as i32;
    if ret_tag == JS_TAG_EXCEPTION {
      // TODO Get the exception and handle and return to host?...
      //
      println!("Exception from js!");
      // Signal error for now.
      return 900;
    }

    if ret_tag != JS_TAG_OBJECT {
      println!("Return from run was not an object!");
      // Signal error for now.
      return 999;
    }

    let retvec = js_to_vec(*context, ret);
    let (ptr, len) = set_buffer(retvec);
    return pack_uint32(ptr, len);
  }
}
