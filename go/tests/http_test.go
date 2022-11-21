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

//import (
//	"bytes"
//	"context"
//	"fmt"
//	adapter "github.com/loopholelabs/scale/go/adapters/http"
//	"github.com/loopholelabs/scale/go/runtime"
//	"github.com/loopholelabs/scale/scalefile"
//	"github.com/loopholelabs/scale/scalefunc"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"os"
//	"path"
//	"testing"
//)
//
//func TestHTTPMiddleware(t *testing.T) {
//	module, err := os.ReadFile(path.Join("modules", fmt.Sprintf("%s.wasm", "http-middleware")))
//	assert.NoError(t, err)
//
//	scaleFunc := scalefunc.ScaleFunc{
//		ScaleFile: scalefile.ScaleFile{
//			Name: "http-middleware",
//			Build: scalefile.Build{
//				Language: "go",
//			},
//			Middleware: true,
//		},
//		Function: module,
//	}
//
//	r, err := runtime.New(context.Background(), []scalefunc.ScaleFunc{scaleFunc})
//	require.NoError(t, err)
//
//	httpAdapter := adapter.New(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
//		w.Header().Set("NEXT", "TRUE")
//		w.WriteHeader(http.StatusOK)
//		_, _ = w.Write([]byte("Hello World"))
//	}), r)
//
//	server := httptest.NewServer(httpAdapter)
//	defer server.Close()
//
//	req, err := http.NewRequest("GET", server.URL, nil)
//	assert.NoError(t, err)
//
//	res, err := http.DefaultClient.Do(req)
//	assert.NoError(t, err)
//
//	body, err := io.ReadAll(res.Body)
//	assert.NoError(t, err)
//	assert.Equal(t, "Hello World", string(body))
//	assert.Equal(t, "TRUE", res.Header.Get("MIDDLEWARE"))
//	assert.Equal(t, "TRUE", res.Header.Get("NEXT"))
//}
//
//func TestHTTPEndpoint(t *testing.T) {
//	module, err := os.ReadFile(path.Join("modules", fmt.Sprintf("%s.wasm", "http-endpoint")))
//	assert.NoError(t, err)
//
//	scaleFunc := scalefunc.ScaleFunc{
//		ScaleFile: scalefile.ScaleFile{
//			Name: "http-endpoint",
//			Build: scalefile.Build{
//				Language: "go",
//			},
//		},
//		Function: module,
//	}
//
//	r, err := runtime.New(context.Background(), []scalefunc.ScaleFunc{scaleFunc})
//	require.NoError(t, err)
//
//	httpAdapter := adapter.New(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
//		t.Fatal("HTTP handler next should not be called")
//	}), r)
//
//	server := httptest.NewServer(httpAdapter)
//	defer server.Close()
//
//	req, err := http.NewRequest("GET", server.URL, bytes.NewBufferString("Hello World"))
//	assert.NoError(t, err)
//
//	res, err := http.DefaultClient.Do(req)
//	assert.NoError(t, err)
//
//	body, err := io.ReadAll(res.Body)
//	assert.NoError(t, err)
//	assert.Equal(t, "Hello World", string(body))
//}
//
//func TestHTTPChain(t *testing.T) {
//	middlewareModule, err := os.ReadFile(path.Join("modules", fmt.Sprintf("%s.wasm", "http-middleware")))
//	assert.NoError(t, err)
//
//	endpointModule, err := os.ReadFile(path.Join("modules", fmt.Sprintf("%s.wasm", "http-endpoint")))
//	assert.NoError(t, err)
//
//	middlewareScaleFunc := scalefunc.ScaleFunc{
//		ScaleFile: scalefile.ScaleFile{
//			Name: "http-middleware",
//			Build: scalefile.Build{
//				Language: "go",
//			},
//			Middleware: true,
//		},
//		Function: middlewareModule,
//	}
//
//	endpointScaleFunc := scalefunc.ScaleFunc{
//		ScaleFile: scalefile.ScaleFile{
//			Name: "http-endpoint",
//			Build: scalefile.Build{
//				Language: "go",
//			},
//		},
//		Function: endpointModule,
//	}
//
//	r, err := runtime.New(context.Background(), []scalefunc.ScaleFunc{middlewareScaleFunc, endpointScaleFunc})
//	require.NoError(t, err)
//
//	httpAdapter := adapter.New(nil, r)
//
//	server := httptest.NewServer(httpAdapter)
//	defer server.Close()
//
//	req, err := http.NewRequest("GET", server.URL, bytes.NewBufferString("Hello World"))
//	assert.NoError(t, err)
//
//	res, err := http.DefaultClient.Do(req)
//	assert.NoError(t, err)
//
//	assert.Equal(t, "TRUE", res.Header.Get("MIDDLEWARE"))
//	body, err := io.ReadAll(res.Body)
//	assert.NoError(t, err)
//	assert.Equal(t, "Hello World", string(body))
//}
