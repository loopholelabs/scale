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
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/loopholelabs/scale/signature"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	V1AlphaVersion = "v1alpha"
)

var (
	ErrInvalidName         = errors.New("invalid name")
	ErrInvalidFunctionName = errors.New("invalid function name")
	ErrInvalidTag          = errors.New("invalid tag")
	ErrNoInstanceId        = errors.New("Extension has no int32 InstanceId defined")
)

var (
	ValidLabel    = regexp.MustCompile(`^[A-Za-z0-9]*$`)
	InvalidString = regexp.MustCompile(`[^A-Za-z0-9-.]`)
)

var (
	TitleCaser = cases.Title(language.Und, cases.NoLower)
)

type Schema struct {
	Version            string                   `hcl:"version,attr"`
	Name               string                   `hcl:"name,attr"`
	Tag                string                   `hcl:"tag,attr"`
	Interfaces         []*InterfaceSchema       `hcl:"interface,block"`
	Functions          []*FunctionSchema        `hcl:"function,block"`
	Enums              []*signature.EnumSchema  `hcl:"enum,block"`
	Models             []*signature.ModelSchema `hcl:"model,block"`
	hasLimitValidator  bool
	hasLengthValidator bool
	hasRegexValidator  bool
	hasCaseModifier    bool
}

func ReadSchema(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	s := new(Schema)
	return s, s.Decode(data)
}

func (s *Schema) Decode(data []byte) error {
	file, diag := hclsyntax.ParseConfig(data, "", hcl.Pos{Line: 1, Column: 1})
	if diag.HasErrors() {
		return diag.Errs()[0]
	}

	diag = gohcl.DecodeBody(file.Body, nil, s)
	if diag.HasErrors() {
		return diag.Errs()[0]
	}

	return nil
}

func (s *Schema) Encode() ([]byte, error) {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(s, f.Body())
	return f.Bytes(), nil
}

