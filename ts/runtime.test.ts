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

import { ScaleFunc, V1Alpha, Go } from "@loopholelabs/scalefile";

import * as signature from "./runtime/signature" ;
import { New } from "./runtime";

window.TextEncoder = TextEncoder;
window.TextDecoder = TextDecoder as typeof window["TextDecoder"];

describe("sigruntime", () => {
  const modPassthrough = fs.readFileSync("./runtime/modules/passthrough-TestRuntime.wasm");
  const modNext = fs.readFileSync("./runtime/modules/next-TestRuntime.wasm");
  const modFile = fs.readFileSync("./runtime/modules/file-TestRuntime.wasm");
  const modNetwork = fs.readFileSync("./runtime/modules/network-TestRuntime.wasm");
  const modPanic = fs.readFileSync("./runtime/modules/panic-TestRuntime.wasm");
  const modBadSignature = fs.readFileSync("./runtime/modules/bad-signature-TestRuntime.wasm");

  it("Can run passthrough", async () => {

    const scalefnPassthrough = new ScaleFunc(V1Alpha, "Test.Passthrough", "Test.TestTag", "ExampleName@ExampleVersion", Go, modPassthrough);

    const r = await New(signature.New, [scalefnPassthrough]);

    const i = await r.Instance(null);
    i.Context().Data = "Test Data";

    expect(() => {
      i.Run();
    }).not.toThrowError();

    expect(i.Context().Data).toBe("Test Data");
  });

  it("Can run next", async () => {

    const scalefnNext = new ScaleFunc(V1Alpha, "Test.Next", "Test.TestTag", "ExampleName@ExampleVersion", Go, modNext);

    const r = await New(signature.New, [scalefnNext]);

    const nextfn = (ctx: signature.Context): signature.Context => {
      ctx.Data = "Hello, world!";
      return ctx;
    }

    const i = await r.Instance(nextfn);

    i.Context().Data = "Test Data";

    expect(() => {
      i.Run();
    }).not.toThrowError();

    expect(i.Context().Data).toBe("Hello, world!");
  });

  it("Can run next error", async () => {

    const scalefnNext = new ScaleFunc(V1Alpha, "Test.Next", "Test.TestTag", "ExampleName@ExampleVersion", Go, modNext);

    const r = await New(signature.New, [scalefnNext]);

    const nextfn = (ctx: signature.Context): signature.Context => {
      throw new Error("Hello error");
    }

    const i = await r.Instance(nextfn);

    i.Context().Data = "Test Data";

    expect(() => {
      i.Run();
    }).toThrow("Hello error");

  });

  it("Can run file error", async () => {

    const scalefnFile = new ScaleFunc(V1Alpha, "Test.File", "Test.TestTag", "ExampleName@ExampleVersion", Go, modFile);

    const r = await New(signature.New, [scalefnFile]);

    const i = await r.Instance(null);

    expect(() => {
      i.Run();
    }).toThrowError();

  });

  it("Can run network error", async () => {

    const scalefnNetwork = new ScaleFunc(V1Alpha, "Test.Network", "Test.TestTag", "ExampleName@ExampleVersion", Go, modNetwork);

    const r = await New(signature.New, [scalefnNetwork]);

    const i = await r.Instance(null);

    expect(() => {
      i.Run();
    }).toThrowError();

  });

  it("Can run panic error", async () => {

    const scalefnPanic = new ScaleFunc(V1Alpha, "Test.Panic", "Test.TestTag", "ExampleName@ExampleVersion", Go, modPanic);

    const r = await New(signature.New, [scalefnPanic]);

    const i = await r.Instance(null);

    expect(() => {
      i.Run();
    }).toThrowError();

  });

  it("Can run bad-signature error", async () => {

    const scalefnBadSignature = new ScaleFunc(V1Alpha, "Test.BadSig", "Test.TestTag", "ExampleName@ExampleVersion", Go, modBadSignature);

    const r = await New(signature.New, [scalefnBadSignature]);

    const i = await r.Instance(null);

    expect(() => {
      i.Run();
    }).toThrowError();

  });
});
