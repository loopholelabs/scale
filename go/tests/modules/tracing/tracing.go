//go:build tinygo || js || wasm

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
	"fmt"
	"time"
	"unsafe"

	signature "github.com/loopholelabs/scale/go/tests/signature/example-signature"
)

func Scale(ctx *signature.Context) (*signature.Context, error) {

	invocationID := make([]byte, 16)
	// Fill in the invocation ID
	get_invocation_id(uint32(uintptr(unsafe.Pointer(unsafe.SliceData(invocationID)))))

	length := get_service_name_len()
	serviceName := make([]byte, length)
	get_service_name(uint32(uintptr(unsafe.Pointer(unsafe.SliceData(serviceName)))))

	now := time.Now()

	// Send a single trace item with the invocationID and serviceName provided by the host.
	data := fmt.Sprintf("{\"invocationID\":\"%x\", \"serviceName\":\"%s\", \"timestamp\": %d}", invocationID, serviceName, now.UnixNano())

	ptr := unsafe.Pointer(unsafe.StringData(data))
	send_otel_trace_json(uint32(uintptr(ptr)), uint32(len(data)))

	return ctx.Next()
}

//go:wasm-module scale
//export get_invocation_id
func get_invocation_id(ptr uint32)

//go:wasm-module scale
//export get_service_name_len
func get_service_name_len() uint32

//go:wasm-module scale
//export get_service_name
func get_service_name(ptr uint32)

//go:wasm-module scale
//export send_otel_trace_json
func send_otel_trace_json(ptr uint32, len uint32)