func (s *Schema) Validate() error {
	switch s.Version {
	case V1AlphaVersion:
		if !ValidLabel.MatchString(s.Name) {
			return ErrInvalidName
		}

		if InvalidString.MatchString(s.Tag) {
			return ErrInvalidTag
		}

		// Transform all model names and references to TitleCase (e.g. "myModel" -> "MyModel")
		for _, model := range s.Models {
			model.Normalize()
		}

		// Transform all model names and references to TitleCase (e.g. "myModel" -> "MyModel")
		for _, enum := range s.Enums {
			enum.Normalize()
		}

		// Validate all models
		knownModels := make(map[string]struct{})
		for _, model := range s.Models {
			err := model.Validate(knownModels, s.Enums)
			if err != nil {
				return err
			}
		}

		// Validate all enums
		knownEnums := make(map[string]struct{})
		for _, enum := range s.Enums {
			err := enum.Validate(knownEnums)
			if err != nil {
				return err
			}
		}

		// Ensure all model and enum references are valid
		for _, model := range s.Models {
			for _, modelReference := range model.Models {
				if _, ok := knownModels[modelReference.Reference]; !ok {
					return fmt.Errorf("unknown %s.%s.reference: %s", model.Name, modelReference.Name, modelReference.Reference)
				}
			}

			for _, modelReferenceArray := range model.ModelArrays {
				if _, ok := knownModels[modelReferenceArray.Reference]; !ok {
					return fmt.Errorf("unknown %s.%s.reference: %s", model.Name, modelReferenceArray.Name, modelReferenceArray.Reference)
				}
			}

			for _, str := range model.Strings {
				if str.LengthValidator != nil {
					s.hasLengthValidator = true
				}
				if str.RegexValidator != nil {
					s.hasRegexValidator = true
				}
				if str.CaseModifier != nil {
					s.hasCaseModifier = true
				}
			}

			for _, strMap := range model.StringMaps {
				if !signature.ValidPrimitiveType(strMap.Value) {
					if _, ok := knownModels[strMap.Value]; !ok {
						return fmt.Errorf("unknown %s.%s.value: %s", model.Name, strMap.Name, strMap.Value)
					}
				}
			}

			for _, i32 := range model.Int32s {
				if i32.LimitValidator != nil {
					s.hasLimitValidator = true
				}
			}

			for _, i32Map := range model.Int32Maps {
				if !signature.ValidPrimitiveType(i32Map.Value) {
					if _, ok := knownModels[i32Map.Value]; !ok {
						return fmt.Errorf("unknown %s.%s.value: %s", model.Name, i32Map.Name, i32Map.Value)
					}
				}
			}

			for _, i64 := range model.Int64s {
				if i64.LimitValidator != nil {
					s.hasLimitValidator = true
				}
			}

			for _, i64Map := range model.Int64Maps {
				if !signature.ValidPrimitiveType(i64Map.Value) {
					if _, ok := knownModels[i64Map.Value]; !ok {
						return fmt.Errorf("unknown %s.%s.value: %s", model.Name, i64Map.Name, i64Map.Value)
					}
				}
			}

			for _, u32 := range model.Uint32s {
				if u32.LimitValidator != nil {
					s.hasLimitValidator = true
				}
			}

			for _, u32Map := range model.Uint32Maps {
				if !signature.ValidPrimitiveType(u32Map.Value) {
					if _, ok := knownModels[u32Map.Value]; !ok {
						return fmt.Errorf("unknown %s.%s.value: %s", model.Name, u32Map.Name, u32Map.Value)
					}
				}
			}

			for _, u64 := range model.Uint64s {
				if u64.LimitValidator != nil {
					s.hasLimitValidator = true
				}
			}

			for _, u64Map := range model.Uint64Maps {
				if !signature.ValidPrimitiveType(u64Map.Value) {
					if _, ok := knownModels[u64Map.Value]; !ok {
						return fmt.Errorf("unknown %s.%s.value: %s", model.Name, u64Map.Name, u64Map.Value)
					}
				}
			}

			for _, f32 := range model.Float32s {
				if f32.LimitValidator != nil {
					s.hasLimitValidator = true
				}
			}

			for _, f64 := range model.Float64s {
				if f64.LimitValidator != nil {
					s.hasLimitValidator = true
				}
			}

			for _, enumReference := range model.Enums {
				if _, ok := knownEnums[enumReference.Reference]; !ok {
					return fmt.Errorf("unknown %s.%s.reference: %s", model.Name, enumReference.Name, enumReference.Reference)
				}
			}

			for _, enumReferenceArray := range model.EnumArrays {
				if _, ok := knownEnums[enumReferenceArray.Reference]; !ok {
					return fmt.Errorf("unknown %s.%s.reference: %s", model.Name, enumReferenceArray.Name, enumReferenceArray.Reference)
				}
			}

			for _, enumReferenceMap := range model.EnumMaps {
				if _, ok := knownEnums[enumReferenceMap.Reference]; !ok {
					return fmt.Errorf("unknown %s.%s.reference: %s", model.Name, enumReferenceMap.Name, enumReferenceMap.Reference)
				}

				if !signature.ValidPrimitiveType(enumReferenceMap.Value) {
					if _, ok := knownModels[enumReferenceMap.Value]; !ok {
						return fmt.Errorf("unknown %s.%s.value: %s", model.Name, enumReferenceMap.Name, enumReferenceMap.Value)
					}
				}
			}
		}

		// Map of interfaces, and check for name collisions.
		knownInterfaces := make(map[string]struct{})
		for _, inter := range s.Interfaces {
			_, dupe := knownModels[inter.Name]
			if dupe {
				return fmt.Errorf("interface name collides with a model %s", inter.Name)
			}
			_, dupe = knownInterfaces[inter.Name]
			if dupe {
				return fmt.Errorf("interface name collides with an interface %s", inter.Name)
			}
			knownInterfaces[inter.Name] = struct{}{}
		}

		for _, inter := range s.Interfaces {
			for _, f := range inter.Functions {
				// Make sure the function name is ok
				if !ValidLabel.MatchString(f.Name) {
					return ErrInvalidFunctionName
				}

				// Make sure the params exist as model.
				if f.Params != "" {
					f.Params = TitleCaser.String(f.Params)
					if _, ok := knownModels[f.Params]; !ok {
						return fmt.Errorf("unknown params in function %s: %s", f.Name, f.Params)
					}
				}

				// Return can either be a model or interface
				if f.Return != "" {
					f.Return = TitleCaser.String(f.Return)
					_, foundModel := knownModels[f.Return]
					_, foundInterface := knownInterfaces[f.Return]
					if !foundModel && !foundInterface {
						return fmt.Errorf("unknown return in function %s: %s", f.Name, f.Return)
					}
				}
			}
		}

		// Check any global functions
		for _, f := range s.Functions {
			// Make sure the function name is ok
			if !ValidLabel.MatchString(f.Name) {
				return ErrInvalidFunctionName
			}

			// Make sure the params exist as model.
			if f.Params != "" {
				f.Params = TitleCaser.String(f.Params)
				if _, ok := knownModels[f.Params]; !ok {
					return fmt.Errorf("unknown params in function %s: %s", f.Name, f.Params)
				}
			}

			// Return can either be a model or interface
			if f.Return != "" {
				f.Return = TitleCaser.String(f.Return)
				_, foundModel := knownModels[f.Return]
				_, foundInterface := knownInterfaces[f.Return]
				if !foundModel && !foundInterface {
					return fmt.Errorf("unknown return in function %s: %s", f.Name, f.Return)
				}
			}
		}

		return nil
	default:
		return fmt.Errorf("unknown schema version: %s", s.Version)
	}

}

