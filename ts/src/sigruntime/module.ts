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

import { Signature } from "../signature/signature";

import { Instance } from "./instance";
import { Runtime } from "./runtime";
import { SFunction } from "./sfunction";

import { v4 as uuidv4 } from 'uuid';

export class Module<T extends Signature> {
  private runtime: Runtime<T>;
  public sfunction: SFunction<T>;

  private moduleName: string = "nothing";
  private waInstance: undefined | WebAssembly.Instance;
  public run: undefined | Function;
  public resize: undefined | Function;
  public memory: undefined | WebAssembly.Memory;

  constructor(f: SFunction<T>, r: Runtime<T>) {
    this.sfunction = f;
    this.runtime = r;
  }

  init(i: Instance<T>) {
    this.waInstance = this.runtime.instantiate(this.sfunction.mod, this, i);

    this.run = this.waInstance.exports.run as Function;
    this.resize = this.waInstance.exports.resize as Function;
    this.memory = this.waInstance.exports.memory as WebAssembly.Memory;
  }

}