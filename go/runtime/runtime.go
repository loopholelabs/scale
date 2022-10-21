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
	"errors"
	"fmt"
	"github.com/loopholelabs/scale/go/scalefunc"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"sync"
)

var (
	NextFunctionRequiredError = errors.New("next function required when the scale function chain only contains middleware")
)

// Next is the next function in the middleware chain. It's meant to be implemented
// by whatever adapter is being used.
type Next func(ctx *Context) *Context

// Runtime is the Scale Runtime. It is responsible for initializing
// and managing the WASM runtime as well as the scale function chain.
type Runtime struct {
	runtime      wazero.Runtime
	moduleConfig wazero.ModuleConfig

	functions  []*Function
	instanceMu sync.RWMutex
	instances  map[string]*Instance
}

func New(ctx context.Context, functions []scalefunc.ScaleFunc) (*Runtime, error) {
	r := &Runtime{
		runtime:      wazero.NewRuntime(ctx),
		moduleConfig: wazero.NewModuleConfig(),
		instances:    make(map[string]*Instance),
	}

	module := r.runtime.NewHostModuleBuilder("env")
	module = module.ExportFunction("next", r.next)

	compiled, err := module.Compile(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to compile env: %w", err)
	}

	_, err = r.runtime.InstantiateModule(ctx, compiled, r.moduleConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate env: %w", err)
	}

	_, err = wasi_snapshot_preview1.Instantiate(ctx, r.runtime)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate wasi: %w", err)
	}

	for _, f := range functions {
		err = r.registerFunction(ctx, f)
		if err != nil {
			return nil, fmt.Errorf("failed to register function '%s': %w", f.ScaleFile.Name, err)
		}
	}

	return r, nil
}
