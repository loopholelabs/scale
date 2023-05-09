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
	"crypto/rand"
	"errors"

	signature "github.com/loopholelabs/scale-signature"
)

type Instance[T signature.Signature] struct {
	next       Next[T]
	runtime    *Runtime[T]
	ctx        T
	runtimeCtx signature.RuntimeContext
}

func (r *Runtime[T]) Instance(next ...Next[T]) (*Instance[T], error) {
	ctx := r.new()
	i := &Instance[T]{
		runtime:    r,
		ctx:        ctx,
		runtimeCtx: ctx.RuntimeContext(),
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

func (i *Instance[T]) Context() T {
	return i.ctx
}

func (i *Instance[T]) runtimeContext() signature.RuntimeContext {
	return i.runtimeCtx
}

func (i *Instance[T]) Run(ctx context.Context) error {
	if i.runtime.head == nil {
		return errors.New("no compiled functions found in runtime")
	}

	// Create a random InvocationID for this Run. This will be shared between functions in the chain.
	rand.Read(i.runtime.InvocationID)

	function := i.runtime.head
	return function.Run(ctx, i)
}
