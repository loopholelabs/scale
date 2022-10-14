// Package runtime implements the Scale Runtime in Go.
package runtime

import (
	"context"
	"github.com/loopholelabs/scale-go/scalefunc"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type Runtime struct {
	Next          Next
	ctx           context.Context
	runtime       wazero.Runtime
	compileConfig wazero.CompileConfig
	moduleConfig  wazero.ModuleConfig

	functions []*Function
}

func New(ctx context.Context, next Next, functions []scalefunc.ScaleFunc) (*Runtime, error) {
	r := &Runtime{
		Next:          next,
		ctx:           ctx,
		runtime:       wazero.NewRuntimeWithConfig(wazero.NewRuntimeConfig().WithWasmCore2()),
		compileConfig: wazero.NewCompileConfig(),
		moduleConfig:  wazero.NewModuleConfig(),
	}

	module := r.runtime.NewModuleBuilder("env")
	//module = module.ExportFunction("hostReceiveRequestForNext", receiveRequestForNext)
	//module = module.ExportFunction("hostReceiveRequest", receiveRequest)
	//module = module.ExportFunction("__Next", __Next)
	//module = module.ExportFunction("debugWasm", debugWasm)

	compiled, err := module.Compile(r.ctx, r.compileConfig)
	if err != nil {
		return nil, err
	}

	_, err = r.runtime.InstantiateModule(r.ctx, compiled, r.moduleConfig)
	if err != nil {
		return nil, err
	}

	var parent api.Module
	var function *Function
	for _, f := range functions {
		function, err = r.registerFunction(f, parent)
		if err != nil {
			return nil, err
		}
		parent = function.Parent
	}

	return r, nil
}
