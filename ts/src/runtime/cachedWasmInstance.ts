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

import { WasiContext } from "./runtime";

export class CachedWasmInstance {
  private wasiBuilder: () => WasiContext;

  private instance: undefined | WebAssembly.Instance;
  private nextFn: undefined | Function;

  constructor(wasiBuilder: () => WasiContext) {
    this.wasiBuilder = wasiBuilder;
  }

  // Initialize the instance
  async init(m: WebAssembly.Module) {
    const wasi = this.wasiBuilder();
    const importObject = {
      wasi_snapshot_preview1: wasi.getImportObject(),
      env: {
        next: this.next.bind(this),
      },
    };

    this.instance = await WebAssembly.instantiate(m, importObject);
    wasi.start(this.instance);
  }

  next(ptr: number, len: number): number {
    if (this.nextFn === undefined) {
      console.log("nextFn was not set!");
      return 0;   // TODO
    } else {
      return this.nextFn(ptr, len);
    }
  }

  // Set the next function
  setNext(fn: Function) {
    this.nextFn = fn;
  }

  getInstance(): WebAssembly.Instance {
    if (this.instance === undefined) {
      throw new Error("Instance wasn't created correctly");
    }
    return this.instance;
  }
}