/*
	Copyright 2022 Loophole Labs

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		   http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package main

import (
	"github.com/loopholelabs/scale/go/compile/scale"
	"github.com/loopholelabs/scale/go/context"
	"github.com/loopholelabs/scale/go/utils"
)

// needed to satisfy compiler
func main() {}

//export run
func run() uint64 {
	ctx := context.New()
	if ctx.FromReadBuffer() != nil {
		return 0
	}
	ctx = scale.Scale(ctx)
	return utils.PackUint32(ctx.ToWriteBuffer())
}

//export resize
func resize(size uint32) uint32 {
	return context.Resize(size)
}
