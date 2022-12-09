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
import { HttpContext } from "../http-signature/HttpContext";

import { Runtime } from "../sigruntime/runtime";
import {
  Context, Request, Response, StringList
} from "../http-signature/generated/generated";

export class HttpAdapter {
  private _runtime: Runtime<HttpContext>;

  constructor(runt: Runtime<HttpContext>) {
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

      const i = this._runtime.Instance(null);
      i.Context().ctx = c;
      i.Run();
      const newc = i.Context().ctx;
  
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
    const ctxresp = ctx.Response;

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
    const reqheaders = new Map<string, StringList>();

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
      reqheaders.set(k, new StringList(sl));
    }

    const preq = new Request(
      method,
      BigInt(body.length),
      protocol,
      ip,
      body,
      reqheaders
    );
    const presp = new Response(
      resp.statusCode,
      new Uint8Array(),
      new Map<string, StringList>()
    )
    const c = new Context(preq, presp);
    return c;
  }
}