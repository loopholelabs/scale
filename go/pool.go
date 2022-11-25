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

package runtime

import (
	"context"
	signature "github.com/loopholelabs/scale-signature"
	"sync"
)

type Pool[T signature.Signature] struct {
	pool sync.Pool
	new  func() (*Module[T], error)
}

func NewPool[T signature.Signature](ctx context.Context, f *Function[T], r *Runtime[T]) *Pool[T] {
	return &Pool[T]{
		new: func() (*Module[T], error) {
			return NewModule[T](ctx, f, r)
		},
	}
}

func (p *Pool[T]) Put(module *Module[T]) {
	if module != nil {
		p.pool.Put(module)
	}
}

func (p *Pool[T]) Get() (*Module[T], error) {
	rv, ok := p.pool.Get().(*Module[T])
	if ok && rv != nil {
		return rv, nil
	}

	return p.new()
}
