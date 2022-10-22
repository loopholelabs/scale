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
	"sync"
)

type Pool struct {
	pool sync.Pool
	new  func() (*Module, error)
}

func NewPool(ctx context.Context, f *Function, r *Runtime) *Pool {
	return &Pool{
		new: func() (*Module, error) {
			return NewModule(ctx, f, r)
		},
	}
}

func (p *Pool) Put(module *Module) {
	if module != nil {
		p.pool.Put(module)
	}
}

func (p *Pool) Get() (*Module, error) {
	rv, ok := p.pool.Get().(*Module)
	if ok && rv != nil {
		return rv, nil
	}

	return p.new()
}
