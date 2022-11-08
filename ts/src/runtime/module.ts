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

  private _wasi: WasiContext;

  private _next: Module | null = null;

  private _wasmInstance: WebAssembly.Instance | null = null;

  private _mem: WebAssembly.Memory | null = null;

  private _allocFn: Function | null = null;

  private _runFn: Function | null = null;

  private _importObject: any;

  constructor(code: Buffer, w: WasiContext) {
    this._code = code;
    this._wasi = w;

    const nextFn = (
      (thisModule: Module) =>
      (ptr: number, len: number): BigInt => {
        if (thisModule._wasmInstance != null && thisModule._allocFn != null) {        
          const mem = thisModule._wasmInstance.exports.memory as WebAssembly.Memory;
          const c = Context.readFrom(mem, ptr, len);

          if (thisModule._next != null) {
            const rc = thisModule._next.run(c);
            if (rc == null) {
              console.log("Next module didn't seem to run correctly.");
              return Host.packMemoryRef(ptr, len);
            }
            const v = rc.writeTo(mem, thisModule._allocFn);
            return Host.packMemoryRef(v.ptr, v.len);
          }
        }
        return Host.packMemoryRef(ptr, len);
      }
    )(this);

    this._importObject = {
      wasi_snapshot_preview1: w.getImportObject(),
      env: {
        next: nextFn,
      },
    };
  }

  async init() { 
    // This is done possibly in the background...
    const i = await WebAssembly.instantiate(this._code, this._importObject);
    this._wasmInstance = i.instance;

    //const wasmMod = new WebAssembly.Module(this._code);
    //this._wasmInstance = new WebAssembly.Instance(wasmMod, this._importObject);

    this._allocFn = this._wasmInstance.exports.malloc as Function;
    // If the module has a 'resize', use that instead of 'malloc'.
    if ("resize" in this._wasmInstance.exports) {
      this._allocFn = this._wasmInstance.exports.resize as Function;
    }

    this._mem = this._wasmInstance.exports.memory as WebAssembly.Memory;
    this._runFn = this._wasmInstance.exports.run as Function;

    this._wasi.start(this._wasmInstance);
  }

  // Set the next module to run
  setNext(mod: Module | null) {
    this._next = mod;
  }

  // Run this module, with an optional next module
  run(context: Context): Context | null {
    if (this._mem != null && this._allocFn != null && this._runFn != null) {
      const v = context.writeTo(this._mem, this._allocFn);
      const packed = this._runFn(v.ptr, v.len);
      if (packed === 0) {
        return null;
      }
      const [outContextPtr, outContextLen] = Host.unpackMemoryRef(packed);
      return Context.readFrom(this._mem, outContextPtr, outContextLen);
    }
    return null;
  }
}
