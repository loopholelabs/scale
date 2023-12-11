/*
	Copyright 2023 Loophole Labs

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

import * as signature from "signature";

function delay(ms: number) {
  return new Promise( resolve => setTimeout(resolve, ms) );
}

function blah() {
  console.log("BLAH!");
}

export async function example(ctx?: signature.ModelWithAllFieldTypes): Promise<signature.ModelWithAllFieldTypes | undefined> {
    console.log("This is a Typescript Function");

    let ctime = (new Date()).getTime();

    let v150 = setTimeout(()=>{
      console.log("TIMEOUT 150");
    }, 150)

    let v10 = setTimeout(()=>{
      console.log("TIMEOUT 10");
    }, 10)

    let v20 = setTimeout(blah, 20)

    let vi10 = setInterval(()=>{
      console.log("INTERVAL 10");
    }, 10)

    let vi20 = setInterval(()=>{
      console.log("INTERVAL 20");
    }, 20)

    await delay(100);

    let dtime1 = (new Date()).getTime() - ctime;
    console.log("After small delay of " + dtime1 + "ms...");

    clearInterval(vi10);
    clearInterval(vi20);
    clearTimeout(v150);

    await delay(100);

    let dtime = (new Date()).getTime() - ctime;
    console.log("After small delay of " + dtime + "ms...");

    if (typeof ctx !== "undefined") {
        ctx.stringField = "This is a Typescript Function"
    }
    return signature.Next(ctx);
}

