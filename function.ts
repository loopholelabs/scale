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
import {Module, NewModule} from "./module";
import {Instance} from "./instance";
import {Template} from "./template";

export class Func<T extends Signature> {
    private readonly ready: Promise<void>;
    public instance: Instance<T>;
    private template: Template<T>;
    public next: Func<T> | undefined;
    public module: Module<T> | undefined;

    constructor(instance: Instance<T>, template: Template<T>) {
        this.instance = instance;
        this.template = template;

        this.ready = new Promise(async (resolve) => { // eslint-disable-line no-async-promise-executor
            if (template.modulePool === undefined) {
                this.module = await NewModule<T>(template);
                this.module.Register(this);
            }
            resolve();
        });
    }

    public async Ready(): Promise<void> {
        await this.ready;
    }

    public async GetModule(signature: T): Promise<Module<T>> {
        if (this.module !== undefined) {
            this.module.SetSignature(signature);
            return this.module;
        }
        if (this.template.modulePool === undefined) {
            throw new Error(`cannot get module from pool for function ${this.template.identifier}: module pool is undefined`);
        }
        const m = await this.template.modulePool.Get();
        m.Register(this);
        m.SetSignature(signature);
        return m;
    }

    public PutModule(m: Module<T>){
        if (this.template.modulePool !== undefined) {
            m.Cleanup()
            this.template.modulePool.Put(m);
        }
    }
}

export async function NewFunc<T extends Signature>(instance: Instance<T>, template: Template<T>): Promise<Func<T>> {
    const fn = new Func<T>(instance, template);
    await fn.Ready();
    return fn;
}