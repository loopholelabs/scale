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

import { Signature, New } from "@loopholelabs/scale-signature-interfaces";
import { ScaleFunc } from "./scalefunc/scalefunc";

const envStringRegex = /[^A-Za-z0-9_]/;

class ConfigFunction {
    function: ScaleFunc;
    env: { [key: string]: string } | undefined;

    constructor(fn: ScaleFunc, env?: { [key: string]: string }) {
        this.function = fn;
        this.env = env;
    }
}

export class Config<T extends Signature> {
    newSignature: New<T>;
    functions: ConfigFunction[] = [];

    constructor(newSignature: New<T>) {
        this.newSignature = newSignature;
    }

    validate(): Error | null {
        if (!this) {
            return new Error("no config provided");
        }
        if (this.functions.length === 0) {
            return new Error("no functions provided");
        }

        for (const f of this.functions) {
            if (!f.function) {
                return new Error("invalid function");
            }
            for (const k in f.env) {
                if (!validEnv(k)) {
                    return new Error("invalid environment variable");
                }
            }
        }

        return null;
    }

    withFunction(func: ScaleFunc, env?: { [key: string]: string }): Config<T> {
        const f = new ConfigFunction(func, env);
        f.function = func;

        if (env) {
            f.env = env;
        }

        this.functions.push(f);
        return this;
    }

    withFunctions(functions: ScaleFunc[], env?: { [key: string]: string }): Config<T> {
        for (const func of functions) {
            this.withFunction(func, env);
        }
        return this;
    }
}

// Helper function
function validEnv(str: string): boolean {
    return !envStringRegex.test(str);
}