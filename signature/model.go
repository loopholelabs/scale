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
	"strings"
)

type ModelSchema struct {
	Name        string `hcl:"name,label"`
	Description string `hcl:"description,optional"`

	Models      []*ModelReferenceSchema      `hcl:"model,block"`
	ModelArrays []*ModelReferenceArraySchema `hcl:"model_array,block"`

	Strings      []*StringSchema      `hcl:"string,block"`
	StringArrays []*StringArraySchema `hcl:"string_array,block"`
	StringMaps   []*StringMapSchema   `hcl:"string_map,block"`

	Bools      []*BoolSchema      `hcl:"bool,block"`
	BoolArrays []*BoolArraySchema `hcl:"bool_array,block"`

	Bytes       []*BytesSchema      `hcl:"bytes,block"`
	BytesArrays []*BytesArraySchema `hcl:"bytes_array,block"`

	Enums      []*EnumReferenceSchema `hcl:"enum,block"`
	EnumArrays []*EnumArraySchema     `hcl:"enum_array,block"`
	EnumMaps   []*EnumMapSchema       `hcl:"enum_map,block"`

	Int32s      []*NumberSchema[int32]      `hcl:"int32,block"`
	Int32Arrays []*NumberArraySchema[int32] `hcl:"int32_array,block"`
	Int32Maps   []*NumberMapSchema[int32]   `hcl:"int32_map,block"`

	Int64s      []*NumberSchema[int64]      `hcl:"int64,block"`
	Int64Arrays []*NumberArraySchema[int64] `hcl:"int64_array,block"`
	Int64Maps   []*NumberMapSchema[int64]   `hcl:"int64_map,block"`

	Uint32s      []*NumberSchema[uint32]      `hcl:"uint32,block"`
	Uint32Arrays []*NumberArraySchema[uint32] `hcl:"uint32_array,block"`
	Uint32Maps   []*NumberMapSchema[uint32]   `hcl:"uint32_map,block"`

	Uint64s      []*NumberSchema[uint64]      `hcl:"uint64,block"`
	Uint64Arrays []*NumberArraySchema[uint64] `hcl:"uint64_array,block"`
	Uint64Maps   []*NumberMapSchema[uint64]   `hcl:"uint64_map,block"`

	Float32s      []*NumberSchema[float32]      `hcl:"float32,block"`
	Float32Arrays []*NumberArraySchema[float32] `hcl:"float32_array,block"`

	Float64s      []*NumberSchema[float64]      `hcl:"float64,block"`
	Float64Arrays []*NumberArraySchema[float64] `hcl:"float64_array,block"`
}

