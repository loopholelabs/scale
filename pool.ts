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

import {Signature} from "@loopholelabs/scale-signature-interfaces";
import {Template} from "./template";
import {NewModule, Module} from "./module";

export class ModulePool<T extends Signature> {
	private readonly objects: Module<T>[] = [];
	private readonly allocator: () => Promise<Module<T>>;

	constructor(template: Template<T>) {
		this.allocator = async (): Promise<Module<T>> => {
			return await NewModule(template);
		}
	}

	public async Get(): Promise<Module<T>> {
		if (this.objects.length > 0) {
			const t = this.objects.pop();
			if (typeof t !== "undefined") {
				return t;
			}
		}
		return await this.allocator()
	}

	public Put(t: Module<T>): void {
		this.objects.push(t);
	}
}