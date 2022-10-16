package main

import (
	"github.com/loopholelabs/scale-go/compile/scale"
	"github.com/loopholelabs/scale-go/context"
	"github.com/loopholelabs/scale-go/utils"
)

// needed to satisfy compiler
func main() {}

//export run
func run(ptr uint32, size uint32) uint64 {
	ctx := context.New()
	ctx.Deserialize(ptr, size)

	ctx = scale.Scale(ctx)

	return utils.PackUint32(ctx.Serialize())
}
