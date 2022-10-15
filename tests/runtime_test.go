package tests

import (
	"context"
	adapter "github.com/loopholelabs/scale-go/adapters/http"
	"github.com/loopholelabs/scale-go/runtime"
	"github.com/loopholelabs/scale-go/scalefile"
	"github.com/loopholelabs/scale-go/scalefunc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

func TestRuntime(t *testing.T) {

	wasm, err := os.ReadFile("example.wasm")
	assert.NoError(t, err)

	scaleFunc := &scalefunc.ScaleFunc{
		ScaleFile: scalefile.ScaleFile{
			Name: "example",
			Build: scalefile.Build{
				Language: "go",
			},
		},
		Function: wasm,
	}
	next := func(ctx *runtime.Context) *runtime.Context {
		ctx.Context.Request.Method = "POST"
		return ctx
	}

	r, err := runtime.New(context.Background(), next, []scalefunc.ScaleFunc{*scaleFunc})
	require.NoError(t, err)

	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	assert.NoError(t, err)

	ctx := runtime.NewContext()
	adapter.Serialize(ctx, req)

	assert.Equal(t, "GET", ctx.Context.Request.Method)

	err = r.Run(ctx)
	assert.NoError(t, err)

	assert.Equal(t, "POST", ctx.Context.Request.Method)
}
