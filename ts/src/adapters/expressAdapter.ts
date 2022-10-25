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
import express from 'express';
import { Context } from '../runtime/context';
import { Module } from '../runtime/module';
import { Context as pgContext, Request as pgRequest, Response as pgResponse, StringList as pgStringList } from "../runtime/generated/generated";

export class ExpressAdapter {
    private _module: Module;

    constructor(mod: Module) {
        this._module = mod;
    }

    handler(req:express.Request, res:express.Response, next:express.NextFunction) {
        let c = ExpressAdapter.toContext(req, res);
        let newc = this._module.run(c);

        // Now write it back out...
        ExpressAdapter.fromContext(newc, res);

//        next();
    }

    static fromContext(ctx: Context, resp: express.Response) {
        let response = ctx.context().Response;
        for(let k of response.Headers.keys()) {
            let vals = response.Headers.get(k);
            if (vals!==undefined) {
                let s = vals.Value;
                for(let v of s.values()) {
                    resp.setHeader(k, v);
                }
            }
        }

        let respBody = new TextDecoder().decode(response.Body);
        resp.status(response.StatusCode).send(respBody);
    }

    static toContext(req: express.Request, resp: express.Response): Context {

        let reqheaders = new Map<string, pgStringList>();

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
            reqheaders.set(k, new pgStringList(sl));
        }

        let bodylen = 0;
        let body = new Uint8Array();
        if (req.body) {
            if (req.body.length!==undefined) {
                bodylen = req.body.length;
            }
            body = req.body;
        }

        let preq = new pgRequest(req.method, BigInt(bodylen), req.protocol, req.ip, body, reqheaders);

/*
        let respheaders = new Map<string, pgStringList>();
//
        for(let k in resp.getHeaders()) {
            let vals = resp.getHeaders()[k];
            let sl: string[] = [];
            if (typeof vals === 'string') {
                sl.push(vals);  // Single value
            } else if (typeof vals === 'number' ) {
                sl.push('' + vals);
            } else if (vals===undefined) {
                // Just empty values
            } else {
                sl = vals;  // Multiple values
            }
            respheaders.set(k, new pgStringList(sl));
        }

        var enc = new TextEncoder();
        let respBody = enc.encode("TODO: Response body");
*/
        let presp = new pgResponse(resp.statusCode, new Uint8Array(), new Map<string, pgStringList>()); //respBody, respheaders);
        return new Context(new pgContext(preq, presp));
    }

}