package context

import "github.com/loopholelabs/scale-go/runtime/generated"

type Response struct {
	value *generated.Response
}

func (ctx *Context) Response() *Response {
	return &Response{value: ctx.generated.Response}
}

func (res *Response) Body() []byte {
	return res.value.Body
}

func (res *Response) StatusCode() int32 {
	return res.value.StatusCode
}

func (res *Response) SetStatusCode(s int32) int32 {
	res.value.StatusCode = s
	return res.value.StatusCode
}

func (res *Response) SetBody(body string) []byte {
	res.value.Body = []byte(body)
	return res.value.Body
}

func (res *Response) SetBodyBytes(body []byte) []byte {
	res.value.Body = body
	return res.value.Body
}

type ResponseHeaders struct {
	value generated.ResponseHeadersMap
}

func (res *Response) Headers() *ResponseHeaders {
	return &ResponseHeaders{
		value: res.value.Headers,
	}
}
