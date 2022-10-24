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

type Module struct {
	module   api.Module
	function *Function
	runtime  *Runtime
	run      api.Function
	resize   api.Function

	instance *Instance
}

func NewModule(ctx context.Context, f *Function, r *Runtime) (*Module, error) {
	module, err := r.runtime.InstantiateModule(ctx, f.compiled, r.moduleConfig.WithName(fmt.Sprintf("%s.%s", f.scaleFunc.ScaleFile.Name, uuid.New().String())))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate function '%s': %w", f.scaleFunc.ScaleFile.Name, err)
	}

	run := module.ExportedFunction("run")
	resize := module.ExportedFunction("resize")
	if run == nil || resize == nil {
		return nil, fmt.Errorf("failed to find run or resize implementations for function %s", f.scaleFunc.ScaleFile.Name)
	}

	return &Module{
		module:   module,
		function: f,
		runtime:  r,
		run:      run,
		resize:   resize,
	}, nil
}

func (m *Module) init(i *Instance) {
	m.instance = i
	m.runtime.modulesMu.Lock()
	m.runtime.modules[m.module.Name()] = m
	m.runtime.modulesMu.Unlock()
}

func (m *Module) reset() {
	m.instance = nil
	m.runtime.modulesMu.Lock()
	delete(m.runtime.modules, m.module.Name())
	m.runtime.modulesMu.Unlock()
}
