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
	"regexp"
)

type StringRegexValidatorSchema struct {
	Expression string `hcl:"expression,attr"`
}

type StringLengthValidatorSchema struct {
	Minimum *uint `hcl:"min,optional"`
	Maximum *uint `hcl:"max,optional"`
}

type StringCaseModifierSchema struct {
	Kind string `hcl:"kind,attr"`
}

type StringSchema struct {
	Name            string                       `hcl:"name,label"`
	Default         string                       `hcl:"default,attr"`
	Accessor        *bool                        `hcl:"accessor,optional"`
	RegexValidator  *StringRegexValidatorSchema  `hcl:"regex_validator,block"`
	LengthValidator *StringLengthValidatorSchema `hcl:"length_validator,block"`
	CaseModifier    *StringCaseModifierSchema    `hcl:"case_modifier,block"`
}

func (s *StringSchema) Validate(model *ModelSchema) error {
	if !ValidLabel.MatchString(s.Name) {
		return fmt.Errorf("invalid %s.string name: %s", model.Name, s.Name)
	}

	if s.LengthValidator != nil {
		if s.LengthValidator.Maximum != nil {
			if *s.LengthValidator.Maximum == 0 {
				return fmt.Errorf("invalid %s.%s.length_validator: maximum length cannot be zero", model.Name, s.Name)
			}

			if s.LengthValidator.Minimum != nil {
				if *s.LengthValidator.Minimum > *s.LengthValidator.Maximum {
					return fmt.Errorf("invalid %s.%s.length_validator: minimum length cannot be greater than maximum length", model.Name, s.Name)
				}
			}
		}
	}

	if s.RegexValidator != nil {
		regex, err := regexp.Compile(s.RegexValidator.Expression)
		if err != nil {
			return fmt.Errorf("invalid %s.%s.regex_validator: %w", model.Name, s.Name, err)
		}
		if !regex.MatchString(s.Default) {
			return fmt.Errorf("invalid %s.%s.default: does not match regex", model.Name, s.Name)
		}
	}

	if s.CaseModifier != nil {
		switch s.CaseModifier.Kind {
		case "upper", "lower", "none":
		default:
			return fmt.Errorf("invalid %s.%s.caseModifier: kind must be upper, lower or none", model.Name, s.Name)
		}
	}

	if s.Accessor != nil {
		if *s.Accessor == false && (s.LengthValidator != nil || s.RegexValidator != nil || s.CaseModifier != nil) {
			return fmt.Errorf("invalid %s.%s.accessor: cannot be false while using validators or modifiers", model.Name, s.Name)
		}
	} else {
		if s.LengthValidator != nil || s.RegexValidator != nil || s.CaseModifier != nil {
			s.Accessor = new(bool)
			*s.Accessor = true
		} else {
			s.Accessor = new(bool)
			*s.Accessor = false
		}
	}

	return nil
}

type StringArraySchema struct {
	Name            string                       `hcl:"name,label"`
	InitialSize     uint32                       `hcl:"initial_size,attr"`
	Accessor        *bool                        `hcl:"accessor,optional"`
	RegexValidator  *StringRegexValidatorSchema  `hcl:"regex_validator,block"`
	LengthValidator *StringLengthValidatorSchema `hcl:"length_validator,block"`
	CaseModifier    *StringCaseModifierSchema    `hcl:"caseModifier,block"`
}

func (s *StringArraySchema) Validate(model *ModelSchema) error {
	if !ValidLabel.MatchString(s.Name) {
		return fmt.Errorf("invalid %s.string_array name: %s", model.Name, s.Name)
	}

	if s.LengthValidator != nil {
		if s.LengthValidator.Maximum != nil {
			if *s.LengthValidator.Maximum == 0 {
				return fmt.Errorf("invalid %s.%s.length_validator: maximum length cannot be zero", model.Name, s.Name)
			}

			if s.LengthValidator.Minimum != nil {
				if *s.LengthValidator.Minimum > *s.LengthValidator.Maximum {
					return fmt.Errorf("invalid %s.%s.length_validator: minimum length cannot be greater than maximum length", model.Name, s.Name)
				}
			}
		}
	}

	if s.RegexValidator != nil {
		if _, err := regexp.Compile(s.RegexValidator.Expression); err != nil {
			return fmt.Errorf("invalid %s.%s.regex_validator: %w", model.Name, s.Name, err)
		}
	}

	if s.CaseModifier != nil {
		switch s.CaseModifier.Kind {
		case "upper", "lower", "none":
		default:
			return fmt.Errorf("invalid %s.%s.caseModifier: kind must be upper, lower or none", model.Name, s.Name)
		}
	}

	if s.Accessor != nil {
		if *s.Accessor == false && (s.LengthValidator != nil || s.RegexValidator != nil || s.CaseModifier != nil) {
			return fmt.Errorf("invalid %s.%s.accessor: cannot be false while using validators or modifiers", model.Name, s.Name)
		}
	} else {
		if s.LengthValidator != nil || s.RegexValidator != nil || s.CaseModifier != nil {
			s.Accessor = new(bool)
			*s.Accessor = true
		} else {
			s.Accessor = new(bool)
			*s.Accessor = false
		}
	}

	return nil
}

type StringMapSchema struct {
	Name            string                       `hcl:"name,label"`
	Value           string                       `hcl:"value,attr"`
	Accessor        *bool                        `hcl:"accessor,optional"`
	RegexValidator  *StringRegexValidatorSchema  `hcl:"regex_validator,block"`
	LengthValidator *StringLengthValidatorSchema `hcl:"length_validator,block"`
	CaseModifier    *StringCaseModifierSchema    `hcl:"caseModifier,block"`
}

func (s *StringMapSchema) Validate(model *ModelSchema) error {
	if !ValidLabel.MatchString(s.Name) {
		return fmt.Errorf("invalid %s.string_map name: %s", model.Name, s.Name)
	}

	if s.LengthValidator != nil {
		if s.LengthValidator.Maximum != nil {
			if *s.LengthValidator.Maximum == 0 {
				return fmt.Errorf("invalid %s.%s.length_validator: maximum length cannot be zero", model.Name, s.Name)
			}

			if s.LengthValidator.Minimum != nil {
				if *s.LengthValidator.Minimum > *s.LengthValidator.Maximum {
					return fmt.Errorf("invalid %s.%s.length_validator: minimum length cannot be greater than maximum length", model.Name, s.Name)
				}
			}
		}
	}

	if s.RegexValidator != nil {
		if _, err := regexp.Compile(s.RegexValidator.Expression); err != nil {
			return fmt.Errorf("invalid %s.%s.regex_validator: %w", model.Name, s.Name, err)
		}
	}

	if s.CaseModifier != nil {
		switch s.CaseModifier.Kind {
		case "upper", "lower":
		default:
			return fmt.Errorf("invalid %s.%s.caseModifier: kind must be upper or lowere", model.Name, s.Name)
		}
	}

	if s.Accessor != nil {
		if *s.Accessor == false && (s.LengthValidator != nil || s.RegexValidator != nil || s.CaseModifier != nil) {
			return fmt.Errorf("invalid %s.%s.accessor: cannot be false while using validators or modifiers", model.Name, s.Name)
		}
	} else {
		if s.LengthValidator != nil || s.RegexValidator != nil || s.CaseModifier != nil {
			s.Accessor = new(bool)
			*s.Accessor = true
		} else {
			s.Accessor = new(bool)
			*s.Accessor = false
		}
	}

	return nil
}
