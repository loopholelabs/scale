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

	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
	"github.com/tetratelabs/wazero"
)

// Function is the runtime representation of a scale function.
type Function[T signature.Signature] struct {
	identifier string
	compiled   wazero.CompiledModule
	scaleFunc  *scalefunc.Schema
	next       *Function[T]
	modulePool *Pool[T]
}

func (f *Function[T]) Run(ctx context.Context, signature T, i *Instance[T]) error {
	module, err := f.modulePool.Get()
	if err != nil {
		return fmt.Errorf("failed to get module from pool for function %s: %w", f.scaleFunc.Name, err)
	}

	module.init(signature, i)
	buf := signature.Write()
	ctxBufferLength := uint64(len(buf))
	writeBuffer, err := module.resize.Call(ctx, ctxBufferLength)
	if err != nil {
		return fmt.Errorf("failed to allocate memory for function '%s': %w", f.scaleFunc.Name, err)
	}

	if !module.module.Memory().Write(uint32(writeBuffer[0]), buf) {
		return fmt.Errorf("failed to write memory for function '%s'", f.scaleFunc.Name)
	}

	packed, err := module.run.Call(ctx)
	if err != nil {
		return fmt.Errorf("failed to run function '%s': %w", f.scaleFunc.Name, err)
	}
	if packed[0] == 0 {
		return fmt.Errorf("failed to run function '%s'", f.scaleFunc.Name)
	}

	offset, length := unpackUint32(packed[0])
	buf, ok := module.module.Memory().Read(offset, length)
	if !ok {
		return fmt.Errorf("failed to read memory for function '%s'", f.scaleFunc.Name)
	}

	err = signature.Read(buf)
	if err != nil {
		return fmt.Errorf("error while running function '%s': %w", f.scaleFunc.Name, err)
	}

	module.reset()
	f.modulePool.Put(module)
	return nil
}
