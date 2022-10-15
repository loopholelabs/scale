package context

import "github.com/loopholelabs/scale-go/runtime/generated"

// Request is the HTTP Request object
type Request struct {
	value *generated.Request
}

// Request returns the Request object for the Context
func (ctx *Context) Request() *Request {
	return &Request{value: ctx.generated.Request}
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
	return req.value.Ip
}

// RequestHeaders is are the headers in the request
type RequestHeaders struct {
	value generated.RequestHeadersMap
}

// Headers returns the headers of the request
func (req *Request) Headers() *RequestHeaders {
	return &RequestHeaders{
		value: req.value.Headers,
	}
}

// Get returns the value of the header
func (h *RequestHeaders) Get(k string) []string {
	return h.value[k].Value
}

// Set sets the value of the header
func (h *RequestHeaders) Set(k string, v []string) {
	h.value[k].Value = v
}
