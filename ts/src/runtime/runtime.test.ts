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

import { Context } from "./signature/runtime";
import { GetRuntime, Runtime } from "./runtime";

window.TextEncoder = TextEncoder;
window.TextDecoder = TextDecoder as typeof window["TextDecoder"];

describe("sigruntime", () => {
  const modPassthrough = fs.readFileSync("./src/runtime/modules/passthrough-TestRuntime.wasm");
  const modNext = fs.readFileSync("./src/runtime/modules/next-TestRuntime.wasm");
  const modFile = fs.readFileSync("./src/runtime/modules/file-TestRuntime.wasm");
  const modNetwork = fs.readFileSync("./src/runtime/modules/network-TestRuntime.wasm");
  const modPanic = fs.readFileSync("./src/runtime/modules/panic-TestRuntime.wasm");
  const modBadSignature = fs.readFileSync("./src/runtime/modules/bad-signature-TestRuntime.wasm");

  const signatureFactory = () => {
    return new Context();        // TODO: Should be signature encapsulating context really...
  }

  it("Can run passthrough", async () => {

    const scalefnPassthrough = new ScaleFunc(V1Alpha, "Test.Passthrough", "ExampleName@ExampleVersion", Go, [], modPassthrough);

    const r = await GetRuntime(signatureFactory, [scalefnPassthrough]);

    const i = await r.Instance(null);
    i.Context().Data = "Test Data";

    expect(() => {
      i.Run();
    }).not.toThrowError();

    expect(i.Context().Data).toBe("Test Data");
  });

  it("Can run next", async () => {

    const scalefnNext = new ScaleFunc(V1Alpha, "Test.Next", "ExampleName@ExampleVersion", Go, [], modNext);

    const r = await GetRuntime(signatureFactory, [scalefnNext]);

    const nextfn = (ctx: Context): Context => {
      console.log("HERE");
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

    const scalefnNext = new ScaleFunc(V1Alpha, "Test.Next", "ExampleName@ExampleVersion", Go, [], modNext);

    const r = await GetRuntime(signatureFactory, [scalefnNext]);

    const nextfn = (ctx: Context): Context => {
      throw new Error("Hello error");
    }

    const i = await r.Instance(nextfn);

    i.Context().Data = "Test Data";

    expect(() => {
      i.Run();
    }).toThrow("Hello error");

  });

  it("Can run file error", async () => {

    const scalefnFile = new ScaleFunc(V1Alpha, "Test.File", "ExampleName@ExampleVersion", Go, [], modFile);

    const r = await GetRuntime(signatureFactory, [scalefnFile]);

    const i = await r.Instance(null);

    expect(() => {
      i.Run();
    }).toThrowError();

  });

  it("Can run network error", async () => {

    const scalefnNetwork = new ScaleFunc(V1Alpha, "Test.Network", "ExampleName@ExampleVersion", Go, [], modNetwork);

    const r = await GetRuntime(signatureFactory, [scalefnNetwork]);

    const i = await r.Instance(null);

    expect(() => {
      i.Run();
    }).toThrowError();

  });

  it("Can run panic error", async () => {

    const scalefnPanic = new ScaleFunc(V1Alpha, "Test.Panic", "ExampleName@ExampleVersion", Go, [], modPanic);

    const r = await GetRuntime(signatureFactory, [scalefnPanic]);

    const i = await r.Instance(null);

    expect(() => {
      i.Run();
    }).toThrowError();

  });

  it("Can run bad-signature error", async () => {

    const scalefnBadSignature = new ScaleFunc(V1Alpha, "Test.BadSig", "ExampleName@ExampleVersion", Go, [], modBadSignature);

    const r = await GetRuntime(signatureFactory, [scalefnBadSignature]);

    const i = await r.Instance(null);

    expect(() => {
      i.Run();
    }).toThrowError();

  });
});
