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

import {
  Headers,
  Request,
  Response,
} from 'node-fetch';

if (!global.fetch) {
//  (global as any).fetch = fetch;
  (global as any).Headers = Headers;
  (global as any).Request = Request;
  (global as any).Response = Response;
}

import { TextEncoder, TextDecoder } from "util";

window.TextEncoder = TextEncoder;
window.TextDecoder = TextDecoder as typeof window["TextDecoder"];

import { Host } from "../runtime/host";
import { NextAdapter } from "./nextAdapter";

import { NextRequest, NextResponse } from 'next/server';
import { Context } from "../runtime/context";
import { Context as PgContext, Request as PgRequest, Response as PgResponse, StringList as PgStringList } from "../runtime/generated/generated";

describe("nextAdapter", () => {

  it("Can convert Request to Context", async () => {
    const bodyData = '{"foo": "bar"}';
    const request = new NextRequest('https://example.com', {method: 'POST', body: bodyData});

    const ctx = await NextAdapter.toContext(request);

    Host.showContext(ctx.context());

    if (request.body != null ) {
      expect(ctx.context().Request.Method).toBe(request.method);
      expect(ctx.context().Request.Protocol).toBe((new URL(request.url)).protocol);
      expect(Number(ctx.context().Request.ContentLength)).toBe(bodyData.length);
      const reqBody = new TextDecoder().decode(ctx.context().Request.Body);
      expect(reqBody).toBe(bodyData);
    }
  });

  it("Can convert Context to Response", async () => {
    const req = new PgRequest("GET", BigInt(0), "http", "1.2.3.4", new Uint8Array(), new Map<string, PgStringList>);
    
    const body = new TextEncoder().encode("Hello world");
    const headers = new Map<string, PgStringList>;

    headers.set("MIDDLEWARE", new PgStringList(["Hello"]));    

    const resp = new PgResponse(200, body, headers);
    const c = new PgContext(req, resp);
    const ctx = new Context(c);

    const response = NextAdapter.fromContext(ctx);
  });

});
