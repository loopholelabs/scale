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
	"sync"

	interfaces "github.com/loopholelabs/scale-signature-interfaces"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"

	"github.com/loopholelabs/scale/extension"
)

// Next is the next function in the middleware chain. It's meant to be implemented
// by whatever adapter is being used.
type Next[T interfaces.Signature] func(ctx T) (T, error)

// Scale is the Scale Runtime. It is responsible for initializing
// and managing the WASM runtime as well as the scale function chain.
type Scale[T interfaces.Signature] struct {
	runtime      wazero.Runtime
	moduleConfig wazero.ModuleConfig

	config *Config[T]

	head *template[T]
	tail *template[T]

	activeModulesMu sync.RWMutex
	activeModules   map[string]*module[T]

	TraceDataCallback func(data string)
}

func New[T interfaces.Signature](config *Config[T]) (*Scale[T], error) {
	r := &Scale[T]{
		runtime:       wazero.NewRuntimeWithConfig(config.context, wazero.NewRuntimeConfig().WithCloseOnContextDone(true)),
		moduleConfig:  wazero.NewModuleConfig().WithSysNanotime().WithSysWalltime().WithRandSource(rand.Reader),
		activeModules: make(map[string]*module[T]),
		config:        config,
	}

	return r, r.init()
}

// Instance returns a new instance of a Scale Function chain
// with the provided and optional next function.
func (r *Scale[T]) Instance(next ...Next[T]) (*Instance[T], error) {
	return newInstance(r.config.context, r, next...)
}

func (r *Scale[T]) init() error {
	err := r.config.validate()
	if err != nil {
		return err
	}

	if r.config.stdout != nil && r.config.rawOutput {
		r.moduleConfig = r.moduleConfig.WithStdout(r.config.stdout)
	}

	if r.config.stderr != nil && r.config.rawOutput {
		r.moduleConfig = r.moduleConfig.WithStderr(r.config.stderr)
	}

	envModule := r.runtime.NewHostModuleBuilder("env")

	// Install any extensions...
	for name, fn := range r.config.extensions {
		fmt.Printf("Installing module [%s]\n", name)
		wfn := func(n string, f extension.InstallableFunc) func(context.Context, api.Module, []uint64) {
			return func(ctx context.Context, mod api.Module, params []uint64) {
				fmt.Printf("HOST FUNCTION CALLED %s\n", n)
				mem := mod.Memory()
				resize := func(name string, size uint64) (uint64, error) {
					w, err := mod.ExportedFunction(name).Call(context.Background(), size)
					return w[0], err
				}
				f(mem, resize, params)
			}
		}(name, fn)

		envModule.NewFunctionBuilder().
			WithGoModuleFunction(api.GoModuleFunc(wfn), []api.ValueType{api.ValueTypeI64, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI64}).
			WithParameterNames("instance", "pointer", "length").Export(name)
	}

	envHostModuleBuilder := envModule.
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

	testSignature := r.config.newSignature()
	for _, sf := range r.config.functions {
		if testSignature.Hash() != "" && testSignature.Hash() != sf.function.SignatureHash {
			return fmt.Errorf("passed in function '%s:%s' has an invalid signatures", sf.function.Name, sf.function.Tag)
		}

		t, err := newTemplate(r.config.context, r, sf.function, sf.env)
		if err != nil {
			return fmt.Errorf("failed to pre-compile function '%s:%s': %w", sf.function.Name, sf.function.Tag, err)
		}

		if r.head == nil {
			r.head = t
		}

		if r.tail != nil {
			r.tail.next = t
		}

		r.tail = t
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

	if m.function.next != nil {
		nextModule, err := m.function.next.getModule(m.signature)
		if err == nil {
			err = nextModule.run(ctx)
			if err == nil {
				m.function.next.putModule(nextModule)
			}
		}
	} else {
		m.signature, err = m.function.instance.next(m.signature)
	}
	if err != nil {
		buf = m.signature.Error(err)
	} else {
		buf = m.signature.Write()
	}

	writeBuffer, err := m.resizeFunction.Call(ctx, uint64(len(buf)))
	if err != nil {
		return
	}
	m.instantiatedModule.Memory().Write(uint32(writeBuffer[0]), buf)
}
