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
	"strings"
)

func (r *Runtime) next(ctx context.Context, module api.Module, offset uint32, length uint32) {
	i := r.instances[strings.Split(module.Name(), ".")[0]]
	if i == nil {
		return
	}

	im := i.instantiated[module.Name()]
	if im == nil {
		return
	}

	buf, ok := module.Memory().Read(ctx, offset, length)
	if !ok {
		return
	}

	err := i.Context().Read(buf)
	if err != nil {
		return
	}

	if im.m.next == nil {
		i.ctx = i.next(i.Context())
	} else {
		err = im.m.next.Run(ctx)
		if err != nil {
			return
		}
	}

	ctxBuffer := i.Context().Write()
	ctxBufferLength := uint64(len(ctxBuffer))
	writeBuffer, err := im.resize.Call(ctx, ctxBufferLength)
	if err != nil {
		return
	}
	module.Memory().Write(ctx, uint32(writeBuffer[0]), ctxBuffer)
}
