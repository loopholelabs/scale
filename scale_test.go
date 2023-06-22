//go:build !tinygo && !js && !wasm

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

package scale

//import (
//	"context"
//	"encoding/json"
//	"errors"
//	"fmt"
//	"os"
//	"testing"
//	"time"
//
//	httpSignature "github.com/loopholelabs/scale-signature-http"
//	"github.com/loopholelabs/scale/go/tests/harness"
//	signature "github.com/loopholelabs/scale/go/tests/signature/example-signature"
//	"github.com/loopholelabs/scalefile"
//	"github.com/loopholelabs/scalefile/scalefunc"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//)
//
//type TestCase struct {
//	Name   string
//	Module *harness.Module
//	Run    func(*scalefunc.Schema, *testing.T)
//}
//
//func TestRuntimeGo(t *testing.T) {
//	passthroughModule := &harness.Module{
//		Name:      "passthrough",
//		Path:      "tests/modules/passthrough/passthrough.go",
//		Signature: "github.com/loopholelabs/scale/go/tests/signature/example-signature",
//	}
//
//	modifyModule := &harness.Module{
//		Name:      "modify",
//		Path:      "tests/modules/modify/modify.go",
//		Signature: "github.com/loopholelabs/scale/go/tests/signature/example-signature",
//	}
//
//	nextModule := &harness.Module{
//		Name:      "next",
//		Path:      "tests/modules/next/next.go",
//		Signature: "github.com/loopholelabs/scale/go/tests/signature/example-signature",
//	}
//
//	modifyNextModule := &harness.Module{
//		Name:      "modifynext",
//		Path:      "tests/modules/modifynext/modifynext.go",
//		Signature: "github.com/loopholelabs/scale/go/tests/signature/example-signature",
//	}
//
//	fileModule := &harness.Module{
//		Name:      "file",
//		Path:      "tests/modules/file/file.go",
//		Signature: "github.com/loopholelabs/scale/go/tests/signature/example-signature",
//	}
//
//	networkModule := &harness.Module{
//		Name:      "network",
//		Path:      "tests/modules/network/network.go",
//		Signature: "github.com/loopholelabs/scale/go/tests/signature/example-signature",
//	}
//
//	panicModule := &harness.Module{
//		Name:      "panic",
//		Path:      "tests/modules/panic/panic.go",
//		Signature: "github.com/loopholelabs/scale/go/tests/signature/example-signature",
//	}
//
//	badSignatureModule := &harness.Module{
//		Name:      "bad-signature",
//		Path:      "tests/modules/bad-signature/bad-signature.go",
//		Signature: "github.com/loopholelabs/scale/go/tests/signature/bad-signature",
//	}
//
//	modules := []*harness.Module{passthroughModule, modifyModule, nextModule, modifyNextModule, fileModule, networkModule, panicModule, badSignatureModule}
//
//	generatedModules := harness.GoSetup(t, modules, "github.com/loopholelabs/scale/go/tests/modules")
//
//	var testCases = []TestCase{
//		{
//			Name:   "Passthrough",
//			Module: passthroughModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance()
//				require.NoError(t, err)
//
//				i.Context().Data = "Test Data"
//
//				err = i.Run(context.Background())
//				assert.NoError(t, err)
//
//				assert.Equal(t, "Test Data", i.Context().Data)
//			},
//		},
//		{
//			Name:   "Modify",
//			Module: modifyModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance()
//				require.NoError(t, err)
//
//				i.Context().Data = "Test Data"
//
//				err = i.Run(context.Background())
//				assert.NoError(t, err)
//
//				assert.Equal(t, "modified", i.Context().Data)
//			},
//		},
//		{
//			Name:   "Next",
//			Module: nextModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				next := func(ctx *signature.Context) (*signature.Context, error) {
//					ctx.Data = "Hello, World!"
//					return ctx, nil
//				}
//
//				r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance(next)
//				require.NoError(t, err)
//
//				i.Context().Data = "Test Data"
//
//				err = i.Run(context.Background())
//				assert.NoError(t, err)
//
//				assert.Equal(t, "Hello, World!", i.Context().Data)
//			},
//		},
//		{
//			Name:   "ModifyNext",
//			Module: modifyNextModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				next := func(ctx *signature.Context) (*signature.Context, error) {
//					ctx.Data = ctx.Data + "-next"
//					return ctx, nil
//				}
//
//				r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance(next)
//				require.NoError(t, err)
//
//				i.Context().Data = "Test Data"
//
//				err = i.Run(context.Background())
//				assert.NoError(t, err)
//
//				assert.Equal(t, "modified-next", i.Context().Data)
//			},
//		},
//		{
//			Name:   "NextError",
//			Module: nextModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				next := func(ctx *signature.Context) (*signature.Context, error) {
//					return nil, errors.New("next error")
//				}
//
//				r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance(next)
//				require.NoError(t, err)
//
//				err = i.Run(context.Background())
//				require.ErrorIs(t, err, errors.New("next error"))
//			},
//		},
//		{
//			Name:   "File",
//			Module: fileModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance(nil)
//				require.NoError(t, err)
//
//				err = i.Run(context.Background())
//				require.Error(t, err)
//			},
//		},
//		{
//			Name:   "Network",
//			Module: networkModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance(nil)
//				require.NoError(t, err)
//
//				err = i.Run(context.Background())
//				require.Error(t, err)
//			},
//		},
//		{
//			Name:   "Panic",
//			Module: panicModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance(nil)
//				require.NoError(t, err)
//
//				err = i.Run(context.Background())
//				require.Error(t, err)
//			},
//		},
//		{
//			Name:   "BadSignature",
//			Module: badSignatureModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance(nil)
//				require.NoError(t, err)
//
//				err = i.Run(context.Background())
//				assert.Error(t, err)
//			},
//		},
//	}
//
//	for _, testCase := range testCases {
//		t.Run(testCase.Name, func(t *testing.T) {
//			module, err := os.ReadFile(generatedModules[testCase.Module])
//			require.NoError(t, err)
//
//			scaleFunc := &scalefunc.Schema{
//				Version:   scalefunc.V1Alpha,
//				Name:      "TestName",
//				Tag:       "TestTag",
//				Signature: "ExampleName@ExampleVersion",
//				Language:  scalefunc.Go,
//				Function:  module,
//			}
//			testCase.Run(scaleFunc, t)
//		})
//	}
//}
//
//func TestRuntimeGoTracing(t *testing.T) {
//	tracingModule := &harness.Module{
//		Name:      "tracing",
//		Path:      "tests/modules/tracing/tracing.go",
//		Signature: "github.com/loopholelabs/scale/go/tests/signature/example-signature",
//	}
//
//	modules := []*harness.Module{tracingModule}
//	generatedModules := harness.GoSetup(t, modules, "github.com/loopholelabs/scale/go/tests/modules")
//
//	// Now run it and see...
//	module, err := os.ReadFile(generatedModules[tracingModule])
//	require.NoError(t, err)
//
//	scaleFunc1 := &scalefunc.Schema{
//		Version:   scalefunc.V1Alpha,
//		Name:      "TestName1",
//		Tag:       "TestTag1",
//		Signature: "ExampleName@ExampleVersion",
//		Language:  scalefunc.Go,
//		Function:  module,
//	}
//
//	scaleFunc2 := &scalefunc.Schema{
//		Version:   scalefunc.V1Alpha,
//		Name:      "TestName2",
//		Tag:       "TestTag2",
//		Signature: "ExampleName@ExampleVersion",
//		Language:  scalefunc.Go,
//		Function:  module,
//	}
//
//	r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc1, scaleFunc2})
//	require.NoError(t, err)
//
//	traceData := make([]string, 0)
//
//	r.TraceDataCallback = func(data string) {
//		traceData = append(traceData, data)
//
//		// Slightly hacky, but we want to make sure time has moved on
//		t1 := time.Now().UnixNano()
//		t2 := t1
//		for {
//			if t2 > t1 {
//				break
//			}
//			t2 = time.Now().UnixNano()
//		}
//
//	}
//
//	i, err := r.Instance(nil)
//	require.NoError(t, err)
//
//	err = i.Run(context.Background())
//	assert.NoError(t, err)
//
//	// Assert that the trace data is there and as expected...
//	assert.Equal(t, 2, len(traceData))
//
//	data1 := traceData[0]
//	trace1 := &TraceData{}
//	err = json.Unmarshal([]byte(data1), trace1)
//	require.NoError(t, err)
//	// Json decode, and check it
//
//	assert.Equal(t, "TestName1:TestTag1", trace1.ServiceName)
//	assert.Equal(t, fmt.Sprintf("%x", r.InvocationID), trace1.InvocationId)
//
//	data2 := traceData[1]
//	trace2 := &TraceData{}
//	err = json.Unmarshal([]byte(data2), trace2)
//	require.NoError(t, err)
//	// Json decode, and check it
//
//	assert.Equal(t, "TestName2:TestTag2", trace2.ServiceName)
//	assert.Equal(t, fmt.Sprintf("%x", r.InvocationID), trace2.InvocationId)
//
//	// Make sure we can get time
//	assert.True(t, trace2.Timestamp > trace1.Timestamp)
//
//	// Run it again and make sure it's a new invocationID
//	i, err = r.Instance(nil)
//	require.NoError(t, err)
//
//	err = i.Run(context.Background())
//	assert.NoError(t, err)
//
//	// Assert that the trace data is there and as expected but a different invocationID from previous Run
//	assert.Equal(t, 4, len(traceData))
//
//	data3 := traceData[2]
//	trace3 := &TraceData{}
//	err = json.Unmarshal([]byte(data3), trace3)
//	require.NoError(t, err)
//	// Json decode, and check it
//
//	assert.NotEqual(t, trace3.InvocationId, trace2.InvocationId)
//}
//
//type TraceData struct {
//	ServiceName  string `json:"serviceName"`
//	InvocationId string `json:"invocationId"`
//	Timestamp    uint64 `json:"timestamp"`
//}
//
//func TestRuntimeHTTPSignatureGo(t *testing.T) {
//	passthroughModule := &harness.Module{
//		Name:      "http-passthrough",
//		Path:      "tests/modules/http-passthrough/http-passthrough.go",
//		Signature: "github.com/loopholelabs/scale-signature-http",
//	}
//
//	handlerModule := &harness.Module{
//		Name:      "http-handler",
//		Path:      "tests/modules/http-handler/http-handler.go",
//		Signature: "github.com/loopholelabs/scale-signature-http",
//	}
//
//	nextModule := &harness.Module{
//		Name:      "http-next",
//		Path:      "tests/modules/http-next/http-next.go",
//		Signature: "github.com/loopholelabs/scale-signature-http",
//	}
//
//	modules := []*harness.Module{passthroughModule, handlerModule, nextModule}
//
//	generatedModules := harness.GoSetup(t, modules, "github.com/loopholelabs/scale/go/tests/modules")
//
//	var testCases = []TestCase{
//		{
//			Name:   "Passthrough",
//			Module: passthroughModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				r, err := New(context.Background(), []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance()
//				require.NoError(t, err)
//
//				i.Context().Response.Body = []byte("Test Data")
//				err = i.Run(context.Background())
//				assert.NoError(t, err)
//
//				assert.Equal(t, "Test Data", string(i.Context().Response.Body))
//			},
//		},
//		{
//			Name:   "Handler",
//			Module: handlerModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				r, err := New(context.Background(), []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance()
//				require.NoError(t, err)
//
//				i.Context().Response.Body = []byte("Test Data")
//				err = i.Run(context.Background())
//				assert.NoError(t, err)
//
//				assert.Equal(t, "Test Data-modified", string(i.Context().Response.Body))
//			},
//		},
//		{
//			Name:   "Next",
//			Module: nextModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				r, err := New(context.Background(), []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance(func(ctx *httpSignature.Context) (*httpSignature.Context, error) {
//					ctx.Response.Body = append(ctx.Response.Body, []byte("-next")...)
//					return ctx, nil
//				})
//				require.NoError(t, err)
//
//				i.Context().Response.Body = []byte("Test Data")
//				err = i.Run(context.Background())
//				assert.NoError(t, err)
//
//				assert.Equal(t, "Test Data-modified-next", string(i.Context().Response.Body))
//			},
//		},
//	}
//
//	for _, testCase := range testCases {
//		t.Run(testCase.Name, func(t *testing.T) {
//			module, err := os.ReadFile(generatedModules[testCase.Module])
//			require.NoError(t, err)
//
//			scaleFunc := &scalefunc.Schema{
//				Version:   scalefunc.V1Alpha,
//				Name:      "TestName",
//				Tag:       "TestTag",
//				Signature: "ExampleName@ExampleVersion",
//				Language:  scalefunc.Go,
//				Function:  module,
//			}
//			testCase.Run(scaleFunc, t)
//		})
//	}
//}
//
//func TestRuntimeRs(t *testing.T) {
//	passthroughModule := &harness.Module{
//		Name:          "passthrough",
//		Path:          "../rust/tests/modules/passthrough/passthrough.rs",
//		Signature:     "example_signature",
//		SignaturePath: "../../../signature/example-signature",
//	}
//
//	modifyModule := &harness.Module{
//		Name:          "modify",
//		Path:          "../rust/tests/modules/modify/modify.rs",
//		Signature:     "example_signature",
//		SignaturePath: "../../../signature/example-signature",
//	}
//
//	nextModule := &harness.Module{
//		Name:          "next",
//		Path:          "../rust/tests/modules/next/next.rs",
//		Signature:     "example_signature",
//		SignaturePath: "../../../signature/example-signature",
//	}
//
//	modifyNextModule := &harness.Module{
//		Name:          "modifynext",
//		Path:          "../rust/tests/modules/modifynext/modifynext.rs",
//		Signature:     "example_signature",
//		SignaturePath: "../../../signature/example-signature",
//	}
//
//	fileModule := &harness.Module{
//		Name:          "file",
//		Path:          "../rust/tests/modules/file/file.rs",
//		Signature:     "example_signature",
//		SignaturePath: "../../../signature/example-signature",
//	}
//
//	networkModule := &harness.Module{
//		Name:          "network",
//		Path:          "../rust/tests/modules/network/network.rs",
//		Signature:     "example_signature",
//		SignaturePath: "../../../signature/example-signature",
//	}
//
//	panicModule := &harness.Module{
//		Name:          "panic",
//		Path:          "../rust/tests/modules/panic/panic.rs",
//		Signature:     "example_signature",
//		SignaturePath: "../../../signature/example-signature",
//	}
//
//	badSignatureModule := &harness.Module{
//		Name:          "bad_signature",
//		Path:          "../rust/tests/modules/bad_signature/bad_signature.rs",
//		Signature:     "bad_signature",
//		SignaturePath: "../../../signature/bad-signature",
//	}
//
//	modules := []*harness.Module{passthroughModule, modifyModule, nextModule, modifyNextModule, fileModule, networkModule, panicModule, badSignatureModule}
//
//	dependencies := []*scalefile.Dependency{
//		{
//			Name:    "scale_signature",
//			Version: "0.2.11",
//		},
//	}
//
//	generatedModules := harness.RustSetup(t, modules, dependencies)
//
//	var testCases = []TestCase{
//		{
//			Name:   "Passthrough",
//			Module: passthroughModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance()
//				require.NoError(t, err)
//
//				i.Context().Data = "Test Data"
//
//				err = i.Run(context.Background())
//				assert.NoError(t, err)
//
//				assert.Equal(t, "Test Data", i.Context().Data)
//			},
//		},
//		{
//			Name:   "Modify",
//			Module: modifyModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance()
//				require.NoError(t, err)
//
//				i.Context().Data = "Test Data"
//
//				err = i.Run(context.Background())
//				assert.NoError(t, err)
//
//				assert.Equal(t, "modified", i.Context().Data)
//			},
//		},
//		{
//			Name:   "Next",
//			Module: nextModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				next := func(ctx *signature.Context) (*signature.Context, error) {
//					ctx.Data = "Hello, World!"
//					return ctx, nil
//				}
//
//				r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance(next)
//				require.NoError(t, err)
//
//				i.Context().Data = "Test Data"
//
//				err = i.Run(context.Background())
//				assert.NoError(t, err)
//
//				assert.Equal(t, "Hello, World!", i.Context().Data)
//			},
//		},
//		{
//			Name:   "ModifyNext",
//			Module: modifyNextModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				next := func(ctx *signature.Context) (*signature.Context, error) {
//					ctx.Data = ctx.Data + "-next"
//					return ctx, nil
//				}
//
//				r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance(next)
//				require.NoError(t, err)
//
//				i.Context().Data = "Test Data"
//
//				err = i.Run(context.Background())
//				assert.NoError(t, err)
//
//				assert.Equal(t, "modified-next", i.Context().Data)
//			},
//		},
//		{
//			Name:   "NextError",
//			Module: nextModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				next := func(ctx *signature.Context) (*signature.Context, error) {
//					return nil, errors.New("next error")
//				}
//
//				r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance(next)
//				require.NoError(t, err)
//
//				err = i.Run(context.Background())
//				require.ErrorIs(t, err, errors.New("next error"))
//			},
//		},
//		{
//			Name:   "File",
//			Module: fileModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance(nil)
//				require.NoError(t, err)
//
//				err = i.Run(context.Background())
//				require.ErrorContains(t, err, "operation not supported on this platform")
//			},
//		},
//		{
//			Name:   "Network",
//			Module: networkModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance(nil)
//				require.NoError(t, err)
//
//				err = i.Run(context.Background())
//				require.ErrorContains(t, err, "operation not supported on this platform")
//			},
//		},
//		{
//			Name:   "Panic",
//			Module: panicModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance(nil)
//				require.NoError(t, err)
//
//				err = i.Run(context.Background())
//				require.Error(t, err)
//			},
//		},
//		{
//			Name:   "BadSignature",
//			Module: badSignatureModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				r, err := NewWithSignature(context.Background(), signature.New, []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance(nil)
//				require.NoError(t, err)
//
//				err = i.Run(context.Background())
//				assert.Error(t, err)
//			},
//		},
//	}
//
//	for _, testCase := range testCases {
//		t.Run(testCase.Name, func(t *testing.T) {
//			module, err := os.ReadFile(generatedModules[testCase.Module])
//			require.NoError(t, err)
//
//			scaleFunc := &scalefunc.Schema{
//				Version:   scalefunc.V1Alpha,
//				Name:      "TestName",
//				Tag:       "TestTag",
//				Signature: "ExampleName@ExampleVersion",
//				Language:  scalefunc.Rust,
//				Function:  module,
//			}
//			testCase.Run(scaleFunc, t)
//		})
//	}
//}
//
//func TestRuntimeHTTPSignatureRs(t *testing.T) {
//	passthroughModule := &harness.Module{
//		Name:      "http_passthrough",
//		Path:      "../rust/tests/modules/http_passthrough/http_passthrough.rs",
//		Signature: "scale_signature_http",
//	}
//
//	handlerModule := &harness.Module{
//		Name:      "http_handler",
//		Path:      "../rust/tests/modules/http_handler/http_handler.rs",
//		Signature: "scale_signature_http",
//	}
//
//	nextModule := &harness.Module{
//		Name:      "http_next",
//		Path:      "../rust/tests/modules/http_next/http_next.rs",
//		Signature: "scale_signature_http",
//	}
//
//	modules := []*harness.Module{passthroughModule, handlerModule, nextModule}
//
//	dependencies := []*scalefile.Dependency{
//		{
//			Name:    "scale_signature",
//			Version: "0.2.11",
//		},
//		{
//			Name:    "scale_signature_http",
//			Version: "0.3.8",
//		},
//	}
//
//	generatedModules := harness.RustSetup(t, modules, dependencies)
//
//	var testCases = []TestCase{
//		{
//			Name:   "Passthrough",
//			Module: passthroughModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				r, err := New(context.Background(), []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance()
//				require.NoError(t, err)
//
//				i.Context().Response.Body = []byte("Test Data")
//				err = i.Run(context.Background())
//				assert.NoError(t, err)
//
//				assert.Equal(t, "Test Data", string(i.Context().Response.Body))
//			},
//		},
//		{
//			Name:   "Handler",
//			Module: handlerModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				r, err := New(context.Background(), []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance()
//				require.NoError(t, err)
//
//				i.Context().Response.Body = []byte("Test Data")
//				err = i.Run(context.Background())
//				assert.NoError(t, err)
//
//				assert.Equal(t, "Test Data-modified", string(i.Context().Response.Body))
//			},
//		},
//		{
//			Name:   "Next",
//			Module: nextModule,
//			Run: func(scaleFunc *scalefunc.Schema, t *testing.T) {
//				r, err := New(context.Background(), []*scalefunc.Schema{scaleFunc})
//				require.NoError(t, err)
//
//				i, err := r.Instance(func(ctx *httpSignature.Context) (*httpSignature.Context, error) {
//					ctx.Response.Body = append(ctx.Response.Body, []byte("-next")...)
//					return ctx, nil
//				})
//				require.NoError(t, err)
//
//				i.Context().Response.Body = []byte("Test Data")
//				err = i.Run(context.Background())
//				assert.NoError(t, err)
//
//				assert.Equal(t, "Test Data-modified-next", string(i.Context().Response.Body))
//			},
//		},
//	}
//
//	for _, testCase := range testCases {
//		t.Run(testCase.Name, func(t *testing.T) {
//			module, err := os.ReadFile(generatedModules[testCase.Module])
//			require.NoError(t, err)
//
//			scaleFunc := &scalefunc.Schema{
//				Version:   scalefunc.V1Alpha,
//				Name:      "TestName",
//				Tag:       "TestTag",
//				Signature: "ExampleName@ExampleVersion",
//				Language:  scalefunc.Go,
//				Function:  module,
//			}
//			testCase.Run(scaleFunc, t)
//		})
//	}
//}