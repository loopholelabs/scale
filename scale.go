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

// Package scale implements the Scale Runtime in Go.
package scale

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"sync"

	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

var (
	NoFunctionsError = errors.New("no available functions")
)

// Next is the next function in the middleware chain. It's meant to be implemented
// by whatever adapter is being used.
type Next[T signature.Signature] func(ctx T) (T, error)

// Scale is the Scale Runtime. It is responsible for initializing
// and managing the WASM runtime as well as the scale function chain.
type Scale[T signature.Signature] struct {
	runtime      wazero.Runtime
	moduleConfig wazero.ModuleConfig

	new signature.New[T]

	functions []*Function[T]
	head      *Function[T]
	tail      *Function[T]

	modulesMu sync.RWMutex
	modules   map[string]*Module[T]

	TraceDataCallback func(data string)
}

func New[T signature.Signature](ctx context.Context, sig signature.New[T], functions []*scalefunc.Schema) (*Scale[T], error) {
	if len(functions) == 0 {
		return nil, NoFunctionsError
	}

	r := &Scale[T]{
		runtime:      wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().WithCloseOnContextDone(true)),
		moduleConfig: wazero.NewModuleConfig().WithSysNanotime().WithSysWalltime().WithRandSource(rand.Reader),
		modules:      make(map[string]*Module[T]),
		new:          sig,
	}

	module := r.runtime.NewHostModuleBuilder("env").
		NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(r.next), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{}).
		WithParameterNames("pointer", "length").Export("next")

	compiled, err := module.Compile(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to compile env: %w", err)
	}

	_, err = r.runtime.InstantiateModule(ctx, compiled, r.moduleConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate host module env: %w", err)
	}

	// Tracing Functions
	builder := r.runtime.NewHostModuleBuilder("scale")
	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(r.getFunctionNameLen), []api.ValueType{}, []api.ValueType{api.ValueTypeI32}).
		WithParameterNames("pointer").Export("get_function_name_len")
	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(r.getFunctionName), []api.ValueType{api.ValueTypeI32}, []api.ValueType{}).
		WithParameterNames("pointer").Export("get_function_name")
	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(r.getInstanceID), []api.ValueType{api.ValueTypeI32}, []api.ValueType{}).
		WithParameterNames("pointer").Export("get_instance_id")
	builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(r.OTELTraceJSON), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{}).
		WithParameterNames("pointer", "length").Export("otel_trace_json")

	_, err = builder.Instantiate(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate host module scale: %w", err)
	}

	_, err = wasi_snapshot_preview1.Instantiate(ctx, r.runtime)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate host module wasi: %w", err)
	}

	for _, f := range functions {
		sf, err := r.CompileFunction(ctx, f)
		if err != nil {
			return nil, fmt.Errorf("failed to pre-compile function '%s': %w", f.Name, err)
		}
		if r.head == nil {
			r.head = sf
		}
		if r.tail != nil {
			r.tail.next = sf
		}
		r.tail = sf
	}

	return r, nil
}

func (r *Scale[T]) CompileFunction(ctx context.Context, scaleFunc *scalefunc.Schema) (*Function[T], error) {
	compiled, err := r.runtime.CompileModule(ctx, scaleFunc.Function)
	if err != nil {
		return nil, fmt.Errorf("failed to compile module '%s': %w", scaleFunc.Name, err)
	}

	f := &Function[T]{
		identifier: fmt.Sprintf("%s:%s", scaleFunc.Name, scaleFunc.Tag),
		scaleFunc:  scaleFunc,
		compiled:   compiled,
	}

	f.modulePool = NewPool[T](ctx, f, r)

	r.functions = append(r.functions, f)
	return f, nil
}
