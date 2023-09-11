package example

import (
	"fmt"
	"signature"
)

func Example(ctx *signature.ModelWithAllFieldTypes) (*signature.ModelWithAllFieldTypes, error) {
	fmt.Printf("This is a Golang Function")
	if ctx != nil {
		ctx.StringField = "This is a Golang Function"
	}
	return signature.Next(ctx)
}
