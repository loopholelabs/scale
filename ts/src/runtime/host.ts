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

import { Context, Request, Response, StringList } from "./generated/generated";

export class Host {
    constructor() {

    }

    // Pack a pointer and length into a single 64bit
    public static packMemoryRef(ptr: number, len: number): BigInt {
        if (ptr > 0xffffffff || len > 0xffffffff) {
            // Error! We can't do it.
        }
        return (BigInt(ptr) << BigInt(32)) | BigInt(len);
    }

    // Unpack a memory ref from 64bit to 2x32bits
    public static unpackMemoryRef(packed: bigint): [number, number] {
        let ptr = Number((packed >> BigInt(32)) & BigInt(0xffffffff));
        let len = Number((packed & BigInt(0xffffffff)));
        return [ptr, len];
    }

    public static showContext(c: Context) {
        let req = c.Request;
        console.log("== Context ==");
        console.log("Request method=" + req.Method + ", proto=" + req.Protocol + ", ip=" + req.IP + ", len=" + req.ContentLength);
        console.log(" Headers: " + Host.stringHeaders(req.Headers));
        let reqBody = new TextDecoder().decode(req.Body);
        console.log(" Body: " + reqBody);
        let resp = c.Response;
        console.log("Response code=" + resp.StatusCode);
        console.log(" Headers: " + Host.stringHeaders(resp.Headers));
        let respBody = new TextDecoder().decode(resp.Body);
        console.log(" Body: " + respBody);
    }

    public static stringHeaders(h: Map<string, StringList>): string {
        let r = '';
        for (let k of h.keys()) {
            let values = h.get(k);
            if (values!=undefined) {
                for (let i of values.Value.values()) {
                    r = r + " " + k + "=" + i;
                }
            } else {
                r = r + " " + k;
            }
        }
        return r;
    }
}