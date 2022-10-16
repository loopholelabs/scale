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

// New creates a new empty Context that must be initialized with the Read method
func New() *Context {
	return &Context{
		generated: generated.NewContext(),
		buffer:    polyglot.NewBuffer(),
	}
}

// Write writes the Context into a pointer and returns the pointer and its size
func (ctx *Context) Write() (uint32, uint32) {
	ctx.buffer.Reset()
	ctx.generated.Encode(ctx.buffer)
	underlying := ctx.buffer.Bytes()
	ptr := &underlying[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(ctx.buffer.Len())
}

// Read takes a pointer and size and reads the data into the Context
func (ctx *Context) Read(ptr uint32, size uint32) {
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
	ctx.Read(utils.UnpackUint32(next(ctx.Write())))
	return ctx
}
