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

import {v4 as uuidv4} from 'uuid';
import { ScaleFunc } from "../signature/scaleFunc";
import { Signature } from "../signature/signature";

import { Instance } from "./instance";
import { Pool } from "./pool";

export class SFunction<T extends Signature> {
  private scalefn: ScaleFunc;
  public mod: WebAssembly.Module;
  public next: undefined | SFunction<T>;
  public id: string;

  public modulePool: undefined | Pool<T>;

  constructor(scalefn: ScaleFunc, mod: WebAssembly.Module) {
    this.scalefn = scalefn;
    this.mod = mod;
    this.id = uuidv4();
  }

  Run(i: Instance<T>) {
    if (this.modulePool === undefined) {
      throw new Error("modulePool not set");
    }
    const module = this.modulePool.Get();

    module.init(i);

    const encoded = i.RuntimeContext().Write();

    if (module.resize === undefined || module.run === undefined || module.memory === undefined) {
      throw new Error("Module doesn't have resize/run/memory");
    }

    const encPtr = module.resize(encoded.length);

    const memData = new Uint8Array(module.memory.buffer);
    memData.set(encoded, encPtr);

    const packed = module.run(encPtr, encoded.length);    // Backward compat. TODO: Remove args

    const [ptr, len] = SFunction.unpackMemoryRef(packed);
    const memDataOut = new Uint8Array(module.memory.buffer);
    const inContextBuff = memDataOut.slice(ptr, ptr + len);

    i.RuntimeContext().Read(inContextBuff);   // NB This may throw an Error if the inContextBuff contains one.
    // TODO: Put module back in the pool?...
  }

  // Pack a pointer and length into a single 64bit
  public static packMemoryRef(ptr: number, len: number): BigInt {
    if (ptr > 0xffffffff || len > 0xffffffff) {
      // Error! We can't do it.
    }
    return (BigInt(ptr) << BigInt(32)) | BigInt(len);
  }

  // Unpack a memory ref from 64bit to 2x32bits
  public static unpackMemoryRef(packed: bigint): [number, number] {
    const ptr = Number((packed >> BigInt(32)) & BigInt(0xffffffff));
    const len = Number(packed & BigInt(0xffffffff));
    return [ptr, len];
  }

}