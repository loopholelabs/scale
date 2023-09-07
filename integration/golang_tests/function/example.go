package example

import (
	"signature"
)

func Example(ctx *signature.ModelWithAllFieldTypes) (*signature.ModelWithAllFieldTypes, error) {
	ctx.StringField = "This is a Golang Function"
	return signature.Next(ctx)
}
