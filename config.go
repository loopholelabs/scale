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

package scale

import (
	"context"
	"errors"
	"io"
	"regexp"

	interfaces "github.com/loopholelabs/scale-signature-interfaces"
	"github.com/loopholelabs/scale/scalefunc"
)

var (
	ErrNoConfig        = errors.New("no config provided")
	ErrNoFunctions     = errors.New("no functions provided")
	ErrInvalidFunction = errors.New("invalid function")
	ErrInvalidEnv      = errors.New("invalid environment variable")
)

var (
	envStringRegex = regexp.MustCompile(`[^A-Za-z0-9_]`)
)

type configFunction struct {
	function *scalefunc.Schema
	env      map[string]string
}

// Config is the configuration for a Scale Runtime
type Config[T interfaces.Signature] struct {
	newSignature interfaces.New[T]
	functions    []configFunction
	context      context.Context
	Stdout       io.Writer
	Stderr       io.Writer
	Stdin        io.Reader
}

// NewConfig returns a new Scale Runtime Config
func NewConfig[T interfaces.Signature](newSignature interfaces.New[T]) *Config[T] {
	return &Config[T]{
		newSignature: newSignature,
	}
}

func (c *Config[T]) validate() error {
	if c == nil {
		return ErrNoConfig
	}
	if len(c.functions) == 0 {
		return ErrNoFunctions
	}

	if c.context == nil {
		c.context = context.Background()
	}

	for _, f := range c.functions {
		if f.function == nil {
			return ErrInvalidFunction
		}
		for k := range f.env {
			if !validEnv(k) {
				return ErrInvalidEnv
			}
		}
	}

	return nil
}

func (c *Config[T]) WithSignature(newSignature interfaces.New[T]) *Config[T] {
	c.newSignature = newSignature
	return c
}

func (c *Config[T]) WithFunction(function *scalefunc.Schema, env ...map[string]string) *Config[T] {
	f := configFunction{
		function: function,
	}

	if len(env) > 0 {
		f.env = env[0]
	}

	c.functions = append(c.functions, f)
	return c
}

func (c *Config[T]) WithFunctions(function []*scalefunc.Schema, env ...map[string]string) *Config[T] {
	for _, f := range function {
		c.WithFunction(f, env...)
	}
	return c
}

func (c *Config[T]) WithContext(ctx context.Context) *Config[T] {
	c.context = ctx
	return c
}

func (c *Config[T]) WithStdout(w io.Writer) *Config[T] {
	c.Stdout = w
	return c
}

func (c *Config[T]) WithStderr(w io.Writer) *Config[T] {
	c.Stderr = w
	return c
}

func (c *Config[T]) WithStdin(r io.Reader) *Config[T] {
	c.Stdin = r
	return c
}

// validEnv returns true if the string is valid for use as an environment variable
func validEnv(str string) bool {
	return !envStringRegex.MatchString(str)
}
