// Package runtime implements the Scale Runtime in Go.
package runtime

import (
	"context"
	"fmt"
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

	c         *Context
	functions []*Function
	modules   map[api.Module]*Function
}

func New(ctx context.Context, next Next, functions []scalefunc.ScaleFunc) (*Runtime, error) {
	r := &Runtime{
		Next:          next,
		ctx:           ctx,
		runtime:       wazero.NewRuntimeWithConfig(wazero.NewRuntimeConfig().WithWasmCore2()),
		compileConfig: wazero.NewCompileConfig(),
		moduleConfig:  wazero.NewModuleConfig(),
		c:             NewContext(),
		modules:       make(map[api.Module]*Function),
	}

	module := r.runtime.NewModuleBuilder("env")
	module = module.ExportFunction("next", r.next)

	compiled, err := module.Compile(r.ctx, r.compileConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to compile module: %w", err)
	}

	_, err = r.runtime.InstantiateModule(r.ctx, compiled, r.moduleConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate module: %w", err)
	}

	module = r.runtime.NewModuleBuilder("wasi_snapshot_preview1")
	module = module.ExportFunction("fd_write", r.fd_write)

	compiled, err = module.Compile(r.ctx, r.compileConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to compile module: %w", err)
	}

	_, err = r.runtime.InstantiateModule(r.ctx, compiled, r.moduleConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate module: %w", err)
	}

	for _, f := range functions {
		err = r.registerFunction(f)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}
