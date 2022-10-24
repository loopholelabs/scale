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
	"github.com/loopholelabs/polyglot-go"
	"github.com/loopholelabs/scale/generated"
	"net/http"
)

// Context is a wrapper around generated.Context
// and it facilitates the exchange of data between
// the runtime and the scale functions
type Context struct {
	Context *generated.Context
	Buffer  *polyglot.Buffer
}

// NewContext returns a new Context
func NewContext() *Context {
	c := &Context{
		Context: generated.NewContext(),
		Buffer:  polyglot.NewBuffer(),
	}
	c.Context.Request.Headers = generated.NewRequestHeadersMap(8)
	c.Context.Response.Headers = generated.NewResponseHeadersMap(8)

	c.Context.Response.StatusCode = http.StatusOK

	return c
}

// Read reads the context from the given byte slice
func (c *Context) Read(b []byte) error {
	return c.Context.Decode(b)
}

// Write writes the context into a byte slice and returns it
func (c *Context) Write() []byte {
	c.Buffer.Reset()
	c.Context.Encode(c.Buffer)
	return c.Buffer.Bytes()
}
