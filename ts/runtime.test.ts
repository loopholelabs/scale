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

import {ScaleFunc, V1Alpha, Go, Rust} from "@loopholelabs/scalefile/scalefunc";

import * as signature from "./tests/signature" ;
import * as httpSignature from "@loopholelabs/scale-signature-http";

import { New, NewFromSignature } from "./index";

describe("TestRuntimeTs", () => {
    const modPassthroughGo = fs.readFileSync("./ts/tests/modules/passthrough-TestRuntimeGo.wasm");
    const modPassthroughRs = fs.readFileSync("./ts/tests/modules/passthrough-TestRuntimeRs.wasm");
    const modModifyGo = fs.readFileSync("./ts/tests/modules/modify-TestRuntimeGo.wasm");
    const modModifyRs = fs.readFileSync("./ts/tests/modules/modify-TestRuntimeRs.wasm");
    const modNextGo = fs.readFileSync("./ts/tests/modules/next-TestRuntimeGo.wasm");
    const modNextRs = fs.readFileSync("./ts/tests/modules/next-TestRuntimeRs.wasm");
    const modModifyNextGo = fs.readFileSync("./ts/tests/modules/modifynext-TestRuntimeGo.wasm");
    const modModifyNextRs = fs.readFileSync("./ts/tests/modules/modifynext-TestRuntimeRs.wasm");
    const modFileGo = fs.readFileSync("./ts/tests/modules/file-TestRuntimeGo.wasm");
    const modFileRs = fs.readFileSync("./ts/tests/modules/file-TestRuntimeRs.wasm");
    const modNetworkGo = fs.readFileSync("./ts/tests/modules/network-TestRuntimeGo.wasm");
    const modNetworkRs = fs.readFileSync("./ts/tests/modules/network-TestRuntimeRs.wasm");
    const modPanicGo = fs.readFileSync("./ts/tests/modules/panic-TestRuntimeGo.wasm");
    const modPanicRs = fs.readFileSync("./ts/tests/modules/panic-TestRuntimeRs.wasm");
    const modBadSignatureGo = fs.readFileSync("./ts/tests/modules/bad-signature-TestRuntimeGo.wasm");
    const modBadSignatureRs = fs.readFileSync("./ts/tests/modules/bad-signature-TestRuntimeRs.wasm");

    const modHttpPassthroughGo = fs.readFileSync("./ts/tests/modules/http-passthrough-TestRuntimeHTTPSignatureGo.wasm");
    const modHttpPassthroughRs = fs.readFileSync("./ts/tests/modules/http-passthrough-TestRuntimeHTTPSignatureRs.wasm");
    const modHttpHandlerGo = fs.readFileSync("./ts/tests/modules/http-handler-TestRuntimeHTTPSignatureGo.wasm");
    const modHttpHandlerRs = fs.readFileSync("./ts/tests/modules/http-handler-TestRuntimeHTTPSignatureRs.wasm");
    const modHttpNextGo = fs.readFileSync("./ts/tests/modules/http-next-TestRuntimeHTTPSignatureGo.wasm");
    const modHttpNextRs = fs.readFileSync("./ts/tests/modules/http-next-TestRuntimeHTTPSignatureRs.wasm");

    const modTracingGo = fs.readFileSync("./ts/tests/modules/tracing-TestRuntimeGoTracing.wasm");

    it("Tracing", async () => {
      const sfn1 = new ScaleFunc(V1Alpha, `Test.Tracing1`, "Test.TestTag1", "ExampleName@ExampleVersion", Go, modTracingGo);
      const sfn2 = new ScaleFunc(V1Alpha, `Test.Tracing2`, "Test.TestTag2", "ExampleName@ExampleVersion", Go, modTracingGo);
      const r = await NewFromSignature(signature.New, [sfn1, sfn2]);

      let traceData: string[] = []

      r.TraceDataCallback = (s: string) => {
        traceData.push(s);
      }

      const i = await r.Instance(null);
      expect(() => {
        i.Run();
      }).not.toThrowError();

      // Check we got some tracing data...
      expect(traceData.length).toBe(2);

      let d1 = JSON.parse(traceData[0])
      expect(d1.serviceName).toBe("Test.Tracing1:Test.TestTag1");
      let iid1 = [...r.InvocationId];
      let hexID1 = Array.from(iid1, function(byte: number) {
        return ('0' + (byte & 0xFF).toString(16)).slice(-2);
      }).join('');
      expect(d1.invocationID).toBe(hexID1);

      let d2 = JSON.parse(traceData[1])
      expect(d2.serviceName).toBe("Test.Tracing2:Test.TestTag2");
      let iid2 = [...r.InvocationId];
      let hexID2 = Array.from(iid2, function(byte: number) {
        return ('0' + (byte & 0xFF).toString(16)).slice(-2);
      }).join('');
      expect(d2.invocationID).toBe(hexID2);

      // Run it again and check the invocationID is now changed

      const i2 = await r.Instance(null);
      expect(() => {
        i2.Run();
      }).not.toThrowError();

      // Check we got some tracing data...
      expect(traceData.length).toBe(4);

      let d3 = JSON.parse(traceData[2])
      expect(d3.invocationID).not.toBe(d1.invocationID);

    });

  const passthrough = [
    { name: "Passthrough", module: modPassthroughGo, language: Go },
    { name: "Passthrough", module: modPassthroughRs, language: Rust },
  ]

  passthrough.forEach((fn) => {
    it(`${fn.name} ${fn.language}`, async () => {
      const sfn = new ScaleFunc(V1Alpha, `Test.${fn.name}`, "Test.TestTag", "ExampleName@ExampleVersion", fn.language, fn.module);
      const r = await NewFromSignature(signature.New, [sfn]);
      const i = await r.Instance(null);
      i.Context().Data = "Test Data";
      expect(() => {
        i.Run();
      }).not.toThrowError();
      expect(i.Context().Data).toBe("Test Data");
    });
  })

  const modify = [
    { name: "Modify", module: modModifyGo, language: Go },
    { name: "Modify", module: modModifyRs, language: Rust },
  ]

  modify.forEach((fn) => {
    it(`${fn.name} ${fn.language}`, async () => {
      const sfn = new ScaleFunc(V1Alpha, `Test.${fn.name}`, "Test.TestTag", "ExampleName@ExampleVersion", fn.language, fn.module);
      const r = await NewFromSignature(signature.New, [sfn]);
      const i = await r.Instance(null);
      i.Context().Data = "Test Data";
      expect(() => {
        i.Run();
      }).not.toThrowError();
      expect(i.Context().Data).toBe("modified");
    });
  })

  const next = [
    { name: "Next", module: modNextGo, language: Go },
    { name: "Next", module: modNextRs, language: Rust },
  ]

  next.forEach((fn) => {
    it(`${fn.name} ${fn.language}`, async () => {
      const sfn = new ScaleFunc(V1Alpha, `Test.${fn.name}`, "Test.TestTag", "ExampleName@ExampleVersion", fn.language, fn.module);
      const r = await NewFromSignature(signature.New, [sfn]);
      const i = await r.Instance(function (ctx: signature.Context): signature.Context {
        ctx.Data = "Hello, world!";
        return ctx;
      });
      i.Context().Data = "Test Data";
      expect(() => {
        i.Run();
      }).not.toThrowError();
      expect(i.Context().Data).toBe("Hello, world!");
    });
  })

  const modifynext = [
    { name: "ModifyNext", module: modModifyNextGo, language: Go },
    { name: "ModifyNext", module: modModifyNextRs, language: Rust },
  ]

  modifynext.forEach((fn) => {
      it(`${fn.name} ${fn.language}`, async () => {
          const sfn = new ScaleFunc(V1Alpha, `Test.${fn.name}`, "Test.TestTag", "ExampleName@ExampleVersion", fn.language, fn.module);
          const r = await NewFromSignature(signature.New, [sfn]);
          const i = await r.Instance(function (ctx: signature.Context): signature.Context {
            ctx.Data += "-next";
            return ctx;
          });
          i.Context().Data = "Test Data";
          expect(() => {
          i.Run();
          }).not.toThrowError();
          expect(i.Context().Data).toBe("modified-next");
      });
  })

  next.forEach((fn) => {
    it(`${fn.name} error ${fn.language}`, async () => {
        const sfn = new ScaleFunc(V1Alpha, `Test.${fn.name}`, "Test.TestTag", "ExampleName@ExampleVersion", fn.language, fn.module);
        const r = await NewFromSignature(signature.New, [sfn]);
        const i = await r.Instance(function (ctx: signature.Context): signature.Context {
          throw new Error("Hello error");
        });
        i.Context().Data = "Test Data";
        expect(() => {
          i.Run();
        }).toThrow("Hello error");
    });
  })

  const file = [
    { name: "File", module: modFileGo, language: Go },
    { name: "File", module: modFileRs, language: Rust },
  ]

    file.forEach((fn) => {
        it(`${fn.name} ${fn.language}`, async () => {
            const sfn = new ScaleFunc(V1Alpha, `Test.${fn.name}`, "Test.TestTag", "ExampleName@ExampleVersion", fn.language, fn.module);
            const r = await NewFromSignature(signature.New, [sfn]);
            const i = await r.Instance(null);
            i.Context().Data = "Test Data";
            expect(() => {
              i.Run();
            }).toThrowError();
        });
    })


   const network = [
    { name: "Network", module: modNetworkGo, language: Go },
    { name: "Network", module: modNetworkRs, language: Rust },
    ]

    network.forEach((fn) => {
        it(`${fn.name} ${fn.language}`, async () => {
            const sfn = new ScaleFunc(V1Alpha, `Test.${fn.name}`, "Test.TestTag", "ExampleName@ExampleVersion", fn.language, fn.module);
            const r = await NewFromSignature(signature.New, [sfn]);
            const i = await r.Instance(null);
            expect(() => {
            i.Run();
            }).toThrowError();
        });
    })

    const panic = [
    { name: "Panic", module: modPanicGo, language: Go },
    { name: "Panic", module: modPanicRs, language: Rust },
    ]

    panic.forEach((fn) => {
        it(`${fn.name} ${fn.language}`, async () => {
            const sfn = new ScaleFunc(V1Alpha, `Test.${fn.name}`, "Test.TestTag", "ExampleName@ExampleVersion", fn.language, fn.module);
            const r = await NewFromSignature(signature.New, [sfn]);
            const i = await r.Instance(null);
            expect(() => {
            i.Run();
            }).toThrowError();
        });
    })

    const badSignature = [
    { name: "BadSignature", module: modBadSignatureGo, language: Go },
    { name: "BadSignature", module: modBadSignatureRs, language: Rust },
    ]

    badSignature.forEach((fn) => {
        it(`${fn.name} ${fn.language}`, async () => {
            const sfn = new ScaleFunc(V1Alpha, `Test.${fn.name}`, "Test.TestTag", "ExampleName@ExampleVersion", fn.language, fn.module);
            const r = await NewFromSignature(signature.New, [sfn]);
            const i = await r.Instance(null);
            expect(() => {
            i.Run();
            }).toThrowError();
        });
    })

    const httpPassthrough = [
    { name: "HttpPassthrough", module: modHttpPassthroughGo, language: Go },
    { name: "HttpPassthrough", module: modHttpPassthroughRs, language: Rust },
    ]

    httpPassthrough.forEach((fn) => {
        it(`${fn.name} ${fn.language}`, async () => {
            const sfn = new ScaleFunc(V1Alpha, `Test.${fn.name}`, "Test.TestTag", "ExampleName@ExampleVersion", fn.language, fn.module);
            const r = await New([sfn]);
            const i = await r.Instance(null)
            i.Context().Response.Body = new TextEncoder().encode("Test Data");
            expect(() => {
                i.Run();
            }).not.toThrowError();
            const textDecoder = new TextDecoder();
            expect(textDecoder.decode(i.Context().Response.Body)).toBe("Test Data");
        });
    })

    const httpHandler = [
    { name: "HttpHandler", module: modHttpHandlerGo, language: Go },
    { name: "HttpHandler", module: modHttpHandlerRs, language: Rust },
    ]

    httpHandler.forEach((fn) => {
        it(`${fn.name} ${fn.language}`, async () => {
            const sfn = new ScaleFunc(V1Alpha, `Test.${fn.name}`, "Test.TestTag", "ExampleName@ExampleVersion", fn.language, fn.module);
            const r = await New([sfn]);
            const i = await r.Instance(null)
            i.Context().Response.Body = new TextEncoder().encode("Test Data");
            expect(() => {
                i.Run();
            }).not.toThrowError();
            const textDecoder = new TextDecoder();
            expect(textDecoder.decode(i.Context().Response.Body)).toBe("Test Data-modified");
        });
    })

    const httpNext = [
    { name: "HttpNext", module: modHttpNextGo, language: Go },
    { name: "HttpNext", module: modHttpNextRs, language: Rust },
    ]

    httpNext.forEach((fn) => {
        it(`${fn.name} ${fn.language}`, async () => {
            const sfn = new ScaleFunc(V1Alpha, `Test.${fn.name}`, "Test.TestTag", "ExampleName@ExampleVersion", fn.language, fn.module);
            const r = await New([sfn]);
            const i = await r.Instance(function (ctx: httpSignature.Context): httpSignature.Context {
                ctx.Response.Body = new TextEncoder().encode(new TextDecoder().decode(ctx.Response.Body) + "-next");
                return ctx;
            });
            i.Context().Response.Body = new TextEncoder().encode("Test Data");
            expect(() => {
                i.Run();
            }).not.toThrowError();
            const textDecoder = new TextDecoder();
            expect(textDecoder.decode(i.Context().Response.Body)).toBe("Test Data-modified-next");
        });
    })
});
