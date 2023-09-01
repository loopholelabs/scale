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
	"sync"

	interfaces "github.com/loopholelabs/scale-signature-interfaces"
)

type modulePool[T interfaces.Signature] struct {
	pool sync.Pool
	new  func() (*moduleInstance[T], error)
}

func newModulePool[T interfaces.Signature](ctx context.Context, r *Scale[T], f *function[T]) *modulePool[T] {
	return &modulePool[T]{
		new: func() (*moduleInstance[T], error) {
			return newModuleInstance[T](ctx, r, f, nil)
		},
	}
}

func (p *modulePool[T]) Put(module *moduleInstance[T]) {
	if module != nil {
		p.pool.Put(module)
	}
}

func (p *modulePool[T]) Get() (*moduleInstance[T], error) {
	rv, ok := p.pool.Get().(*moduleInstance[T])
	if ok && rv != nil {
		return rv, nil
	}
	return p.new()
}
