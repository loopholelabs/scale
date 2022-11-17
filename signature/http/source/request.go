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

// Request is the HTTP Request object
type Request struct {
	value *HttpRequest
}

// Request returns the Request object for the Context
func (x *Context) Request() *Request {
	return &Request{value: x.generated.Request}
}

// Method returns the method of the request
func (req *Request) Method() string {
	return req.value.Method
}

// SetMethod sets the method of the request
func (req *Request) SetMethod(method string) string {
	req.value.Method = method
	return req.value.Method
}

// RemoteIP returns the remote IP of the request
func (req *Request) RemoteIP() string {
	return req.value.IP
}

func (req *Request) Body() []byte {
	return req.value.Body
}

func (req *Request) SetBody(body string) {
	req.value.Body = []byte(body)
}

func (req *Request) SetBodyBytes(body []byte) {
	req.value.Body = body
}

// RequestHeaders is are the headers in the request
type RequestHeaders struct {
	value HttpRequestHeadersMap
}

// Headers returns the headers of the request
func (req *Request) Headers() *RequestHeaders {
	return &RequestHeaders{
		value: req.value.Headers,
	}
}

// Get returns the value of the header
func (h *RequestHeaders) Get(k string) []string {
	v := h.value[k]
	if v == nil {
		return nil
	}
	return v.Value
}

// Set sets the value of the header
func (h *RequestHeaders) Set(k string, v []string) {
	h.value[k] = &HttpStringList{
		Value: v,
	}
}
