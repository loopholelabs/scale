package runtime

import (
	"errors"
)

var (
	MemoryWriteError = errors.New("memory write error")
)

func (r *Runtime) Run(ctx *Context) error {
	rootModule := r.functions[0].Module
	run := rootModule.ExportedFunction("run")
	malloc := rootModule.ExportedFunction("malloc")
	free := rootModule.ExportedFunction("free")

	length := uint64(ctx.Buffer.Len())
	buffer, err := malloc.Call(r.ctx, length)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = free.Call(r.ctx, buffer[0])
	}()

	if !rootModule.Memory().Write(r.ctx, uint32(buffer[0]), ctx.Buffer.Bytes()) {
		return MemoryWriteError
	}

	_, err = run.Call(r.ctx, buffer[0], length)
	return err
}
