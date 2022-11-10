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

export class NextAdapter {
  private _runtime: Runtime;

  constructor(runt: Runtime) {
    this._runtime = runt;
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

  static async toContext(req: Request, resp: Response) {
    const reqheaders = new Map<string, PgStringList>();

    let method = "unknown";
    if (req.method != undefined) {
      method = req.method;
    }

    // TODO: Find protocol
    // TODO: Find IP

    for(let k in req.headers) {
      let vals = req.headers.get(k);
      let sl: string[] = [];
      if (typeof vals === 'string') {
        sl.push(vals);  // Single value
      }
      reqheaders.set(k, new PgStringList(sl));
    }

    // Read all of body
    
    let body = new Uint8Array();

    if (req.body!=null) {
      // Read all of req.body
      let b = await (await req.blob()).arrayBuffer();
      body = new Uint8Array(b);
    }

    const preq = new PgRequest(
      method,
      BigInt(body.length),
      "http",
      "1.2.3.4",
      body,
      reqheaders
    );
    const presp = new PgResponse(
      resp.status,
      new Uint8Array(),
      new Map<string, PgStringList>()
    )
    const c = new Context(new PgContext(preq, presp));
    return c;
  }
}