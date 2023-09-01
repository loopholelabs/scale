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

import { v4 as uuidv4 } from "uuid";
import { Signature } from "@loopholelabs/scale-signature-interfaces";
import { ScaleFunc } from "./scalefunc/scalefunc";
import { Scale } from "./scale";
import { DisabledWASI } from "./wasi";
import {ModuleInstance} from "./module";

export class Func<T extends Signature> {
    private runtime: Scale<T>

    private identifier: string;

    private compiled: WebAssembly.Module;

    private scaleFunc: ScaleFunc;

    public next: undefined | Func<T>;

    private env: { [key: string]: string } | undefined;

    constructor(r: Scale<T>, scaleFunc: ScaleFunc, compiled: WebAssembly.Module, env?: { [key: string]: string }) {
        this.runtime = r;
        this.identifier = `${scaleFunc.Name}:${scaleFunc.Tag}`;
        this.scaleFunc = scaleFunc;
        this.env = env;
        this.compiled = compiled;
    }

    public RunWithModule(moduleInstance: ModuleInstance<T>) {
        if (moduleInstance.signature === undefined) {
            throw new Error("module instance doesn't have signature");
        }
        const encoded = moduleInstance.signature.Write();

        const encPtr = moduleInstance.resize(encoded.length);

        const memData = new Uint8Array(module.memory.buffer);
        memData.set(encoded, encPtr);

        const packed = module.run(encPtr, encoded.length);

        const [ptr, len] = Func.unpackMemoryRef(packed);
        const memDataOut = new Uint8Array(module.memory.buffer);
        const inContextBuff = memDataOut.slice(ptr, ptr + len);

        const err = i.RuntimeContext().Read(inContextBuff);
        if (err !== undefined) {
            throw err;
        }
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

export async function NewFunc<T extends Signature>(r: Scale<T>, scaleFunc: ScaleFunc, moduleConfig: WebAssembly.Imports, env?: { [key: string]: string }): Promise<Func<T>> {
    let compiled: WebAssembly.WebAssemblyInstantiatedSource;
    try {
        compiled = await WebAssembly.instantiate(scaleFunc.Function, moduleConfig);
    } catch (e) {
        throw new Error(`failed to compile wasm module: ${e}`)
    }
    return new Func<T>(r, scaleFunc, compiled.module, env);
}