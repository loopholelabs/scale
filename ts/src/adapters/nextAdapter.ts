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

import { HttpContext } from "../http-signature/HttpContext";

import { Runtime } from "../runtime/runtime";
import {
  Context, Request, Response, StringList
} from "../http-signature/generated/generated";

import { NextRequest, NextResponse } from 'next/server';

// https://vercel.com/docs/concepts/functions/edge-functions#creating-edge-functions

export class NextAdapter {
  private _runtime: Runtime<HttpContext>;

  constructor(runt: Runtime<HttpContext>) {
    this._runtime = runt;
  }

  getHandler() {
    return async (req: NextRequest) => {

      const c = await NextAdapter.toContext(req);

      const i = await this._runtime.Instance(null);
      i.Context().ctx = c;
      i.Run();
      const newc = i.Context().ctx;
  
      if (newc!=null) {
        return NextAdapter.fromContext(newc);
      }

      // Return server error
      return new NextResponse("Internal Server Error", {status: 500});
    };
  }

  static fromContext(ctx: Context) {
    const ctxresp = ctx.Response;

    const headers = new Headers();

    for(let k of ctxresp.Headers.keys()) {
      let vals = ctxresp.Headers.get(k);
      if (vals!==undefined) {
        let s = vals.Value;
        for(let v of s.values()) {
            headers.set(k, v);
        }
      }
    }

    const resp = new NextResponse(ctxresp.Body, {
      status: ctxresp.StatusCode,
      headers: headers
    });
    
    return resp;
  }

  static async toContext(req: NextRequest): Promise<Context> {
    const reqheaders = new Map<string, StringList>();

    let method = "unknown";
    if (req.method != undefined) {
      method = req.method;
    }

    const protocol = (new URL(req.url)).protocol;

    // TODO: Find IP
    const ip = "1.2.3.4";

    for(let k in req.headers) {
      let vals = req.headers.get(k);
      let sl: string[] = [];
      if (typeof vals === 'string') {
        sl.push(vals);  // Single value
      }
      reqheaders.set(k, new StringList(sl));
    }

    // Read all of body
    
    let body = new Uint8Array();

    if (req.body!=null) {
      // Read all of req.body
      let b = await (await req.blob()).arrayBuffer();
      body = new Uint8Array(b);
    }

    const preq = new Request(
      method,
      BigInt(body.length),
      protocol,
      ip,
      body,
      reqheaders
    );

    // Dummy response filled in here...
    let status = 200;
    const presp = new Response(
      status,
      new Uint8Array(),
      new Map<string, StringList>()
    )
    const c = new Context(preq, presp);
    return c;
  }
}