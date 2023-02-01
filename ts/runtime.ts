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

import { Signature, SignatureFactory } from "@loopholelabs/scale-signature";
import { ScaleFunc } from "@loopholelabs/scalefile";

import { Func } from "./func";
import { Instance } from "./instance";
import { Pool } from "./pool";
import { Module } from "./module";
import { Cache } from "./cache";

export type NextFn<T extends Signature> = (ctx: T) => T;

export async function New<T extends Signature>(newSignature: SignatureFactory<T>, functions: ScaleFunc[]): Promise<Runtime<T>> {
  const r = new Runtime(newSignature, functions);
  await r.Ready();
  return r;
}

export class Runtime<T extends Signature> {
  public NewSignature: SignatureFactory<T>;
  private readonly ready: Promise<any>;
  private functions: Func<T>[];
  public head: undefined | Func<T>;
  public tail: undefined | Func<T>;

  constructor(newSignature: SignatureFactory<T>, functions: ScaleFunc[]) {
    this.NewSignature = newSignature;
    this.functions = [];

    this.ready = new Promise(async (resolve) => {
      for (let i = 0; i < functions.length; i++) {
        const fn = functions[i];
        const mod = await WebAssembly.compile(fn.Function as Buffer);

        const f = new Func<T>(fn, mod);
        f.modulePool = new Pool<T>(f, this);
        this.functions.push(f);

        if (this.head === undefined) {
          this.head = f;
        }
        if (this.tail !== undefined) {
          this.tail.next = f;
        }
        this.tail = f;
      }

      resolve(true);
    });
  }

  async Ready() {
    await this.ready;
  }

  async Instance(next: null | NextFn<T>): Promise<Instance<T>> {
    const i = new Instance<T>(this, next);
    for (let a = 0; a < this.functions.length; a++) {
      const mod = this.functions[a].mod;
      const id = this.functions[a].id;
      const cache = new Cache();
      await cache.Initialize(mod);
      i.SetInstance(id, cache);
    }

    return i;
  }

  Instantiate(fnid: string, mod: Module<T>, i: Instance<T>): WebAssembly.Instance {
    const nextFunction = ((): Function => {
      return (ptr: number, len: number): BigInt => {
        if (mod.memory === undefined || mod.resize === undefined) {
          return BigInt(0);
        }
        let buff: Uint8Array;
        try {
          const memDataOut = new Uint8Array(mod.memory.buffer);
          const inContextBuff = memDataOut.slice(ptr, ptr + len);

          i.RuntimeContext().Read(inContextBuff);

          if (mod.sfunction.next === undefined) {
            i.ctx = i.next(i.Context());
            buff = i.RuntimeContext().Write();
          } else {
            mod.sfunction.next.Run(i);
            buff = i.RuntimeContext().Write();
          }
        } catch (e) {
          buff = i.RuntimeContext().Error(e as Error);
        }

        const encPtr = mod.resize(buff.length);
        const memData = new Uint8Array(mod.memory.buffer);
        memData.set(buff, encPtr);
        return Func.packMemoryRef(encPtr, buff.length);
      };
    })();

    const cached = i.GetInstance(fnid);
    cached.SetNext(nextFunction);
    return cached.getInstance();
  }
}