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

import { Signature, RuntimeContext } from "../signature/signature";
import { Runtime, NextFn } from "./runtime";

export class Instance<T extends Signature> {
  private runtime: Runtime<T>;
  public next: NextFn<T>;
  public ctx: T;

  constructor(r: Runtime<T>, n: null | NextFn<T>) {
    this.runtime = r;
    this.ctx = r.signature;

    if (n === null) {
      this.next = (ctx: T) => ctx;
    } else {
      this.next = n;
    }
  }

  Context(): T {
    return this.ctx;
  }

  RuntimeContext(): RuntimeContext {
    return this.runtime.signature.RuntimeContext();
  }

  async Run() {
    if (this.runtime.head === undefined) {
      throw (new Error("no compiled functions found in runtime"));
    }
    const fn = this.runtime.head;
    fn.Run(this);
  }
}