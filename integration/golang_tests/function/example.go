package example

import (
	"signature"
)

func Example(ctx *signature.Example) (*signature.Example, error) {
	ctx.Data = "Hello World!"
	return signature.Next(ctx)
}
