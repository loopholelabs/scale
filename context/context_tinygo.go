//go:build tinygo
// +build tinygo

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

package context

import (
	"github.com/loopholelabs/polyglot-go"
	"github.com/loopholelabs/scale-go/runtime/generated"
	"github.com/loopholelabs/scale-go/utils"
	"reflect"
	"unsafe"
)

// Context is a context object for an incoming request. It is meant to be used
// inside the Scale function.
type Context struct {
	generated *generated.Context
	buffer    *polyglot.Buffer
}

// New creates a new empty Context that must be initialized with the Deserialize method
func New() *Context {
	return &Context{
		generated: generated.NewContext(),
		buffer:    polyglot.NewBuffer(),
	}
}

// Serialize serializes the Context into a pointer and returns the pointer and its size
func (ctx *Context) Serialize() (uint32, uint32) {
	ctx.buffer.Reset()
	ctx.generated.Encode(ctx.buffer)
	underlying := ctx.buffer.Bytes()
	ptr := &underlying[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(ctx.buffer.Len())
}

// Deserialize takes a pointer and size and deserializes the data into the Context
func (ctx *Context) Deserialize(ptr uint32, size uint32) {
	if ptr == 0 || size == 0 {
		panic("context: invalid pointer or size")
	}
	buf := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  uintptr(size), // Tinygo requires this, it's not an error.
		Cap:  uintptr(size), // ^^ See https://github.com/tinygo-org/tinygo/issues/1284
	}))
	_ = ctx.generated.Decode(buf)
}

// Next calls the next host function after writing the Context,
// then it reads the result back into the Context
func (ctx *Context) Next() *Context {
	ctx.Deserialize(utils.UnpackUint32(next(ctx.Serialize())))
	return ctx
}

//func Debug(str string) {
//	data := []byte(str)
//	ptr := &data[0]
//	unsafePtr := uintptr(unsafe.Pointer(ptr))
//	debug(uint32(unsafePtr), uint32(len(data)))
//}
