package runtime

import (
	"context"
	"fmt"
	"github.com/loopholelabs/scale-go/utils"
	"github.com/tetratelabs/wazero/api"
)

type Module struct {
	module   api.Module
	instance *Instance
	function *Function
	next     *Module
	run      api.Function
	malloc   api.Function
	free     api.Function
}

func (m *Module) Run(ctx context.Context) error {
	bufLength := uint64(m.instance.ctx.Buffer.Len())
	buffer, err := m.malloc.Call(ctx, bufLength)
	if err != nil {
		return fmt.Errorf("failed to allocate memory for function '%s': %w", m.function.ScaleFunc.ScaleFile.Name, err)
	}

	if !m.module.Memory().Write(ctx, uint32(buffer[0]), m.instance.ctx.Buffer.Bytes()) {
		return fmt.Errorf("failed to write memory for function '%s'", m.function.ScaleFunc.ScaleFile.Name)
	}

	packed, err := m.run.Call(ctx, buffer[0], bufLength)
	if err != nil {
		return fmt.Errorf("failed to run function '%s': %w", m.function.ScaleFunc.ScaleFile.Name, err)
	}

	offset, length := utils.UnpackUint32(packed[0])
	buf, ok := m.module.Memory().Read(ctx, offset, length)
	if !ok {
		return fmt.Errorf("failed to read memory for function '%s'", m.function.ScaleFunc.ScaleFile.Name)
	}

	return m.instance.ctx.Deserialize(buf)
}

type Instance struct {
	id      string
	runtime *Runtime
	ctx     *Context
	head    *Module
	tail    *Module
	modules map[api.Module]*Module
}

func (i *Instance) Context() *Context {
	return i.ctx
}

func (i *Instance) Run(ctx context.Context) error {
	if i.head == nil {
		return fmt.Errorf("no functions registered for instance %s", i.id)
	}
	module := i.head
	return module.Run(ctx)
}

func (i *Instance) initialize(ctx context.Context) error {
	for _, f := range i.runtime.functions {
		module, err := i.runtime.runtime.InstantiateModule(ctx, f.Compiled, i.runtime.moduleConfig.WithName(i.id))
		if err != nil {
			return fmt.Errorf("failed to instantiate function '%s' for instance %s: %w", f.ScaleFunc.ScaleFile.Name, i.id, err)
		}

		run := module.ExportedFunction("run")
		malloc := module.ExportedFunction("malloc")
		free := module.ExportedFunction("free")

		m := &Module{
			module:   module,
			function: f,
			instance: i,
			run:      run,
			malloc:   malloc,
			free:     free,
		}
		i.modules[module] = m
		if i.head == nil {
			i.head = m
		}
		if i.tail != nil {
			i.tail.next = m
		}
		i.tail = m
	}
	return nil
}
