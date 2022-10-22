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
)

type Module struct {
	pool     *Pool
	instance *Instance
	function *Function
	next     *Module
}

func (m *Module) Run(ctx context.Context) error {
	im, err := m.pool.Get()
	if err != nil {
		return fmt.Errorf("failed to get instance from pool: %w", err)
	}

	defer func() {
		m.instance.instantiatedMu.Lock()
		delete(m.instance.instantiated, im.module.Name())
		m.instance.instantiatedMu.Unlock()
		m.pool.Put(im)
	}()

	m.instance.instantiatedMu.Lock()
	m.instance.instantiated[im.module.Name()] = im
	m.instance.instantiatedMu.Unlock()

	ctxBuffer := m.instance.Context().Write()
	ctxBufferLength := uint64(len(ctxBuffer))
	writeBuffer, err := im.resize.Call(ctx, ctxBufferLength)
	if err != nil {
		return fmt.Errorf("failed to allocate memory for function '%s': %w", m.function.ScaleFunc.ScaleFile.Name, err)
	}

	if !im.module.Memory().Write(ctx, uint32(writeBuffer[0]), ctxBuffer) {
		return fmt.Errorf("failed to write memory for function '%s'", m.function.ScaleFunc.ScaleFile.Name)
	}

	packed, err := im.run.Call(ctx)
	if err != nil {
		return fmt.Errorf("failed to run function '%s': %w", m.function.ScaleFunc.ScaleFile.Name, err)
	}
	if packed[0] == 0 {
		return fmt.Errorf("failed to run function '%s'", m.function.ScaleFunc.ScaleFile.Name)
	}

	offset, length := utils.UnpackUint32(packed[0])
	readBuffer, ok := im.module.Memory().Read(ctx, offset, length)
	if !ok {
		return fmt.Errorf("failed to read memory for function '%s'", m.function.ScaleFunc.ScaleFile.Name)
	}

	err = m.instance.Context().Read(readBuffer)
	if err != nil {
		return fmt.Errorf("failed to deserialize context for function '%s': %w", m.function.ScaleFunc.ScaleFile.Name, err)
	}

	return nil
}
