package example

import (
	"signature"
)

func Example(ctx *signature.ModelWithAllFieldTypes) (*signature.ModelWithAllFieldTypes, error) {
	if ctx != nil {
		ctx.StringField = "This is a Golang Function"
	}
	return signature.Next(ctx)
}
