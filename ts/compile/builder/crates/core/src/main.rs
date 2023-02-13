
extern crate quickjs_wasm_sys;
extern crate once_cell;

use quickjs_wasm_sys::{
  ext_js_exception, ext_js_undefined, size_t as JS_size_t, JSCFunctionData,
  JSContext, JSValue, JS_Eval, JS_FreeCString, JS_GetGlobalObject,
  JS_Call, JS_NewCFunctionData, JS_NewContext, JS_NewInt32_Ext,
  JS_NewInt64_Ext, JS_NewObject, JS_NewRuntime, JS_ToCStringLen2, JS_EVAL_TYPE_GLOBAL,
  JS_GetPropertyStr, JS_GetPropertyUint32, JS_DefinePropertyValueStr, JS_PROP_C_W_E,
  JS_TAG_EXCEPTION, JS_GetArrayBuffer, JS_BigIntToUint64
};
use std::os::raw::{c_int, c_void};

use std::io::{self, Read, Write};
use std::ffi::CString;

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
fn nextwrap(context: *mut JSContext, _jsval1: JSValue, _int1: c_int, jsval2: *mut JSValue, _int2: c_int) -> JSValue {
  unsafe {
    // The args are [ptr / len]
    // TODO: Make sure it's an array of 2 vals?
    let ptr = JS_GetPropertyUint32(context, *jsval2, 0);    
    let len = JS_GetPropertyUint32(context, *jsval2, 1);    

    let packed = _next(ptr as u32, len as u32);
    return JS_NewInt64_Ext(context, packed as i64);    
  }
}

// Get the address of a javascript ArrayBuffer
fn getaddrwrap(context: *mut JSContext, _jsval1: JSValue, _int1: c_int, jsval2: *mut JSValue, _int2: c_int) -> JSValue {
  unsafe {
    let mut len = 0;    
    let addr = JS_GetArrayBuffer(context, &mut len, *jsval2) as i32;    
    return JS_NewInt32_Ext(context, addr);
  }
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

        let resize_key = CString::new("resize").unwrap();
        let resize_fn = JS_GetPropertyStr(context, exports, resize_key.as_ptr());
        ENTRY_RESIZE.set(resize_fn).unwrap();

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

        // Setup a function called getaddr() in the global object
        set_callback(context, global, "getaddr", &getaddrwrap);

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
        let _ret = JS_Call(*context, *main, *exports, args.len() as i32, args.as_slice().as_ptr() as *mut JSValue);
    }
}


#[export_name = "run"]
fn run() -> u64 {
  unsafe {
    let context = JS_CONTEXT.get().unwrap();
    let exports = ENTRY_EXPORTS.get().unwrap();
    let runfn = ENTRY_RUN.get().unwrap();

    let args: Vec<JSValue> = Vec::new();
    let ret = JS_Call(*context, *runfn, *exports, args.len() as i32, args.as_slice().as_ptr() as *mut JSValue);

    let ret_tag = (ret >> 32) as i32;
    if ret_tag == JS_TAG_EXCEPTION {
      // TODO Get the exception and handle and return to host?...
      //
      println!("Rust/js: Exception from js!");
      // Signal error for now.
      return 999;
    }

    let mut valret = 0_u64;
    let err = JS_BigIntToUint64(*context, &mut valret, ret);
    if err < 0 {
      // TODO: Return a better error maybe...
      println!("Rust/js: Error converting run return value");
      return 999;
    }      
    return valret;
  }
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "resize")]
#[no_mangle]
pub unsafe extern "C" fn resize(size: u32) -> *mut u8 {
  unsafe {
    let context = JS_CONTEXT.get().unwrap();
    let exports = ENTRY_EXPORTS.get().unwrap();
    let resizefn = ENTRY_RESIZE.get().unwrap();

    let mut args: Vec<JSValue> = Vec::new();
    let jval = JS_NewInt32_Ext(*context, size as i32);
    args.push(jval);

    let ret = JS_Call(*context, *resizefn, *exports, args.len() as i32, args.as_slice().as_ptr() as *mut JSValue);

    let ret_tag = (ret >> 32) as i32;
    if ret_tag == JS_TAG_EXCEPTION {
      // TODO Get the exception and handle and return to host?...
      //
      println!("Rust/js: Exception from js!");
      // Signal error for now.
      return 999 as *mut u8;
    }

    return ret as *mut u8;
  }
}