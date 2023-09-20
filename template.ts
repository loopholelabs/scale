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

import { Signature  } from "@loopholelabs/scale-signature-interfaces";
import { Scale } from "./scale";
import {ScaleFunc} from "./scalefunc/scalefunc";
import {ModulePool} from "./pool";

export async function NewTemplate<T extends Signature>(runtime: Scale<T>, scaleFunc: ScaleFunc, env?: { [key: string]: string }): Promise<Template<T>> {
    const t = new Template<T>(runtime, scaleFunc, env);
    await t.Ready();
    return t;
}

export class Template<T extends Signature> {
    private readonly ready: Promise<void>;

    public runtime: Scale<T>
    public identifier: string;
    public compiled: WebAssembly.Module | undefined;
    public next: undefined | Template<T>;

    public modulePool: ModulePool<T> | undefined;
    public env: { [key: string]: string } | undefined;

    constructor(runtime: Scale<T>, scaleFunc: ScaleFunc, env?: { [key: string]: string }) {
        this.runtime = runtime;
        this.identifier = `${scaleFunc.Name}:${scaleFunc.Tag}`;
        this.env = env;

        if (scaleFunc.Stateless) {
            this.modulePool = new ModulePool<T>(this);
        }

        this.ready = new Promise(async (resolve) => { // eslint-disable-line no-async-promise-executor
            this.compiled = await WebAssembly.compile(scaleFunc.Function);
            resolve();
        });
    }

    async Ready() {
        await this.ready;
    }
}