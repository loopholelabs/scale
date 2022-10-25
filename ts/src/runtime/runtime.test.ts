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

import { TextEncoder, TextDecoder } from 'util';
import * as fs from 'fs';
import { Module } from './module';
import { Context, Request, Response, StringList } from "./generated/generated";
import { Context as ourContext} from './context';

window.TextEncoder = TextEncoder;
window.TextDecoder = TextDecoder as typeof window["TextDecoder"];

describe("runtime", () => {
    it("Can run a simple e2e", () => {

        // Create a context to send in...
        let enc = new TextEncoder();
        let body = enc.encode("Hello world this is a request body");
        let headers = new Map<string, StringList>();
        headers.set('content', new StringList(['hello']));
        let req1 = new Request('GET', BigInt(100), 'https', '1.2.3.4', body, headers);
        let respBody = enc.encode("Response body");
        let respHeaders = new Map<string, StringList>();        
        const resp1 = new Response(200, respBody, respHeaders);        
        const context = new Context(req1, resp1);

        // Now we can use context with a couple of wasm modules...

        const modHttpEndpoint = fs.readFileSync('./example_modules/http-endpoint.wasm');
        const modHttpMiddleware = fs.readFileSync('./example_modules/http-middleware.wasm');
        let moduleHttpEndpoint = new Module(modHttpEndpoint, null);
        let moduleHttpMiddleware = new Module(modHttpMiddleware, moduleHttpEndpoint);

        // Run the modules...

        let ctx = new ourContext(context);

        let retContext = moduleHttpMiddleware.run(ctx);

        // check the returns...

        let dec = new TextDecoder();
        let bodyText = dec.decode(retContext.context().Response.Body);

        // The http-endpoint.wasm module copies the request body to the response body.
        expect(bodyText).toBe("Hello world this is a request body");        

        // The http-middleware.wasm adds a header
        let middle = retContext.context().Response.Headers.get("MIDDLEWARE");
        expect(middle).toBeDefined();
        let vals = middle?.Value;
        if (vals!==undefined) {
            expect(vals.length).toBe(1);
            expect(vals[0]).toBe("TRUE");
        }
    })
});