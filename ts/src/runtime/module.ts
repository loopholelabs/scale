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
import { WASI } from "wasi";
import { Host } from "./host";
import { Context } from "./context";

export class Module {
  private _code: Buffer;

  private _next: Module | null;

  private _wasmInstance: WebAssembly.Instance;

  private _mem: WebAssembly.Memory;

  private _allocFn: Function;

  private _runFn: Function;

  constructor(code: Buffer, next: Module | null) {
    this._code = code;
    this._next = next;
    const wasmMod = new WebAssembly.Module(this._code);

    const wasi = new WASI({
      args: [],
      env: {},
    });

    let wasmModule: WebAssembly.Instance;
    let allocFn: Function;

    const nextModule = this._next;

    const importObject = {
      wasi_snapshot_preview1: wasi.wasiImport,
      env: {
        next: (ptr: number, len: number): BigInt => {
          const mem = wasmModule.exports.memory as WebAssembly.Memory;
          const c = Context.readFrom(mem, ptr, len);

          if (nextModule != null) {
            const rc = nextModule.run(c);
            if (rc == null) {
              console.log("Next module didn't seem to run correctly.");
              return Host.packMemoryRef(ptr, len);
            }
            const v = rc.writeTo(mem, allocFn);
            return Host.packMemoryRef(v.ptr, v.len);
          }
          return Host.packMemoryRef(ptr, len);
        },
      },
    };

    wasmModule = new WebAssembly.Instance(wasmMod, importObject);

    allocFn = wasmModule.exports.malloc as Function;
    // If the module has a 'resize', use that instead of 'malloc'.
    if ("resize" in wasmModule.exports) {
      allocFn = wasmModule.exports.resize as Function;
    }

    this._wasmInstance = wasmModule;
    this._mem = wasmModule.exports.memory as WebAssembly.Memory;
    this._runFn = wasmModule.exports.run as Function;
    this._allocFn = allocFn;

    wasi.start(wasmModule);
  }

  // Run this module, with an optional next module
  run(context: Context): Context | null {
    const v = context.writeTo(this._mem, this._allocFn);
    const packed = this._runFn(v.ptr, v.len);
    if (packed === 0) {
      return null;
    }
    const [outContextPtr, outContextLen] = Host.unpackMemoryRef(packed);
    return Context.readFrom(this._mem, outContextPtr, outContextLen);
  }
}
