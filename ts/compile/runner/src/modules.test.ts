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

// import {TextEncoder as te2, TextDecoder as td2} from "fastestsmallesttextencoderdecoder";

import { ScaleFunc, V1Alpha, Go } from "@loopholelabs/scalefile/scalefunc";
import { New } from "@loopholelabs/scale";
import * as fs from "fs";

describe("modules", () => {
  let mods = [
    {path: "wasm/module_middleware.wasm", description: "middleware"},
    {path: "wasm/module_middleware_opt.wasm", description: "middleware (speed optimized)"},
    {path: "wasm/module_middleware_gz.wasm", description: "middleware (size optimized)"}
  ];

  for (let k of mods) {
    let fileSizeInMegabytes = "";

    try {
      var stats = fs.statSync(k.path);
      var fileSizeInBytes = stats.size;
      // Convert the file size to megabytes (optional)
      fileSizeInMegabytes = (fileSizeInBytes / (1024*1024)).toFixed(2);
    } catch(e) {}

    it("Can run module " + k.description + " " + fileSizeInMegabytes + "Mb" , async () => {
      const modNext = fs.readFileSync(k.path);
      
      const fn = new ScaleFunc(V1Alpha, "Test.Middleware", "Test.Tag", "ExampleName@ExampleVersion", Go, modNext);
      const r = await New([fn]);
      const i = await r.Instance(null);
      i.Run();
      let header = i.Context().Response.Headers.get("FROM_TYPESCRIPT");
      expect(header).not.toBeUndefined();
      if (header != undefined) {
        let bits = header.Value;
        expect(bits.length).toBe(1);
        expect(bits[0]).toBe("TRUE");
      }
    });
  }

  it("Can run module error" , async () => {
    const modError = fs.readFileSync("wasm/module_error.wasm");    
    const fn = new ScaleFunc(V1Alpha, "Test.Error", "Test.Tag", "ExampleName@ExampleVersion", Go, modError);
    const r = await New([fn]);
    const i = await r.Instance(null);

    expect(() => {
      i.Run();
    }).toThrow("Something went wrong");

  });

  it("Can run module endpoint" , async () => {
    const modEndpoint = fs.readFileSync("wasm/module_endpoint.wasm");    
    const fn = new ScaleFunc(V1Alpha, "Test.Error", "Test.Tag", "ExampleName@ExampleVersion", Go, modEndpoint);
    const r = await New([fn]);
    const i = await r.Instance(null);

    i.Run();

    let decoder = new TextDecoder();
    let body = decoder.decode(i.Context().Response.Body);

    expect(body).toBe("Hello from typescript!");
  });

});