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
	"crypto/rand"
	"github.com/loopholelabs/scale/signature"
)

// Instance is a single instance of a Scale Function chain
type Instance[T signature.Signature] struct {
	runtime    *Scale[T]
	instanceID []byte

	head *moduleInstance[T]

	next Next[T]
}

func newInstance[T signature.Signature](r *Scale[T], next ...Next[T]) (*Instance[T], error) {
	i := &Instance[T]{
		runtime:    r,
		instanceID: make([]byte, 16),
	}

	_, err := rand.Read(i.instanceID)
	if err != nil {
		return nil, err
	}

	if len(next) > 0 && next[0] != nil {
		i.next = next[0]
	} else {
		i.next = func(ctx T) (T, error) {
			return ctx, nil
		}
	}

	return i, nil
}

func newPersistentInstance[T signature.Signature](ctx context.Context, r *Scale[T], next ...Next[T]) (*Instance[T], error) {
	i, err := newInstance[T](r, next...)
	if err != nil {
		return nil, err
	}

	i.head, err = newModuleInstance(ctx, i.runtime, i.runtime.head, i)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (i *Instance[T]) Run(ctx context.Context, signature T) error {
	if i.head != nil {
		i.head.setSignature(signature)
		return i.head.function.runWithModule(ctx, signature, i.head)
	}

	return i.runtime.head.run(ctx, signature, i)
}

func (i *Instance[T]) Cleanup() {
	if i.head != nil {
		i.head.cleanup(i.runtime)
	}
}
