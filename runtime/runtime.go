/*
	Copyright 2022 Loophole Labs

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		   http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

// Package runtime implements the Scale Runtime in Go.
package runtime

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/loopholelabs/scale-go/scalefunc"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// Next is the next function in the middleware chain. It's meant to be implemented
// by whatever adapter is being used.
type Next func(ctx *Context) *Context

// Runtime is the Scale Runtime. It is responsible for initializing
// and managing the WASM runtime as well as the scale function chain.
type Runtime struct {
	runtime       wazero.Runtime
	compileConfig wazero.CompileConfig
	moduleConfig  wazero.ModuleConfig

	functions []*Function
	instances map[string]*Instance
}

func New(ctx context.Context, functions []scalefunc.ScaleFunc) (*Runtime, error) {
	r := &Runtime{
		runtime:       wazero.NewRuntimeWithConfig(wazero.NewRuntimeConfig().WithWasmCore2()),
		compileConfig: wazero.NewCompileConfig(),
		moduleConfig:  wazero.NewModuleConfig(),
		instances:     make(map[string]*Instance),
	}

	module := r.runtime.NewModuleBuilder("env")
	module = module.ExportFunction("next", r.next)
	//module = module.ExportFunction("debug", r.debug)

	compiled, err := module.Compile(ctx, r.compileConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to compile module: %w", err)
	}

	_, err = r.runtime.InstantiateModule(ctx, compiled, r.moduleConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate module: %w", err)
	}

	module = r.runtime.NewModuleBuilder("wasi_snapshot_preview1")
	module = module.ExportFunction("fd_write", r.fd_write)

	compiled, err = module.Compile(ctx, r.compileConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to compile module: %w", err)
	}

	_, err = r.runtime.InstantiateModule(ctx, compiled, r.moduleConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate module: %w", err)
	}

	for _, f := range functions {
		err = r.registerFunction(ctx, f)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r *Runtime) Instance(ctx context.Context, next Next) (*Instance, error) {
	i := &Instance{
		id:      uuid.New().String(),
		next:    next,
		runtime: r,
		modules: make(map[api.Module]*Module),
		ctx:     NewContext(),
	}

	r.instances[i.id] = i

	return i, i.initialize(ctx)
}
