package scale

import (
	"github.com/loopholelabs/scale-go/context"
)

func Scale(ctx *context.Context) *context.Context {
	res := ctx.Response()
	res.Headers().Set("X-Test", []string{"test"})
	return ctx.Next()
}
