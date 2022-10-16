package http

import (
	"fmt"
	"github.com/loopholelabs/scale-go/runtime"
	"github.com/loopholelabs/scale-go/runtime/generated"
	"net/http"
	"strings"
)

// SerializeResponse serializes the *ResponseWriter object into a runtime.Context
func SerializeResponse(ctx *runtime.Context, w *ResponseWriter) {
	ctx.Context.Response.StatusCode = int32(w.statusCode)
	for k, v := range w.headers {
		ctx.Context.Response.Headers[k] = &generated.StringList{
			Value: v,
		}
	}
	ctx.Context.Response.Body = w.buffer.Bytes()
}

// DeserializeResponse deserializes the runtime.Context object into the http.ResponseWriter
func DeserializeResponse(ctx *runtime.Context, w http.ResponseWriter) error {
	for k, v := range ctx.Context.Response.Headers {
		w.Header().Set(k, strings.Join(v.Value, ","))
	}
	w.WriteHeader(int(ctx.Context.Response.StatusCode))

	_, err := w.Write(ctx.Context.Response.Body)
	if err != nil {
		return fmt.Errorf("error writing response body: %w", err)
	}
	return nil
}