// Hash returns the SHA256 hash of the schema
func (s *Schema) Hash() ([]byte, error) {
	d, err := s.Encode()
	if err != nil {
		return nil, err
	}

	h := sha256.New()
	if _, err = h.Write(d); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// Clone returns a deep copy of the schema
func (s *Schema) Clone() (*Schema, error) {
	clone := new(Schema)
	encoded, err := s.Encode()
	if err != nil {
		return nil, err
	}
	if err = clone.Decode(encoded); err != nil {
		return nil, err
	}
	return clone, nil
}

// CloneWithDisabledAccessorsValidatorsAndModifiers returns a clone of the
// schema with all accessors, validators, and modifiers disabled
func (s *Schema) CloneWithDisabledAccessorsValidatorsAndModifiers() (*Schema, error) {
	clone, err := s.Clone()
	if err != nil {
		return nil, err
	}
	clone.hasCaseModifier = false
	clone.hasLimitValidator = false
	clone.hasRegexValidator = false
	clone.hasLengthValidator = false
	for _, model := range clone.Models {
		for _, modelReference := range model.Models {
			modelReference.Accessor = false
		}

		for _, modelReferenceArray := range model.ModelArrays {
			modelReferenceArray.Accessor =
				false
		}

		for _, str := range model.Strings {
			var accessorValue bool
			str.Accessor = &accessorValue
			str.CaseModifier = nil
			str.LengthValidator = nil
			str.RegexValidator = nil
		}

		for _, strArray := range model.StringArrays {
			var accessorValue bool
			strArray.Accessor = &accessorValue
		}

		for _, strMap := range model.StringMaps {
			var accessorValue bool
			strMap.Accessor = &accessorValue
		}

		for _, i32 := range model.Int32s {
			var accessorValue bool
			i32.Accessor = &accessorValue
			i32.LimitValidator = nil
		}

		for _, i32Array := range model.Int32Arrays {
			var accessorValue bool
			i32Array.Accessor = &accessorValue
		}

		for _, i32Map := range model.Int32Maps {
			var accessorValue bool
			i32Map.Accessor = &accessorValue
		}

		for _, i64 := range model.Int64s {
			var accessorValue bool
			i64.Accessor = &accessorValue
			i64.LimitValidator = nil
		}

		for _, i64Array := range model.Int64Arrays {
			var accessorValue bool
			i64Array.Accessor = &accessorValue
		}

		for _, i64Map := range model.Int64Maps {
			var accessorValue bool
			i64Map.Accessor = &accessorValue
		}

		for _, u32 := range model.Uint32s {
			var accessorValue bool
			u32.Accessor = &accessorValue
			u32.LimitValidator = nil
		}

		for _, u32Array := range model.Uint32Arrays {
			var accessorValue bool
			u32Array.Accessor = &accessorValue
		}

		for _, u32Map := range model.Uint32Maps {
			var accessorValue bool
			u32Map.Accessor = &accessorValue
		}

		for _, u64 := range model.Uint64s {
			var accessorValue bool
			u64.Accessor = &accessorValue
			u64.LimitValidator = nil
		}

		for _, u64Array := range model.Uint64Arrays {
			var accessorValue bool
			u64Array.Accessor = &accessorValue
		}

		for _, u64Map := range model.Uint64Maps {
			var accessorValue bool
			u64Map.Accessor = &accessorValue
		}

		for _, f32 := range model.Float32s {
			var accessorValue bool
			f32.Accessor = &accessorValue
			f32.LimitValidator = nil
		}

		for _, f32Array := range model.Float32Arrays {
			var accessorValue bool
			f32Array.Accessor = &accessorValue
		}

		for _, f64 := range model.Float64s {
			var accessorValue bool
			f64.Accessor = &accessorValue
			f64.LimitValidator = nil
		}

		for _, f64Array := range model.Float64Arrays {
			var accessorValue bool
			f64Array.Accessor = &accessorValue
		}

		for _, boolean := range model.Bools {
			boolean.Accessor = false
		}

		for _, booleanArray := range model.BoolArrays {
			booleanArray.Accessor = false
		}

		for _, b := range model.Bytes {
			b.Accessor = false
		}

		for _, bytesArray := range model.BytesArrays {
			bytesArray.Accessor = false
		}

		for _, enumReference := range model.Enums {
			enumReference.Accessor = false
		}

		for _, enumReferenceArray := range model.EnumArrays {
			enumReferenceArray.Accessor = false
		}

		for _, enumReferenceMap := range model.EnumMaps {
			enumReferenceMap.Accessor = false
		}
	}

	return clone, clone.validateAndNormalize()
}

// validateAndNormalize validates the Schema and normalizes it
//
// Note: This function modifies the Schema in-place
func (s *Schema) validateAndNormalize() error {

	// TODO...

	return nil
}

func (s *Schema) HasLimitValidator() bool {
	return s.hasLimitValidator
}

func (s *Schema) HasLengthValidator() bool {
	return s.hasLengthValidator
}

func (s *Schema) HasRegexValidator() bool {
	return s.hasRegexValidator
}

func (s *Schema) HasCaseModifier() bool {
	return s.hasCaseModifier
}

func ValidPrimitiveType(t string) bool {
	switch t {
	case "string", "int32", "int64", "uint32", "uint64", "float32", "float64", "bool", "bytes":
		return true
	default:
		return false
	}
}

const MasterTestingSchema = `
version = "v1alpha"
name = "HttpFetch"
tag = "alpha"

function New {
	params = "HttpConfig"
	return = "HttpConnector"	
}

model HttpConfig {
	int32 timeout {
		default = 60
	}
}

model HttpResponse {
	string_map Headers {
		value = "StringList"
	}
	int32 StatusCode {
		default = 0
	}
	bytes Body {
		initial_size = 0
	}
}

model StringList {
  string_array Values {
    initial_size = 0
  }
}

model ConnectionDetails {
	string url {
		default = "https://google.com"
	}
}

interface HttpConnector {
	function Fetch {
		params = "ConnectionDetails"
		return = "HttpResponse"
	}
}

`
