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
	pool    sync.Pool
	maxSize uint32
	new     func(*module[T]) (*module[T], error)
	close   func(*module[T])
	ch      chan *module[T] // for buffered channel impl
}

func newModulePool[T interfaces.Signature](ctx context.Context, template *template[T], maxSize uint32) *modulePool[T] {
	if maxSize == 0 {
		// If size = 0, use a standard sync.Pool + golang's runtime finalizers
		// to make sure the close function gets called eventually.
		return &modulePool[T]{
			maxSize: maxSize,
			new: func(m *module[T]) (*module[T], error) {
				m.SetFinalizer()
				return newModule[T](ctx, template)
			},
			close: func(m *module[T]) {
				m.Close()
			},
		}
	}

	if maxSize > 0 {
		// if size > 0 then we use buffered channel implementation.
		return &modulePool[T]{
			maxSize: maxSize,
			new: func(m *module[T]) (*module[T], error) {
				return newModule[T](ctx, template)
			},
			close: func(m *module[T]) {
				m.Close()
			},
			ch: make(chan *module[T], maxSize),
		}
	}

	return nil
}

func (p *modulePool[T]) Put(m *module[T]) {
	if m == nil {
		return
	}

	if p.maxSize == 0 {
		p.pool.Put(m)
	} else {
		select {
		case p.ch <- m:
		default:
			// Channel is full, call the close function
			p.close(m)
		}
	}
}

// if you call pool.Get and its a success

//otherwise return m = p.New
//OR
//on line 40, in the New function run Close (change to SetFinalizers),

func (p *modulePool[T]) Get() (*module[T], error) {
	if p.maxSize == 0 {
		m, ok := p.pool.Get().(*module[T])
		if ok && m != nil {
			return m, nil
		}
		return p.new(m)
	}

	// Use buffered channel
	select {
	case m := <-p.ch:
		return m, nil
	default:
		// Channel is empty, create a new module?
		m, _ := p.pool.Get().(*module[T])
		return m, nil
	}
}
