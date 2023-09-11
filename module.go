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

	interfaces "github.com/loopholelabs/scale-signature-interfaces"

	"github.com/loopholelabs/polyglot"

	"github.com/google/uuid"
	"github.com/tetratelabs/wazero/api"
)

type module[T interfaces.Signature] struct {
	// template is the template that the module
	// was created from
	template *template[T]

	// instantiatedModule is the instantiated wasm module
	instantiatedModule api.Module

	// runFunction is the exported `run` function
	runFunction api.Function

	// resizeFunction is the exported `resize` function
	resizeFunction api.Function

	// function is set during the initialization of the module
	function *function[T]

	// signature is set during the initialization of the module
	signature T
}

// newModule creates a new module
func newModule[T interfaces.Signature](ctx context.Context, template *template[T]) (*module[T], error) {
	config := template.runtime.moduleConfig.WithName(fmt.Sprintf("%s.%s", template.identifier, uuid.New().String()))
	for k, v := range template.env {
		config = config.WithEnv(k, v)
	}

	instantiatedModule, err := template.runtime.runtime.InstantiateModule(ctx, template.compiled, config)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate module '%s': %w", template.identifier, err)
	}

	run := instantiatedModule.ExportedFunction("run")
	resize := instantiatedModule.ExportedFunction("resize")
	initialize := instantiatedModule.ExportedFunction("initialize")
	if run == nil || resize == nil || initialize == nil {
		return nil, fmt.Errorf("failed to find run, resize, or initialize implementations for function %s", template.identifier)
	}

	packed, err := initialize.Call(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to run initialize function for '%s': %w", template.identifier, err)
	}

	if packed[0] != 0 {
		offset, length := unpackUint32(packed[0])
		buf, ok := instantiatedModule.Memory().Read(offset, length)
		if !ok {
			return nil, fmt.Errorf("failed to read memory for function '%s'", template.identifier)
		}

		dec := polyglot.GetDecoder(buf)
		valErr, err := dec.Error()
		polyglot.ReturnDecoder(dec)
		if err != nil {
			return nil, fmt.Errorf("failed to decode error for function '%s': %w", template.identifier, err)
		}
		return nil, fmt.Errorf("failed to initialize function '%s': %s", template.identifier, valErr)
	}

	return &module[T]{
		template:           template,
		instantiatedModule: instantiatedModule,
		runFunction:        run,
		resizeFunction:     resize,
	}, nil
}

// run runs the module
//
// The signature of the module must be set before calling this function
func (m *module[T]) run(ctx context.Context) error {
	buf := m.signature.Write()
	writeBuffer, err := m.resizeFunction.Call(ctx, uint64(len(buf)))
	if err != nil {
		return fmt.Errorf("failed to allocate memory for function '%s': %w", m.template.identifier, err)
	}

	if !m.instantiatedModule.Memory().Write(uint32(writeBuffer[0]), buf) {
		return fmt.Errorf("failed to write memory for function '%s'", m.template.identifier)
	}

	packed, err := m.runFunction.Call(ctx)
	if err != nil {
		return fmt.Errorf("failed to run function '%s': %w", m.template.identifier, err)
	}
	if packed[0] == 0 {
		return fmt.Errorf("failed to run function '%s'", m.template.identifier)
	}

	ptr, length := unpackUint32(packed[0])
	buf, ok := m.instantiatedModule.Memory().Read(ptr, length)
	if !ok {
		return fmt.Errorf("failed to read memory for function '%s'", m.template.identifier)
	}

	err = m.signature.Read(buf)
	if err != nil {
		return fmt.Errorf("error while running function '%s': %w", m.template.identifier, err)
	}
	return nil
}

// register sets the module's instance field and registers it as an active module with the runtime
func (m *module[T]) register(function *function[T]) {
	m.function = function
	m.template.runtime.activeModulesMu.Lock()
	m.template.runtime.activeModules[m.instantiatedModule.Name()] = m
	m.template.runtime.activeModulesMu.Unlock()
}

// cleanup removes the module from the runtime's active modules map
func (m *module[T]) cleanup() {
	m.function = nil
	m.template.runtime.activeModulesMu.Lock()
	delete(m.template.runtime.activeModules, m.instantiatedModule.Name())
	m.template.runtime.activeModulesMu.Unlock()
}

// setSignature sets the module's signature
func (m *module[T]) setSignature(signature T) {
	m.signature = signature
}
