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
	"errors"
	"github.com/loopholelabs/polyglot-go"
	"github.com/loopholelabs/scale-go/runtime/generated"
	"unsafe"
)

var (
	InvalidPointerError = errors.New("invalid pointer")
)

// global buffer to copy data into
var Buffer *byte

// Context is a context object for an incoming request. It is meant to be used
// inside the Scale function.
type Context struct {
	generated *generated.Context
	buffer    *polyglot.Buffer
}

// New creates a new empty Context that must be initialized with the FromPointer method
func New() *Context {
	return &Context{
		generated: generated.NewContext(),
		buffer:    polyglot.NewBuffer(),
	}
}

// ToPointer serializes the Context into a pointer and returns the pointer and its size
func (ctx *Context) ToBuffer() (uint32, uint32) {
	ctx.buffer.Reset()
	ctx.generated.Encode(ctx.buffer)
	Buffer = &ctx.buffer.Bytes()[0]
	//ptr := &underlying[0]
	unsafePtr := uintptr(unsafe.Pointer(Buffer))
	return uint32(unsafePtr), uint32(ctx.buffer.Len())
}

// FromPointer takes a pointer and size and deserializes the data into the Context
func (ctx *Context) FromBuffer(length uint32) error {
	if length == 0 {
		return InvalidPointerError
	}
	buf := unsafe.Slice(Buffer, length)
	//buf := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
	//	Data: uintptr(buffer),
	//	Len:  uintptr(size), // Tinygo requires this, it's not an error.
	//	Cap:  uintptr(size), // ^^ See https://github.com/tinygo-org/tinygo/issues/1284
	//}))
	return ctx.generated.Decode(buf)
}

// Next calls the next host function after writing the Context,
// then it reads the result back into the Context
func (ctx *Context) Next() *Context {
	ctx.FromBuffer(next(ctx.ToBuffer()))
	return ctx
}
