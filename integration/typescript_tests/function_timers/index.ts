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

let output = "";

function delay(ms: number) {
  return new Promise( resolve => setTimeout(resolve, ms) );
}

function blah() {
  output = output + " BLAH";
}

export async function example(ctx?: signature.ModelWithAllFieldTypes): Promise<signature.ModelWithAllFieldTypes | undefined> {
    console.log("This is a Typescript Function");

    let ctime = (new Date()).getTime();

    let v150 = setTimeout(()=>{
      output = output + " TIMEOUT 150";
    }, 150)

    let v151 = setTimeout(()=>{
      output = output + " TIMEOUT 151";
    }, 151)

    let v10 = setTimeout(()=>{
      output = output + " TIMEOUT 50";
    }, 50)

    let v20 = setTimeout(blah, 20)

    let vi10 = setInterval(()=>{
      output = output + " INTERVAL 19";
    }, 19)

    let vi20 = setInterval(()=>{
      output = output + " INTERVAL 18"
    }, 18)

    await delay(100);

    let dtime1 = (new Date()).getTime() - ctime;
    // round it to within 20
    dtime1 = Math.round(dtime1 / 20) * 20;
    output = output + " DELAY " + dtime1;

    clearInterval(vi10);
    clearInterval(vi20);
    clearTimeout(v150);

    await delay(100);

    let dtime = (new Date()).getTime() - ctime;
    dtime = Math.round(dtime / 20) * 20;
    output = output + " DELAY " + dtime;

    if (typeof ctx !== "undefined") {
        ctx.stringField = "This is a Typescript Function. " + output;
    }
    return signature.Next(ctx);
}

