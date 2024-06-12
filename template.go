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

	"github.com/tetratelabs/wazero"

	interfaces "github.com/loopholelabs/scale-signature-interfaces"
	"github.com/loopholelabs/scale/scalefunc"
)

// template is a template for creating runnable scale functions
type template[T interfaces.Signature] struct {
	// runtime is the scale runtime that the template belongs to
	runtime *Scale[T]

	// identifier is the identifier for the template
	identifier string

	// compiled is the compiled module source
	compiled wazero.CompiledModule

	// next is the (optional) next function in the chain
	next *template[T]

	// modulePool is the pool of modules for the template
	modulePool *modulePool[T]

	// env contains the (optional) environment variables for the template
	env map[string]string
}

// newTemplate creates a new template from a scale function schema
func newTemplate[T interfaces.Signature](ctx context.Context, runtime *Scale[T], scaleFunc *scalefunc.V1BetaSchema, env map[string]string) (*template[T], error) {
	compiled, err := runtime.runtime.CompileModule(ctx, scaleFunc.Function)
	if err != nil {
		return nil, fmt.Errorf("failed to compile wasm module '%s': %w", scaleFunc.Name, err)
	}

	templ := &template[T]{
		runtime:    runtime,
		identifier: fmt.Sprintf("%s:%s", scaleFunc.Name, scaleFunc.Tag),
		compiled:   compiled,
		env:        env,
	}

  var maxSize uint32 = 0;

	if scaleFunc.Stateless {
		templ.modulePool = newModulePool[T](ctx, templ, maxSize)
	}

	return templ, nil
}
