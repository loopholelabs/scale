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

/* eslint no-bitwise: off */

import { Context, StringList } from "./generated/generated";

export class Host {
  // Pack a pointer and length into a single 64bit
  public static packMemoryRef(ptr: number, len: number): BigInt {
    if (ptr > 0xffffffff || len > 0xffffffff) {
      // Error! We can't do it.
    }
    return (BigInt(ptr) << BigInt(32)) | BigInt(len);
  }

  // Unpack a memory ref from 64bit to 2x32bits
  public static unpackMemoryRef(packed: bigint): [number, number] {
    const ptr = Number((packed >> BigInt(32)) & BigInt(0xffffffff));
    const len = Number(packed & BigInt(0xffffffff));
    return [ptr, len];
  }

  public static stringHeaders(h: Map<string, StringList>): string {
    let r = "";
    for (let k of h.keys()) {
      let values = h.get(k);
      if (values != undefined) {
        for (let i of values.Value.values()) {
          r = r + " " + k + "=" + i;
        }
      } else {
        r = r + " " + k;
      }
    }
    return r;
  }

  public static showContext(c: Context) {
    const req = c.Request;
    const reqBody = new TextDecoder().decode(req.Body);
    console.log(`== Context ==
Request method=${req.Method}, proto=${req.Protocol}, ip=${req.IP}, len=${req.ContentLength}
 Headers: ${Host.stringHeaders(req.Headers)}
 Body: ${reqBody}`);

    const resp = c.Response;
    const respBody = new TextDecoder().decode(resp.Body);
    console.log(`Response code=${resp.StatusCode}
 Headers: ${Host.stringHeaders(resp.Headers)}
 Body: ${respBody}`);
  }
}
