//go:build !tinygo

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

type Response struct{}

// Response returns the Response object for the Context
func (ctx *Context) Response() *Response {
	return &Response{}
}

func (res *Response) Body() []byte {
	return nil
}

func (res *Response) StatusCode() int32 {
	return 0
}

func (res *Response) SetStatusCode(int32) {}

func (res *Response) SetBody(string) {}

func (res *Response) SetBodyBytes([]byte) {}

type ResponseHeaders struct{}

func (res *Response) Headers() *ResponseHeaders {
	return &ResponseHeaders{}
}

func (h *ResponseHeaders) Get(string) []string {
	return nil
}

func (h *ResponseHeaders) Set(k string, v []string) {}
