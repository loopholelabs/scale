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
	"errors"
	"github.com/loopholelabs/scale-signature"
)

type Instance struct {
	next    Next
	runtime *Runtime
	ctx     signature.RuntimeContext
}

func (r *Runtime) Instance(next Next) (*Instance, error) {
	i := &Instance{
		next:    next,
		runtime: r,
		ctx:     r.signature.RuntimeContext(),
	}

	if i.next == nil {
		i.next = func(ctx signature.RuntimeContext) signature.RuntimeContext {
			return ctx
		}
	}

	return i, nil
}

func (i *Instance) Context() signature.RuntimeContext {
	return i.ctx
}

func (i *Instance) Run(ctx context.Context) error {
	if i.runtime.head == nil {
		return errors.New("no compiled functions found in runtime")
	}
	function := i.runtime.head
	return function.Run(ctx, i)
}
