
package serverless

import (
	"signature"
)

func Scale(ctx *signature.Context) (*signature.Context, error) {
	return signature.Next(ctx)
}
