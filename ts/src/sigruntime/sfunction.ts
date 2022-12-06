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

import { ScaleFunc } from "../signature/scaleFunc";
import { Signature } from "../signature/signature";

import { Instance } from "./instance";

export class SFunction<T extends Signature> {
  private scalefn: ScaleFunc;
  private ins: WebAssembly.Instance;
  public next: undefined | SFunction<T>;

  constructor(scalefn: ScaleFunc, ins: WebAssembly.Instance) {
    this.scalefn = scalefn;
    this.ins = ins;

    console.log("SFunction setup for ", scalefn);
    console.log("SFunction setup for WebAssembly.Instance ", ins);
  }

  Run(i: Instance<T>) {
    console.log("TODO: SFunction Run");
  }
}