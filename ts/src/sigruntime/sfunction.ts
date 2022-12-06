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
    const encoded = i.RuntimeContext().Write();
    // TODO: Call resize...
    console.log("RuntimeContext is ", encoded);

    const resizeFn = this.ins.exports.resize as Function;
    const mem = this.ins.exports.memory as WebAssembly.Memory;

    const encPtr = resizeFn(encoded.length);

    console.log("ResizeFn returned", encPtr);

    const memData = new Uint8Array(mem.buffer);
    memData.set(encoded, encPtr); // Writes the context into memory

    // Now run the function...
    const runFn = this.ins.exports.run as Function;

    const packed = runFn();

    console.log("Return from runFn was ", packed);

    const [ptr, len] = SFunction.unpackMemoryRef(packed);

    console.log("PTR = " + ptr + ", LEN = " + len);

    const memDataOut = new Uint8Array(mem.buffer);
    const inContextBuff = memDataOut.slice(ptr, ptr + len);

    console.log("Data is ", inContextBuff);
    i.RuntimeContext().Read(inContextBuff);
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