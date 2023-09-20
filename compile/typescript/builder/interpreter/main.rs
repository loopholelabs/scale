/*
    Copyright 2022 Loophole Labs

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

           http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

pub mod helpers;

use once_cell::sync::OnceCell;
use quickjs_wasm_sys::{
    ext_js_undefined, JSContext, JSRuntime, JSValue, JS_BigIntToUint64, JS_Call,
    JS_DefinePropertyValueStr, JS_Eval, JS_GetArrayBuffer, JS_GetException, JS_GetGlobalObject,
    JS_GetPropertyStr, JS_GetPropertyUint32, JS_IsError, JS_NewContext, JS_NewInt32_Ext,
    JS_NewObject, JS_NewRuntime, JS_NewUint32_Ext, JS_EVAL_TYPE_GLOBAL, JS_PROP_C_W_E,
    JS_TAG_EXCEPTION, JS_TAG_UNDEFINED,
};

use std::ffi::CString;
use std::io::{self, Read};
use std::os::raw::c_int;
use std::str;

#[cfg(any(feature = "runtime_source", feature = "embedded_source"))]
use flate2::read::GzDecoder;

static mut JS_INITIALIZED: bool = false;
static mut JS_CONTEXT: OnceCell<*mut JSContext> = OnceCell::new();

static mut ENTRY_EXPORTS: OnceCell<JSValue> = OnceCell::new();

static mut ENTRY_MAIN: OnceCell<JSValue> = OnceCell::new();
static mut ENTRY_INITIALIZE: OnceCell<JSValue> = OnceCell::new();
static mut ENTRY_RUN: OnceCell<JSValue> = OnceCell::new();
static mut ENTRY_RESIZE: OnceCell<JSValue> = OnceCell::new();

static SCRIPT_NAME: &str = "index.js";

// If the wizer_opt feature is enabled, we will export a function called wizer.initialize
// so wizer knows what entrypoint to use.
#[cfg(feature = "wizer_opt")]
#[export_name = "wizer.initialize"]
#[no_mangle]
pub extern "C" fn init() {
    initialize_runtime();
}

#[cfg(all(not(feature = "embedded_source"), feature = "runtime_source"))]
#[link(wasm_import_module = "env")]
extern "C" {
    // get_js_source_len is a function that is imported from the host
    // that returns the runtime js source length as a u32
    #[link_name = "get_js_source_len"]
    fn get_js_source_len() -> u32;

    // get_js_source is a function that is imported from the host
    // that returns the runtime js source as a byte array
    #[link_name = "get_js_source"]
    fn get_js_source(ptr: *mut u8) -> u32;
}

// initialize_runtime creates the JS runtime and context
// and prepares the script for execution
fn initialize_runtime() {
    unsafe {
        let runtime: *mut JSRuntime = JS_NewRuntime();
        if runtime.is_null() {
            panic!("unable to create js runtime");
        }
        let context: *mut JSContext = JS_NewContext(runtime);
        if context.is_null() {
            panic!("unable to create js execution context");
        }

        let global = JS_GetGlobalObject(context);

        let log_cb = helpers::console_log_to(io::stdout());
        let error_cb = helpers::console_log_to(io::stderr());

        let console = JS_NewObject(context);
        helpers::set_callback(context, console, "log", log_cb);
        helpers::set_callback(context, console, "error", error_cb);

        let console_name = CString::new("console").unwrap();
        JS_DefinePropertyValueStr(
            context,
            global,
            console_name.as_ptr(),
            console,
            JS_PROP_C_W_E as i32,
        );

        helpers::set_callback(
            context,
            global,
            scale_signature_interfaces::TYPESCRIPT_NEXT,
            &next_wrap,
        );
        helpers::set_callback(
            context,
            global,
            scale_signature_interfaces::TYPESCRIPT_ADDRESS_OF,
            &address_of_wrap,
        );

        let mut js_source: String = String::new();

        // if the `embedded_source` feature is enabled,
        // read in the source from the `js_source` file
        #[cfg(all(not(feature = "runtime_source"), feature = "embedded_source"))]
        {
            println!("Reading embedded JS Source...");
            let source_data = include_bytes!(concat!(env!("OUT_DIR"), "/js_source"));
            if source_data.len() < 3 {
                panic!("js_source is empty");
            }
            if source_data[0] == 0x1f && source_data[1] == 0x8b {
                let mut gz = GzDecoder::new(&source_data[..]);
                gz.read_to_string(&mut js_source).unwrap();
            } else {
                js_source = str::from_utf8(source_data).unwrap().to_string();
            }
        }

        // if the `runtime_source` feature is enabled,
        // read in the source from the host using the `get_js_source` and `get_js_source_len` host
        // functions
        #[cfg(all(not(feature = "embedded_source"), feature = "runtime_source"))]
        {
            println!("Reading JS Source from Runtime...");
            let source_len = get_js_source_len() as usize;
            let source_buffer = vec![0; source_len];
            get_js_source(source_buffer.as_ptr() as *mut u8);
            let source_data = source_buffer.as_slice();
            if source_data.len() < 3 {
                panic!("js_source is empty");
            }
            if source_data[0] == 0x1f && source_data[1] == 0x8b {
                let mut gz = GzDecoder::new(&source_data[..]);
                gz.read_to_string(&mut js_source).unwrap();
            } else {
                js_source = str::from_utf8(source_data).unwrap().to_string();
            }
        }

        // if neither the `runtime_source` or `embedded_source` features are enabled,
        // read in the source from stdin
        #[cfg(all(not(feature = "runtime_source"), not(feature = "embedded_source")))]
        {
            println!("Reading JS Source from stdin...");
            io::stdin().read_to_string(&mut js_source).unwrap();
        }

        let js_source_len = js_source.len() - 1;
        let js_source_input = CString::new(js_source).unwrap();
        let js_source_script_name = CString::new(SCRIPT_NAME).unwrap();

        let ret = JS_Eval(
            context,
            js_source_input.as_ptr(),
            js_source_len as _,
            js_source_script_name.as_ptr(),
            JS_EVAL_TYPE_GLOBAL as i32,
        );

        if (ret >> 32) as i32 == JS_TAG_EXCEPTION {
            let e = JS_GetException(context);
            let exception = helpers::to_exception(context, e)
                .expect("getting exception during evaluation failed");

            let mut stack = None;
            let is_error = JS_IsError(context, e) != 0;
            if is_error {
                let cstring_key = CString::new("stack")
                    .expect("getting new CString for JS stack during evaluation failed");
                let raw = JS_GetPropertyStr(context, e, cstring_key.as_ptr());
                if (raw >> 32) as i32 != JS_TAG_UNDEFINED {
                    stack.replace(helpers::to_exception(context, raw));
                }
            }
            let mut err = format!("exception from js runtime during evaluation: {}", exception);
            if let Some(Ok(stack)) = stack {
                if stack.len() > 0 {
                    err.push_str(&format!("\nstack:\n{stack}"));
                }
            }
            panic!("{}", err);
        }

        let exports_key = CString::new("exports").unwrap();
        let exports = JS_GetPropertyStr(context, global, exports_key.as_ptr());

        let main_key = CString::new("main").unwrap();
        let main_fn = JS_GetPropertyStr(context, exports, main_key.as_ptr());
        ENTRY_MAIN.set(main_fn).unwrap();

        let initialize_key = CString::new("initialize").unwrap();
        let initialize_fn = JS_GetPropertyStr(context, exports, initialize_key.as_ptr());
        ENTRY_INITIALIZE.set(initialize_fn).unwrap();

        let run_key = CString::new("run").unwrap();
        let run_fn = JS_GetPropertyStr(context, exports, run_key.as_ptr());
        ENTRY_RUN.set(run_fn).unwrap();

        let resize_key = CString::new("resize").unwrap();
        let resize_fn = JS_GetPropertyStr(context, exports, resize_key.as_ptr());
        ENTRY_RESIZE.set(resize_fn).unwrap();

        ENTRY_EXPORTS.set(exports).unwrap();
        JS_CONTEXT.set(context).unwrap();
        JS_INITIALIZED = true;
    }
}

fn main() {
    unsafe {
        if !JS_INITIALIZED {
            initialize_runtime();
        }

        let context = JS_CONTEXT.get().unwrap();
        let exports = ENTRY_EXPORTS.get().unwrap();
        let mainfn = ENTRY_MAIN.get().unwrap();

        let args: Vec<JSValue> = Vec::new();
        let ret = JS_Call(
            *context,
            *mainfn,
            *exports,
            args.len() as i32,
            args.as_slice().as_ptr() as *mut JSValue,
        );

        if (ret >> 32) as i32 == JS_TAG_EXCEPTION {
            let e = JS_GetException(*context);
            let exception =
                helpers::to_exception(*context, e).expect("getting exception during main failed");

            let mut stack = None;
            let is_error = JS_IsError(*context, e) != 0;
            if is_error {
                let cstring_key = CString::new("stack")
                    .expect("getting new CString for JS stack during main failed");
                let raw = JS_GetPropertyStr(*context, e, cstring_key.as_ptr());
                if (raw >> 32) as i32 != JS_TAG_UNDEFINED {
                    stack.replace(helpers::to_exception(*context, raw));
                }
            }
            let mut err = format!("exception from js runtime during main: {}", exception);
            if let Some(Ok(stack)) = stack {
                if stack.len() > 0 {
                    err.push_str(&format!("\nstack:\n{stack}"));
                }
            }
            panic!("{}", err);
        }
    }
}

#[export_name = "initialize"]
#[no_mangle]
pub extern "C" fn initialize() -> u64 {
    unsafe {
        let context = JS_CONTEXT.get().expect("js context not initialized");
        let exports = ENTRY_EXPORTS.get().expect("exports not initialized");
        let initfn = ENTRY_INITIALIZE.get().expect("initialize not initialized");

        let args: Vec<JSValue> = Vec::new();
        let ret = JS_Call(
            *context,
            *initfn,
            *exports,
            args.len() as i32,
            args.as_slice().as_ptr() as *mut JSValue,
        );

        if (ret >> 32) as i32 == JS_TAG_EXCEPTION {
            let e = JS_GetException(*context);
            let exception = helpers::to_exception(*context, e)
                .expect("getting exception during initialize failed");

            let mut stack = None;
            let is_error = JS_IsError(*context, e) != 0;
            if is_error {
                let cstring_key = CString::new("stack")
                    .expect("getting new CString for JS stack during initialize failed");
                let raw = JS_GetPropertyStr(*context, e, cstring_key.as_ptr());
                if (raw >> 32) as i32 != JS_TAG_UNDEFINED {
                    stack.replace(helpers::to_exception(*context, raw));
                }
            }
            let mut err = format!("exception from js runtime during initialize: {}", exception);
            if let Some(Ok(stack)) = stack {
                if stack.len() > 0 {
                    err.push_str(&format!("\nstack:\n{stack}"));
                }
            }
            return helpers::global_err(err.into());
        }

        let mut valret = 0_u64;
        let err = JS_BigIntToUint64(*context, &mut valret, ret);
        if err < 0 {
            return helpers::global_err("error converting return value from initialize".into());
        }
        valret
    }
}

#[export_name = "run"]
#[no_mangle]
pub extern "C" fn run() -> u64 {
    unsafe {
        let context = JS_CONTEXT.get().unwrap();
        let exports = ENTRY_EXPORTS.get().unwrap();
        let runfn = ENTRY_RUN.get().unwrap();

        let args: Vec<JSValue> = Vec::new();
        let ret = JS_Call(
            *context,
            *runfn,
            *exports,
            args.len() as i32,
            args.as_slice().as_ptr() as *mut JSValue,
        );

        if (ret >> 32) as i32 == JS_TAG_EXCEPTION {
            let e = JS_GetException(*context);
            let exception =
                helpers::to_exception(*context, e).expect("getting exception during run failed");

            let mut stack = None;
            let is_error = JS_IsError(*context, e) != 0;
            if is_error {
                let cstring_key = CString::new("stack")
                    .expect("getting new CString for JS stack during run failed");
                let raw = JS_GetPropertyStr(*context, e, cstring_key.as_ptr());
                if (raw >> 32) as i32 != JS_TAG_UNDEFINED {
                    stack.replace(helpers::to_exception(*context, raw));
                }
            }
            let mut err = format!("exception from js runtime during run: {}", exception);
            if let Some(Ok(stack)) = stack {
                if stack.len() > 0 {
                    err.push_str(&format!("\nstack:\n{stack}"));
                }
            }
            return helpers::global_err(err.into());
        }

        let mut valret = 0_u64;
        let err = JS_BigIntToUint64(*context, &mut valret, ret);
        if err < 0 {
            return helpers::global_err("error converting return value from run".into());
        }
        valret
    }
}

#[export_name = "resize"]
#[no_mangle]
pub extern "C" fn resize(size: u32) -> *mut u8 {
    unsafe {
        let context = JS_CONTEXT.get().unwrap();
        let exports = ENTRY_EXPORTS.get().unwrap();
        let resizefn = ENTRY_RESIZE.get().unwrap();

        let mut args: Vec<JSValue> = Vec::new();
        let jval = JS_NewInt32_Ext(*context, size as i32);
        args.push(jval);

        let ret = JS_Call(
            *context,
            *resizefn,
            *exports,
            args.len() as i32,
            args.as_slice().as_ptr() as *mut JSValue,
        );

        if (ret >> 32) as i32 == JS_TAG_EXCEPTION {
            let e = JS_GetException(*context);
            let exception =
                helpers::to_exception(*context, e).expect("getting exception during resize failed");

            let mut stack = None;
            let is_error = JS_IsError(*context, e) != 0;
            if is_error {
                let cstring_key = CString::new("stack")
                    .expect("getting new CString for JS stack during resize failed");
                let raw = JS_GetPropertyStr(*context, e, cstring_key.as_ptr());
                if (raw >> 32) as i32 != JS_TAG_UNDEFINED {
                    stack.replace(helpers::to_exception(*context, raw));
                }
            }
            let mut err = format!("exception from js runtime during resize: {}", exception);
            if let Some(Ok(stack)) = stack {
                if stack.len() > 0 {
                    err.push_str(&format!("\nstack:\n{stack}"));
                }
            }
            panic!("{}", err);
        }

        ret as *mut u8
    }
}

#[link(wasm_import_module = "env")]
extern "C" {
    #[link_name = "next"]
    fn _next(ptr: u32, size: u32);
}

// Wrap the exported next function so it can be called from the js runtime
fn next_wrap(
    context: *mut JSContext,
    _: JSValue,
    _: c_int,
    js_value: *mut JSValue,
    _: c_int,
) -> JSValue {
    unsafe {
        let ptr = JS_GetPropertyUint32(context, *js_value, 0) as u32;
        let len = JS_GetPropertyUint32(context, *js_value, 1) as u32;
        _next(ptr, len);
        ext_js_undefined
    }
}

// Wrap the exported address_of function so it can be called from the js runtime
fn address_of_wrap(
    context: *mut JSContext,
    _: JSValue,
    _: c_int,
    js_value: *mut JSValue,
    _: c_int,
) -> JSValue {
    unsafe {
        let mut len = 0;
        let addr = JS_GetArrayBuffer(context, &mut len, *js_value) as u32;
        return JS_NewUint32_Ext(context, addr);
    }
}
