package http

import (
	"bytes"
	"fmt"
	"github.com/loopholelabs/scale-go/runtime"
	"github.com/loopholelabs/scale-go/runtime/generated"
	"net/http"
	"strings"
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

// SerializeResponse serializes the *ResponseWriter object into a runtime.Context
func SerializeResponse(ctx *runtime.Context, w *ResponseWriter) {
	ctx.Context.Response.StatusCode = int32(w.statusCode)
	for k, v := range w.headers {
		ctx.Context.Response.Headers[k] = &generated.StringList{
			Value: v,
		}
	}
	ctx.Context.Response.Body = w.buffer.Bytes()
	ctx.Serialize()
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
