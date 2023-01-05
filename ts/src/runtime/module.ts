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

import { Signature } from "@loopholelabs/scale-signature";

import { Instance } from "./instance";
import { Runtime } from "./runtime";
import { Func } from "./func";

export class Module<T extends Signature> {
  private runtime: Runtime<T>;
  public sfunction: Func<T>;

  private wasmInstance: undefined | WebAssembly.Instance;
  public run: undefined | Function;
  public resize: undefined | Function;
  public memory: undefined | WebAssembly.Memory;

  constructor(f: Func<T>, r: Runtime<T>) {
    this.sfunction = f;
    this.runtime = r;
  }

  init(i: Instance<T>) {
    this.wasmInstance = this.runtime.instantiate(this.sfunction.id, this, i);

    this.run = this.wasmInstance.exports.run as Function;
    this.resize = this.wasmInstance.exports.resize as Function;
    if (this.resize === undefined) this.resize = this.wasmInstance.exports.malloc as Function;    // Backward compat. TODO: Remove
    this.memory = this.wasmInstance.exports.memory as WebAssembly.Memory;
  }

}