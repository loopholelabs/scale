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

package signature

import "fmt"

type BytesSchema struct {
	Name        string `hcl:"name,label"`
	InitialSize uint32 `hcl:"initial_size,attr"`
	Accessor    bool   `hcl:"accessor,optional"`
}

func (s *BytesSchema) Validate(model *ModelSchema) error {
	if !ValidLabel.MatchString(s.Name) {
		return fmt.Errorf("invalid %s.bytes name: %s", model.Name, s.Name)
	}

	return nil
}

type BytesArraySchema struct {
	Name        string `hcl:"name,label"`
	InitialSize uint32 `hcl:"initial_size,attr"`
	Accessor    bool   `hcl:"accessor,optional"`
}

func (s *BytesArraySchema) Validate(model *ModelSchema) error {
	if !ValidLabel.MatchString(s.Name) {
		return fmt.Errorf("invalid %s.bytes_array name: %s", model.Name, s.Name)
	}

	return nil
}
