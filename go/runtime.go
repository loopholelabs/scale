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

// Package runtime implements the Scale Runtime in Go.
package runtime

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/loopholelabs/scale-signature"
	"github.com/loopholelabs/scalefile/scalefunc"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"sync"
)

var (
	NoFunctionsError           = errors.New("no functions found in runtime")
	InvalidScaleFunctionError  = errors.New("invalid scale function")
	IncompatibleSignatureError = errors.New("incompatible signature")
	InvalidSignatureError      = errors.New("invalid signature")
)

// Next is the next function in the middleware chain. It's meant to be implemented
// by whatever adapter is being used.
type Next[T signature.Signature] func(ctx T) (T, error)

// Runtime is the Scale Runtime. It is responsible for initializing
// and managing the WASM runtime as well as the scale function chain.
type Runtime[T signature.Signature] struct {
	runtime      wazero.Runtime
	moduleConfig wazero.ModuleConfig

	signature T

	functions []*Function[T]
	head      *Function[T]
	tail      *Function[T]

	modulesMu sync.RWMutex
	modules   map[string]*Module[T]
}

func New[T signature.Signature](ctx context.Context, sig T, functions []*scalefunc.ScaleFunc) (*Runtime[T], error) {
	if len(functions) == 0 {
		return nil, NoFunctionsError
	}

	sigName := sig.Name()
	if sigName == "" {
		return nil, InvalidSignatureError
	}
	sigVersion := sig.Version()
	if sigVersion == "" {
		return nil, InvalidSignatureError
	}
	for _, f := range functions {
		if f.Name == "" {
			return nil, InvalidScaleFunctionError
		}
		if f.Signature == "" {
			return nil, InvalidScaleFunctionError
		}
		if f.Function == nil {
			return nil, InvalidScaleFunctionError
		}

		_, name, version := signature.ParseSignature(f.Signature)
		if sigName != name || sigVersion != version {
			return nil, IncompatibleSignatureError
		}
	}

	r := &Runtime[T]{
		runtime:      wazero.NewRuntime(ctx),
		moduleConfig: wazero.NewModuleConfig().WithSysNanotime().WithSysWalltime().WithRandSource(rand.Reader),
		modules:      make(map[string]*Module[T]),
		signature:    sig,
	}

	module := r.runtime.NewHostModuleBuilder("env").
		NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(r.next), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{}).
		WithParameterNames("pointer", "length").Export("next")

	compiled, err := module.Compile(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to compile env: %w", err)
	}

	_, err = r.runtime.InstantiateModule(ctx, compiled, r.moduleConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate env: %w", err)
	}

	_, err = wasi_snapshot_preview1.Instantiate(ctx, r.runtime)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate wasi: %w", err)
	}

	for _, f := range functions {
		sf, err := r.compileFunction(ctx, f)
		if err != nil {
			return nil, fmt.Errorf("failed to compile function '%s': %w", f.Name, err)
		}
		if r.head == nil {
			r.head = sf
		}
		if r.tail != nil {
			r.tail.next = sf
		}
		r.tail = sf
	}

	return r, nil
}

func (r *Runtime[T]) compileFunction(ctx context.Context, scaleFunc *scalefunc.ScaleFunc) (*Function[T], error) {
	compiled, err := r.runtime.CompileModule(ctx, scaleFunc.Function)
	if err != nil {
		return nil, fmt.Errorf("failed to compile function '%s': %w", scaleFunc.Name, err)
	}

	f := &Function[T]{
		scaleFunc: scaleFunc,
		compiled:  compiled,
	}

	f.modulePool = NewPool[T](ctx, f, r)

	r.functions = append(r.functions, f)
	return f, nil
}
