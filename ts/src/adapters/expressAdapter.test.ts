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
import express from "express";
import bodyParser from "body-parser";
import * as fs from "fs";
import request from "supertest";
import { WASI } from "wasi";

import { ScaleFunc } from "../signature/scaleFunc";
import { HttpContext, HttpContextFactory } from "../http-signature/HttpContext";
import { Runtime as SigRuntime, WasiContext } from "../sigruntime/runtime";

import { ExpressAdapter } from "./expressAdapter";

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


describe("expressAdapter", () => {
  const app = express();

  it("Can run a simple e2e", async () => {
    const modHttpEndpoint = fs.readFileSync(
      "./example_modules/http-endpoint.wasm"
    );
    const modHttpMiddleware = fs.readFileSync(
      "./example_modules/http-middleware.wasm"
    );

    const scalefnEndpoint = new ScaleFunc();
    scalefnEndpoint.Version = "TestVersion";
    scalefnEndpoint.Name = "Test.HttpEndpoint";
    scalefnEndpoint.Signature = "ExampleName@ExampleVersion";
    scalefnEndpoint.Language = "go";
    scalefnEndpoint.Function = modHttpEndpoint;

    const scalefnMiddle = new ScaleFunc();
    scalefnMiddle.Version = "TestVersion";
    scalefnMiddle.Name = "Test.HttpEndpoint";
    scalefnMiddle.Signature = "ExampleName@ExampleVersion";
    scalefnMiddle.Language = "go";
    scalefnMiddle.Function = modHttpMiddleware;

    const signatureFactory = HttpContextFactory;

    const r = new SigRuntime<HttpContext>(getNewWasi, signatureFactory, [scalefnMiddle, scalefnEndpoint]);
    await r.Ready;
    

    const adapter = new ExpressAdapter(r);

    app.use(
      bodyParser.raw({
        type: () => true,
      })
    );

    app.use(adapter.getHandler());

    const res = await request(app).post("/blah").send("HELLO WORLD");

    // Make sure everything worked as expected.
    expect(res.statusCode).toEqual(200);
    expect(res.text).toBe("HELLO WORLD");
    expect(res.headers.middleware).toBe("TRUE");
  });
});
