package http

import (
	"github.com/loopholelabs/scale-go/runtime/context"
	"net/http"
	"strings"
)

func Deserialize(ctx *context.Context, w http.ResponseWriter) error {
	w.WriteHeader(int(ctx.Context.Response.StatusCode))
	for k, v := range ctx.Context.Response.Headers {
		w.Header().Set(k, strings.Join(v.Value, ","))
	}

	_, err := w.Write(ctx.Context.Response.Body)
	return err
}
