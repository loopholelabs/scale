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
import { CallbackFunction as TraceCallbackFunction } from "./tracing";
import {NewTemplate, Template} from "./template";
import {Module} from "./module";
import {Instance, NewInstance} from "./instance";

export type NextFn<T extends Signature> = (ctx: T) => T;

export async function New<T extends Signature>(config: Config<T>): Promise<Scale<T>> {
    const r = new Scale(config);
    await r.Ready();
    return r;
}

export class Scale<T extends Signature> {
    private readonly ready: Promise<void>;
    public readonly config: Config<T>;
    public head: undefined | Template<T>;
    private tail: undefined | Template<T>;

    public TraceDataCallback: TraceCallbackFunction | undefined;

    constructor(config: Config<T>) {
        this.config = config;
        this.ready = new Promise(async (resolve) => { // eslint-disable-line no-async-promise-executor
            this.config.validate()

            const testSignature = this.config.newSignature();
            for (let i = 0; i < this.config.functions.length; i++) {
                const sf = this.config.functions[i];
                if (testSignature.Hash() != sf.function.Signature.Hash) {
                    throw new Error(`passed in function ${sf.function.Name}:${sf.function.Tag} has an invalid signature`);
                }
                let t: Template<T>;
                try {
                    t = await NewTemplate<T>(this, sf.function, sf.env);
                } catch (e) {
                    throw new Error(`failed to pre-compile function ${sf.function.Name}:${sf.function.Tag}: ${e}`);
                }

                if (typeof this.head === "undefined") {
                    this.head = t;
                }

                if (typeof this.tail !== "undefined") {
                    this.tail.next = t;
                }

                this.tail = t;
            }
            resolve();
        });
    }

    async Ready() {
        await this.ready;
    }

    async Instance(next?: NextFn<T>): Promise<Instance<T>> {
        return NewInstance(this, next);
    }

    public Next(m: Module<T>): (ptr: number, len: number) => void {
        return (ptr: number, len: number): void => {
            if (typeof m.memory === "undefined") {
                throw new Error("no memory found in module");
            }
            if (typeof m.resize === "undefined") {
                throw new Error("no resize function found in module");
            }
            if (typeof m.signature === "undefined") {
                throw new Error("no signature found in module");
            }
            if (typeof m.function === "undefined") {
                throw new Error("no function found in module");
            }
            let buf: Uint8Array;
            try {
                const memDataOut = new Uint8Array(m.memory.buffer);
                const inContextBuff = memDataOut.slice(ptr, ptr + len);
                m.signature.Read(inContextBuff);

                if (typeof m.function.next !== "undefined") {
                    m.function.next.GetModule(m.signature).then((nextModule) => {
                      nextModule.Run();
                      if (typeof m.function !== "undefined" && typeof m.function.next !== "undefined") {
                          m.function.next.PutModule(nextModule);
                      }
                    });
                } else {
                    m.signature = m.function.instance.next(m.signature)
                }
                buf = m.signature.Write();
            } catch (e) {
                buf = m.signature.Error(e as Error);
            }
            const writeBufferPointer = m.resize(buf.length);
            const writeBuffer = new Uint8Array(m.memory.buffer);
            writeBuffer.set(buf, writeBufferPointer);
        }
    }
}
