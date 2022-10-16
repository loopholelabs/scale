package runtime

import (
	"github.com/loopholelabs/polyglot-go"
	"github.com/loopholelabs/scale-go/runtime/generated"
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

	return c
}

// Deserialize deserializes the context from the given byte slice
func (c *Context) Deserialize(b []byte) error {
	return c.Context.Decode(b)
}

// Serialize serializes the context into a byte slice and returns it
func (c *Context) Serialize() []byte {
	c.Buffer.Reset()
	c.Context.Encode(c.Buffer)
	return c.Buffer.Bytes()
}
