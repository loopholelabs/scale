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
use polyglot_rs::Encoder;
use quickjs_wasm_sys::{
    ext_js_undefined, size_t as JS_size_t, JSContext, JSValue, JS_DefinePropertyValueStr,
    JS_FreeCString, JS_ThrowInternalError, JS_ToCStringLen2, JS_PROP_C_W_E,
};
use quickjs_wasm_sys::{JSCFunctionData, JS_NewCFunctionData};
use std::error::Error;
use std::ffi::CString;
use std::io::{Cursor, Write};
use std::os::raw::c_int;
use std::os::raw::c_void;

pub static mut GLOBAL_ERROR: Vec<u8> = Vec::new();

pub fn set_callback<F>(context: *mut JSContext, global: JSValue, fn_name: impl Into<Vec<u8>>, f: F)
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

pub fn console_log_to<T>(
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
                let format = CString::new("string pointer is null").unwrap();
                return unsafe { JS_ThrowInternalError(ctx, format.as_ptr()) };
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

pub fn global_err(err: Box<dyn Error>) -> u64 {
    let mut cursor = Cursor::new(Vec::new());
    cursor.encode_error(err).unwrap();

    unsafe {
        GLOBAL_ERROR = cursor.into_inner();

        // Now pack the address and length and return...
        let ptr = GLOBAL_ERROR.as_ptr() as u64;
        let len = GLOBAL_ERROR.len() as u64;
        return (ptr << 32) | len;
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
