package scale

import (
	"github.com/loopholelabs/scale-go/context"
)

func Scale(ctx *context.Context) *context.Context {
	ctx.Response().SetBody("Hello, World!")
	return ctx
}
