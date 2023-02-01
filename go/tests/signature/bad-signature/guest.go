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

package bad

import (
	"github.com/loopholelabs/polyglot-go"
	"github.com/loopholelabs/scale-signature"
	"unsafe"
)

var _ signature.GuestContext = (*GuestContext)(nil)
var _ signature.Context = (*Context)(nil)

var (
	writeBuffer = polyglot.NewBuffer()
	readBuffer  []byte
)

type GuestContext Context

// Context is a context object for an incoming request. It is meant to be used
// inside the Scale function.
type Context struct {
	*BadContext
}

// New creates a new empty Context
func New() *Context {
	return &Context{
		BadContext: NewBadContext(),
	}
}

// GuestContext converts the given Context to a GuestContext.
func (x *Context) GuestContext() signature.GuestContext {
	return (*GuestContext)(x)
}

// ToWriteBuffer serializes the Context into the global writeBuffer and returns the pointer to the buffer and its size
//
// This method should only be used to read the Context from the Scale Runtime.
// Users should not use this method.
func (x *GuestContext) ToWriteBuffer() (uint32, uint32) {
	writeBuffer.Reset()
	x.internalEncode(writeBuffer)
	underlying := writeBuffer.Bytes()
	ptr := &underlying[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(writeBuffer.Len())
}

// FromReadBuffer deserializes the data into the Context from the global readBuffer
//
// It assumes that the readBuffer has been filled with the data from the Scale Runtime after
// a call to the Resize method
func (x *GuestContext) FromReadBuffer() error {
	return x.internalDecode(readBuffer)
}

// ErrorWriteBuffer serializes an error into the global writeBuffer and returns a pointer to the buffer and its size
//
// This method should only be used to write an error to the Scale Runtime, in place of the ToWriteBuffer method.
// Users should not use this method.
func (x *GuestContext) ErrorWriteBuffer(err error) (uint32, uint32) {
	writeBuffer.Reset()
	x.error(writeBuffer, err)
	underlying := writeBuffer.Bytes()
	ptr := &underlying[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(writeBuffer.Len())
}

// Next calls the next host function after writing the Context into the global writeBuffer,
// then it reads the result from the global readBuffer back into the Context
func (x *Context) Next() (*Context, error) {
	next(x.GuestContext().ToWriteBuffer())
	return x, x.GuestContext().FromReadBuffer()
}

func Resize(size uint32) uint32 {
	if uint32(cap(readBuffer)) < size {
		readBuffer = append(make([]byte, 0, uint32(len(readBuffer))+size), readBuffer...)
	}
	readBuffer = readBuffer[:size]
	return uint32(uintptr(unsafe.Pointer(&readBuffer[0])))
}

//export next
//go:linkname next
func next(offset uint32, length uint32)
