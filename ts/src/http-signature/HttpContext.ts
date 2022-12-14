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

/* eslint no-bitwise: off */

import { Context, StringList } from "./generated/generated";

import { Signature, RuntimeContext } from "@loopholelabs/scale-signature";

import { Kind, encodeError, decodeError } from "@loopholelabs/polyglot-ts";

export function HttpContextFactory():HttpContext {
  return new HttpContext();
}

export class HttpContext implements Signature, RuntimeContext {
  public ctx: undefined | Context;

  constructor() {
  }

  RuntimeContext(): RuntimeContext {
    return this;
  }

  Read(d: Uint8Array) {
    if (d.length > 0 && d[0]===Kind.Error) {
      const e = decodeError(d).value;
      throw(e);
    }
    this.ctx = Context.decode(d).value;
  }

  Write(): Uint8Array {
    if (this.ctx === undefined) throw (new Error("ctx undefined"));
    return this.ctx.encode(new Uint8Array());
  }

  Error(e: Error): Uint8Array {
    return encodeError(new Uint8Array(), e);
  }

  private static stringHeaders(h: Map<string, StringList>): string {
    let r = "";
    for (let k of h.keys()) {
      let values = h.get(k);
      if (values != undefined) {
        for (let i of values.Value.values()) {
          r = r + " " + k + "=" + i;
        }
      } else {
        r = r + " " + k;
      }
    }
    return r;
  }

  public show() {
    if (this.ctx === undefined) {
      console.log("== Context undefined ==");
    } else {
      const req = this.ctx.Request;
      const reqBody = new TextDecoder().decode(req.Body);
      console.log(`== Context ==
    Request method=${req.Method}, proto=${req.Protocol}, ip=${req.IP}, len=${req.ContentLength}
    Headers: ${HttpContext.stringHeaders(req.Headers)}
    Body: ${reqBody}`);

      const resp = this.ctx.Response;
      const respBody = new TextDecoder().decode(resp.Body);
      console.log(`Response code=${resp.StatusCode}
    Headers: ${HttpContext.stringHeaders(resp.Headers)}
    Body: ${respBody}`);
    }
  }
}