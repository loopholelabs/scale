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
  ValidString,
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
    const v1Alpha = new V1AlphaSchema("", "", "", Buffer.from(encodedSignature), "Test Signature Hash", Go, false, Buffer.from(encodedFunction), Buffer.from(encodedManifest));
    const encoded = v1Alpha.Encode();

    const decoded = V1AlphaSchema.Decode(encoded);

    expect(decoded.Language).toBe(v1Alpha.Language);
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
