//go:build !tinygo && !js && !wasm

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
)

var _ signature.Signature = (*Context)(nil)
var _ signature.RuntimeContext = (*RuntimeContext)(nil)

type RuntimeContext Context

// Context is a context object for an incoming request. It is meant to be used
// inside the Scale function.
type Context struct {
	*BadContext
	buffer *polyglot.Buffer
}

// New creates a new empty Context
func New() *Context {
	return &Context{
		BadContext: NewBadContext(),
		buffer:     polyglot.NewBuffer(),
	}
}

// RuntimeContext converts a Context into a RuntimeContext.
func (x *Context) RuntimeContext() signature.RuntimeContext {
	return (*RuntimeContext)(x)
}

// Read reads the context from the given byte slice and returns an error if one occurred
//
// This method is meant to be used by the Scale Runtime to deserialize the Context
func (x *RuntimeContext) Read(b []byte) error {
	return x.internalDecode(b)
}

// Write writes the context into a byte slice and returns it
func (x *RuntimeContext) Write() []byte {
	x.buffer.Reset()
	x.internalEncode(x.buffer)
	return x.buffer.Bytes()
}

// Error writes the context into a byte slice and returns it
func (x *RuntimeContext) Error(err error) []byte {
	x.buffer.Reset()
	x.error(x.buffer, err)
	return x.buffer.Bytes()
}
