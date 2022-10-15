package tests

import (
	"context"
	"fmt"
	adapter "github.com/loopholelabs/scale-go/adapters/http"
	"github.com/loopholelabs/scale-go/runtime"
	"github.com/loopholelabs/scale-go/scalefile"
	"github.com/loopholelabs/scale-go/scalefunc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"os/exec"
	"path"
	"testing"
)

type TestCase struct {
	Name   string
	Module string
	Run    func(scalefunc.ScaleFunc, *testing.T)
}

var TestCases = []TestCase{
	{
		Name:   "Passthrough",
		Module: "passthrough",
		Run: func(scaleFunc scalefunc.ScaleFunc, t *testing.T) {
			next := func(ctx *runtime.Context) *runtime.Context {
				ctx.Context.Request.Method = "POST"
				return ctx
			}

			r, err := runtime.New(context.Background(), next, []scalefunc.ScaleFunc{scaleFunc})
			require.NoError(t, err)

			i, err := r.Instance(context.Background())
			require.NoError(t, err)

			req, err := http.NewRequest("GET", "http://localhost:8080", nil)
			assert.NoError(t, err)

			adapter.Serialize(i.Context(), req)
			assert.Equal(t, "GET", i.Context().Context.Request.Method)

			err = i.Run(context.Background())
			assert.NoError(t, err)

			assert.Equal(t, "POST", i.Context().Context.Request.Method)
		},
	},
	{
		Name:   "File Read",
		Module: "fileread",
		Run: func(scaleFunc scalefunc.ScaleFunc, t *testing.T) {
			r, err := runtime.New(context.Background(), nil, []scalefunc.ScaleFunc{scaleFunc})
			require.NoError(t, err)

			_, err = r.Instance(context.Background())
			require.Error(t, err)
		},
	},
	{
		Name:   "Network",
		Module: "network",
		Run: func(scaleFunc scalefunc.ScaleFunc, t *testing.T) {
			r, err := runtime.New(context.Background(), nil, []scalefunc.ScaleFunc{scaleFunc})
			require.NoError(t, err)

			_, err = r.Instance(context.Background())
			require.Error(t, err)
		},
	},
	{
		Name:   "Panic",
		Module: "panic",
		Run: func(scaleFunc scalefunc.ScaleFunc, t *testing.T) {
			r, err := runtime.New(context.Background(), nil, []scalefunc.ScaleFunc{scaleFunc})
			require.NoError(t, err)

			i, err := r.Instance(context.Background())
			require.NoError(t, err)

			req, err := http.NewRequest("GET", "http://localhost:8080", nil)
			assert.NoError(t, err)
			adapter.Serialize(i.Context(), req)
			assert.Equal(t, "GET", i.Context().Context.Request.Method)

			err = i.Run(context.Background())
			assert.Error(t, err)
		},
	},
}

func TestRuntime(t *testing.T) {
	err := exec.Command("sh", "compile.sh").Run()
	require.NoError(t, err)

	for _, testCase := range TestCases {
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

	err = exec.Command("sh", "cleanup.sh").Run()
	assert.NoError(t, err)
}
