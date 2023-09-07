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
import {Signature} from "@loopholelabs/scale-signature-interfaces";
import {Func} from "./function";
import {Template} from "./template";
import {UnpackUint32} from "./utils";
import {Decoder} from "@loopholelabs/polyglot";
import {DisabledWASI} from "./wasi";

export async function NewModule<T extends Signature>(template: Template<T>): Promise<Module<T>> {
    const m = new Module(template);
    await m.Ready();
    return m;
}

export class Module<T extends Signature> {
    private readonly ready: Promise<void>;

    private template: Template<T>;

    private readonly wasi: DisabledWASI;

    private instantiatedModule: WebAssembly.Instance | undefined;
    private run: undefined | ((ptr: number, len: number) => bigint);
    public resize: undefined | ((len: number) => number);
    private initialize: undefined | (() => bigint);

    public memory: undefined | WebAssembly.Memory;

    public function: Func<T> | undefined;
    public signature: T | undefined;

    constructor(template: Template<T>) {
        this.template = template;
        this.wasi = new DisabledWASI(this.template.env);

        const moduleConfig = {
            wasi_snapshot_preview1: this.wasi.GetImports(),
            scale: this.template.tracing.GetImports(),
            env: {
                next: this.template.runtime.Next(this),
            },
        }

        this.ready = new Promise(async (resolve) => { // eslint-disable-line no-async-promise-executor
            if (this.template.compiled !== undefined) {
                this.instantiatedModule = await WebAssembly.instantiate(this.template.compiled, moduleConfig);
                this.wasi.SetInstance(this.instantiatedModule);

                this.run = this.instantiatedModule.exports.run as ((ptr: number, len: number) => bigint) | undefined;
                if (this.run === undefined) {
                    throw new Error("no run function found in module");
                }
                this.resize = this.instantiatedModule.exports.resize as ((len: number) => number) | undefined;
                if (this.resize === undefined) {
                    throw new Error("no resize function found in module");
                }
                this.initialize = this.instantiatedModule.exports.initialize as (() => bigint) | undefined;
                if (this.initialize === undefined) {
                    throw new Error("no initialize function found in module");
                }
                this.memory = this.instantiatedModule.exports.memory as WebAssembly.Memory | undefined;
                if (this.memory === undefined) {
                    throw new Error("no memory found in module");
                }

                const packed = this.initialize();
                if (packed != BigInt(0)) {
                    const [ptr, len] = UnpackUint32(packed);
                    const readData = new Uint8Array(this.memory.buffer);
                    const readBuffer = readData.slice(ptr, ptr + len);

                    const dec = new Decoder(readBuffer)
                    throw dec.error();
                }

                resolve();
            }
            throw new Error("no compiled module found in template");
        });
    }

    public async Ready() {
        return await this.ready;
    }

    public Run() {
        if (this.signature === undefined) {
            throw new Error("no signature found in module");
        }

        if (this.resize === undefined) {
            throw new Error("no resize function found in module");
        }

        if (this.run === undefined) {
            throw new Error("no run function found in module");
        }

        if (this.memory === undefined) {
            throw new Error("no memory found in module");
        }


        const buffer = this.signature.Write();
        const writeBufferPointer = this.resize(buffer.length);
        const writeBuffer = new Uint8Array(this.memory.buffer);
        writeBuffer.set(buffer, writeBufferPointer);

        const packed = this.run(writeBufferPointer, buffer.length);
        const [ptr, len] = UnpackUint32(packed);
        const readBufferPointer = new Uint8Array(this.memory.buffer);
        const readBuffer = readBufferPointer.slice(ptr, ptr + len);

        const err = this.signature.Read(readBuffer);
        if (err !== undefined) {
            throw err;
        }
    }

    public Register(fn: Func<T>) {
        this.function = fn;
    }

    public Cleanup() {
        this.function = undefined;
    }

    public SetSignature(signature: T) {
        this.signature = signature;
    }

}