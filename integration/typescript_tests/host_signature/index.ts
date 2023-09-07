// Code generated by scale-signature 0.3.20, DO NOT EDIT.
// output: local-example-latest-host

/* eslint no-bitwise: off */

import { Signature as SignatureInterface } from "@loopholelabs/scale-signature-interfaces";
import { Decoder, Encoder, Kind } from "@loopholelabs/polyglot";

export * from "./types";
import { ModelWithAllFieldTypes } from "./types";

const hash = "3a592aa345d412faa2e6285ee048ca2ab5aa64b0caa2f9ca67b2c1e0792101e5"

// New returns a new signature and tells the Scale Runtime how to use it
//
// This function should be passed into the scale runtime config as an argument
export function New(): Signature {
  return new Signature();
}

// Signature is the host representation of the signature
//
// Users should not use this type directly, but instead pass the New() function
// to the Scale Runtime
export class Signature implements SignatureInterface {
  public context: ModelWithAllFieldTypes;

  constructor() {
    this.context = new ModelWithAllFieldTypes();
  }

  // Read reads the context from the given Uint8Array and returns an error if one occurred
  //
  // This method is meant to be used by the Scale Runtime to deserialize the Signature
  Read(b: Uint8Array): Error | undefined {
    const dec = new Decoder(b);
    try {
      Object.assign(this.context, ModelWithAllFieldTypes.decode(dec));
    } catch (err) {
      return err;
    }
    return undefined;
  }

  // Write writes the signature into a Uint8Array and returns it
  //
  // This method is meant to be used by the Scale Runtime to serialize the Signature
  Write(): Uint8Array {
    const enc = new Encoder();
    this.context.encode(enc);
    return enc.bytes;
  }

  // Error writes the signature into a Uint8Array and returns it
  //
  // This method is meant to be used by the Scale Runtime to return an error
  Error(err: Error): Uint8Array {
    const enc = new Encoder();
    enc.error(err);
    return enc.bytes;
  }

  // Hash returns the hash of the signature
  //
  // This method is meant to be used by the Scale Runtime to validate Signature and Function compatibility
  Hash(): string {
    return hash;
  }
}
