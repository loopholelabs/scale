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
func (res *Response) SetStatusCode(s int32) int32 {
	res.value.StatusCode = s
	return res.value.StatusCode
}

// SetBody sets the body of the response
func (res *Response) SetBody(body string) []byte {
	res.value.Body = []byte(body)
	return res.value.Body
}

// SetBodyBytes sets the body of the response in bytes
func (res *Response) SetBodyBytes(body []byte) []byte {
	res.value.Body = body
	return res.value.Body
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
	return h.value[k].Value
}

// Set sets the value of the header at key k to value v
func (h *ResponseHeaders) Set(k string, v []string) {
	h.value[k].Value = v
}
