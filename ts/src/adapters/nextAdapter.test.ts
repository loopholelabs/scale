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

import { Host } from "../runtime/host";
import { NextAdapter } from "./nextAdapter";

describe("nextAdapter", () => {

  it("Can convert Request to Context", async () => {
    const request = new Request('https://example.com', {method: 'POST', body: '{"foo": "bar"}'});
    const response = new Response();

    const ctx = await NextAdapter.toContext(request, response);

    Host.showContext(ctx.context());

//    expect(res.headers.middleware).toBe("TRUE");

  });
});
