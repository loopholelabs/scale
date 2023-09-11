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

use std::path::PathBuf;
use structopt::StructOpt;

#[derive(Debug, StructOpt)]
#[structopt(
    name = "js_builder",
    about = "A CLI for compiling JS Source into a Scale Function"
)]
pub struct Options {
    #[structopt(short = "o", parse(from_os_str), default_value = "scale.wasm")]
    pub output: PathBuf,
}
