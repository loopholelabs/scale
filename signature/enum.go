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

import (
	"fmt"
)

type EnumSchema struct {
	Name   string   `hcl:"name,label"`
	Values []string `hcl:"values,attr"`
}

func (s *EnumSchema) Normalize() {
	s.Name = TitleCaser.String(s.Name)
}

func (s *EnumSchema) Validate(knownEnums map[string]struct{}) error {
	if !ValidLabel.MatchString(s.Name) {
		return fmt.Errorf("invalid enum name: %s", s.Name)
	}

	if _, ok := knownEnums[s.Name]; ok {
		return fmt.Errorf("duplicate enum name: %s", s.Name)
	} else {
		knownEnums[s.Name] = struct{}{}
	}

	visitedValues := make(map[string]struct{}, 0)
	for i := 0; i < len(s.Values); i++ {
		if _, ok := visitedValues[s.Values[i]]; ok {
			return fmt.Errorf("duplicate value in %s: %s", s.Name, s.Values[i])
		} else {
			visitedValues[s.Values[i]] = struct{}{}
		}
	}

	return nil
}

type EnumReferenceSchema struct {
	Name      string `hcl:"name,label"`
	Default   string `hcl:"default,attr"`
	Reference string `hcl:"reference,attr"`
	Accessor  bool   `hcl:"accessor,optional"`
}

func (s *EnumReferenceSchema) Validate(model *ModelSchema, enums []*EnumSchema) error {
	if !ValidLabel.MatchString(s.Name) {
		return fmt.Errorf("invalid %s.enum name: %s", model.Name, s.Name)
	}

	if !ValidLabel.MatchString(s.Reference) {
		return fmt.Errorf("invalid %s.%s.reference: %s", model.Name, s.Name, s.Reference)
	}

	for _, enum := range enums {
		if enum.Name != s.Reference {
			continue
		}

		for _, value := range enum.Values {
			if value == s.Default {
				return nil
			}
		}
	}

	return fmt.Errorf("invalid %s.default: %s is not a valid value", s.Name, s.Default)
}

type EnumArraySchema struct {
	Name        string `hcl:"name,label"`
	Reference   string `hcl:"reference,attr"`
	InitialSize uint32 `hcl:"initial_size,attr"`
	Accessor    bool   `hcl:"accessor,optional"`
}

func (s *EnumArraySchema) Validate(model *ModelSchema) error {
	if !ValidLabel.MatchString(s.Name) {
		return fmt.Errorf("invalid %s.enum_array name: %s", model.Name, s.Name)
	}

	if !ValidLabel.MatchString(s.Reference) {
		return fmt.Errorf("invalid %s.%s.reference: %s", model.Name, s.Name, s.Reference)
	}

	return nil
}

type EnumMapSchema struct {
	Name      string `hcl:"name,label"`
	Reference string `hcl:"reference,attr"`
	Value     string `hcl:"value,attr"`
	Accessor  bool   `hcl:"accessor,optional"`
}

func (s *EnumMapSchema) Validate(model *ModelSchema) error {
	if !ValidLabel.MatchString(s.Name) {
		return fmt.Errorf("invalid %s.enum_map name: %s", model.Name, s.Name)
	}

	if !ValidLabel.MatchString(s.Reference) {
		return fmt.Errorf("invalid %s.%s.reference: %s", model.Name, s.Name, s.Reference)
	}

	return nil
}
