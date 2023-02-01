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

import { NewWASI } from "./wasi";

export class Cache {
  private readonly wasiBuilder = NewWASI
  private instance: undefined | WebAssembly.Instance;
  private nextFunction: undefined | Function;
  constructor() {}

  async Initialize(m: WebAssembly.Module) {
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

  private next(ptr: number, len: number): number {
    if (this.nextFunction === undefined) {
      return 0;
    } else {
      return this.nextFunction(ptr, len);
    }
  }

  // Set the next function
  SetNext(fn: Function) {
    this.nextFunction = fn;
  }

  getInstance(): WebAssembly.Instance {
    if (this.instance === undefined) {
      throw new Error("Instance wasn't created correctly");
    }
    return this.instance;
  }
}