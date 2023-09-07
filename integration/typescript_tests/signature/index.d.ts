// Code generated by scale-signature 0.3.20, DO NOT EDIT.
// output: local-example-latest-guest

export * from "./types";
import { ModelWithAllFieldTypes } from "./types";

// Write serializes the signature into the global WRITE_BUFFER and returns the pointer to the buffer and its size
//
// Users should not use this method.
export declare function Write(ctx: ModelWithAllFieldTypes): number[];

// Read deserializes signature from the global READ_BUFFER
//
// Users should not use this method.
export declare function Read(): ModelWithAllFieldTypes | undefined;

// Error serializes an error into the global writeBuffer and returns a pointer to the buffer and its size
//
// Users should not use this method.
export function Error(err: Error): number[];

// Resize resizes the global READ_BUFFER to the given size and returns the pointer to the buffer
//
// Users should not use this method.
export function Resize(size: number): number;

// Hash returns the hash of the Scale Signature
//
// Users should not use this method.
export function Hash(): number[];

// Next calls the next function in the Scale Function Chain
export function Next(ctx: ModelWithAllFieldTypes): ModelWithAllFieldTypes | undefined;
