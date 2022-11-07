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

import { TextEncoder, TextDecoder } from "util";

import { createServer } from "http";
import request from "supertest";

import { Context } from "../runtime/context";
import { Module } from "../runtime/module";
import {
  Context as PgContext,
  Request as PgRequest,
  Response as PgResponse,
  StringList as PgStringList,
} from "../runtime/generated/generated";
import { time } from "console";

window.TextEncoder = TextEncoder;
window.TextDecoder = TextDecoder as typeof window["TextDecoder"];

describe("nextAdapter", () => {
  it("Can run a simple e2e", async () => {
/*
    const port = 3000;

    console.log("Starting nextAdapter example");

    createServer(async (req, res) => {
      console.log(req);
      res.writeHead(200, { "Content-Type": "text/plain" });
      res.write("Hello World!");
      res.end();
    }).listen(port, () => {
      console.log(`> Ready on ${port}`);
    });

    // Now fire off a request and check it worked...


    
    const delay = (ms: number) =>
      new Promise((resolve) => setTimeout(resolve, ms));

    console.log("Waiting...");
    await delay(30 * 1000);
*/
  });
});

