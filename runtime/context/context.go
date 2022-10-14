package context

import (
	"github.com/loopholelabs/polyglot-go"
	"github.com/loopholelabs/scale-go/runtime/context/generated"
	"io"
)

type Context struct {
	Context      *generated.Context
	Buffer       *polyglot.Buffer
	RawBuffer    []byte
	StreamBuffer *io.Reader
}

func New() *Context {
	return &Context{
		Context: generated.NewContext(),
		Buffer:  polyglot.NewBuffer(),
	}
}

func (c *Context) Deserialize(b []byte) error {
	return c.Context.Decode(b)
}

func (c *Context) Encode() {
	c.Buffer.Reset()
	c.Context.Encode(c.Buffer)
}
