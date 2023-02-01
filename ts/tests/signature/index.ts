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

import { Kind, encodeError, decodeError, encodeString, decodeString } from "@loopholelabs/polyglot-ts";

import { RuntimeContext } from "@loopholelabs/scale-signature";

export function New(): Context {
  return new Context();
}

export class Context implements RuntimeContext {
  public Data: string

  constructor() {
    this.Data = "";
  }

  RuntimeContext(): RuntimeContext {
    return this;
  }

  Read(d: Uint8Array) {
    if (d.length > 0 && d[0] === Kind.Error) {
      const e = decodeError(d).value;
      throw (e);
    }
    this.Data = decodeString(d).value;
  }

  Write(): Uint8Array {
    return encodeString(new Uint8Array(), this.Data);
  }

  Error(e: Error): Uint8Array {
    return encodeError(new Uint8Array(), e);
  }
}

export default {
    New: New,
    Context: Context,
};