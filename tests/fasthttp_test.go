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
	adapter "github.com/loopholelabs/scale-go/adapters/fasthttp"
	"github.com/loopholelabs/scale-go/runtime"
	"github.com/loopholelabs/scale-go/scalefile"
	"github.com/loopholelabs/scale-go/scalefunc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
	"net"
	"os"
	"path"
	"testing"
)

func TestFastHTTPEndpoint(t *testing.T) {
	module, err := os.ReadFile(path.Join("modules", fmt.Sprintf("%s.wasm", "http-endpoint")))
	assert.NoError(t, err)

	scaleFunc := scalefunc.ScaleFunc{
		ScaleFile: scalefile.ScaleFile{
			Name: "http-endpoint",
			Build: scalefile.Build{
				Language: "go",
			},
		},
		Function: module,
	}

	r, err := runtime.New(context.Background(), []scalefunc.ScaleFunc{scaleFunc})
	require.NoError(t, err)

	fasthttpAdapter := adapter.New(r)

	listener := fasthttputil.NewInmemoryListener()
	defer func() {
		err := listener.Close()
		assert.NoError(t, err)
	}()

	go func() {
		err := fasthttp.Serve(listener, fasthttpAdapter.Handle)
		assert.NoError(t, err)
	}()

	client := &fasthttp.Client{
		Name: "test-client",
		Dial: func(_ string) (net.Conn, error) {
			return listener.Dial()
		},
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://test.com")
	req.Header.SetMethod("GET")
	req.SetBodyString("Hello World")

	res := fasthttp.AcquireResponse()

	err = client.Do(req, res)
	assert.NoError(t, err)

	assert.Equal(t, "Hello World", string(res.Body()))
}

func TestFastHTTPChain(t *testing.T) {
	middlewareModule, err := os.ReadFile(path.Join("modules", fmt.Sprintf("%s.wasm", "http-middleware")))
	assert.NoError(t, err)

	endpointModule, err := os.ReadFile(path.Join("modules", fmt.Sprintf("%s.wasm", "http-endpoint")))
	assert.NoError(t, err)

	middlewareScaleFunc := scalefunc.ScaleFunc{
		ScaleFile: scalefile.ScaleFile{
			Name: "http-middleware",
			Build: scalefile.Build{
				Language: "go",
			},
			Middleware: true,
		},
		Function: middlewareModule,
	}

	endpointScaleFunc := scalefunc.ScaleFunc{
		ScaleFile: scalefile.ScaleFile{
			Name: "http-endpoint",
			Build: scalefile.Build{
				Language: "go",
			},
		},
		Function: endpointModule,
	}

	r, err := runtime.New(context.Background(), []scalefunc.ScaleFunc{middlewareScaleFunc, endpointScaleFunc})
	require.NoError(t, err)

	fasthttpAdapter := adapter.New(r)

	listener := fasthttputil.NewInmemoryListener()
	defer func() {
		err := listener.Close()
		assert.NoError(t, err)
	}()

	go func() {
		err := fasthttp.Serve(listener, fasthttpAdapter.Handle)
		assert.NoError(t, err)
	}()

	client := &fasthttp.Client{
		Name: "test-client",
		Dial: func(_ string) (net.Conn, error) {
			return listener.Dial()
		},
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://test.com")
	req.Header.SetMethod("GET")
	req.SetBodyString("Hello World")

	res := fasthttp.AcquireResponse()

	err = client.Do(req, res)
	assert.NoError(t, err)

	assert.Equal(t, "Hello World", string(res.Body()))
	assert.Equal(t, "test", string(res.Header.Peek("X-Test")))
}
