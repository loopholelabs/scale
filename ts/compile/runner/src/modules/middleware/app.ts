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

import "../../runner";
import { GuestContext, HttpStringList } from "@loopholelabs/scale-signature-http";


function scalefn(inContext: GuestContext) {
  inContext.Next();

  // Lets add a header to show things are working...
  inContext.Response.SetHeader("FROM_TYPESCRIPT", new HttpStringList(["TRUE"]));
}

(global as any).scale = scalefn;