package scale

import (
	"github.com/loopholelabs/scale-go/context"
)

func Scale(ctx *context.Context) *context.Context {
	req := ctx.Request()
	res := ctx.Response()
	res.SetBodyBytes(req.Body())
	return ctx
}
