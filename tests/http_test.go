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
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestHTTPAdapter(t *testing.T) {
	module, err := os.ReadFile(path.Join("modules", fmt.Sprintf("%s.wasm", "http-middleware")))
	assert.NoError(t, err)

	scaleFunc := scalefunc.ScaleFunc{
		ScaleFile: scalefile.ScaleFile{
			Name: "http-middleware",
			Build: scalefile.Build{
				Language: "go",
			},
		},
		Function: module,
	}

	r, err := runtime.New(context.Background(), []scalefunc.ScaleFunc{scaleFunc})
	require.NoError(t, err)

	httpAdapter := adapter.New(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("NEXT", "true")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello World"))
	}), r)

	server := httptest.NewServer(httpAdapter)
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Hello World", string(body))
	assert.Equal(t, "test", res.Header.Get("X-Test"))
	assert.Equal(t, "true", res.Header.Get("NEXT"))
}
