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
	"github.com/loopholelabs/scale-go/scalefunc"
	"github.com/tetratelabs/wazero"
)

// Function is the runtime representation of a scale function.
type Function struct {
	ScaleFunc scalefunc.ScaleFunc
	Compiled  wazero.CompiledModule
}

func (r *Runtime) registerFunction(ctx context.Context, scaleFunc scalefunc.ScaleFunc) error {
	compiled, err := r.runtime.CompileModule(ctx, scaleFunc.Function)
	if err != nil {
		return fmt.Errorf("failed to compile function '%s': %w", scaleFunc.ScaleFile.Name, err)
	}

	f := &Function{
		ScaleFunc: scaleFunc,
		Compiled:  compiled,
	}

	r.functions = append(r.functions, f)
	return nil
}
