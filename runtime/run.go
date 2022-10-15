package runtime

import (
	"errors"
)

var (
	MemoryWriteError = errors.New("memory write error")
)

func (r *Runtime) Run(ctx *Context) error {
	rootFunc := r.functions[0]
	r.c = ctx
	return rootFunc.Run(r)
}
