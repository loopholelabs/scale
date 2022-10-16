// Package http provides a Scale Runtime Adapter for the Standard net/http library.
package http

import (
	"github.com/loopholelabs/scale-go/runtime"
	"net/http"
)

// This compiler guard ensures that the HTTP adapter implements the net/http.Handler interface.
var _ http.Handler = (*HTTP)(nil)

type HTTP struct {
	next    http.Handler
	runtime *runtime.Runtime
}

func New(next http.Handler, runtime *runtime.Runtime) *HTTP {
	return &HTTP{
		next:    next,
		runtime: runtime,
	}
}

func (h *HTTP) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	i, err := h.runtime.Instance(req.Context(), h.Next(req))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = SerializeRequest(i.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = i.Run(req.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	err = DeserializeResponse(i.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *HTTP) Next(req *http.Request) runtime.Next {
	return func(ctx *runtime.Context) *runtime.Context {
		DeserializeRequest(ctx, req)
		w := NewResponseWriter()
		h.next.ServeHTTP(w, req)
		err := SerializeRequest(ctx, req)
		if err != nil {
			ctx.Error(err)
			return ctx
		}
		SerializeResponse(ctx, w)
		return ctx
	}
}
