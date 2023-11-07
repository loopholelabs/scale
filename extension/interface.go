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

type InterfaceSchema struct {
	Name        string            `hcl:"name,label"`
	Description string            `hcl:"description,optional"`
	Functions   []*FunctionSchema `hcl:"function,block"`
}

func (s *InterfaceSchema) Normalize() {
	s.Name = signature.TitleCaser.String(s.Name)
	for _, function := range s.Functions {
		function.Normalize()
	}
}

func (s *InterfaceSchema) Validate(knownInterfaces map[string]map[string]struct{}) error {
	if !signature.ValidLabel.MatchString(s.Name) {
		return fmt.Errorf("invalid interface name: %s", s.Name)
	}

	if _, ok := knownInterfaces[s.Name]; ok {
		return fmt.Errorf("duplicate interface name: %s", s.Name)
	}
	knownInterfaces[s.Name] = make(map[string]struct{})
	for _, function := range s.Functions {
		if err := function.Validate(knownInterfaces[s.Name]); err != nil {
			return err
		}
	}

	return nil
}
