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
	"github.com/loopholelabs/scale/go/utils"
	"github.com/tetratelabs/wazero/api"
)

type Module struct {
	module   api.Module
	instance *Instance
	function *Function
	next     *Module
	run      api.Function
	resize   api.Function
}

func (m *Module) Run(ctx context.Context) error {
	ctxBuffer := m.instance.Context().Write()
	ctxBufferLength := uint64(len(ctxBuffer))
	writeBuffer, err := m.resize.Call(ctx, ctxBufferLength)
	if err != nil {
		return fmt.Errorf("failed to allocate memory for function '%s': %w", m.function.ScaleFunc.ScaleFile.Name, err)
	}

	if !m.module.Memory().Write(ctx, uint32(writeBuffer[0]), ctxBuffer) {
		return fmt.Errorf("failed to write memory for function '%s'", m.function.ScaleFunc.ScaleFile.Name)
	}

	packed, err := m.run.Call(ctx, writeBuffer[0], ctxBufferLength)
	if err != nil {
		return fmt.Errorf("failed to run function '%s': %w", m.function.ScaleFunc.ScaleFile.Name, err)
	}
	if packed[0] == 0 {
		return fmt.Errorf("failed to run function '%s'", m.function.ScaleFunc.ScaleFile.Name)
	}

	offset, length := utils.UnpackUint32(packed[0])
	readBuffer, ok := m.module.Memory().Read(ctx, offset, length)
	if !ok {
		return fmt.Errorf("failed to read memory for function '%s'", m.function.ScaleFunc.ScaleFile.Name)
	}

	err = m.instance.Context().Read(readBuffer)
	if err != nil {
		return fmt.Errorf("failed to deserialize context for function '%s': %w", m.function.ScaleFunc.ScaleFile.Name, err)
	}

	return nil
}
