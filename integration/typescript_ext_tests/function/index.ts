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

import * as extension from "extension";
import * as types from "extension/types"

export function example(ctx?: signature.ModelWithAllFieldTypes): signature.ModelWithAllFieldTypes | undefined {
    console.log("This is a Typescript Function, which calls an extension.");

    const s = new types.Stringval();
    s.value = "dummy";
    const ex = extension.New(s);
    const hello = ex.Hello(s);
    const world = extension.World(s);

    if (typeof ctx !== "undefined") {
        ctx.stringField = "This is a Typescript Function. Extension New().Hello()=" + hello.value + " World()=" + world.value;
    }
    return signature.Next(ctx);
}