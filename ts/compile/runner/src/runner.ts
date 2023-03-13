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

declare var scale: ScaleFunction;     // This should be defined in the global scope.


import { TextEncoder, TextDecoder } from "./text";

//import { TextEncoder, TextDecoder } from "text-encoding";
(global as any).TextEncoder = TextEncoder;
(global as any).TextDecoder = TextDecoder;

import { Resize, GuestContext, HttpContext, HttpRequest, HttpResponse, HttpStringList } from "@loopholelabs/scale-signature-http";

export type ScaleFunction = (a: GuestContext) => GuestContext;

function mainFunction() {
//  console.log("Console.log Main function called");
//  console.error("Console.error Main function called");
}

// Create a new dummy GuestContext.
var ctx = new GuestContext(new HttpContext(
  new HttpRequest("GET", "http://example.com", BigInt(0), "http", "1.2.3.4", new Uint8Array(), new Map<string, HttpStringList>),
  new HttpResponse(200, new Uint8Array(), new Map<string, HttpStringList>)
));

// Our own run function wrapper
function runFunction(): bigint {
  ctx.FromReadBuffer();   // Read from the read buffer.

  try {
    scale(ctx);

    let [ptr, len] = ctx.ToWriteBuffer();
    return BigInt(ptr) << BigInt(32) | BigInt(len);
  } catch(e) {
    let [ptr, len] = ctx.ErrorWriteBuffer(e as Error);
    return BigInt(ptr) << BigInt(32) | BigInt(len);
  }

}

// Route the resize through to the signature
function resizeFunction(size: number): number {
  return Resize(size);
}

(global as any).Exports = {
  main: mainFunction,
  run: runFunction,
  resize: resizeFunction
}
