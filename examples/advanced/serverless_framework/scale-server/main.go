package main

import (
  "net/http"
  "context"
  "fmt"
  "log"
  signature "signature/signature"
  scale "github.com/loopholelabs/scale"
  scalefunc "github.com/loopholelabs/scale/scalefunc"
	interfaces "github.com/loopholelabs/scale-signature-interfaces"
)

func handler(w http.ResponseWriter, r *http.Request) {
    sf, err := scalefunc.Read("local-serverless-latest.scale")
    if err != nil {
        return
    }

    var newSig interfaces.New[*signature.Signature] = signature.New

    config := scale.NewConfig(newSig).WithFunction(sf)

	  runtime, _ := scale.New(config)
    i, _ := runtime.Instance()

    ctx := context.Background()

    i.Run(ctx, signature.New())
    fmt.Fprintf(w, "Ran the function!")
}

func main() {
    http.Handle("/", http.HandlerFunc(handler))
    log.Fatal(http.ListenAndServe(":3000", nil))
}

