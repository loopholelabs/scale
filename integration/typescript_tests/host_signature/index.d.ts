// Code generated by scale-signature 0.3.20, DO NOT EDIT.
// output: local-example-latest-host

import { Signature as SignatureInterface } from "@loopholelabs/scale-signature-interfaces";

export * from "./types";

// New returns a new signature and tells the Scale Runtime how to use it
//
// This function should be passed into the scale runtime config as an argument
export declare function New();

// Signature is the host representation of the signature
//
// Users should not use this type directly, but instead pass the New() function
// to the Scale Runtime
export declare class Signature implements SignatureInterface {
  public context: ModelWithAllFieldTypes;

  constructor();

  // Read reads the context from the given Uint8Array and returns an error if one occurred
  //
  // This method is meant to be used by the Scale Runtime to deserialize the Signature
  Read(b: Uint8Array): Error | undefined;

  // Write writes the signature into a Uint8Array and returns it
  //
  // This method is meant to be used by the Scale Runtime to serialize the Signature
  Write(): Uint8Array;

  // Error writes the signature into a Uint8Array and returns it
  //
  // This method is meant to be used by the Scale Runtime to return an error
  Error(err: Error): Uint8Array;

  // Hash returns the hash of the signature
  //
  // This method is meant to be used by the Scale Runtime to validate Signature and Function compatibility
  Hash(): string;
}
