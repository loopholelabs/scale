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
import { argv, env } from 'node:process';
import { Host } from './host';
import { Context } from "./context"

const WASI = require("wasi");

export class Module {
    private _code: Buffer;
    private _wasmMod: WebAssembly.Module;
    private _next: Module | null;

    constructor(code: Buffer, next: Module | null) {
        this._code = code;
        this._next = next;
        this._wasmMod = new WebAssembly.Module(this._code);
    }

    // Run this module, with an optional next module
    run(context: Context): Context | null{
        const wasi = new WASI({
            args: argv,
            env,
            preopens: {
                '/sandbox': '/some/real/path/that/wasm/can/access'
            }
        });

        let wasmModule: WebAssembly.Instance;
        let allocFn: Function;

        let nextModule = this._next;

        const importObject = {
            wasi_snapshot_preview1: wasi.exports,
            env: {
                next: function(ptr: number, len: number): BigInt {
                    const mem = wasmModule.exports.memory as WebAssembly.Memory;
                    let c = Context.readFrom(mem, ptr, len);

                    if (nextModule != null) {
                        let rc = nextModule.run(c);
                        if (rc==null) {
                            console.log("Next module didn't seem to run correctly.");
                            return Host.packMemoryRef(ptr, len);
                        }
                        let v = rc.writeTo(mem, allocFn);
                        return Host.packMemoryRef(v.ptr, v.len);
                    } else {
                        return Host.packMemoryRef(ptr, len);
                    }
                }
            }
        };

        wasmModule = new WebAssembly.Instance(this._wasmMod, importObject);

        const mem = wasmModule.exports.memory as WebAssembly.Memory;
        allocFn = wasmModule.exports.malloc as Function;

        // If the module has a 'resize', use that instead of 'malloc'.
        if ("resize" in wasmModule.exports) {
            allocFn = wasmModule.exports.resize as Function;
        }

        let v = context.writeTo(mem, allocFn);

        const runfn = wasmModule.exports.run as Function;
        let packed = runfn(v.ptr, v.len);

        if (packed ==0) {
            return null;
        }

        let [outContextPtr, outContextLen] = Host.unpackMemoryRef(packed);
        return Context.readFrom(mem, outContextPtr, outContextLen);
    }
}
