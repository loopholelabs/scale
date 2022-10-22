/*
	Copyright 2022 Loophole Labs

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		   http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package tests

import (
	"context"
	"fmt"
	adapter "github.com/loopholelabs/scale/go/adapters/http"
	"github.com/loopholelabs/scale/go/runtime"
	"github.com/loopholelabs/scale/go/runtime/generated"
	"github.com/loopholelabs/scale/go/scalefile"
	"github.com/loopholelabs/scale/go/scalefunc"
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

				i, err := r.Instance(next)
				require.NoError(t, err)

				req, err := http.NewRequest("GET", "http://localhost:8080", nil)
				assert.NoError(t, err)

				err = adapter.FromRequest(i.Context(), req)
				assert.NoError(t, err)
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

				i, err := r.Instance(next)
				require.NoError(t, err)

				req, err := http.NewRequest("GET", "http://localhost:8080", nil)
				assert.NoError(t, err)

				err = adapter.FromRequest(i.Context(), req)
				assert.NoError(t, err)
				assert.Equal(t, "GET", i.Context().Context.Request.Method)

				err = i.Run(context.Background())
				assert.NoError(t, err)

				assert.Equal(t, "POST", i.Context().Context.Request.Method)
				assert.Equal(t, &generated.StringList{Value: []string{"TRUE"}}, i.Context().Context.Response.Headers["MIDDLEWARE"])
			},
		},
		{
			Name:   "File Read",
			Module: "fileread",
			Run: func(scaleFunc scalefunc.ScaleFunc, t *testing.T) {
				next := func(ctx *runtime.Context) *runtime.Context {
					t.Fatal("next should not be called")
					return ctx
				}

				r, err := runtime.New(context.Background(), []scalefunc.ScaleFunc{scaleFunc})
				require.NoError(t, err)

				i, err := r.Instance(next)
				require.NoError(t, err)

				req, err := http.NewRequest("GET", "http://localhost:8080", nil)
				assert.NoError(t, err)

				err = adapter.FromRequest(i.Context(), req)
				assert.NoError(t, err)

				err = i.Run(context.Background())
				assert.NoError(t, err)

				assert.Equal(t, "error reading file", string(i.Context().Context.Response.Body))
			},
		},
		{
			Name:   "Network",
			Module: "network",
			Run: func(scaleFunc scalefunc.ScaleFunc, t *testing.T) {
				next := func(ctx *runtime.Context) *runtime.Context {
					t.Fatal("next should not be called")
					return ctx
				}

				r, err := runtime.New(context.Background(), []scalefunc.ScaleFunc{scaleFunc})
				require.NoError(t, err)

				i, err := r.Instance(next)
				require.NoError(t, err)

				req, err := http.NewRequest("GET", "http://localhost:8080", nil)
				assert.NoError(t, err)

				err = adapter.FromRequest(i.Context(), req)
				assert.NoError(t, err)

				err = i.Run(context.Background())
				assert.NoError(t, err)

				assert.Equal(t, "error opening connection", string(i.Context().Context.Response.Body))
			},
		},
		{
			Name:   "Panic",
			Module: "panic",
			Run: func(scaleFunc scalefunc.ScaleFunc, t *testing.T) {
				r, err := runtime.New(context.Background(), []scalefunc.ScaleFunc{scaleFunc})
				require.NoError(t, err)

				i, err := r.Instance(func(ctx *runtime.Context) *runtime.Context {
					return ctx
				})
				require.NoError(t, err)

				req, err := http.NewRequest("GET", "http://localhost:8080", nil)
				assert.NoError(t, err)
				err = adapter.FromRequest(i.Context(), req)
				assert.NoError(t, err)
				assert.Equal(t, "GET", i.Context().Context.Request.Method)

				err = i.Run(context.Background())
				assert.Error(t, err)
			},
		},
		{
			Name:   "Next Function Required",
			Module: "passthrough",
			Run: func(scaleFunc scalefunc.ScaleFunc, t *testing.T) {
				r, err := runtime.New(context.Background(), []scalefunc.ScaleFunc{scaleFunc})
				require.NoError(t, err)

				_, err = r.Instance(nil)
				require.ErrorIs(t, err, runtime.NextFunctionRequiredError)
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
					Middleware: true,
				},
				Function: module,
			}
			testCase.Run(scaleFunc, t)
		})
	}
}
