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
	"fmt"
	"strings"
	"sync"

	"github.com/loopholelabs/scale/signature"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// Next is the next function in the middleware chain. It's meant to be implemented
// by whatever adapter is being used.
type Next[T signature.Signature] func(ctx T) (T, error)

// Scale is the Scale Runtime. It is responsible for initializing
// and managing the WASM runtime as well as the scale function chain.
type Scale[T signature.Signature] struct {
	runtime      wazero.Runtime
	moduleConfig wazero.ModuleConfig

	config *Config[T]

	head *function[T]
	tail *function[T]

	activeModulesMu sync.RWMutex
	activeModules   map[string]*moduleInstance[T]

	TraceDataCallback func(data string)
}

func New[T signature.Signature](config *Config[T]) (*Scale[T], error) {
	r := &Scale[T]{
		runtime:       wazero.NewRuntimeWithConfig(config.context, wazero.NewRuntimeConfig().WithCloseOnContextDone(true)),
		moduleConfig:  wazero.NewModuleConfig().WithSysNanotime().WithSysWalltime().WithRandSource(rand.Reader),
		activeModules: make(map[string]*moduleInstance[T]),
		config:        config,
	}

	return r, r.init()
}

// Instance returns a new instance of a Scale Function chain
// with the provided and optional next function.
func (r *Scale[T]) Instance(next ...Next[T]) (*Instance[T], error) {
	if r.config.pooling {
		return newPoolingInstance(r, next...)
	}

	return newInstance(r.config.context, r, next...)
}

func (r *Scale[T]) init() error {
	err := r.config.validate()
	if err != nil {
		return err
	}

	envHostModuleBuilder := r.runtime.NewHostModuleBuilder("env").
		NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(r.next), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{}).
		WithParameterNames("pointer", "length").Export("next")

	_, err = envHostModuleBuilder.Instantiate(r.config.context)
	if err != nil {
		return fmt.Errorf("failed to instantiate env host module: %w", err)
	}

	// Tracing Functions
	scaleHostModuleBuilder := r.runtime.NewHostModuleBuilder("scale")
	scaleHostModuleBuilder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(r.getFunctionNameLen), []api.ValueType{}, []api.ValueType{api.ValueTypeI32}).
		WithParameterNames("pointer").Export("get_function_name_len")
	scaleHostModuleBuilder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(r.getFunctionName), []api.ValueType{api.ValueTypeI32}, []api.ValueType{}).
		WithParameterNames("pointer").Export("get_function_name")
	scaleHostModuleBuilder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(r.getInstanceID), []api.ValueType{api.ValueTypeI32}, []api.ValueType{}).
		WithParameterNames("pointer").Export("get_instance_id")
	scaleHostModuleBuilder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(r.otelTraceJSON), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{}).
		WithParameterNames("pointer", "length").Export("otel_trace_json")

	_, err = scaleHostModuleBuilder.Instantiate(r.config.context)
	if err != nil {
		return fmt.Errorf("failed to instantiate scale host module: %w", err)
	}

	_, err = wasi_snapshot_preview1.Instantiate(r.config.context, r.runtime)
	if err != nil {
		return fmt.Errorf("failed to instantiate host module wasi: %w", err)
	}

	for _, sf := range r.config.functions {
		f, err := newFunction(r.config.context, r, sf.function, sf.env)
		if err != nil {
			return fmt.Errorf("failed to pre-compile function '%s:%s': %w", sf.function.Name, sf.function.Tag, err)
		}

		if r.head == nil {
			r.head = f
		}

		if r.head.scaleFunc.SignatureHash != sf.function.SignatureHash {
			return fmt.Errorf("function '%s:%s' and '%s:%s' have mismatching signatures", r.head.scaleFunc.Name, r.head.scaleFunc.Tag, sf.function.Name, sf.function.Tag)
		}

		if r.tail != nil {
			r.tail.next = f
		}

		r.tail = f
	}

	return nil
}

func (r *Scale[T]) next(ctx context.Context, module api.Module, params []uint64) {
	r.activeModulesMu.RLock()
	m := r.activeModules[module.Name()]
	r.activeModulesMu.RUnlock()
	if m == nil {
		return
	}

	pointer := uint32(params[0])
	length := uint32(params[1])
	buf, ok := m.instantiatedModule.Memory().Read(pointer, length)
	if !ok {
		return
	}

	err := m.signature.Read(buf)
	if err != nil {
		return
	}

	if m.nextInstance != nil {
		m.nextInstance.setSignature(m.signature)
		err = m.nextInstance.function.runWithModule(ctx, m.nextInstance)
	} else if m.function.next == nil {
		m.signature, err = m.instance.next(m.signature)
	} else {
		err = m.function.next.run(ctx, m.signature, m.instance)
	}
	if err != nil {
		buf = m.signature.Error(err)
	} else {
		buf = m.signature.Write()
	}

	writeBuffer, err := m.resize.Call(ctx, uint64(len(buf)))
	if err != nil {
		return
	}
	module.Memory().Write(uint32(writeBuffer[0]), buf)
}

type Parsed struct {
	Organization string
	Name         string
	Tag          string
}

// Parse parses a function or signature name of the form <org>/<name>:<tag> into its organization, name, and tag
func Parse(name string) *Parsed {
	orgSplit := strings.Split(name, "/")
	if len(orgSplit) == 1 {
		orgSplit = []string{"", name}
	}
	tagSplit := strings.Split(orgSplit[1], ":")
	if len(tagSplit) == 1 {
		tagSplit = []string{tagSplit[0], ""}
	}
	return &Parsed{
		Organization: orgSplit[0],
		Name:         tagSplit[0],
		Tag:          tagSplit[1],
	}
}
