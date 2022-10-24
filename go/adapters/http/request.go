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

package http

import (
	"bytes"
	"github.com/loopholelabs/scale/go/runtime"
	"github.com/loopholelabs/scale/go/runtime/generated"
	"io"
	"net/http"
)

const (
	BodyLimit = 1024 * 1024 * 10
)

// FromRequest serializes http.Request object into a runtime.Context
func FromRequest(ctx *runtime.Context, req *http.Request) error {
	for k, v := range req.Header {
		ctx.Context.Request.Headers[k] = &generated.StringList{
			Value: v,
		}
	}
	ctx.Context.Request.Method = req.Method
	ctx.Context.Request.ContentLength = req.ContentLength
	ctx.Context.Request.Protocol = req.Proto
	ctx.Context.Request.IP = req.RemoteAddr

	if req.ContentLength != 0 {
		var err error
		ctx.Context.Request.Body, err = io.ReadAll(io.LimitReader(req.Body, BodyLimit))
		if err != nil {
			return err
		}
	} else {
		ctx.Context.Request.Body = nil
	}

	return nil
}

// ToRequest deserializes the runtime.Context object into an existing http.Request
func ToRequest(ctx *runtime.Context, req *http.Request) {
	req.Method = ctx.Context.Request.Method
	req.ContentLength = ctx.Context.Request.ContentLength
	req.Proto = ctx.Context.Request.Protocol
	req.RemoteAddr = ctx.Context.Request.IP

	for k, v := range ctx.Context.Request.Headers {
		req.Header[k] = v.Value
	}

	if ctx.Context.Request.ContentLength != 0 {
		req.Body = io.NopCloser(bytes.NewReader(ctx.Context.Request.Body))
	} else {
		req.Body = io.NopCloser(bytes.NewReader(nil))
	}
}
