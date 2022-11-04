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

import next from "next";
import { createServer } from "http";

import { Context } from "../runtime/context";
import { Module } from "../runtime/module";
import {
  Context as PgContext,
  Request as PgRequest,
  Response as PgResponse,
  StringList as PgStringList,
} from "../runtime/generated/generated";


const dev = process.env.NODE_ENV !== 'production'
const hostname = 'localhost'
const port = 3000
// when using middleware `hostname` and `port` must be provided below
const app = next({ dev, hostname, port })
const handle = app.getRequestHandler()

console.log("Starting nextAdapter example");

app.prepare().then(() => {
  createServer(async (req, res) => {
    console.log(req);

  }).listen(port, () => {
    console.log(`> Ready on http://${hostname}:${port}`)
  })
/*
export class NextAdapter {
  private _module: Module;

  constructor(mod: Module) {
    this._module = mod;
  }

  handler(
    req: next.Request,
    res: express.Response,
    next: express.NextFunction
  ) {
    const c = ExpressAdapter.toContext(req, res);
    const newc = this._module.run(c);

    if (newc != null) {
      // Now write it back out...
      ExpressAdapter.fromContext(newc, res);
    }
    // next();
  }

  static fromContext(ctx: Context, resp: express.Response) {
    const response = ctx.context().Response;
    for(let k of response.Headers.keys()) {
      let vals = response.Headers.get(k);
      if (vals!==undefined) {
        let s = vals.Value;
        for(let v of s.values()) {
            resp.setHeader(k, v);
        }
      }
    }

    const respBody = new TextDecoder().decode(response.Body);
    resp.status(response.StatusCode).send(respBody);
  }

  static toContext(req: express.Request, resp: express.Response): Context {
    const reqheaders = new Map<string, PgStringList>();

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

    let bodylen = 0;
    let body = new Uint8Array();
    if (req.body) {
      if (req.body.length !== undefined) {
        bodylen = req.body.length;
      }
      body = req.body;
    }

    const preq = new PgRequest(
      req.method,
      BigInt(bodylen),
      req.protocol,
      req.ip,
      body,
      reqheaders
    );

    const presp = new PgResponse(
      resp.statusCode,
      new Uint8Array(),
      new Map<string, PgStringList>()
    ); // respBody, respheaders);
    return new Context(new PgContext(preq, presp));
  }
}
*/