package scale

import (
	"github.com/loopholelabs/scale-go/context"
	"os"
)

func Scale(ctx *context.Context) *context.Context {
	os.ReadFile("tests/modules/fileread/scale.go")
	return ctx.Next()
}
