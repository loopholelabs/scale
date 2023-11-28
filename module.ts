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
import {v4 as uuid} from "uuid";
import {NamedLogger} from "./log";
import {Tracing} from "./tracing";

export async function NewModule<T extends Signature>(template: Template<T>): Promise<Module<T>> {
    const m = new Module(template);
    await m.Ready();
    return m;
}

export class Module<T extends Signature> {
    private readonly ready: Promise<void>;

    public template: Template<T>;

    private readonly wasi: DisabledWASI;
    private readonly tracing: Tracing<T>;

    private instantiatedModule: WebAssembly.Instance | undefined;
    private run: undefined | ((ptr: number, len: number) => bigint);
    public resize: undefined | ((len: number) => number);
    private initialize: undefined | (() => bigint);

    public memory: undefined | WebAssembly.Memory;

    public function: Func<T> | undefined;
    public signature: T | undefined;

    constructor(template: Template<T>) {
        this.template = template;

        const name = `${this.template.identifier}.${uuid()}`

        let stdout = this.template.runtime.config.stdout;
        let stderr = this.template.runtime.config.stderr;

        if (typeof stdout !== "undefined" && !this.template.runtime.config.rawOutput) {
            stdout = NamedLogger(name, stdout);
        }

        if (typeof stderr !== "undefined" && !this.template.runtime.config.rawOutput) {
            stderr = NamedLogger(name, stderr);
        }

        this.wasi = new DisabledWASI(this.template.env, stdout, stderr);
        this.tracing = new Tracing(this, this.template.runtime.TraceDataCallback);

        const envValue:WebAssembly.ModuleImports = {};

        envValue["next"] = this.template.runtime.Next(this);

        const mem = function(m: Module<T>) {
          return {
            Write: (offset:number, v:Uint8Array): void => {
              if (m.memory==undefined) {
                throw new Error("no memory found in module");
              }
              const writeData = new Uint8Array(m.memory.buffer);
              writeData.set(v, offset);
            },
            Read: (offset:number, byteCount:number): Uint8Array => {
              if (m.memory==undefined) {
                throw new Error("no memory found in module");
              }
              const readData = new Uint8Array(m.memory.buffer);
              return readData.slice(offset, offset + byteCount);
            },
          }
        }(this)

        const resize = function(m: Module<T>) {
            return (name:string, size:number): number => {
              if (m.instantiatedModule==undefined) {
                throw new Error("no resize function found in module " + name);
              }
              const resize_function = m.instantiatedModule.exports[name] as ((len: number) => number) | undefined;
              if (resize_function==undefined) {
                throw new Error("no resize function found in module " + name);
              }
              return resize_function(size);
          }
        }(this);

        // Add any extensions...
        for (const ext of this.template.runtime.config.extensions) {
          const fns = ext.Init();
          for (const [n, fn] of fns) {
            envValue[n] = (instance: number, ptr: number, len: number): bigint => {
              const params: number[] = [instance, ptr, len];
              fn(mem, resize, params);
              return BigInt(params[0]);
            };
          }
        }

        const moduleConfig = {
          wasi_snapshot_preview1: this.wasi.GetImports(),
          scale: this.tracing.GetImports(),
          env: envValue,
        }

        this.ready = new Promise(async (resolve) => { // eslint-disable-line no-async-promise-executor
            if (typeof this.template.compiled !== "undefined") {
                this.instantiatedModule = await WebAssembly.instantiate(this.template.compiled, moduleConfig);
                this.wasi.SetInstance(this.instantiatedModule);

                this.run = this.instantiatedModule.exports.run as ((ptr: number, len: number) => bigint) | undefined;
                if (typeof this.run === "undefined") {
                    throw new Error("no run function found in module");
                }
                this.resize = this.instantiatedModule.exports.resize as ((len: number) => number) | undefined;
                if (typeof this.resize === "undefined") {
                    throw new Error("no resize function found in module");
                }
                this.initialize = this.instantiatedModule.exports.initialize as (() => bigint) | undefined;
                if (typeof this.initialize === "undefined") {
                    throw new Error("no initialize function found in module");
                }
                this.memory = this.instantiatedModule.exports.memory as WebAssembly.Memory | undefined;
                if (typeof this.memory === "undefined") {
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
            } else {
                throw new Error("no compiled module found in template");
            }
        });
    }

    public async Ready() {
        return await this.ready;
    }

    public Run() {
        if (typeof this.signature === "undefined") {
            throw new Error("no signature found in module");
        }

        if (typeof this.resize === "undefined") {
            throw new Error("no resize function found in module");
        }

        if (typeof this.run === "undefined") {
            throw new Error("no run function found in module");
        }

        if (typeof this.memory === "undefined") {
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
        if (typeof err !== "undefined") {
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