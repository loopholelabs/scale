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

package scale

import (
	"context"
	"fmt"

	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
	"github.com/tetratelabs/wazero"
)

// function is the runtime representation of a scale function with
// a compiled module source.
type function[T signature.Signature] struct {
	// runtime is the scale runtime that the function belongs to
	runtime *Scale[T]

	// identifier is the identifier for the function
	identifier string

	// compiled is the compiled module source
	compiled wazero.CompiledModule

	// scaleFunc is the scale function definition schema
	scaleFunc *scalefunc.Schema

	// next is the (optional) next function in the chain
	next *function[T]

	// modulePool is the pool of modules for the function
	modulePool *modulePool[T]

	// env contains the (optional) environment variables for the function
	env map[string]string
}

func (f *function[T]) runWithModule(ctx context.Context, signature T, moduleInstance *moduleInstance[T]) error {
	buf := signature.Write()
	ctxBufferLength := uint64(len(buf))
	writeBuffer, err := moduleInstance.resize.Call(ctx, ctxBufferLength)
	if err != nil {
		return fmt.Errorf("failed to allocate memory for function '%s': %w", f.scaleFunc.Name, err)
	}

	if !moduleInstance.instantiatedModule.Memory().Write(uint32(writeBuffer[0]), buf) {
		return fmt.Errorf("failed to write memory for function '%s'", f.scaleFunc.Name)
	}

	packed, err := moduleInstance.run.Call(ctx)
	if err != nil {
		return fmt.Errorf("failed to run function '%s': %w", f.scaleFunc.Name, err)
	}
	if packed[0] == 0 {
		return fmt.Errorf("failed to run function '%s'", f.scaleFunc.Name)
	}

	offset, length := unpackUint32(packed[0])
	buf, ok := moduleInstance.instantiatedModule.Memory().Read(offset, length)
	if !ok {
		return fmt.Errorf("failed to read memory for function '%s'", f.scaleFunc.Name)
	}

	err = signature.Read(buf)
	if err != nil {
		return fmt.Errorf("error while running function '%s': %w", f.scaleFunc.Name, err)
	}
	return nil
}

func (f *function[T]) run(ctx context.Context, signature T, instance *Instance[T]) error {
	var module *moduleInstance[T]
	var err error
	if f.scaleFunc.Stateless {
		module, err = f.modulePool.Get()
		if err != nil {
			return fmt.Errorf("failed to get module from pool for function %s: %w", f.scaleFunc.Name, err)
		}
	} else {
		module, err = newModuleInstance[T](ctx, f.runtime, f, instance)
		if err != nil {
			return fmt.Errorf("failed to create module for function %s: %w", f.scaleFunc.Name, err)
		}
	}

	module.init(f.runtime, instance)
	module.setSignature(signature)

	err = f.runWithModule(ctx, signature, module)
	if f.scaleFunc.Stateless {
		module.cleanup(f.runtime)
	}
	if err != nil {
		return fmt.Errorf("error while running function '%s': %w", f.scaleFunc.Name, err)
	}

	return nil
}

func newFunction[T signature.Signature](ctx context.Context, r *Scale[T], scaleFunc *scalefunc.Schema, env map[string]string) (*function[T], error) {
	compiled, err := r.runtime.CompileModule(ctx, scaleFunc.Function)
	if err != nil {
		return nil, fmt.Errorf("failed to compile wasm module '%s': %w", scaleFunc.Name, err)
	}

	f := &function[T]{
		runtime:    r,
		identifier: fmt.Sprintf("%s:%s", scaleFunc.Name, scaleFunc.Tag),
		compiled:   compiled,
		scaleFunc:  scaleFunc,
		env:        env,
	}

	f.modulePool = newModulePool[T](ctx, r, f)

	return f, nil
}
