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

mod builder;
mod options;

use std::{env, fs};
use structopt::StructOpt;

fn main() {
    let opts = options::Options::from_args();
    let as_stdin = env::var("AS_STDIN");

    if as_stdin.is_ok() {
        let interpreter: &[u8] = include_bytes!(concat!(env!("OUT_DIR"), "/interpreter.wasm"));
        builder::Builder::new(interpreter)
            .build_interpreter(true, opts.output)
            .unwrap();

        env::remove_var("AS_STDIN");
    } else {
        let js_source = fs::File::open(&opts.input).expect("Failed to open input file");
        let self_cmd = env::args().next().unwrap();
        env::set_var("AS_STDIN", "true");
        let status = std::process::Command::new(self_cmd)
            .arg(&opts.input)
            .arg("-o")
            .arg(&opts.output)
            .stdin(js_source)
            .status()
            .expect("Failed to execute child process");
        if !status.success() {
            panic!("Couldn't create wasm from input");
        }
    }
}
