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
  ScaleFunc,
  Go,
  VersionErr,
  LanguageErr,
  ChecksumErr,
  V1Alpha,
  ValidString,
} from "./scalefunc";

// eslint-disable-next-line @typescript-eslint/no-var-requires
const Buffer = require("buffer/").Buffer;

window.TextEncoder = TextEncoder;
window.TextDecoder = TextDecoder as typeof window["TextDecoder"];

describe("scaleFunc", () => {
  const enc = new TextEncoder(); // always utf-8

  const someFunction = enc.encode("Hello world some function here");

  it("Can encode and decode", () => {
    expect(() => {
      const sfInvalid = new ScaleFunc(
        "invalid",
        "name",
        "tag",
        "signature",
        Buffer.from(someFunction),
        "signatureHash",
        Go,
        false,
        [],
        Buffer.from(someFunction)
      );
      const b = sfInvalid.Encode();
      ScaleFunc.Decode(b);
    }).toThrow(VersionErr);

    expect(() => {
      const sfInvalid = new ScaleFunc(
        V1Alpha,
        "name",
        "tag",
        "signature",
        Buffer.from(someFunction),
        "signatureHash",
        "invalid",
        true,
        [],
        Buffer.from(someFunction)
      );
      const b = sfInvalid.Encode();
      ScaleFunc.Decode(b);
    }).toThrow(LanguageErr);

    const sf = new ScaleFunc(
      V1Alpha,
      "Test name",
      "Test tag",
      "Test Signature",
        Buffer.from(someFunction),
        "signatureHash",
      Go,
      false,
      [],
      Buffer.from(someFunction)
    );

    const buff = sf.Encode();
    const sf2 = ScaleFunc.Decode(buff);

    expect(sf.Version).toBe(sf2.Version);
    expect(sf.Name).toBe(sf2.Name);
    expect(sf.Tag).toBe(sf2.Tag);
    expect(sf.SignatureName).toBe(sf2.SignatureName);
    expect(sf.Language).toBe(sf2.Language);
    expect(sf.Function).toStrictEqual(sf2.Function);

    if (typeof sf2.Size !== "undefined" && typeof sf2.Hash !== "undefined") {
      buff[buff.length - 1]++; // This increments the last byte of the hash
      // Now try to decode again with a bad checksum...
      expect(() => {
        ScaleFunc.Decode(buff);
      }).toThrow(ChecksumErr);
    } else {
      throw new Error("Size or Checksum were not set!");
    }
  });
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
