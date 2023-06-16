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
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/loopholelabs/scale/signature"
	"github.com/tetratelabs/wazero/api"
)

var (
	EnvStringRegex = regexp.MustCompile(`[^A-Za-z0-9_]`)
)

// ValidEnv returns true if the string is valid for use as an environment variable
func ValidEnv(str string) bool {
	return !EnvStringRegex.MatchString(str)
}

type Module[T signature.Signature] struct {
	module   api.Module
	function *Function[T]
	runtime  *Scale[T]

	run    api.Function
	resize api.Function

	instance  *Instance[T]
	signature T
}

func NewModule[T signature.Signature](ctx context.Context, f *Function[T], r *Scale[T]) (*Module[T], error) {
	config := r.moduleConfig.WithName(fmt.Sprintf("%s.%s", f.identifier, uuid.New().String()))
	for k, v := range f.scaleFunc.GetEnv() {
		if ValidEnv(k) && len(k) > 0 {
			config = config.WithEnv(k, v)
		}
	}

	module, err := r.runtime.InstantiateModule(ctx, f.compiled, config)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate module '%s': %w", f.scaleFunc.Name, err)
	}

	run := module.ExportedFunction("run")
	resize := module.ExportedFunction("resize")
	if run == nil || resize == nil {
		return nil, fmt.Errorf("failed to find run or resize implementations for function %s", f.scaleFunc.Name)
	}

	return &Module[T]{
		module:   module,
		function: f,
		runtime:  r,
		run:      run,
		resize:   resize,
	}, nil
}

func (m *Module[T]) init(signature T, i *Instance[T]) {
	m.signature = signature
	m.instance = i
	m.runtime.modulesMu.Lock()
	m.runtime.modules[m.module.Name()] = m
	m.runtime.modulesMu.Unlock()
}

func (m *Module[T]) reset() {
	m.runtime.modulesMu.Lock()
	delete(m.runtime.modules, m.module.Name())
	m.runtime.modulesMu.Unlock()
}
