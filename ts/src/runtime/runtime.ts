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
import { CachedWasmInstance } from "./cachedWasmInstance";

import { getNewWasi } from "./wasi";

export interface WasiContext {
  start(instance: WebAssembly.Instance): void;
  getImportObject(): any;
}

export type NextFn<T extends Signature> = (ctx: T) => T;

export async function GetRuntime<T extends Signature>(sigfac: SignatureFactory<T>, fns: ScaleFunc[]) {
  const wasiBuilder = getNewWasi;
  const r = new Runtime(wasiBuilder, sigfac, fns);
  await r.Ready;
  return r;
}

export class Runtime<T extends Signature> {
  public Ready: Promise<any>;

  public signatureFactory: SignatureFactory<T>;
  private fns: Func<T>[];
  public head: undefined | Func<T>;
  private tail: undefined | Func<T>;

  private wasiBuilder: () => WasiContext;

  constructor(wasiBuilder: () => WasiContext, sigfac: SignatureFactory<T>, fns: ScaleFunc[]) {
    this.signatureFactory = sigfac;
    this.fns = [];

    this.wasiBuilder = wasiBuilder;

    // We compile the modules async...
    // After creating a Runtime you should then do 'await runtime.Ready' or equivalent.
    this.Ready = new Promise(async (resolve, reject) => {

      for (let i = 0; i < fns.length; i++) {
        const fn = fns[i];
        const mod = await WebAssembly.compile(fn.Function as Buffer);

        const f = new Func<T>(fn, mod);
        f.modulePool = new Pool<T>(f, this);
        this.fns.push(f);

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

  async Instance(next: null | NextFn<T>): Promise<Instance<T>> {
    const i = new Instance<T>(this, next);
    // Create instances for the functions and save them whithin the Instance.
    for (let a = 0; a < this.fns.length; a++) {
      const mod = this.fns[a].mod;
      const id = this.fns[a].id;
      const cached = new CachedWasmInstance(this.wasiBuilder);
      await cached.init(mod);
      i.setInstance(id, cached);    // Store the instance here
    }

    return i;
  }

  instantiate(fnid: string, mod: Module<T>, i: Instance<T>): WebAssembly.Instance {

    // NB This closure captures i and mod
    const nextFn = ((runtimeThis: Runtime<T>): Function => {
      return (ptr: number, len: number): BigInt => {
        if (mod.memory === undefined || mod.resize === undefined) {
          // Critical unrecoverable error
          // NB This would only ever happen if init() wasn't called on the Module.
          return BigInt(0);
        }

        let buff: Uint8Array = new Uint8Array();
        try {
          const memDataOut = new Uint8Array(mod.memory.buffer);
          const inContextBuff = memDataOut.slice(ptr, ptr + len);
          i.RuntimeContext().Read(inContextBuff);

          // Now call next...
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

        // Write it back out
        const encPtr = mod.resize(buff.length);
        const memData = new Uint8Array(mod.memory.buffer);
        memData.set(buff, encPtr);
        return Func.packMemoryRef(encPtr, buff.length);
      };
    })(this);

    //
    const cached = i.getInstance(fnid);
    cached.setNext(nextFn);
    return cached.getInstance();
  }
}