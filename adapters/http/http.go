// Package http provides a Scale Runtime Adapter for the Standard net/http library.
package http

import (
	"github.com/loopholelabs/scale-go/runtime"
	"net/http"
)

// This compiler guard ensures that the HTTP adapter implements the net/http.Handler interface.
var _ http.Handler = (*HTTP)(nil)

type next struct {
	next http.Handler
	w    http.ResponseWriter
	r    *http.Request
}

type HTTP struct {
	next http.Handler
}

func New(next http.Handler, runtime *runtime.Runtime) *HTTP {
	return &HTTP{
		next: next,
	}
}

func (h *HTTP) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//run, err := runtime.NewContext(req.Context(), h.Next, nil)
	//
	//ctx := context.NewContext()
	//Serialize(ctx, req)
	//ctx = h.runtime.Run(ctx)
	//err := Deserialize(ctx, w)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//}
}

func (h *HTTP) Next(ctx *runtime.Context) {
	//h.next.ServeHTTP(ctx, ctx)
}
