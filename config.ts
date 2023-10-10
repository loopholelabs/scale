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
import { Extension } from "@loopholelabs/scale-extension-interfaces";

const envStringRegex = /[^A-Za-z0-9_]/;

export interface Writer {
    (message: string): void;
}

class ConfigFunction {
    function: ScaleFunc;
    env: { [key: string]: string } | undefined;
    constructor(fn: ScaleFunc, env?: { [key: string]: string }) {
        this.function = fn;
        this.env = env;
    }
}

export class Config<T extends Signature> {
    public newSignature: New<T>;
    public functions: ConfigFunction[] = [];
    stdout: Writer | undefined;
    stderr: Writer | undefined;
    rawOutput: boolean = false;
    public extensions: Extension[] = [];

    constructor(newSignature: New<T>) {
        this.newSignature = newSignature;
    }

    public validate() {
        if (!this) {
           throw new Error("no config provided");
        }
        if (this.functions.length === 0) {
            throw new Error("no functions provided");
        }

        for (const f of this.functions) {
            if (!f.function) {
                throw new Error("invalid function");
            }
            for (const k in f.env) {
                if (!validEnv(k)) {
                    throw new Error("invalid environment variable");
                }
            }
        }
    }

    public WithSignature(newSignature: New<T>): Config<T> {
        this.newSignature = newSignature;
        return this;
    }

    public WithFunction(func: ScaleFunc, env?: { [key: string]: string }): Config<T> {
        const f = new ConfigFunction(func, env);
        f.function = func;

        if (env) {
            f.env = env;
        }

        this.functions.push(f);
        return this;
    }

    public WithExtension(e: Extension): Config<T> {
      this.extensions.push(e);
      return this;
    }

    public WithFunctions(functions: ScaleFunc[], env?: { [key: string]: string }): Config<T> {
        for (const func of functions) {
            this.WithFunction(func, env);
        }
        return this;
    }

    public WithStdout(writer: Writer): Config<T> {
        this.stdout = writer
        return this;
    }

    public WithStderr(writer: Writer): Config<T> {
        this.stderr = writer
        return this;
    }

    public WithRawOutput(rawOutput: boolean): Config<T> {
        this.rawOutput = rawOutput
        return this;
    }
}

// Helper function
function validEnv(str: string): boolean {
    return !envStringRegex.test(str);
}