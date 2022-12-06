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

import { TextEncoder, TextDecoder } from "util";
import * as fs from "fs";

import { Signature } from "../signature/signature";
import { ScaleFunc } from "../signature/scaleFunc";

import { Context } from "./signature/runtime";
import { Runtime, WasiContext } from "./runtime";
import { WASI } from "wasi";

window.TextEncoder = TextEncoder;
window.TextDecoder = TextDecoder as typeof window["TextDecoder"];

function getNewWasi(): WasiContext {
  const wasi = new WASI({
    args: [],
    env: {},
  });
  const w: WasiContext = {
    getImportObject: () => wasi.wasiImport,
    start: (instance: WebAssembly.Instance) => {
      wasi.start(instance);
    }
  }
  return w;
}

describe("sigruntime", () => {
  it("Can run a simple signature e2e one module", async () => {    

    const modPassthrough = fs.readFileSync(
      "./src/sigruntime/modules/passthrough-TestRuntime.wasm"
    );

    const modNext = fs.readFileSync(
      "./src/sigruntime/modules/next-TestRuntime.wasm"
    );

    const scalefnPassthrough = new ScaleFunc();
    scalefnPassthrough.Version = "TestVersion";
    scalefnPassthrough.Name = "Test.Passthrough";
    scalefnPassthrough.Signature = "ExampleName@ExampleVersion";
    scalefnPassthrough.Language = "go";
    scalefnPassthrough.Function = modPassthrough;

    const scalefnNext = new ScaleFunc();
    scalefnNext.Version = "TestVersion";
    scalefnNext.Name = "Test.Next";
    scalefnNext.Signature = "ExampleName@ExampleVersion";
    scalefnNext.Language = "go";
    scalefnNext.Function = modNext;

    const signature = new Context();    // TODO: Should be signature encapsulating context really...
    const r = new Runtime(getNewWasi, signature, [scalefnNext, scalefnPassthrough]);
    await r.Ready;

    const i = r.Instance(null);

    i.Context().Data = "Test Data";

    i.Run();

    expect(i.Context().Data).toBe("Test Data");

//    expect(12).toBe(34);
  });

});
