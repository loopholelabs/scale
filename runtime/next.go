package runtime

import (
	"context"
	"fmt"
	"github.com/loopholelabs/scale-go/utils"
	"github.com/tetratelabs/wazero/api"
)

type Next func(ctx *Context) *Context

func (r *Runtime) next(ctx context.Context, module api.Module, offset uint32, length uint32) uint64 {
	malloc := module.ExportedFunction("malloc")
	free := module.ExportedFunction("free")

	buf, ok := module.Memory().Read(ctx, offset, length)
	if !ok {
		panic("failed to read memory")
	}

	err := r.c.Deserialize(buf)
	if err != nil {
		panic("failed to deserialize context")
	}

	f := r.modules[module]
	if f.Next == nil {
		r.c = r.Next(r.c)
	} else {
		err = f.Next.Run(r)
		if err != nil {
			panic(fmt.Errorf("failed to run next function: %w", err))
		}
	}

	r.c.Encode()
	bufLength := uint64(r.c.Buffer.Len())
	buffer, err := malloc.Call(r.ctx, bufLength)
	if err != nil {
		panic(fmt.Errorf("failed to allocate memory: %w", err))
	}
	defer func() {
		_, _ = free.Call(r.ctx, buffer[0])
	}()

	if !module.Memory().Write(r.ctx, uint32(buffer[0]), r.c.Buffer.Bytes()) {
		panic("failed to write memory")
	}

	return utils.PackUint32(uint32(buffer[0]), uint32(bufLength))
}
