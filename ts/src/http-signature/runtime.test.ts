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
import { Context, Request, Response, StringList } from "./generated/generated";
import { WASI } from "wasi";

import { ScaleFunc } from "../signature/scaleFunc";
import { HttpContext, HttpContextFactory } from "./HttpContext";

import { Runtime as SigRuntime, WasiContext } from "../sigruntime/runtime";

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

describe("runtime", () => {
  it("Can run a simple e2e one module", async () => {    
    // Create a context to send in...
    const enc = new TextEncoder();
    const body = enc.encode("Hello world this is a request body");
    const headers = new Map<string, StringList>();
    headers.set("content", new StringList(["hello"]));
    const req1 = new Request(
      "GET",
      BigInt(100),
      "https",
      "1.2.3.4",
      body,
      headers
    );
    const respBody = enc.encode("Response body");
    const respHeaders = new Map<string, StringList>();
    const resp1 = new Response(200, respBody, respHeaders);
    const context = new Context(req1, resp1);

    const modHttpEndpoint = fs.readFileSync(
      "./example_modules/http-endpoint.wasm"
    );

    const scalefnEndpoint = new ScaleFunc();
    scalefnEndpoint.Version = "TestVersion";
    scalefnEndpoint.Name = "Test.HttpEndpoint";
    scalefnEndpoint.Signature = "ExampleName@ExampleVersion";
    scalefnEndpoint.Language = "go";
    scalefnEndpoint.Function = modHttpEndpoint;

    const signatureFactory = HttpContextFactory;

    const r = new SigRuntime<HttpContext>(getNewWasi, signatureFactory, [scalefnEndpoint]);
    await r.Ready;

    const i = await r.Instance(null);
    i.Context().ctx = context;

    i.Run();

    const retContext = i.Context().ctx;

    expect(retContext).not.toBeNull();

    if (retContext != null) {
      // check the returns...

      expect(retContext.Response.StatusCode).toBe(200);

      const dec = new TextDecoder();
      const bodyText = dec.decode(retContext.Response.Body);

      // The http-endpoint.wasm module copies the request body to the response body.
      expect(bodyText).toBe("Hello world this is a request body");
    }
  });

  it("Can run a simple e2e using runtime", async () => {

    // Create a context to send in...
    const enc = new TextEncoder();
    const body = enc.encode("Hello world this is a request body");
    const headers = new Map<string, StringList>();
    headers.set("content", new StringList(["hello"]));
    const req1 = new Request(
      "GET",
      BigInt(100),
      "https",
      "1.2.3.4",
      body,
      headers
    );
    const respBody = enc.encode("Response body");
    const respHeaders = new Map<string, StringList>();
    const resp1 = new Response(200, respBody, respHeaders);
    const context = new Context(req1, resp1);

    // Now we can use context with a couple of wasm modules...

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

    const i = await r.Instance(null);
    i.Context().ctx = context;

    i.Run();

    const retContext = i.Context().ctx;

    expect(retContext).not.toBeNull();

    if (retContext != null) {
      // check the returns...

      expect(retContext.Response.StatusCode).toBe(200);

      const dec = new TextDecoder();
      const bodyText = dec.decode(retContext.Response.Body);

      // The http-endpoint.wasm module copies the request body to the response body.
      expect(bodyText).toBe("Hello world this is a request body");

      // The http-middleware.wasm adds a header
      const middle = retContext.Response.Headers.get("MIDDLEWARE");
      expect(middle).toBeDefined();
      const vals = middle?.Value;
      if (vals !== undefined) {
        expect(vals.length).toBe(1);
        expect(vals[0]).toBe("TRUE");
      }
    }
  });
});
