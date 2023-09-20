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

import { randomBytes } from "crypto";

import { Signature  } from "@loopholelabs/scale-signature-interfaces";
import { NextFn, Scale } from "./scale";
import {Func, NewFunc} from "./function";

const ErrNoCompiledFunctions = new Error("no compiled functions found in runtime");

export async function NewInstance<T extends Signature>(runtime: Scale<T>, next?: NextFn<T>): Promise<Instance<T>> {
  const i = new Instance(runtime, next);
  await i.Ready();
  return i;
}

export class Instance<T extends Signature> {
  private readonly ready: Promise<void>;

  private runtime: Scale<T>;
  public identifier: Buffer;
  private head: undefined | Func<T>;
  next: NextFn<T>;

  constructor(runtime: Scale<T>, next?: NextFn<T>) {
    this.runtime = runtime;
    this.identifier = randomBytes(16);
    if (typeof next !== "undefined") {
      this.next = next;
    } else {
      this.next = (ctx: T) => ctx;
    }

    this.ready = new Promise(async (resolve) => { // eslint-disable-line no-async-promise-executor
      let previousFunction = this.head;
      let nextTemplate = this.runtime.head;

      while (typeof nextTemplate !== "undefined") {
        const fn = await NewFunc(this, nextTemplate);
        if (typeof this.head === "undefined") {
            this.head = fn;
        }
        if (typeof previousFunction !== "undefined") {
          previousFunction.next = fn;
        }
        previousFunction = fn;
        nextTemplate = nextTemplate.next;
      }

      resolve();
    });

  }

  public async Ready() {
    await this.ready;
  }

  public async Run(signature: T) {
    if (typeof this.head === "undefined") {
      throw ErrNoCompiledFunctions;
    }
    const m = await this.head.GetModule(signature);
    m.Run();
    this.head.PutModule(m);
  }
}