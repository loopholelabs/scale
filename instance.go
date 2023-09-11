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
	"crypto/rand"
	"fmt"
	interfaces "github.com/loopholelabs/scale-signature-interfaces"
)

// Instance is a single instance of a Scale Function chain
type Instance[T interfaces.Signature] struct {
	// runtime is the runtime that this instance belongs to
	runtime *Scale[T]

	// identifier is the unique identifier for this instance
	identifier []byte

	// head is the head function in the chain for this instance
	head *function[T]

	// next is the next function in the chain for this instance
	next Next[T]
}

func newInstance[T interfaces.Signature](ctx context.Context, runtime *Scale[T], next ...Next[T]) (*Instance[T], error) {
	instance := &Instance[T]{
		runtime:    runtime,
		identifier: make([]byte, 16),
	}

	_, err := rand.Read(instance.identifier)
	if err != nil {
		return nil, err
	}

	if len(next) > 0 && next[0] != nil {
		instance.next = next[0]
	} else {
		instance.next = func(ctx T) (T, error) {
			return ctx, nil
		}
	}

	previousFunction := instance.head
	nextTemplate := instance.runtime.head

	for nextTemplate != nil {
		fn, err := newFunction(ctx, instance, nextTemplate)
		if err != nil {
			return nil, fmt.Errorf("failed to create function: %w", err)
		}
		if instance.head == nil {
			instance.head = fn
		}
		if previousFunction != nil {
			previousFunction.next = fn
		}
		previousFunction = fn
		nextTemplate = nextTemplate.next
	}

	return instance, nil
}

func (i *Instance[T]) Run(ctx context.Context, signature T) error {
	m, err := i.head.getModule(signature)
	if err != nil {
		return fmt.Errorf("failed to get module for function '%s': %w", i.head.template.identifier, err)
	}
	err = m.run(ctx)
	i.head.putModule(m)
	if err != nil {
		return fmt.Errorf("failed to run function '%s': %w", i.head.template.identifier, err)
	}
	return nil
}
