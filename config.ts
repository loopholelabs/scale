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



const envStringRegex = /[^A-Za-z0-9_]/;

class ConfigFunction {
    function: ScaleFunc.Schema;
    env: { [key: string]: string };
}

export class Config<T extends Signature> {
    newSignature: NewSignature<T>;
    functions: ConfigFunction[] = [];
    context: any = null; // Note: You'll want to import or define a more specific context type.
    pooling: boolean;

    constructor(newSignature: NewSignature<T>) {
        this.newSignature = newSignature;
    }
}