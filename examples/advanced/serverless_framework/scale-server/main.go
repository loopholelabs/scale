package main

import (
  "net/http"
  "context"
  "fmt"
  "signature"
  scale "github.com/loopholelabs/scale"
  scalefunc "github.com/loopholelabs/scale/scalefunc"
)

func handler(w http.ResponseWriter, r *http.Request) {
    sf, err := scalefunc.Read("local-serverless-latest.scale")
    if err != nil {
        return
    }

    config := scale.NewConfig().WithSignature(*signature.Signature).WithFunction(sf)

	  r, _ := scale.New(config)
    i, _ := r.Instance()

    i.Run(context.Context)
    fmt.Fprintf(w, "Ran the function!")
}

func main() {
    http.Handle("/", http.HandlerFunc(handler))
    log.Fatal(http.ListenAndServe(":3000", nil))
}

