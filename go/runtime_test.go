//go:build !tinygo && !js && !wasm
// +build !tinygo,!js,!wasm

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

package runtime

import (
	"context"
	"errors"
	"github.com/loopholelabs/scale/go/tests/harness"
	signature "github.com/loopholelabs/scale/go/tests/signature"
	"github.com/loopholelabs/scalefile/scalefunc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type TestCase struct {
	Name   string
	Module *harness.Module
	Run    func(*scalefunc.ScaleFunc, *testing.T)
}

func TestRuntime(t *testing.T) {
	passthroughModule := &harness.Module{
		Name:      "passthrough",
		Path:      "tests/modules/passthrough/passthrough.go",
		Signature: "github.com/loopholelabs/scale/go/tests/signature",
	}

	nextModule := &harness.Module{
		Name:      "next",
		Path:      "tests/modules/next/next.go",
		Signature: "github.com/loopholelabs/scale/go/tests/signature",
	}

	fileModule := &harness.Module{
		Name:      "file",
		Path:      "tests/modules/file/file.go",
		Signature: "github.com/loopholelabs/scale/go/tests/signature",
	}

	networkModule := &harness.Module{
		Name:      "network",
		Path:      "tests/modules/network/network.go",
		Signature: "github.com/loopholelabs/scale/go/tests/signature",
	}

	panicModule := &harness.Module{
		Name:      "panic",
		Path:      "tests/modules/panic/panic.go",
		Signature: "github.com/loopholelabs/scale/go/tests/signature",
	}

	modules := []*harness.Module{passthroughModule, nextModule, fileModule, networkModule, panicModule}

	generatedModules := harness.Setup(t, modules, "github.com/loopholelabs/scale/go/tests/modules")

	var testCases = []TestCase{
		{
			Name:   "Passthrough",
			Module: passthroughModule,
			Run: func(scaleFunc *scalefunc.ScaleFunc, t *testing.T) {
				r, err := New(context.Background(), signature.New(), []*scalefunc.ScaleFunc{scaleFunc})
				require.NoError(t, err)

				i, err := r.Instance(nil)
				require.NoError(t, err)

				i.Context().Data = "Test Data"

				err = i.Run(context.Background())
				assert.NoError(t, err)

				assert.Equal(t, "Test Data", i.Context().Data)
			},
		},
		{
			Name:   "Next",
			Module: nextModule,
			Run: func(scaleFunc *scalefunc.ScaleFunc, t *testing.T) {
				next := func(ctx *signature.Context) (*signature.Context, error) {
					ctx.Data = "Hello, World!"
					return ctx, nil
				}

				r, err := New(context.Background(), signature.New(), []*scalefunc.ScaleFunc{scaleFunc})
				require.NoError(t, err)

				i, err := r.Instance(next)
				require.NoError(t, err)

				i.Context().Data = "Test Data"

				err = i.Run(context.Background())
				assert.NoError(t, err)

				assert.Equal(t, "Hello, World!", i.Context().Data)
			},
		},
		{
			Name:   "NextError",
			Module: nextModule,
			Run: func(scaleFunc *scalefunc.ScaleFunc, t *testing.T) {
				next := func(ctx *signature.Context) (*signature.Context, error) {
					return nil, errors.New("next error")
				}

				r, err := New(context.Background(), signature.New(), []*scalefunc.ScaleFunc{scaleFunc})
				require.NoError(t, err)

				i, err := r.Instance(next)
				require.NoError(t, err)

				i.Context().Data = "Test Data"

				err = i.Run(context.Background())
				require.ErrorIs(t, err, errors.New("next error"))
			},
		},
		{
			Name:   "File",
			Module: fileModule,
			Run: func(scaleFunc *scalefunc.ScaleFunc, t *testing.T) {
				r, err := New(context.Background(), signature.New(), []*scalefunc.ScaleFunc{scaleFunc})
				require.NoError(t, err)

				i, err := r.Instance(nil)
				require.NoError(t, err)

				i.Context().Data = "Test Data"

				err = i.Run(context.Background())
				require.Error(t, err)
			},
		},
		{
			Name:   "Network",
			Module: networkModule,
			Run: func(scaleFunc *scalefunc.ScaleFunc, t *testing.T) {
				r, err := New(context.Background(), signature.New(), []*scalefunc.ScaleFunc{scaleFunc})
				require.NoError(t, err)

				i, err := r.Instance(nil)
				require.NoError(t, err)

				i.Context().Data = "Test Data"

				err = i.Run(context.Background())
				require.Error(t, err)
			},
		},
		{
			Name:   "Panic",
			Module: panicModule,
			Run: func(scaleFunc *scalefunc.ScaleFunc, t *testing.T) {
				r, err := New(context.Background(), signature.New(), []*scalefunc.ScaleFunc{scaleFunc})
				require.NoError(t, err)

				i, err := r.Instance(nil)
				require.NoError(t, err)

				i.Context().Data = "Test Data"

				err = i.Run(context.Background())
				require.Error(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {

			module, err := os.ReadFile(generatedModules[testCase.Module])
			require.NoError(t, err)

			scaleFunc := &scalefunc.ScaleFunc{
				Version:   "TestVersion",
				Name:      "TestName",
				Signature: "ExampleName@ExampleVersion",
				Language:  "go",
				Function:  module,
			}
			testCase.Run(scaleFunc, t)
		})
	}
}
