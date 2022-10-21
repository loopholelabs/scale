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
	"fmt"
	"github.com/google/uuid"
	"github.com/tetratelabs/wazero/api"
)

type Instance struct {
	id      string
	next    Next
	runtime *Runtime
	ctx     *Context
	head    *Module
	tail    *Module
	modules map[api.Module]*Module
}

func (r *Runtime) Instance(ctx context.Context, next Next) (*Instance, error) {
	i := &Instance{
		id:      uuid.New().String(),
		next:    next,
		runtime: r,
		modules: make(map[api.Module]*Module),
		ctx:     NewContext(),
	}

	if i.next == nil {
		endpoint := false
		for _, f := range r.functions {
			if !f.ScaleFunc.ScaleFile.Middleware {
				endpoint = true
				break
			}
		}
		if !endpoint {
			return nil, NextFunctionRequiredError
		}
		i.next = func(ctx *Context) *Context {
			return ctx
		}
	}

	r.instanceMu.Lock()
	r.instances[i.id] = i
	r.instanceMu.Unlock()

	return i, i.initialize(ctx)
}

func (i *Instance) Context() *Context {
	return i.ctx
}

func (i *Instance) Run(ctx context.Context) error {
	defer func() {
		i.runtime.instanceMu.Lock()
		delete(i.runtime.instances, i.id)
		i.runtime.instanceMu.Unlock()
	}()
	if i.head == nil {
		return fmt.Errorf("no functions registered for instance %s", i.id)
	}
	module := i.head
	return module.Run(ctx)
}

func (i *Instance) initialize(ctx context.Context) error {
	for _, f := range i.runtime.functions {
		module, err := i.runtime.runtime.InstantiateModule(ctx, f.Compiled, i.runtime.moduleConfig.WithName(fmt.Sprintf("%s.%s", i.id, f.ScaleFunc.ScaleFile.Name)))
		if err != nil {
			return fmt.Errorf("failed to instantiate function '%s' for instance %s: %w", f.ScaleFunc.ScaleFile.Name, i.id, err)
		}

		run := module.ExportedFunction("run")
		resize := module.ExportedFunction("resize")
		if run == nil || resize == nil {
			return fmt.Errorf("failed to find run or resize functions for instance %s", i.id)
		}

		m := &Module{
			module:   module,
			function: f,
			instance: i,
			run:      run,
			resize:   resize,
		}
		i.modules[module] = m
		if i.head == nil {
			i.head = m
		}
		if i.tail != nil {
			i.tail.next = m
		}
		i.tail = m
	}
	return nil
}
