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

import { Signature } from "@loopholelabs/scale-signature-interfaces";
import { Config } from "./config";
import { Module as TracingModule, CallbackFunction as TraceCallbackFunction } from "./tracing";
import { Func, NewFunc } from "./function";
import { DisabledWASI } from "./wasi";

export type NextFn<T extends Signature> = (ctx: T) => T;

export async function New<T extends Signature>(config: Config<T>): Promise<Scale<T>> {
    const r = new Scale(config);
    await r.Ready();
    return r;
}

export class Scale<T extends Signature> {
    private readonly ready: Promise<any>;

    public readonly moduleConfig: WebAssembly.Imports;
    private readonly config: Config<T>;

    public head: undefined | Func<T>;
    public tail: undefined | Func<T>;

    public TraceDataCallback: TraceCallbackFunction | undefined;

    constructor(config: Config<T>) {
        this.config = config;
        const wasi = new DisabledWASI();
        const tracing = new TracingModule("unknown", Buffer.alloc(16), undefined);
        this.moduleConfig = {
            wasi_snapshot_preview1: wasi.GetImports(),
            scale: tracing.GetImports(),
            env: {
                next: (ptr: number, len: number): number => {
                    return 0;
                }
            },
        }

        this.ready = new Promise(async (resolve) => {
            const testSignature = this.config.newSignature();

            for (let i = 0; i < this.config.functions.length; i++) {
                const fn = this.config.functions[i];

                if (testSignature.Hash() != fn.function.SignatureHash) {
                    throw new Error(`passed in function ${fn.function.Name}:${fn.function.Tag} has an invalid signature`);
                }

                let f: Func<T>;
                try {
                    f = await NewFunc<T>(this, fn.function, this.moduleConfig, fn.env);
                } catch (e) {
                    throw new Error(`failed to pre-compile function ${fn.function.Name}:${fn.function.Tag}: ${e}`);
                }

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
