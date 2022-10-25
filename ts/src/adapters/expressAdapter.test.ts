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

import express from 'express';
import bodyParser from 'body-parser';

import * as fs from 'fs';
import { Module } from '../runtime/module';
import { ExpressAdapter } from './expressAdapter';

import request from 'supertest';

window.TextEncoder = TextEncoder;
window.TextDecoder = TextDecoder as typeof window["TextDecoder"];

describe("expressAdapter", () => {
    let port = 8090;
    let app = express();

    it("Can run a simple e2e", async () => {
        const modHttpEndpoint = fs.readFileSync('./example_modules/http-endpoint.wasm');
        const modHttpMiddleware = fs.readFileSync('./example_modules/http-middleware.wasm');
        let moduleHttpEndpoint = new Module(modHttpEndpoint, null);
        let moduleHttpMiddleware = new Module(modHttpMiddleware, moduleHttpEndpoint);
        
        var adapter = new ExpressAdapter(moduleHttpMiddleware);
        
        app.use(bodyParser.raw({
            type: (t)=>true,
        }));
        app.use(adapter.handler.bind(adapter));
        
        let res = await request(app).post("/blah").send("HELLO WORLD");

        // Make sure everything worked as expected.
        expect(res.statusCode).toEqual(200);
        expect(res.text).toBe("HELLO WORLD");
        expect(res.headers["middleware"]).toBe("TRUE");
    });
});