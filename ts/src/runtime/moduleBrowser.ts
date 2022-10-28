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

export class ModuleBrowser {
  private _code: Buffer;

  private _wasmMod: WebAssembly.Module;

  private _next: ModuleBrowser | null;

  constructor(code: Buffer, next: ModuleBrowser | null) {
    this._code = code;
    this._next = next;
    this._wasmMod = new WebAssembly.Module(this._code);
  }

  // Run this module, with an optional next module
  run(context: Context): Context {
    let wasmModule: WebAssembly.Instance;

    const nextModule = this._next;

    const importObject = {
      wasi_snapshot_preview1: {
        fd_write: () => {},
      },
      env: {
        next: (ptr: number, len: number): BigInt => {
          const mem = wasmModule.exports.memory as WebAssembly.Memory;
          const c = Context.readFrom(mem, ptr, len);

          if (nextModule != null) {
            const rc = nextModule.run(c);
            const v = rc.writeTo(mem, wasmModule.exports.malloc as Function);
            return Host.packMemoryRef(v.ptr, v.len);
          }
          return Host.packMemoryRef(ptr, len);
        },
      },
    };

    wasmModule = new WebAssembly.Instance(this._wasmMod, importObject);

    const mem = wasmModule.exports.memory as WebAssembly.Memory;
    const v = context.writeTo(mem, wasmModule.exports.malloc as Function);

    const runfn = wasmModule.exports.run as Function;
    const packed = runfn(v.ptr, v.len);
    const [outContextPtr, outContextLen] = Host.unpackMemoryRef(packed);
    return Context.readFrom(mem, outContextPtr, outContextLen);
  }
}
