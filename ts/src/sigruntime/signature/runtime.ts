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
import { Kind, encodeError, decodeError, encodeString, decodeString, decodeNull } from "@loopholelabs/polyglot-ts";

import { Signature, RuntimeContext } from "../../signature/signature";

const NilError = new Error(""); // TODO

export class Context implements RuntimeContext {
  // For now we'll just put the stuff here...
  public Data: string

  constructor() {
    this.Data = "";
  }

  RuntimeContext(): RuntimeContext {
    return this;  //(this as RuntimeContext);
  }

  Read(d: Uint8Array): Error {
    this.Data = decodeString(d).value;
    return NilError;
  }

  Write(): Uint8Array {
    return encodeString(new Uint8Array(), this.Data);
  }

  Error(e: Error): Uint8Array {
    return encodeError(new Uint8Array(), e);
  }
}