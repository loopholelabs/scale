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
	interfaces "github.com/loopholelabs/scale-signature-interfaces"
	"sync"
)

var (
	// CacheAllocatedModules is the number of allocated modules that we will
	// maintain in a cache for reuse. If we need more modules than we have
	// cached, they will be allocated and provided to Scale. All modules are
	// returned first to the cache, or if the cache is full, to the sync.Pool
	// which may choose to de-allocate unused modules at any time. In summary,
	// this knob tunes the maximum number of modules that we will keep around
	// for reuse indefinitely, or in other words is the maximum number of
	// modules available for use immediately without allocation.
	CacheAllocatedModules = 1
)

type modulePool[T interfaces.Signature] struct {
	// primary is the default pool of modules. A module is always pulled from
	// or put back into the primary pool first if possible. This is a buffered
	// channel. If we're trying to put and the buffer is full, then we put into
	// the fallback which is a sync.Pool. If we're trying to get and the channel
	// is empty, then we try to get from the fallback.
	//
	// The intent is to reduce the required module allocation since sync.Pool
	// may empty unused references on each GC cycle.
	//
	// See https://github.com/golang/go/issues/22950.
	primary  chan *module[T]
	fallback sync.Pool
	new      func() (*module[T], error)
}

func newModulePool[T interfaces.Signature](ctx context.Context, template *template[T]) *modulePool[T] {
	return &modulePool[T]{
		primary: make(chan *module[T], CacheAllocatedModules),
		new: func() (*module[T], error) {
			return newModule[T](ctx, template)
		},
	}
}

func (p *modulePool[T]) Put(m *module[T]) {
	if m != nil {
		select {
		case p.primary <- m:
		default:
			p.fallback.Put(m)
		}
	}
}

func (p *modulePool[T]) Get() (*module[T], error) {
	select {
	case m := <-p.primary:
		return m, nil
	default:
		m, ok := p.fallback.Get().(*module[T])
		if m != nil && ok {
			return m, nil
		}
		return p.new()
	}
}
