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
	"crypto/rand"
	"errors"
	"fmt"
	"sync"

	signature "github.com/loopholelabs/scale-signature"
	httpSignature "github.com/loopholelabs/scale-signature-http"
	"github.com/loopholelabs/scalefile/scalefunc"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

var (
	NoFunctionsError = errors.New("no functions found in runtime")
)

// Next is the next function in the middleware chain. It's meant to be implemented
// by whatever adapter is being used.
type Next[T signature.Signature] func(ctx T) (T, error)

// Runtime is the Scale Runtime. It is responsible for initializing
// and managing the WASM runtime as well as the scale function chain.
type Runtime[T signature.Signature] struct {
	runtime      wazero.Runtime
	moduleConfig wazero.ModuleConfig

	InvocationID      []byte
	ServiceName       string
	TraceDataCallback func(data string)

	new signature.NewSignature[T]

	functions []*Function[T]
	head      *Function[T]
	tail      *Function[T]

	modulesMu sync.RWMutex
	modules   map[string]*Module[T]
}

func New(ctx context.Context, functions []*scalefunc.ScaleFunc) (*Runtime[*httpSignature.Context], error) {
	return NewWithSignature(ctx, httpSignature.New, functions)
}

func NewWithSignature[T signature.Signature](ctx context.Context, sig signature.NewSignature[T], functions []*scalefunc.ScaleFunc) (*Runtime[T], error) {
	if len(functions) == 0 {
		return nil, NoFunctionsError
	}

	r := &Runtime[T]{
		runtime:      wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().WithCloseOnContextDone(true)),
		moduleConfig: wazero.NewModuleConfig().WithSysNanotime().WithSysWalltime().WithRandSource(rand.Reader),
		modules:      make(map[string]*Module[T]),
		new:          sig,
		InvocationID: make([]byte, 16),
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
		return nil, fmt.Errorf("failed to instantiate env: %w", err)
	}

	// Add scale trace functions
	scale_builder := r.runtime.NewHostModuleBuilder("scale")

	scale_builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(r.getServiceNameLen), []api.ValueType{}, []api.ValueType{api.ValueTypeI32}).
		WithParameterNames("pointer").Export("get_service_name_len")

	scale_builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(r.getServiceName), []api.ValueType{api.ValueTypeI32}, []api.ValueType{}).
		WithParameterNames("pointer").Export("get_service_name")

	scale_builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(r.getInvocationID), []api.ValueType{api.ValueTypeI32}, []api.ValueType{}).
		WithParameterNames("pointer").Export("get_invocation_id")

	scale_builder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(r.sendOtelTraceJson), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{}).
		WithParameterNames("pointer", "length").Export("send_otel_trace_json")

	_, err = scale_builder.Instantiate(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate scale: %w", err)
	}

	_, err = wasi_snapshot_preview1.Instantiate(ctx, r.runtime)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate wasi: %w", err)
	}

	for _, f := range functions {
		sf, err := r.compileFunction(ctx, f)
		if err != nil {
			return nil, fmt.Errorf("failed to compile function '%s': %w", f.Name, err)
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

func (r *Runtime[T]) compileFunction(ctx context.Context, scaleFunc *scalefunc.ScaleFunc) (*Function[T], error) {
	compiled, err := r.runtime.CompileModule(ctx, scaleFunc.Function)
	if err != nil {
		return nil, fmt.Errorf("failed to compile function '%s': %w", scaleFunc.Name, err)
	}

	f := &Function[T]{
		scaleFunc: scaleFunc,
		compiled:  compiled,
	}

	f.modulePool = NewPool[T](ctx, f, r)

	r.functions = append(r.functions, f)
	return f, nil
}

// Host provided function get service name length
func (r *Runtime[T]) getServiceNameLen(ctx context.Context, mod api.Module, params []uint64) {
	params[0] = uint64(len([]byte(r.ServiceName)))
}

// Host provided function get service name
func (r *Runtime[T]) getServiceName(ctx context.Context, mod api.Module, params []uint64) {
	ptr := uint32(params[0])
	mem := mod.Memory()
	mem.Write(ptr, []byte(r.ServiceName))
}

// Host provided function to get 16 byte Invocation ID.
func (r *Runtime[T]) getInvocationID(ctx context.Context, mod api.Module, params []uint64) {
	ptr := uint32(params[0])
	mem := mod.Memory()
	mem.Write(ptr, r.InvocationID)
}

// Host provided function to receive otel trace data.
func (r *Runtime[T]) sendOtelTraceJson(ctx context.Context, mod api.Module, params []uint64) {
	if r.TraceDataCallback == nil {
		return // Drop it
	}

	ptr := uint32(params[0])
	length := uint32(params[1])
	mem := mod.Memory()
	data, ok := mem.Read(ptr, length)

	if ok {
		// Make a copy of the data
		copied_data := make([]byte, len(data))
		copy(copied_data, data)

		r.TraceDataCallback(string(copied_data))
	}
}
