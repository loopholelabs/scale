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

fn main() {
    #[cfg(all(not(feature = "runtime_source"), feature = "embedded_source"))]
    {
        let output_directory = std::env::var("OUT_DIR").unwrap();
        let js_source_input_path = std::env::var("JS_SOURCE_PATH").unwrap();
        let js_source_output_path = format!("{}/js_source", output_directory);

        fs::copy(&js_source_input_path, &js_source_output_path)
            .expect("unable to copy js source to output directory");
    }
}
