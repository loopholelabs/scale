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
	"github.com/tetratelabs/wazero/api"
)

func (r *Runtime) next(ctx context.Context, module api.Module, pointer uint32, length uint32) {
	r.modulesMu.RLock()
	m := r.modules[module.Name()]
	r.modulesMu.RUnlock()

	if m == nil {
		return
	}

	buf, ok := m.module.Memory().Read(ctx, pointer, length)
	if !ok {
		return
	}

	err := m.instance.Context().Read(buf)
	if err != nil {
		return
	}

	if m.function.next == nil {
		m.instance.ctx, err = m.instance.next(m.instance.Context())
		if err != nil {
			return
		}
	} else {
		err = m.function.next.Run(ctx, m.instance)
		if err != nil {
			return
		}
	}

	ctxBuffer := m.instance.Context().Write()
	ctxBufferLength := uint64(len(ctxBuffer))
	writeBuffer, err := m.resize.Call(ctx, ctxBufferLength)
	if err != nil {
		return
	}
	module.Memory().Write(ctx, uint32(writeBuffer[0]), ctxBuffer)
}
