package runtime

import (
	"github.com/loopholelabs/polyglot-go"
	"github.com/loopholelabs/scale-go/runtime/generated"
	"io"
)

// Context is a wrapper around generated.Context
// and it facilitates the exchange of data between
// the runtime and the scale functions
type Context struct {
	Context      *generated.Context
	Buffer       *polyglot.Buffer
	RawBuffer    []byte
	StreamBuffer *io.Reader
}

// NewContext returns a new Context
func NewContext() *Context {
	return &Context{
		Context: generated.NewContext(),
		Buffer:  polyglot.NewBuffer(),
	}
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
