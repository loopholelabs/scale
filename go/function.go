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

	signature "github.com/loopholelabs/scale-signature"
	"github.com/loopholelabs/scale/go/utils"
	"github.com/loopholelabs/scalefile/scalefunc"
	"github.com/tetratelabs/wazero"
)

// Function is the runtime representation of a scale function.
type Function[T signature.Signature] struct {
	compiled   wazero.CompiledModule
	scaleFunc  *scalefunc.ScaleFunc
	next       *Function[T]
	modulePool *Pool[T]
}

func (f *Function[T]) Run(ctx context.Context, i *Instance[T]) error {
	module, err := f.modulePool.Get()
	if err != nil {
		return fmt.Errorf("failed to get module from pool for function %s: %w", f.scaleFunc.Name, err)
	}

	module.init(i)
	defer func() {
		module.reset()
		f.modulePool.Put(module)
	}()

	ctxBuffer := i.runtimeContext().Write()
	ctxBufferLength := uint64(len(ctxBuffer))
	writeBuffer, err := module.resize.Call(ctx, ctxBufferLength)
	if err != nil {
		return fmt.Errorf("failed to allocate memory for function '%s': %w", f.scaleFunc.Name, err)
	}

	if !module.module.Memory().Write(uint32(writeBuffer[0]), ctxBuffer) {
		return fmt.Errorf("failed to write memory for function '%s'", f.scaleFunc.Name)
	}

	packed, err := module.run.Call(ctx)
	if err != nil {
		return fmt.Errorf("failed to run function '%s': %w", f.scaleFunc.Name, err)
	}
	if packed[0] == 0 {
		return fmt.Errorf("failed to run function '%s'", f.scaleFunc.Name)
	}

	offset, length := utils.UnpackUint32(packed[0])
	readBuffer, ok := module.module.Memory().Read(offset, length)
	if !ok {
		return fmt.Errorf("failed to read memory for function '%s'", f.scaleFunc.Name)
	}

	err = i.runtimeContext().Read(readBuffer)
	if err != nil {
		return fmt.Errorf("failed to deserialize context for function '%s': %w", f.scaleFunc.Name, err)
	}

	return nil
}
