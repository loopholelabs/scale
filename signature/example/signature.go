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

package example

import (
	"github.com/loopholelabs/scale/signature"
	example "github.com/loopholelabs/scale/signature/example/source"
)

var _ signature.Signature = (*Signature)(nil)

type Signature struct{}

func (s *Signature) Version() string {
	return example.VERSION
}

func (s *Signature) RuntimeContext() signature.RuntimeContext {
	return example.New().RuntimeContext()
}

func (s *Signature) Resize(size uint32) uint32 {
	return example.Resize(size)
}
