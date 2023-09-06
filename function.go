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
)

// function is an instantiated function that can be run
type function[T interfaces.Signature] struct {
	// instance is the instance that this function belongs to
	instance *Instance[T]

	// template is the template that the function
	// was created from
	template *template[T]

	// next is the next function in the chain
	next *function[T]

	// module is the optional, stateful, instantiated module for this function
	//
	// If the function is stateless, then this will be nil
	module *module[T]
}

// newFunction creates a new function from an Instance and a template
func newFunction[T interfaces.Signature](ctx context.Context, instance *Instance[T], template *template[T]) (fn *function[T], err error) {
	fn = &function[T]{
		instance: instance,
		template: template,
	}

	if template.modulePool == nil {
		fn.module, err = newModule[T](ctx, fn.template)
		if err != nil {
			return nil, fmt.Errorf("failed to create module for function '%s': %w", fn.template.identifier, err)
		}
		fn.module.register(fn)
	}

	return fn, nil
}

func (f *function[T]) getModule(signature T) (*module[T], error) {
	if f.module != nil {
		f.module.setSignature(signature)
		return f.module, nil
	}
	if f.template.modulePool == nil {
		return nil, fmt.Errorf("cannot get module from pool for function %s: module pool is nil", f.template.identifier)
	}
	m, err := f.template.modulePool.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get module from pool for function %s: %w", f.template.identifier, err)
	}
	m.register(f)
	m.setSignature(signature)
	return m, nil
}

func (f *function[T]) putModule(m *module[T]) {
	if f.template.modulePool != nil {
		m.cleanup()
		f.template.modulePool.Put(m)
	}
}