func (m *ModelSchema) Normalize() {
	m.Name = TitleCaser.String(m.Name)
	for _, modelReference := range m.Models {
		modelReference.Name = TitleCaser.String(modelReference.Name)
		modelReference.Reference = TitleCaser.String(modelReference.Reference)
	}

	for _, modelReferenceArray := range m.ModelArrays {
		modelReferenceArray.Name = TitleCaser.String(modelReferenceArray.Name)
		modelReferenceArray.Reference = TitleCaser.String(modelReferenceArray.Reference)
	}

	for _, enumReference := range m.Enums {
		enumReference.Name = TitleCaser.String(enumReference.Name)
		enumReference.Reference = TitleCaser.String(enumReference.Reference)
	}

	for _, enumReferenceArray := range m.EnumArrays {
		enumReferenceArray.Name = TitleCaser.String(enumReferenceArray.Name)
		enumReferenceArray.Reference = TitleCaser.String(enumReferenceArray.Reference)
	}

	for _, enumReferenceMap := range m.EnumMaps {
		enumReferenceMap.Name = TitleCaser.String(enumReferenceMap.Name)
		enumReferenceMap.Reference = TitleCaser.String(enumReferenceMap.Reference)

		if !ValidPrimitiveType(strings.ToLower(enumReferenceMap.Value)) {
			enumReferenceMap.Value = TitleCaser.String(enumReferenceMap.Value)
		} else {
			enumReferenceMap.Value = strings.ToLower(enumReferenceMap.Value)
		}
	}

	for _, str := range m.Strings {
		str.Name = TitleCaser.String(str.Name)
	}

	for _, strArray := range m.StringArrays {
		strArray.Name = TitleCaser.String(strArray.Name)
	}

	for _, strMap := range m.StringMaps {
		strMap.Name = TitleCaser.String(strMap.Name)

		if !ValidPrimitiveType(strings.ToLower(strMap.Value)) {
			strMap.Value = TitleCaser.String(strMap.Value)
		} else {
			strMap.Value = strings.ToLower(strMap.Value)
		}
	}

	for _, i32 := range m.Int32s {
		i32.Name = TitleCaser.String(i32.Name)
	}

	for _, i32Array := range m.Int32Arrays {
		i32Array.Name = TitleCaser.String(i32Array.Name)
	}

	for _, i32Map := range m.Int32Maps {
		i32Map.Name = TitleCaser.String(i32Map.Name)

		if !ValidPrimitiveType(strings.ToLower(i32Map.Value)) {
			i32Map.Value = TitleCaser.String(i32Map.Value)
		} else {
			i32Map.Value = strings.ToLower(i32Map.Value)
		}
	}

	for _, i64 := range m.Int64s {
		i64.Name = TitleCaser.String(i64.Name)
	}

	for _, i64Array := range m.Int64Arrays {
		i64Array.Name = TitleCaser.String(i64Array.Name)
	}

	for _, i64Map := range m.Int64Maps {
		i64Map.Name = TitleCaser.String(i64Map.Name)

		if !ValidPrimitiveType(strings.ToLower(i64Map.Value)) {
			i64Map.Value = TitleCaser.String(i64Map.Value)
		} else {
			i64Map.Value = strings.ToLower(i64Map.Value)
		}
	}

	for _, u32 := range m.Uint32s {
		u32.Name = TitleCaser.String(u32.Name)
	}

	for _, u32Array := range m.Uint32Arrays {
		u32Array.Name = TitleCaser.String(u32Array.Name)
	}

	for _, u32Map := range m.Uint32Maps {
		u32Map.Name = TitleCaser.String(u32Map.Name)

		if !ValidPrimitiveType(strings.ToLower(u32Map.Value)) {
			u32Map.Value = TitleCaser.String(u32Map.Value)
		} else {
			u32Map.Value = strings.ToLower(u32Map.Value)
		}
	}

	for _, u64 := range m.Uint64s {
		u64.Name = TitleCaser.String(u64.Name)
	}

	for _, u64Array := range m.Uint64Arrays {
		u64Array.Name = TitleCaser.String(u64Array.Name)
	}

	for _, u64Map := range m.Uint64Maps {
		u64Map.Name = TitleCaser.String(u64Map.Name)

		if !ValidPrimitiveType(strings.ToLower(u64Map.Value)) {
			u64Map.Value = TitleCaser.String(u64Map.Value)
		} else {
			u64Map.Value = strings.ToLower(u64Map.Value)
		}
	}

	for _, f32 := range m.Float32s {
		f32.Name = TitleCaser.String(f32.Name)
	}

	for _, f32Array := range m.Float32Arrays {
		f32Array.Name = TitleCaser.String(f32Array.Name)
	}

	for _, f64 := range m.Float64s {
		f64.Name = TitleCaser.String(f64.Name)
	}

	for _, f64Array := range m.Float64Arrays {
		f64Array.Name = TitleCaser.String(f64Array.Name)
	}

	for _, b := range m.Bools {
		b.Name = TitleCaser.String(b.Name)
	}

	for _, bArray := range m.BoolArrays {
		bArray.Name = TitleCaser.String(bArray.Name)
	}

	for _, b := range m.Bytes {
		b.Name = TitleCaser.String(b.Name)
	}

	for _, bArray := range m.BytesArrays {
		bArray.Name = TitleCaser.String(bArray.Name)
	}

}

