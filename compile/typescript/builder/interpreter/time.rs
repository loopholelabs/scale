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

use crate::helpers;
use crate::helpers::set_callback;

use quickjs_wasm_sys::{
    ext_js_undefined, JSContext, JSRuntime, JSValue, JS_BigIntToUint64, JS_Call,
    JS_DefinePropertyValueStr, JS_Eval, JS_GetArrayBuffer, JS_GetException, JS_GetGlobalObject,
    JS_GetPropertyStr, JS_GetPropertyUint32, JS_SetPropertyInt64, JS_IsError, JS_NewContext, JS_NewInt32_Ext,
    JS_NewObject, JS_NewRuntime, JS_NewUint32_Ext, JS_EVAL_TYPE_GLOBAL, JS_PROP_C_W_E,
    JS_TAG_EXCEPTION, JS_TAG_UNDEFINED, JS_ExecutePendingJob, JS_GetOpaque, JS_TAG_OBJECT, JS_ToInt64,
    JS_IsFunction, JS_SetPropertyUint32, JS_IsLiveObject, JS_HackDupValue,
};

use std::time::{SystemTime, Duration};

use std::os::raw::c_int;

// Something to store our timer info in
pub struct timer_info {
  pub id: u64,
  pub callback: JSValue,
  pub repeating: bool,
  pub delay: Duration,
  pub ctime: SystemTime,
}

static mut TIMER_ID:u64 = 0;

// Active timers
pub static mut ACTIVE_TIMERS:Vec<timer_info> = Vec::new();

/**
 * Install time functions (setTimeout, clearTimeout, setInterval, clearInterval) into js runtime.
 *
 */
pub fn install(context: *mut JSContext) {
  unsafe {
    let global = JS_GetGlobalObject(context);

    helpers::set_callback(
      context,
      global,
      "setTimeout".to_string(),
      &set_timeout_wrap,
    );

    helpers::set_callback(
      context,
      global,
      "clearTimeout".to_string(),
      &clear_timeout_wrap,
    );

    helpers::set_callback(
      context,
      global,
      "setInterval".to_string(),
      &set_interval_wrap,
    );

    helpers::set_callback(
      context,
      global,
      "clearInterval".to_string(),
      &clear_interval_wrap,
    );
  }
}

/**
 * Run any pending jobs
 *
 */
pub fn run_pending_jobs(runtime: *mut JSRuntime, context: *mut JSContext) {
  unsafe {
    let global = JS_GetGlobalObject(context);
    let now = SystemTime::now();

//    print!("time:run_pending_jobs {}\n", ACTIVE_TIMERS.len());

    for tim in ACTIVE_TIMERS.iter_mut() {
      let elapsed = now.duration_since(tim.ctime).unwrap();

      if elapsed.as_secs_f64() > tim.delay.as_secs_f64() {
        // Time to trigger it then...

        if JS_IsFunction(context, tim.callback)==1 && JS_IsLiveObject(runtime, tim.callback)==1 {
          // Try calling it...
          let args: Vec<JSValue> = Vec::new();
          let r = JS_Call(context, tim.callback, global, args.len() as i32, args.as_slice().as_ptr() as *mut JSValue);
          if (r >> 32) as i32 == JS_TAG_EXCEPTION {
            // Show the issue...
            let err = helpers::error(context, "time");
            print!("Error {err}\n");
            //
          }
        } else {
          // GC or something...
          print!(" FUNC {:x} func?={} live?={}\n", tim.callback, JS_IsFunction(context, tim.callback), JS_IsLiveObject(runtime, tim.callback));
        }

        // If it's an interval, update the ctime value.
        if tim.repeating {
          tim.ctime = now
        }
      }
    }

    // Remove any timers we don't need anymore
    ACTIVE_TIMERS.retain(|tim| {
      let elapsed = now.duration_since(tim.ctime).unwrap();

      if elapsed.as_secs_f64() > tim.delay.as_secs_f64() {
        return false
      }
      return true
    })

  }
}

/**
 * Wrapper for setTimeout
 *
 */
pub fn set_timeout_wrap(
  context: *mut JSContext,
  _: JSValue,
  argc: c_int,
  argv: *mut JSValue,
  _: c_int,
) -> JSValue {
  unsafe {
    let global = JS_GetGlobalObject(context);
    let now = SystemTime::now();

    let func = *argv.offset(0);
    if JS_IsFunction(context, func)==1 {
      let mut delay:i64 = 0;
      JS_ToInt64(context, &mut delay as *mut i64, *argv.offset(1));

      let func2 = JS_HackDupValue(context, func);

      let t = timer_info{
        id: TIMER_ID,
        delay: Duration::from_millis(delay as u64),
        callback: func2,
        ctime: now,
        repeating: false,
      };
      TIMER_ID+=1;
      ACTIVE_TIMERS.push(t);
    }

    // FIXME: Return an object so they can cancel it...
    ext_js_undefined
  }
}

/**
 * Wrapper for clearTimeout
 *
 */
 pub fn clear_timeout_wrap(
  context: *mut JSContext,
  _: JSValue,
  argc: c_int,
  argv: *mut JSValue,
  _: c_int,
) -> JSValue {
  unsafe {
      print!("ClearTimeout argc={argc}\n");
      ext_js_undefined
  }
}

/**
 * Wrapper for setInterval
 *
 */
 pub fn set_interval_wrap(
  context: *mut JSContext,
  _: JSValue,
  argc: c_int,
  argv: *mut JSValue,
  _: c_int,
) -> JSValue {
  unsafe {
    let now = SystemTime::now();

    let func = *argv.offset(0);
    if JS_IsFunction(context, func)==1 {
      let mut delay:i64 = 0;
      JS_ToInt64(context, &mut delay as *mut i64, *argv.offset(1));

      let func2 = JS_HackDupValue(context, func);

      let t = timer_info{
        id: TIMER_ID,
        delay: Duration::from_millis(delay as u64),
        callback: func2,
        ctime: now,
        repeating: true,
      };
      TIMER_ID+=1;
      ACTIVE_TIMERS.push(t);
    }

    // FIXME: Return an object so they can cancel it...
    ext_js_undefined
  }
}

/**
 * Wrapper for clearInterval
 *
 */
 pub fn clear_interval_wrap(
  context: *mut JSContext,
  _: JSValue,
  argc: c_int,
  argv: *mut JSValue,
  _: c_int,
) -> JSValue {
  unsafe {
      print!("ClearInterval argc={argc}\n");
      ext_js_undefined
  }
}