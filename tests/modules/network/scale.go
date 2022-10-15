package scale

import (
	"github.com/loopholelabs/scale-go/context"
	"net"
)

func Scale(ctx *context.Context) *context.Context {
	net.Dial("tcp", "0.0.0.0:80")
	return ctx.Next()
}
