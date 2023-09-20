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

use binaryen::{CodegenConfig, Module};
use std::path::Path;
use wizer::Wizer;

pub(crate) struct Builder<'a> {
    interpreter: &'a [u8],
}

impl<'a> Builder<'a> {
    pub fn new(interpreter: &'a [u8]) -> Self {
        Self { interpreter }
    }
    pub fn build_interpreter(
        self,
        optimize: bool,
        out_path: impl AsRef<Path>,
    ) -> Result<(), Box<dyn std::error::Error>> {
        let mut interpreter = Wizer::new()
            .wasm_bulk_memory(true)
            .allow_wasi(true)?
            .inherit_stdio(true)
            .run(self.interpreter)?;

        if optimize {
            let codegen_cfg = CodegenConfig {
                optimization_level: 3, // Aggressively optimize for speed.
                shrink_level: 0,       // Don't optimize for size at the expense of performance.
                debug_info: false,     // Disable debug info.
            };

            if let Ok(mut optimized_interpreter) = Module::read(&interpreter) {
                optimized_interpreter.optimize(&codegen_cfg);
                optimized_interpreter
                    .run_optimization_passes(vec!["strip"], &codegen_cfg)
                    .unwrap();
                interpreter = optimized_interpreter.write();
            } else {
                Err("Unable to read interpreter wasm binary for wasm-opt optimizations")?
            }
        }

        std::fs::write(out_path.as_ref(), interpreter)?;
        Ok(())
    }
}
