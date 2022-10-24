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
	"github.com/loopholelabs/scale/generated"
	"github.com/loopholelabs/scale/go/runtime"
	"github.com/valyala/fasthttp"
	"strings"
)

func FromRequestContext(ctx *runtime.Context, fastCTX *fasthttp.RequestCtx) {
	ctx.Context.Request.Protocol = "HTTP/1.1"
	ctx.Context.Request.Method = string(fastCTX.Request.Header.Method())
	ctx.Context.Request.IP = fastCTX.RemoteAddr().String()
	ctx.Context.Request.ContentLength = int64(fastCTX.Request.Header.ContentLength())
	ctx.Context.Request.Body = fastCTX.Request.Body()

	fastCTX.Request.Header.VisitAll(func(key []byte, value []byte) {
		ctx.Context.Request.Headers[string(key)] = &generated.StringList{
			Value: strings.Split(string(value), ","),
		}
	})
}

func ToRequestContext(ctx *runtime.Context, fastCTX *fasthttp.RequestCtx) {
	fastCTX.Request.Header.SetMethod(ctx.Context.Request.Method)
	fastCTX.Request.Header.SetContentLength(int(ctx.Context.Request.ContentLength))
	fastCTX.Request.SetBody(ctx.Context.Request.Body)

	for k, v := range ctx.Context.Request.Headers {
		fastCTX.Request.Header.Set(k, strings.Join(v.Value, ","))
	}
}
