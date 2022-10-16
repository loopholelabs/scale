package http

import (
	"bytes"
	"github.com/loopholelabs/scale-go/runtime"
	"github.com/loopholelabs/scale-go/runtime/generated"
	"io"
	"net/http"
)

// SerializeRequest serializes http.Request object into a runtime.Context
func SerializeRequest(ctx *runtime.Context, req *http.Request) {
	for k, v := range req.Header {
		ctx.Context.Request.Headers[k] = &generated.StringList{
			Value: v,
		}
	}
	ctx.Context.Request.Method = req.Method
	ctx.Context.Request.Length = req.ContentLength
	ctx.Context.Request.Protocol = req.Proto
	ctx.Context.Request.Ip = req.RemoteAddr
	ctx.Context.Request.Body = []byte("")

	ctx.Serialize()
}

// DeserializeRequest deserializes the runtime.Context object into an existing http.Request
func DeserializeRequest(ctx *runtime.Context, req *http.Request) {
	req.Method = ctx.Context.Request.Method
	req.ContentLength = ctx.Context.Request.Length
	req.Proto = ctx.Context.Request.Protocol
	req.RemoteAddr = ctx.Context.Request.Ip

	for k, v := range ctx.Context.Request.Headers {
		req.Header[k] = v.Value
	}

	if ctx.Context.Request.Body != nil {
		req.Body = io.NopCloser(bytes.NewReader(ctx.Context.Request.Body))
	}
}
