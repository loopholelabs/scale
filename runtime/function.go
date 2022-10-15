package runtime

import (
	"fmt"
	"github.com/loopholelabs/scale-go/scalefunc"
	"github.com/loopholelabs/scale-go/utils"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type Function struct {
	ScaleFunc scalefunc.ScaleFunc
	Next      *Function
	Compiled  wazero.CompiledModule
	Module    api.Module
}

func (f *Function) Run(r *Runtime) error {
	module := f.Module
	run := module.ExportedFunction("run")
	malloc := module.ExportedFunction("malloc")
	free := module.ExportedFunction("free")

	bufLength := uint64(r.c.Buffer.Len())
	buffer, err := malloc.Call(r.ctx, bufLength)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = free.Call(r.ctx, buffer[0])
	}()

	if !module.Memory().Write(r.ctx, uint32(buffer[0]), r.c.Buffer.Bytes()) {
		return MemoryWriteError
	}

	packed, err := run.Call(r.ctx, buffer[0], bufLength)
	if err != nil {
		return err
	}

	offset, length := utils.UnpackUint32(packed[0])
	buf, ok := module.Memory().Read(r.ctx, offset, length)
	if !ok {
		panic("failed to read memory")
	}

	return r.c.Deserialize(buf)
}

func (r *Runtime) registerFunction(scaleFunc scalefunc.ScaleFunc) error {
	compiled, err := r.runtime.CompileModule(r.ctx, scaleFunc.Function, r.compileConfig)
	if err != nil {
		return fmt.Errorf("failed to compile function '%s': %w", scaleFunc.ScaleFile.Name, err)
	}
	module, err := r.runtime.InstantiateModule(r.ctx, compiled, r.moduleConfig.WithName(scaleFunc.ScaleFile.Name))
	if err != nil {
		return fmt.Errorf("failed to instantiate function '%s': %w", scaleFunc.ScaleFile.Name, err)
	}

	f := &Function{
		ScaleFunc: scaleFunc,
		Compiled:  compiled,
		Module:    module,
	}

	r.functions = append(r.functions, f)
	r.modules[module] = f

	if len(r.functions) > 1 {
		r.functions[len(r.functions)-2].Next = f
	}

	return nil
}