func (m *ModelSchema) Validate(knownModels map[string]struct{}, enums []*EnumSchema) error {
	if !ValidLabel.MatchString(m.Name) {
		return fmt.Errorf("invalid model name: %s", m.Name)
	}

	if _, ok := knownModels[m.Name]; ok {
		return fmt.Errorf("duplicate model name: %s", m.Name)
	}
	knownModels[m.Name] = struct{}{}

	knownFields := make(map[string]struct{})
	for _, modelReference := range m.Models {
		err := modelReference.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[modelReference.Name]; ok {
			return fmt.Errorf("duplicate %s.model name: %s", m.Name, modelReference.Name)
		}
		knownFields[modelReference.Name] = struct{}{}
	}

	for _, modelReferenceArray := range m.ModelArrays {
		err := modelReferenceArray.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[modelReferenceArray.Name]; ok {
			return fmt.Errorf("duplicate %s.model_array name: %s", m.Name, modelReferenceArray.Name)
		}
		knownFields[modelReferenceArray.Name] = struct{}{}
	}

	for _, str := range m.Strings {
		err := str.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[str.Name]; ok {
			return fmt.Errorf("duplicate %s.string name: %s", m.Name, str.Name)
		}
		knownFields[str.Name] = struct{}{}
	}

	for _, strArray := range m.StringArrays {
		err := strArray.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[strArray.Name]; ok {
			return fmt.Errorf("duplicate %s.string_array name: %s", m.Name, strArray.Name)
		}
		knownFields[strArray.Name] = struct{}{}
	}

	for _, strMap := range m.StringMaps {
		err := strMap.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[strMap.Name]; ok {
			return fmt.Errorf("duplicate %s.string_map name: %s", m.Name, strMap.Name)
		}
		knownFields[strMap.Name] = struct{}{}
	}

	for _, i32 := range m.Int32s {
		err := i32.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[i32.Name]; ok {
			return fmt.Errorf("duplicate %s.i32 name: %s", m.Name, i32.Name)
		}
		knownFields[i32.Name] = struct{}{}
	}

	for _, i32Array := range m.Int32Arrays {
		err := i32Array.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[i32Array.Name]; ok {
			return fmt.Errorf("duplicate %s.i32_array name: %s", m.Name, i32Array.Name)
		}
		knownFields[i32Array.Name] = struct{}{}
	}

	for _, i32Map := range m.Int32Maps {
		err := i32Map.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[i32Map.Name]; ok {
			return fmt.Errorf("duplicate %s.i32_map name: %s", m.Name, i32Map.Name)
		}
		knownFields[i32Map.Name] = struct{}{}
	}

	for _, i64 := range m.Int64s {
		err := i64.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[i64.Name]; ok {
			return fmt.Errorf("duplicate %s.i64 name: %s", m.Name, i64.Name)
		}
		knownFields[i64.Name] = struct{}{}
	}

	for _, i64Array := range m.Int64Arrays {
		err := i64Array.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[i64Array.Name]; ok {
			return fmt.Errorf("duplicate %s.i64_array name: %s", m.Name, i64Array.Name)
		}
		knownFields[i64Array.Name] = struct{}{}
	}

	for _, i64Map := range m.Int64Maps {
		err := i64Map.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[i64Map.Name]; ok {
			return fmt.Errorf("duplicate %s.i64_map name: %s", m.Name, i64Map.Name)
		}
		knownFields[i64Map.Name] = struct{}{}
	}

	for _, u32 := range m.Uint32s {
		err := u32.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[u32.Name]; ok {
			return fmt.Errorf("duplicate %s.u32 name: %s", m.Name, u32.Name)
		}
		knownFields[u32.Name] = struct{}{}
	}

	for _, u32Array := range m.Uint32Arrays {
		err := u32Array.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[u32Array.Name]; ok {
			return fmt.Errorf("duplicate %s.u32_array name: %s", m.Name, u32Array.Name)
		}
		knownFields[u32Array.Name] = struct{}{}
	}

	for _, u32Map := range m.Uint32Maps {
		err := u32Map.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[u32Map.Name]; ok {
			return fmt.Errorf("duplicate %s.u32_map name: %s", m.Name, u32Map.Name)
		}
		knownFields[u32Map.Name] = struct{}{}
	}

	for _, u64 := range m.Uint64s {
		err := u64.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[u64.Name]; ok {
			return fmt.Errorf("duplicate %s.u64 name: %s", m.Name, u64.Name)
		}
		knownFields[u64.Name] = struct{}{}
	}

	for _, u64Array := range m.Uint64Arrays {
		err := u64Array.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[u64Array.Name]; ok {
			return fmt.Errorf("duplicate %s.u64_array name: %s", m.Name, u64Array.Name)
		}
		knownFields[u64Array.Name] = struct{}{}
	}

	for _, u64Map := range m.Uint64Maps {
		err := u64Map.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[u64Map.Name]; ok {
			return fmt.Errorf("duplicate %s.u64_map name: %s", m.Name, u64Map.Name)
		}
		knownFields[u64Map.Name] = struct{}{}
	}

	for _, f32 := range m.Float32s {
		err := f32.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[f32.Name]; ok {
			return fmt.Errorf("duplicate %s.f32 name: %s", m.Name, f32.Name)
		}
		knownFields[f32.Name] = struct{}{}
	}

	for _, f32Array := range m.Float32Arrays {
		err := f32Array.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[f32Array.Name]; ok {
			return fmt.Errorf("duplicate %s.f32_array name: %s", m.Name, f32Array.Name)
		}
		knownFields[f32Array.Name] = struct{}{}
	}

	for _, f64 := range m.Float64s {
		err := f64.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[f64.Name]; ok {
			return fmt.Errorf("duplicate %s.f64 name: %s", m.Name, f64.Name)
		}
		knownFields[f64.Name] = struct{}{}
	}

	for _, f64Array := range m.Float64Arrays {
		err := f64Array.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[f64Array.Name]; ok {
			return fmt.Errorf("duplicate %s.f64_array name: %s", m.Name, f64Array.Name)
		}
		knownFields[f64Array.Name] = struct{}{}
	}

	for _, b := range m.Bools {
		err := b.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[b.Name]; ok {
			return fmt.Errorf("duplicate %s.bool name: %s", m.Name, b.Name)
		}
		knownFields[b.Name] = struct{}{}
	}

	for _, bArray := range m.BoolArrays {
		err := bArray.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[bArray.Name]; ok {
			return fmt.Errorf("duplicate %s.bool_array name: %s", m.Name, bArray.Name)
		}
		knownFields[bArray.Name] = struct{}{}
	}

	for _, b := range m.Bytes {
		err := b.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[b.Name]; ok {
			return fmt.Errorf("duplicate %s.bytes name: %s", m.Name, b.Name)
		}
		knownFields[b.Name] = struct{}{}
	}

	for _, bArray := range m.BytesArrays {
		err := bArray.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[bArray.Name]; ok {
			return fmt.Errorf("duplicate %s.bytes_array name: %s", m.Name, bArray.Name)
		}
		knownFields[bArray.Name] = struct{}{}
	}

	for _, enum := range m.Enums {
		err := enum.Validate(m, enums)
		if err != nil {
			return err
		}

		if _, ok := knownFields[enum.Name]; ok {
			return fmt.Errorf("duplicate %s.enum name: %s", m.Name, enum.Name)
		}
		knownFields[enum.Name] = struct{}{}
	}

	for _, enumArray := range m.EnumArrays {
		err := enumArray.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[enumArray.Name]; ok {
			return fmt.Errorf("duplicate %s.enum_array name: %s", m.Name, enumArray.Name)
		}
		knownFields[enumArray.Name] = struct{}{}
	}

	for _, enumMap := range m.EnumMaps {
		err := enumMap.Validate(m)
		if err != nil {
			return err
		}

		if _, ok := knownFields[enumMap.Name]; ok {
			return fmt.Errorf("duplicate %s.enum_map name: %s", m.Name, enumMap.Name)
		}
		knownFields[enumMap.Name] = struct{}{}
	}

	return nil
}

