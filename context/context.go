package context

import (
	"github.com/loopholelabs/polyglot-go"
	"github.com/loopholelabs/scale-go/runtime/generated"
	"reflect"
	"unsafe"
)

type Context struct {
	generated *generated.Context
	buffer    *polyglot.Buffer
}

func NewContext(ctx *generated.Context) *Context {
	return &Context{
		generated: ctx,
		buffer:    polyglot.NewBuffer(),
	}
}

func (ctx *Context) serialize() (uint32, uint32) {
	ctx.buffer.Reset()
	ctx.generated.Encode(ctx.buffer)
	underlying := ctx.buffer.Bytes()
	ptr := &underlying[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(len(*ctx.buffer))
}

func (ctx *Context) deserialize(ptr uint32, size uint32) {
	buf := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  uintptr(size), // Tinygo requires this, it's not an error.
		Cap:  uintptr(size), // ^^ See https://github.com/tinygo-org/tinygo/issues/1284
	}))
	_ = ctx.generated.Decode(buf)
}

func (ctx *Context) Next() {
	next(ctx.serialize())
	ctx.deserialize(globalOffset, globalLength)
}
