/*
	Copyright 2023 Loophole Labs

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

package extension

import (
	"fmt"

	"github.com/loopholelabs/scale/signature"
)

type FunctionSchema struct {
	Name        string `hcl:"name,label"`
	Description string `hcl:"description,optional"`
	Params      string `hcl:"params,optional"`
	Return      string `hcl:"return,optional"`
}

func (s *FunctionSchema) Normalize() {
	s.Name = signature.TitleCaser.String(s.Name)
	s.Params = signature.TitleCaser.String(s.Params)
	s.Return = signature.TitleCaser.String(s.Return)
}

func (s *FunctionSchema) Validate(knownFunctions map[string]struct{}) error {
	if !signature.ValidLabel.MatchString(s.Name) {
		return fmt.Errorf("invalid function name: %s", s.Name)
	}

	if _, ok := knownFunctions[s.Name]; ok {
		return fmt.Errorf("duplicate function name: %s", s.Name)
	}
	knownFunctions[s.Name] = struct{}{}

	return nil
}
