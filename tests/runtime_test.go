package tests

import (
	"context"
	"fmt"
	adapter "github.com/loopholelabs/scale-go/adapters/http"
	"github.com/loopholelabs/scale-go/runtime"
	"github.com/loopholelabs/scale-go/runtime/generated"
	"github.com/loopholelabs/scale-go/scalefile"
	"github.com/loopholelabs/scale-go/scalefunc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"path"
	"testing"
)

func TestRuntime(t *testing.T) {
	var testCases = []TestCase{
		{
			Name:   "Passthrough",
			Module: "passthrough",
			Run: func(scaleFunc scalefunc.ScaleFunc, t *testing.T) {
				next := func(ctx *runtime.Context) *runtime.Context {
					ctx.Context.Request.Method = "POST"
					return ctx
				}

				r, err := runtime.New(context.Background(), []scalefunc.ScaleFunc{scaleFunc})
				require.NoError(t, err)

				i, err := r.Instance(context.Background(), next)
				require.NoError(t, err)

				req, err := http.NewRequest("GET", "http://localhost:8080", nil)
				assert.NoError(t, err)

				adapter.SerializeRequest(i.Context(), req)
				assert.Equal(t, "GET", i.Context().Context.Request.Method)

				err = i.Run(context.Background())
				assert.NoError(t, err)

				assert.Equal(t, "POST", i.Context().Context.Request.Method)
			},
		},
		{
			Name:   "HTTP Middleware",
			Module: "http-middleware",
			Run: func(scaleFunc scalefunc.ScaleFunc, t *testing.T) {
				next := func(ctx *runtime.Context) *runtime.Context {
					ctx.Context.Request.Method = "POST"
					return ctx
				}

				r, err := runtime.New(context.Background(), []scalefunc.ScaleFunc{scaleFunc})
				require.NoError(t, err)

				i, err := r.Instance(context.Background(), next)
				require.NoError(t, err)

				req, err := http.NewRequest("GET", "http://localhost:8080", nil)
				assert.NoError(t, err)

				adapter.SerializeRequest(i.Context(), req)
				assert.Equal(t, "GET", i.Context().Context.Request.Method)

				err = i.Run(context.Background())
				assert.NoError(t, err)

				assert.Equal(t, "POST", i.Context().Context.Request.Method)
				assert.Equal(t, &generated.StringList{Value: []string{"test"}}, i.Context().Context.Response.Headers["X-Test"])
			},
		},
		{
			Name:   "File Read",
			Module: "fileread",
			Run: func(scaleFunc scalefunc.ScaleFunc, t *testing.T) {
				r, err := runtime.New(context.Background(), []scalefunc.ScaleFunc{scaleFunc})
				require.NoError(t, err)

				_, err = r.Instance(context.Background(), nil)
				require.Error(t, err)
			},
		},
		{
			Name:   "Network",
			Module: "network",
			Run: func(scaleFunc scalefunc.ScaleFunc, t *testing.T) {
				r, err := runtime.New(context.Background(), []scalefunc.ScaleFunc{scaleFunc})
				require.NoError(t, err)

				_, err = r.Instance(context.Background(), nil)
				require.Error(t, err)
			},
		},
		{
			Name:   "Panic",
			Module: "panic",
			Run: func(scaleFunc scalefunc.ScaleFunc, t *testing.T) {
				r, err := runtime.New(context.Background(), []scalefunc.ScaleFunc{scaleFunc})
				require.NoError(t, err)

				i, err := r.Instance(context.Background(), nil)
				require.NoError(t, err)

				req, err := http.NewRequest("GET", "http://localhost:8080", nil)
				assert.NoError(t, err)
				adapter.SerializeRequest(i.Context(), req)
				assert.Equal(t, "GET", i.Context().Context.Request.Method)

				err = i.Run(context.Background())
				assert.Error(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			module, err := os.ReadFile(path.Join("modules", fmt.Sprintf("%s.wasm", testCase.Module)))
			assert.NoError(t, err)

			scaleFunc := scalefunc.ScaleFunc{
				ScaleFile: scalefile.ScaleFile{
					Name: testCase.Name,
					Build: scalefile.Build{
						Language: "go",
					},
				},
				Function: module,
			}
			testCase.Run(scaleFunc, t)
		})
	}
}
