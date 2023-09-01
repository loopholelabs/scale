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

import { Signature, NewSignature } from "@loopholelabs/scale-signature";
import { ScaleFunc } from "@loopholelabs/scalefile/scalefunc";
import * as httpSignature  from "@loopholelabs/scale-signature-http";

import { randomBytes } from "crypto";

import { ScaleMod } from "./scalemod";

import { Func } from "./func";
import { Instance } from "./instance";
import { Pool } from "./pool";
import { Module } from "./module";
import { Cache } from "./cache";
import {DisabledWASI} from "./wasi";

import { TraceCallbackFunc } from "./scalemod";

export * from "./instance"
export * from "./module"
export * from "./func"
export * from "./wasi"

export type NextFn<T extends Signature> = (ctx: T) => T;

export async function New<T extends Signature>(newSignature: NewSignature<T>, functions: (Promise<ScaleFunc>|ScaleFunc|Func<T>)[]): Promise<Runtime<T>> {
    const r = new Runtime(newSignature, functions);
    await r.Ready();
    return r;
}

export class Runtime<T extends Signature> {
    public NewSignature: NewSignature<T>;
    private readonly ready: Promise<any>;
    private functions: Func<T>[];
    public head: undefined | Func<T>;
    public tail: undefined | Func<T>;

    public InvocationId: Buffer;
    public TraceDataCallback: TraceCallbackFunc | undefined;

    constructor(newSignature: NewSignature<T>, functions: (Promise<ScaleFunc>|ScaleFunc|Func<T>)[]) {
        this.NewSignature = newSignature;
        this.functions = [];
        this.InvocationId = Buffer.alloc(16);

        this.ready = new Promise(async (resolve) => {
            for (let i = 0; i < functions.length; i++) {
                const fn = functions[i];
                let f: Func<T>;
                if (fn instanceof ScaleFunc) {
                    const wasi = new DisabledWASI();
                    const scaleMod = new ScaleMod("unknown", Buffer.alloc(16), undefined);
                    const instantiatedSource = await WebAssembly.instantiate(fn.Function, {
                        wasi_snapshot_preview1: wasi.GetImports(),
                        env: {
                            next: (ptr: number, len: number): number => {
                                return 0;
                            }
                        },
                        scale: scaleMod.GetImports(),
                    });
                    f = new Func<T>(fn, instantiatedSource.module);
                } else if (fn instanceof Promise<ScaleFunc>) {
                    const wasi = new DisabledWASI();
                    const scaleMod = new ScaleMod("unknown", Buffer.alloc(16), undefined);
                    const resolvedFn = await fn;
                    const instantiatedSource = await WebAssembly.instantiate(resolvedFn.Function, {
                        wasi_snapshot_preview1: wasi.GetImports(),
                        env: {
                            next: (ptr: number, len: number): number => {
                                return 0;
                            }
                        },
                        scale: scaleMod.GetImports(),
                    });
                    f = new Func<T>(resolvedFn, instantiatedSource.module);
                } else {
                    f = fn
                }

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
        this.InvocationId = randomBytes(16);

        const i = new Instance<T>(this, next);
        for (let a = 0; a < this.functions.length; a++) {
            const module = this.functions[a].wasmModule;
            const id = this.functions[a].id;
            const cache = new Cache();

            let serviceName = this.functions[a].scaleFunc.Name + ":" + this.functions[a].scaleFunc.Tag;
            let scalemod = new ScaleMod(serviceName, this.InvocationId, this.TraceDataCallback);

            await cache.Initialize(module, scalemod);
            i.SetInstance(id, cache);
        }

        return i;
    }

    InstantiateModule(fnid: string, mod: Module<T>, i: Instance<T>): WebAssembly.Instance {
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

                    if (mod.func.next === undefined) {
                        i.ctx = i.next(i.Context());
                        buff = i.RuntimeContext().Write();
                    } else {
                        mod.func.next.Run(i);
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
        return cached.GetInstance();
    }
}
