//go:build tinygo
// +build tinygo

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

package context

import "github.com/loopholelabs/scale-go/runtime/generated"

// Response is the HTTP Response object
type Response struct {
	value *generated.Response
}

// Response returns the Response object for the Context
func (ctx *Context) Response() *Response {
	return &Response{value: ctx.generated.Response}
}

// Body returns the body of the response
func (res *Response) Body() []byte {
	return res.value.Body
}

// StatusCode returns the status code of the response
func (res *Response) StatusCode() int32 {
	return res.value.StatusCode
}

// SetStatusCode sets the status code of the response
func (res *Response) SetStatusCode(s int32) {
	res.value.StatusCode = s
}

// SetBody sets the body of the response
func (res *Response) SetBody(body string) {
	res.value.Body = []byte(body)
}

// SetBodyBytes sets the body of the response in bytes
func (res *Response) SetBodyBytes(body []byte) {
	res.value.Body = body
}

// ResponseHeaders are the headers in the response
type ResponseHeaders struct {
	value generated.ResponseHeadersMap
}

// Headers returns the headers of the response
func (res *Response) Headers() *ResponseHeaders {
	return &ResponseHeaders{
		value: res.value.Headers,
	}
}

// Get returns the value of the header at key k
func (h *ResponseHeaders) Get(k string) []string {
	v := h.value[k]
	if v == nil {
		return nil
	}
	return v.Value
}

// Set sets the value of the header at key k to value v
func (h *ResponseHeaders) Set(k string, v []string) {
	h.value[k] = &generated.StringList{
		Value: v,
	}
}
