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

export interface WasiContext {
  start(instance: WebAssembly.Instance): void;
  getImportObject(): any;
}

export type NextFn<T extends Signature> = (ctx: T) => T;

export class Runtime<T extends Signature> {
  public Ready: Promise<any>;
  
  public signature: T;
  private fns: SFunction<T>[];
  private head: undefined | SFunction<T>;
  private tail: undefined | SFunction<T>;
  
  constructor(wasi: WasiContext, sig: T, fns: ScaleFunc[]) {
    this.signature = sig;
    this.fns = [];

    this.Ready = new Promise(async (resolve, reject) => {
      const nextFn = (ptr: number, len: number): BigInt => {
        console.log("Hello from next()");
        return BigInt(0);
      };

      const importObject = {
        wasi_snapshot_preview1: wasi.getImportObject(),
        env: {
          next: nextFn,
        },
      };

      for(let i=0;i<fns.length;i++) {
        const fn = fns[i];
        const ins = await WebAssembly.instantiate(fn.Function as Buffer, importObject);
        wasi.start(ins.instance);

        const f = new SFunction<T>(fn, ins.instance);
        this.fns.push(f);
      }

      resolve(true);
    });

  }

  Instance(next: null | NextFn<T>): Instance<T> {
    return new Instance<T>(this, next);
  }
}