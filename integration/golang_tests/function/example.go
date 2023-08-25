package example

import (
	"signature"
)

func Example(ctx *signature.ModelWithAllFieldTypes) (*signature.ModelWithAllFieldTypes, error) {
	ctx.StringField = "Hello World!"
	return signature.Next(ctx)
}
