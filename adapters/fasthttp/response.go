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

package fasthttp

import (
	"github.com/loopholelabs/scale-go/runtime"
	"github.com/valyala/fasthttp"
	"strings"
)

func ToResponseContext(ctx *runtime.Context, fastCTX *fasthttp.RequestCtx) {
	fastCTX.Response.SetStatusCode(int(ctx.Context.Response.StatusCode))
	fastCTX.Response.SetBody(ctx.Context.Response.Body)

	for k, v := range ctx.Context.Response.Headers {
		fastCTX.Response.Header.Set(k, strings.Join(v.Value, ","))
	}
}
