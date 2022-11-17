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
import http from "http";
import request from "supertest";
import { HttpAdapter } from "./httpAdapter";
import * as fs from "fs";
import { Runtime } from "../runtime/runtime";
import { WASI } from "wasi";

import { Module, WasiContext } from "../runtime/module";

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

describe("httpAdapter", () => {

  it("Can run a simple e2e", async () => {
    const modHttpEndpoint = fs.readFileSync(
      "./example_modules/http-endpoint.wasm"
    );
    const modHttpMiddleware = fs.readFileSync(
      "./example_modules/http-middleware.wasm"
    );
    
    // Create factory to return initted module
    // module.new(moduleCode)

    // const runtime = new Runtime([module.new(modCode), module.new(modCode2)]);

    const moduleHttpEndpoint = new Module(modHttpEndpoint, getNewWasi());
    await moduleHttpEndpoint.init();
    const moduleHttpMiddleware = new Module(modHttpMiddleware, getNewWasi());
    await moduleHttpMiddleware.init();
    const runtime = new Runtime([moduleHttpMiddleware, moduleHttpEndpoint]);

    var adapter = new HttpAdapter(runtime);

    var server = http.createServer(adapter.getHandler());
  

    const res = await request(server).post("/blah").send("HELLO WORLD");

    console.log("Done POST");

    // Make sure everything worked as expected.
    expect(res.statusCode).toEqual(200);
    expect(res.text).toBe("HELLO WORLD");
    expect(res.headers.middleware).toBe("TRUE");

  });
});
