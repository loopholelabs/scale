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
import { ScaleFunc } from "../signature/scaleFunc";

import { SFunction } from "./sfunction";
import { Instance } from "./instance";
import { Pool } from "./pool";
import { Module } from "./module";

export interface WasiContext {
  start(instance: WebAssembly.Instance): void;
  getImportObject(): any;
}

export type NextFn<T extends Signature> = (ctx: T) => T;

export class Runtime<T extends Signature> {
  public Ready: Promise<any>;
  
  public signature: T;
  private fns: SFunction<T>[];
  public head: undefined | SFunction<T>;
  private tail: undefined | SFunction<T>;

  private wasiBuilder: () => WasiContext;

  public modules: Map<string, Module<T>>;   // Map from unique module ID to module

  constructor(wasiBuilder: () => WasiContext, sig: T, fns: ScaleFunc[]) {
    this.signature = sig;
    this.fns = [];
    this.modules = new Map<string, Module<T>>;

    this.wasiBuilder = wasiBuilder;

    this.Ready = new Promise(async (resolve, reject) => {

      for(let i=0;i<fns.length;i++) {
        const fn = fns[i];
        const mod = await WebAssembly.compile(fn.Function as Buffer);

        const f = new SFunction<T>(fn, mod);
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

  Instance(next: null | NextFn<T>): Instance<T> {
    return new Instance<T>(this, next);
  }

  instantiate(m: WebAssembly.Module, mod: Module<T>, i: Instance<T>): WebAssembly.Instance {
    // NB This closure captures i.
    const nextFn = ((runtimeThis: Runtime<T>): Function => {
      return (ptr: number, len: number): BigInt => {
        if (mod.memory===undefined || mod.resize===undefined) return BigInt(0);   // TODO
        const memDataOut = new Uint8Array(mod.memory.buffer);
        const inContextBuff = memDataOut.slice(ptr, ptr + len);
        i.RuntimeContext().Read(inContextBuff);

        // Now call next...
        let buff: Uint8Array = new Uint8Array();
        if (mod.sfunction.next === undefined) {
          try {
            i.ctx = i.next(i.Context());
            buff = i.RuntimeContext().Write();
          } catch(e) {
            buff = i.RuntimeContext().Error(e as Error);
          }
        } else {
          try {
            // TODO: This should be await?
            mod.sfunction.next.Run(i);
            buff = i.RuntimeContext().Write();
          } catch(e) {
            buff = i.RuntimeContext().Error(e as Error);
          }
        }

        // Write it back out
        const encPtr = mod.resize(buff.length);

        const memData = new Uint8Array(mod.memory.buffer);
        memData.set(buff, encPtr);

        return SFunction.packMemoryRef(encPtr, buff.length);
      };
    })(this);
    
    const wasi = this.wasiBuilder();
    const importObject = {
      wasi_snapshot_preview1: wasi.getImportObject(),
      env: {
        next: nextFn,
      },
    };

    const inst = new WebAssembly.Instance(m, importObject);

    wasi.start(inst);
    return inst;

    //return WebAssembly.instantiate(m, importObject);
  }
}