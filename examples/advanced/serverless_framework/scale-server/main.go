package main

import (
  "fmt"
  "log"
	"embed"
	"io"
  "net/http"
	scale "github.com/loopholelabs/scale/go"
	adapter "github.com/loopholelabs/scale-http-adapters/http"
  "github.com/loopholelabs/scalefile/scalefunc"
)

//go:embed local-serverless-latest.scale
var embeddedFunction []byte

func main() {
	sf := new(ScaleFunc)
	_ = sf.Decode(embeddedFunction)

	r, _ := scale.New(context.Background(), []*scalefunc.ScaleFunc{sf})
    handler := adapter.New(nil, r)

    http.Handle("/", handler)
    log.Fatal(http.ListenAndServe(":3000", nil))
}

