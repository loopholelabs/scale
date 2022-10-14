package context

import "github.com/loopholelabs/scale-go/runtime/generated"

type Request struct {
	value *generated.Request
}

type RequestHeaders struct {
	value generated.RequestHeadersMap
}

func (ctx *Context) Request() *Request {
	return &Request{value: ctx.generated.Request}
}

func (req *Request) Method() string {
	return req.value.Method
}

func (req *Request) SetMethod(method string) string {
	req.value.Method = method
	return req.value.Method
}

func (req *Request) RemoteIP() string {
	return req.value.Ip
}

func (req *Request) Headers() *RequestHeaders {
	return &RequestHeaders{
		value: req.value.Headers,
	}
}

func (h *RequestHeaders) Get(k string) []string {
	return h.value[k].Value
}

func (h *RequestHeaders) Set(k string, v []string) {
	h.value[k].Value = v
}