type ModelReferenceSchema struct {
	Name      string `hcl:"name,label"`
	Reference string `hcl:"reference,attr"`
	Accessor  bool   `hcl:"accessor,optional"`
}

func (m *ModelReferenceSchema) Validate(model *ModelSchema) error {
	if !ValidLabel.MatchString(m.Name) {
		return fmt.Errorf("invalid %s.model name: %s", model.Name, m.Name)
	}

	if !ValidLabel.MatchString(m.Reference) {
		return fmt.Errorf("invalid %s.%s.reference: %s", model.Name, m.Name, m.Reference)
	}

	return nil
}

type ModelReferenceArraySchema struct {
	Name        string `hcl:"name,label"`
	Reference   string `hcl:"reference,attr"`
	InitialSize uint32 `hcl:"initial_size,attr"`
	Accessor    bool   `hcl:"accessor,optional"`
}

func (m *ModelReferenceArraySchema) Validate(model *ModelSchema) error {
	if !ValidLabel.MatchString(m.Name) {
		return fmt.Errorf("invalid %s.model name: %s", model.Name, m.Name)
	}

	if !ValidLabel.MatchString(m.Reference) {
		return fmt.Errorf("invalid %s.%s.reference: %s", model.Name, m.Name, m.Reference)
	}

	return nil
}

type ModelReferenceMapSchema struct {
	Name      string `hcl:"name,label"`
	Reference string `hcl:"reference,attr"`
	Value     string `hcl:"value,attr"`
	Accessor  bool   `hcl:"accessor,optional"`
}

func (m *ModelReferenceMapSchema) Validate(model *ModelSchema) error {
	if !ValidLabel.MatchString(m.Name) {
		return fmt.Errorf("invalid %s.model name: %s", model.Name, m.Name)
	}

	if !ValidLabel.MatchString(m.Reference) {
		return fmt.Errorf("invalid %s.%s.reference: %s", model.Name, m.Name, m.Reference)
	}

	return nil
}
