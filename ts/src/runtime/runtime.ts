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

import { Module } from "./module";
import { Context } from "./context";

export class Runtime {
  private _modules: Module[] = [];

  constructor(mods: Module[]) {
    let prevModule: Module | null = null;
    mods.forEach((m: Module) => {
      m.setNext(null);  // Clear our own next.
      this._modules.push(m);
      if (prevModule != null) {
        prevModule.setNext(m);
      }
      prevModule = m;
    });
    this._modules = mods;
  }

  run(context: Context): Context | null {
    // Run the first module only,
    // which will be already chained to the others
    if (this._modules.length === 0) return context;
    return this._modules[0].run(context);
  }
}
