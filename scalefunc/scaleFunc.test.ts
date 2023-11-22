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

import {
  V1AlphaSchema,
  Go,
  VersionErr,
  ValidString, V1BetaSchema, V1BetaSignature,
} from "./scalefunc";

// eslint-disable-next-line @typescript-eslint/no-var-requires
const Buffer = require("buffer/").Buffer;

window.TextEncoder = TextEncoder;
window.TextDecoder = TextDecoder as typeof window["TextDecoder"];

describe("encode_decode", () => {
  const enc = new TextEncoder(); // always utf-8
  const encodedSignature = enc.encode("encoded signature");
  const encodedFunction = enc.encode("encoded function");
  const encodedManifest = enc.encode("encoded manifest");

  it("v1Alpha", () => {
    expect(() => {
      const v1BetaSignature = new V1BetaSignature("Test Signature", "", "", Buffer.from(encodedSignature), "Test Signature Hash");
      const v1Beta = new V1BetaSchema("", "", v1BetaSignature, [], Go, Buffer.from(encodedManifest), false, Buffer.from(encodedFunction));

      const encoded = v1Beta.Encode();
      V1AlphaSchema.Decode(encoded);
    }).toThrow(VersionErr);
  });

  it("v1Beta", () => {
    const v1Alpha = new V1AlphaSchema("", "", "", Buffer.from(encodedSignature), "Test Signature Hash", Go, false, Buffer.from(encodedFunction), Buffer.from(encodedManifest));
    const encoded = v1Alpha.Encode();

    const v1Beta = V1BetaSchema.Decode(encoded);

    expect(v1Beta.Name).toBe(v1Alpha.Name);
    expect(v1Beta.Tag).toBe(v1Alpha.Tag);
    expect(v1Beta.Signature.Name).toBe(v1Alpha.SignatureName);
    expect(v1Beta.Language).toBe(v1Alpha.Language);
    expect(v1Beta.Function).toStrictEqual(v1Alpha.Function);
    expect(v1Beta.Stateless).toBe(v1Alpha.Stateless);
  })
});


describe("ValidName", () => {
  it("Valid Name", () => {
    expect(ValidString("test")).toBe(true);
    expect(ValidString("test1")).toBe(true);
    expect(ValidString("test.1")).toBe(true);
    expect(ValidString("te---.-1")).toBe(true);
    expect(ValidString("test-1")).toBe(true);
  });

  it("Invalid Name", () => {
    expect(ValidString("test_1")).toBe(false);
    expect(ValidString("test 1")).toBe(false);
    expect(ValidString("test1 ")).toBe(false);
    expect(ValidString(" test1")).toBe(false);
    expect(ValidString("test1_")).toBe(false);
    expect(ValidString("test1?")).toBe(false);
    expect(ValidString("test1!")).toBe(false);
    expect(ValidString("test1@")).toBe(false);
    expect(ValidString("test1#")).toBe(false);
    expect(ValidString("test1$")).toBe(false);
    expect(ValidString("test1%")).toBe(false);
    expect(ValidString("test1^")).toBe(false);
    expect(ValidString("test1&")).toBe(false);
    expect(ValidString("test1*")).toBe(false);
    expect(ValidString("test1(")).toBe(false);
    expect(ValidString("test1-1!")).toBe(false);
  });
});
