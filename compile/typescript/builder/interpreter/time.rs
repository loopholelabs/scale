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

use quickjs_wasm_sys::{
    ext_js_undefined, JSContext, JSRuntime, JSValue, JS_Call,
    JS_GetGlobalObject, JS_NewInt64_Ext, JS_TAG_EXCEPTION, JS_ToInt64,
    JS_IsFunction, JS_IsLiveObject, JS_DuplicateValue,
};

use std::time::{SystemTime, Duration};
use std::ops::Add;

use std::os::raw::c_int;

// Something to store our timer info in
pub struct TimerInfo {
  pub id: u64,                  // Unique ID that represents this timer.
  pub active: bool,             // Is the timer still active, or pending deletion?
  pub callback: JSValue,        // Function to callback
  pub repeating: bool,          // Is the timer repeating? (setInterval)
  pub delay: Duration,          // Delay between ticks for a repeating timer
  pub trigger_time: SystemTime, // Next time the timer will trigger
}

static mut TIMER_ID:u64 = 0;

// Active timers
pub static mut ACTIVE_TIMERS:Vec<TimerInfo> = Vec::new();

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

    for tim in ACTIVE_TIMERS.iter_mut() {
      match tim.trigger_time.elapsed() {
        Ok(..) => {
          if JS_IsFunction(context, tim.callback)==1 && JS_IsLiveObject(runtime, tim.callback)==1 {
            let args: Vec<JSValue> = Vec::new();
            let r = JS_Call(context, tim.callback, global, args.len() as i32, args.as_slice().as_ptr() as *mut JSValue);
            if (r >> 32) as i32 == JS_TAG_EXCEPTION {
              // A trigger function from a timer threw an exception.
              // For now, we will just show it in output.
              let err = helpers::error(context, "time");
              print!("Error from timer function {err}\n");
            }
          } else {
            // Either a timer was called with a non-function argument, or quickjs GC stole it?
            print!("Error from timer function {:x} func?={} live?={}\n", tim.callback, JS_IsFunction(context, tim.callback), JS_IsLiveObject(runtime, tim.callback));
          }
  
          // If it's an interval, update the next trigger time value.
          // If not, it can be marked for removal.
          if tim.repeating {
            tim.trigger_time = now.add(tim.delay);
          } else {
            tim.active = false
          }
        },
        Err(..) => {
          // This timer isn't ready to trigger yet.
        }
      }
    }

    // Remove any timers we don't need anymore
    ACTIVE_TIMERS.retain(|tim| {
      tim.active
    })
  }
}

/**
 * Wrapper for setTimeout(func delay)=>ID
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
    if argc!=2 {
      return ext_js_undefined;
    }
    let now = SystemTime::now();

    let mut ret = ext_js_undefined;

    let func = *argv.offset(0);
    if JS_IsFunction(context, func)==1 {
      let mut delay:i64 = 0;
      JS_ToInt64(context, &mut delay as *mut i64, *argv.offset(1));

      let func2 = JS_DuplicateValue(context, func);

      let delay_duration = Duration::from_millis(delay as u64);
      let t = TimerInfo{
        id: TIMER_ID,
        active: true,
        callback: func2,
        trigger_time: now.add(delay_duration),
        repeating: false,
        delay: delay_duration,
      };

      ret = JS_NewInt64_Ext(context, TIMER_ID as i64);

      TIMER_ID+=1;
      ACTIVE_TIMERS.push(t);
    }

    // Return an object so they can cancel it...
    ret
  }
}

/**
 * Wrapper for clearTimeout(ID)
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
    if argc!=1 {
      return ext_js_undefined
    }
    let mut id:i64 = 0;
    JS_ToInt64(context, &mut id as *mut i64, *argv.offset(0));

    ACTIVE_TIMERS.retain(|tim| {
      tim.id!=id as u64
    });

    ext_js_undefined
  }
}

/**
 * Wrapper for setInterval(func, delay)->ID
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
    if argc!=2 {
      return ext_js_undefined;
    }
    let now = SystemTime::now();

    let mut ret = ext_js_undefined;

    let func = *argv.offset(0);
    if JS_IsFunction(context, func)==1 {
      let mut delay:i64 = 0;
      JS_ToInt64(context, &mut delay as *mut i64, *argv.offset(1));

      let func2 = JS_DuplicateValue(context, func);

      let delay_duration = Duration::from_millis(delay as u64);
      let t = TimerInfo{
        id: TIMER_ID,
        active: true,
        delay: delay_duration,
        callback: func2,
        trigger_time: now.add(delay_duration),
        repeating: true,
      };

      ret = JS_NewInt64_Ext(context, TIMER_ID as i64);

      TIMER_ID+=1;
      ACTIVE_TIMERS.push(t);
    }

    // Return an object so they can cancel it...
    ret
  }
}

/**
 * Wrapper for clearInterval(ID)
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
    if argc!=1 {
      return ext_js_undefined
    }
    let mut id:i64 = 0;
    JS_ToInt64(context, &mut id as *mut i64, *argv.offset(0));

    ACTIVE_TIMERS.retain(|tim| {
      tim.id!=id as u64
    });

    ext_js_undefined
  }
}