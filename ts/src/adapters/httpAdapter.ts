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
import http from "http";
import { Context } from "../runtime/context";
import { Runtime } from "../runtime/runtime";
import {
  Context as PgContext,
  Request as PgRequest,
  Response as PgResponse,
  StringList as PgStringList,
} from "../runtime/generated/generated";
import { Host } from "../runtime/host";

export class HttpAdapter {
  private _runtime: Runtime;

  constructor(runt: Runtime) {
    this._runtime = runt;
  }

  getHandler() {
    return this.handler.bind(this);
  }

  handler(
    req: http.IncomingMessage,
    res: http.ServerResponse
  ) {
    let chunks:Uint8Array[] = [];

    req.on('data', (chunk) => {
        chunks.push(chunk);
    });

    req.on('end', () => {
      let length = 0;
      chunks.forEach(item => {
        length += item.length;
      });
      
      let body = new Uint8Array(length);
      let offset = 0;
      chunks.forEach(item => {
        body.set(item, offset);
        offset += item.length;
      });      

      const c = HttpAdapter.toContext(req, body, res);

      const newc = this._runtime.run(c);

      if (newc != null) {
        HttpAdapter.fromContext(newc, res);        
      } else {
        res.statusCode = 500;
        res.write("Server error");
        res.end();
      }
    });
  }

  static fromContext(ctx: Context, res: http.ServerResponse) {
    // Copy stuff over here...
    const ctxresp = ctx.context().Response;

    res.statusCode = ctxresp.StatusCode;
    for(let k of ctxresp.Headers.keys()) {
      let vals = ctxresp.Headers.get(k);
      if (vals!==undefined) {
        let s = vals.Value;
        for(let v of s.values()) {
            res.setHeader(k, v);
        }
      }
    }

    res.write(ctxresp.Body);
    res.end();
  }

  static toContext(req: http.IncomingMessage, body: Uint8Array, resp: http.ServerResponse) {
    const reqheaders = new Map<string, PgStringList>();

    let method = "unknown";
    if (req.method != undefined) {
      method = req.method;
    }

    // TODO: Find protocol
    // TODO: Find IP
    const protocol = "http";
    const ip = "1.2.3.4";

    for(let k in req.headers) {
      let vals = req.headers[k];
      let sl: string[] = [];
      if (typeof vals === 'string') {
        sl.push(vals);  // Single value
      } else if (vals===undefined) {
        // Just empty values
      } else {
        sl = vals;  // Multiple values
      }
      reqheaders.set(k, new PgStringList(sl));
    }

    const preq = new PgRequest(
      method,
      BigInt(body.length),
      protocol,
      ip,
      body,
      reqheaders
    );
    const presp = new PgResponse(
      resp.statusCode,
      new Uint8Array(),
      new Map<string, PgStringList>()
    )
    const c = new Context(new PgContext(preq, presp));
    return c;
  }
}