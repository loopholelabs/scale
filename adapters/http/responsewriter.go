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
	"net/http"
)

var _ http.ResponseWriter = (*ResponseWriter)(nil)

type ResponseWriter struct {
	headers    http.Header
	statusCode int
	buffer     *bytes.Buffer
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		headers: make(http.Header),
		buffer:  bytes.NewBuffer(nil),
	}
}

func (r *ResponseWriter) Header() http.Header {
	return r.headers
}

func (r *ResponseWriter) Write(bytes []byte) (int, error) {
	return r.buffer.Write(bytes)
}

func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}
