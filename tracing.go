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

	"github.com/tetratelabs/wazero/api"
)

// getFunctionNameLen is the Host function for getting the Function Name Length
func (r *Scale[T]) getFunctionNameLen(_ context.Context, module api.Module, params []uint64) {
	r.activeModulesMu.RLock()
	m := r.activeModules[module.Name()]
	r.activeModulesMu.RUnlock()
	if m == nil {
		return
	}

	params[0] = uint64(len([]byte(m.function.identifier)))
}

// getFunctionName is the Host function for getting the Function Name
func (r *Scale[T]) getFunctionName(_ context.Context, module api.Module, params []uint64) {
	r.activeModulesMu.RLock()
	m := r.activeModules[module.Name()]
	r.activeModulesMu.RUnlock()
	if m == nil {
		return
	}

	ptr := uint32(params[0])
	mem := module.Memory()
	mem.Write(ptr, []byte(m.function.identifier))
}

// getInstanceID is the Host function to get 16 byte Instance ID
func (r *Scale[T]) getInstanceID(_ context.Context, module api.Module, params []uint64) {
	r.activeModulesMu.RLock()
	m := r.activeModules[module.Name()]
	r.activeModulesMu.RUnlock()
	if m == nil {
		return
	}

	ptr := uint32(params[0])
	mem := module.Memory()
	mem.Write(ptr, m.instance.instanceID)
}

// otelTraceJSON is the Host function to receive OTEL Trace data in JSON
// and then call the TraceDataCallback
func (r *Scale[T]) otelTraceJSON(_ context.Context, module api.Module, params []uint64) {
	if r.TraceDataCallback == nil {
		return
	}

	ptr := uint32(params[0])
	length := uint32(params[1])
	mem := module.Memory()
	if data, ok := mem.Read(ptr, length); ok {
		r.TraceDataCallback(string(data))
	}
}
