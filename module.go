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

	"github.com/google/uuid"
	"github.com/tetratelabs/wazero/api"

	"github.com/loopholelabs/scale/signature"
)

type moduleInstance[T signature.Signature] struct {
	// function is the function that this moduleInstance is associated with
	function *function[T]

	// instantiatedModule is the instantiated wasm module
	instantiatedModule api.Module

	// run is the exported `run` function
	run api.Function

	// resize is the exported `resize` function
	resize api.Function

	// instance is set during the initialization of the moduleInstance
	instance *Instance[T]

	// signature is set during the initialization of the moduleInstance
	signature T

	// nextInstance is the nextInstance moduleInstance in the chain
	//
	// This is only set for persistent instances, otherwise it's nil
	nextInstance *moduleInstance[T]
}

// newModuleInstance creates a new moduleInstance
//
// If the instance parameter is given then the entire moduleInstance chain is initialized, with `module.init`
// being called, and the moduleInstance.nextInstance field is set to the nextInstance moduleInstance in the chain.
// In this case, `module.cleanup` must be called when the moduleInstance is no longer needed.
func newModuleInstance[T signature.Signature](ctx context.Context, r *Scale[T], f *function[T], i *Instance[T]) (*moduleInstance[T], error) {
	config := r.moduleConfig.WithName(fmt.Sprintf("%s.%s", f.identifier, uuid.New().String()))
	for k, v := range f.env {
		config = config.WithEnv(k, v)
	}

	instantiatedModule, err := r.runtime.InstantiateModule(ctx, f.compiled, config)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate module '%s': %w", f.identifier, err)
	}

	run := instantiatedModule.ExportedFunction("run")
	resize := instantiatedModule.ExportedFunction("resize")
	if run == nil || resize == nil {
		return nil, fmt.Errorf("failed to find run or resize implementations for function %s", f.identifier)
	}

	m := &moduleInstance[T]{
		function:           f,
		instantiatedModule: instantiatedModule,
		run:                run,
		resize:             resize,
	}

	if i != nil {
		m.init(r, i)
		if f.next != nil {
			m.nextInstance, err = newModuleInstance[T](ctx, r, f.next, i)
			if err != nil {
				return nil, err
			}
		}

	}

	return m, nil
}

// init initializes the moduleInstance and registers it with the Scale runtime
func (m *moduleInstance[T]) init(r *Scale[T], i *Instance[T]) {
	m.instance = i
	r.activeModulesMu.Lock()
	r.activeModules[m.instantiatedModule.Name()] = m
	r.activeModulesMu.Unlock()
}

// cleanup removes the moduleInstance from the Scale runtime
//
// if it has a nextInstance moduleInstance, then it is also cleaned up
func (m *moduleInstance[T]) cleanup(r *Scale[T]) {
	r.activeModulesMu.Lock()
	delete(r.activeModules, m.instantiatedModule.Name())
	r.activeModulesMu.Unlock()
	if m.nextInstance != nil {
		m.nextInstance.cleanup(r)
	}
}

// setSignature sets the signature for the moduleInstance
func (m *moduleInstance[T]) setSignature(signature T) {
	m.signature = signature
}
