package runtime

import (
	"context"
	"github.com/loopholelabs/scale-go/utils"
	"github.com/tetratelabs/wazero/api"
)

type Next func(ctx *Context) *Context

func (r *Runtime) next(ctx context.Context, module api.Module, offset uint32, length uint32) uint64 {
	i := r.instances[module.Name()]
	if i == nil {
		// TODO: have an error field in the context
		return 0
	}

	m := i.modules[module]
	if m == nil {
		// TODO: have an error field in the context
		return 0
	}

	buf, ok := m.module.Memory().Read(ctx, offset, length)
	if !ok {
		// TODO: have an error field in the context
		return 0
	}

	err := i.Context().Deserialize(buf)
	if err != nil {
		// TODO: have an error field in the context
		return 0
	}

	if m.next == nil {
		i.ctx = r.Next(i.Context())
	} else {
		err = m.next.Run(ctx)
		if err != nil {
			// TODO: have an error field in the context
			return 0
		}
	}

	i.Context().Serialize()
	bufLength := uint64(i.Context().Buffer.Len())
	buffer, err := m.malloc.Call(ctx, bufLength)
	if err != nil {
		// TODO: have an error field in the context
		return 0
	}

	if !module.Memory().Write(ctx, uint32(buffer[0]), i.Context().Buffer.Bytes()) {
		// TODO: have an error field in the context
		return 0
	}

	return utils.PackUint32(uint32(buffer[0]), uint32(bufLength))
}
