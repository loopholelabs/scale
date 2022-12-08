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

import { Context as pgContext } from "./generated/generated";

// Wrapper for context
export class Context {
  private _context: pgContext;

  constructor(ctx: pgContext) {
    this._context = ctx;
  }

  context(): pgContext {
    return this._context;
  }

  // Write a context into WebAssembly memory and return a ptr/length
  writeTo(
    mem: WebAssembly.Memory,
    mallocfn: Function
  ): { ptr: number; len: number } {
    const inContextBuff = new Uint8Array();
    const encoded = this._context.encode(inContextBuff);
    const encPtr = mallocfn(encoded.length);

    const memData = new Uint8Array(mem.buffer);
    memData.set(encoded, encPtr); // Writes the context into memory
    return { ptr: encPtr, len: encoded.length };
  }

  // Read a context from WebAssembly memory
  public static readFrom(
    mem: WebAssembly.Memory,
    ptr: number,
    len: number
  ): Context {
    const memData = new Uint8Array(mem.buffer);
    const inContextBuff = memData.slice(ptr, ptr + len);
    const c = pgContext.decode(inContextBuff);
    return new Context(c.value);
  }
}
