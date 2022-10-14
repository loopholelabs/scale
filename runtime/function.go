package runtime

import (
	"github.com/loopholelabs/scale-go/scalefunc"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type Function struct {
	ScaleFunc scalefunc.ScaleFunc
	Parent    api.Module
	Compiled  wazero.CompiledModule
	Module    api.Module
}

func (r *Runtime) registerFunction(scaleFunc scalefunc.ScaleFunc, parent api.Module) (*Function, error) {
	compiled, err := r.runtime.CompileModule(r.ctx, scaleFunc.Function, r.compileConfig)
	if err != nil {
		return nil, err
	}
	module, err := r.runtime.InstantiateModule(r.ctx, compiled, r.moduleConfig.WithName(scaleFunc.ScaleFile.Name))
	if err != nil {
		return nil, err
	}

	f := &Function{
		ScaleFunc: scaleFunc,
		Parent:    parent,
		Compiled:  compiled,
		Module:    module,
	}

	r.functions = append(r.functions, f)

	return f, nil
}
