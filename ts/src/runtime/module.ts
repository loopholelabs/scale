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
import { Host } from "./host";
import { Context } from "./context";

export interface WasiContext {
  start(instance: WebAssembly.Instance): void;
  getImportObject(): any;
}

/*
import { WASI } from "wasi";
const wasi = new WASI({
  args: [],
  env: {},
});

const w = {
  getImportObject: () => wasi.wasiImport;
  start: (instance: WebAssembly.Instance) => {
    wasi.start(instance);
  }
}

*/

export class Module {
  private _code: Buffer;

  private _next: Module | null = null;

  private _wasmInstance: WebAssembly.Instance;

  private _mem: WebAssembly.Memory;

  private _allocFn: Function;

  private _runFn: Function;

  constructor(code: Buffer, w: WasiContext) {
    this._code = code;
    const wasmMod = new WebAssembly.Module(this._code);

    let wasmModule: WebAssembly.Instance;
    let allocFn: Function;

    const nextFn = (
      (thisModule: Module) =>
      (ptr: number, len: number): BigInt => {
        const mem = wasmModule.exports.memory as WebAssembly.Memory;
        const c = Context.readFrom(mem, ptr, len);

        if (thisModule._next != null) {
          const rc = thisModule._next.run(c);
          if (rc == null) {
            console.log("Next module didn't seem to run correctly.");
            return Host.packMemoryRef(ptr, len);
          }
          const v = rc.writeTo(mem, allocFn);
          return Host.packMemoryRef(v.ptr, v.len);
        }
        return Host.packMemoryRef(ptr, len);
      }
    )(this);

    const importObject = {
      wasi_snapshot_preview1: w.getImportObject(),
      env: {
        next: nextFn,
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

    w.start(wasmModule);
  }

  // Set the next module to run
  setNext(mod: Module | null) {
    this._next = mod;
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
