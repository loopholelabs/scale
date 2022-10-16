package runtime

import (
	"context"
	"github.com/loopholelabs/scale-go/utils"
	"github.com/tetratelabs/wazero/api"
)

func (r *Runtime) fd_write(int32, int32, int32, int32) int32 {
	return 0
}

func (r *Runtime) next(ctx context.Context, module api.Module, offset uint32, length uint32) uint64 {
	i := r.instances[module.Name()]
	if i == nil {
		return 0
	}

	m := i.modules[module]
	if m == nil {
		return 0
	}

	buf, ok := m.module.Memory().Read(ctx, offset, length)
	if !ok {
		return 0
	}

	err := i.Context().Read(buf)
	if err != nil {
		// TODO: have an error field in the context
		return 0
	}

	if m.next == nil {
		i.ctx = i.next(i.Context())
	} else {
		err = m.next.Run(ctx)
		if err != nil {
			i.Context().Error(err)
		}
	}

	i.Context().Write()
	bufLength := uint64(i.Context().Buffer.Len())
	buffer, err := m.malloc.Call(ctx, bufLength)
	if err != nil {
		return 0
	}

	if !module.Memory().Write(ctx, uint32(buffer[0]), i.Context().Buffer.Bytes()) {
		return 0
	}

	return utils.PackUint32(uint32(buffer[0]), uint32(bufLength))
}

//func (r *Runtime) debug(ctx context.Context, module api.Module, offset uint32, length uint32) {
//	buf, ok := module.Memory().Deserialize(ctx, offset, length)
//	if !ok {
//		panic("failed to read memory")
//	}
//
//	fmt.Println(string(buf))
//}
