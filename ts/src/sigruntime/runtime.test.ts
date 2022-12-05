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

    const scalefn = new ScaleFunc();
    scalefn.Version = "TestVersion";
    scalefn.Name = "TestName";
    scalefn.Signature = "ExampleName@ExampleVersion";
    scalefn.Language = "go";
    scalefn.Function = modPassthrough;

    const signature = new Context();    // TODO: Should be signature encapsulating context really...
    const r = new Runtime(getNewWasi(), signature, [scalefn]);
    await r.Ready;


    const i = r.Instance(null);

    i.Context().Data = "Test Data";

    i.Run();

    expect(i.Context().Data).toBe("Test Data");

//    expect(12).toBe(34);
  });

});
