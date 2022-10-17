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

// Package fiber provides a Scale Runtime Adapter for the fiber library
package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/loopholelabs/scale-go/adapters/fasthttp"
	"github.com/loopholelabs/scale-go/runtime"
)

type Fiber struct {
	runtime *runtime.Runtime
}

func New(runtime *runtime.Runtime) *Fiber {
	return &Fiber{
		runtime: runtime,
	}
}

func (f *Fiber) Handle(ctx *fiber.Ctx) error {
	i, err := f.runtime.Instance(ctx.Context(), nil)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	fasthttp.FromRequestContext(i.Context(), ctx.Context())
	err = i.Run(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}
	fasthttp.ToResponseContext(i.Context(), ctx.Context())
	return nil
}

func (f *Fiber) Next(fCtx *fiber.Ctx) runtime.Next {
	if fCtx.Next() == nil {
		return nil
	}
	return func(ctx *runtime.Context) *runtime.Context {
		fasthttp.ToRequestContext(ctx, fCtx.Context())
		fasthttp.ToResponseContext(ctx, fCtx.Context())
		err := fCtx.Next()
		if err != nil {
			_ = fiber.DefaultErrorHandler(fCtx, err)
		}
		fasthttp.FromRequestContext(ctx, fCtx.Context())
		fasthttp.FromResponseContext(ctx, fCtx.Context())
		return ctx
	}
}
