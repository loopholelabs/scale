package http

import (
	"github.com/loopholelabs/scale-go/runtime/context"
	"github.com/loopholelabs/scale-go/runtime/context/generated"
	"net/http"
)

func Serialize(ctx *context.Context, r *http.Request) {
	ctx.Context.Request.Headers = generated.NewRequestHeadersMap(uint32(len(r.Header)))
	for k, v := range r.Header {
		ctx.Context.Request.Headers[k] = &generated.StringList{
			Value: v,
		}
	}
	ctx.Context.Request.Method = r.Method
	ctx.Context.Request.Length = r.ContentLength
	ctx.Context.Request.Protocol = r.Proto
	ctx.Context.Request.Ip = r.RemoteAddr
	ctx.Context.Request.Body = []byte("")

	ctx.Encode()
}
