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

import { ScaleFunc } from "../scalefunc";
import { New as NewScale } from "../scale";

import fs from "fs";
import {Config} from "../config";
import { New as NewSignature, Signature } from "./typescript_tests/host_signature";
import {TextDecoder, TextEncoder} from "util";

window.TextEncoder = TextEncoder;
window.TextDecoder = TextDecoder as typeof window["TextDecoder"];

test("test-typescript-host-rust-guest", async () => {
    const file = fs.readFileSync(process.cwd() + "/integration/rust.scale")
    const sf = ScaleFunc.Decode(file);
    const config = new Config<Signature>(NewSignature).WithFunction(sf).WithStdout(console.log).WithStderr(console.error);
    const s = await NewScale(config);

    const i = await s.Instance();
    const sig = NewSignature();

    await i.Run(sig);

    expect(sig.context.stringField).toBe("This is a Rust Function");
});

test("test-typescript-host-golang-guest", async () => {
    const file = fs.readFileSync(process.cwd() + "/integration/golang.scale")
    const sf = ScaleFunc.Decode(file);
    const config = new Config<Signature>(NewSignature).WithFunction(sf).WithStdout(console.log).WithStderr(console.error);
    const s = await NewScale(config);

    const i = await s.Instance();
    const sig = NewSignature();

    await i.Run(sig);

    expect(sig.context.stringField).toBe("This is a Golang Function");
});

test("test-typescript-host-typescript-guest", async () => {
    const file = fs.readFileSync(process.cwd() + "/integration/typescript.scale")
    const sf = ScaleFunc.Decode(file);
    const config = new Config<Signature>(NewSignature).WithFunction(sf).WithStdout(console.log).WithStderr(console.error);
    const s = await NewScale(config);

    const i = await s.Instance();
    const sig = NewSignature();

    await i.Run(sig);

    expect(sig.context.stringField).toBe("This is a Typescript Function");
});
