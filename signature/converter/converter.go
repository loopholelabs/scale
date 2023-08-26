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

// Package converter generates a polyglot-encoded buffer from a signature schema
// and a data payload
package converter

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/loopholelabs/polyglot"
	"github.com/loopholelabs/scale/signature"
	"strconv"
)

var (
	ErrInvalidSchema = errors.New("invalid schema")
	ErrInvalidData   = errors.New("invalid data")
)

// Converter converts a signature schema and data payload into a polyglot-encoded buffer
type Converter struct {
	signature *signature.Schema
	models    map[string]*signature.ModelSchema
	enums     map[string]*signature.EnumSchema
	ctxName   string
	ctxModel  *signature.ModelSchema
}

func Convert(schema *signature.Schema, data map[string]interface{}, encoder *polyglot.BufferEncoder) error {
	p, err := New(schema)
	if err != nil {
		return err
	}
	return p.Convert(data, encoder)
}

func New(schema *signature.Schema) (*Converter, error) {
	p := &Converter{
		signature: schema,
		models:    make(map[string]*signature.ModelSchema),
		enums:     make(map[string]*signature.EnumSchema),
	}

	p.ctxName = p.signature.Context
	if p.ctxName == "" {
		return nil, fmt.Errorf("%w: missing context", ErrInvalidSchema)
	}

	for _, model := range p.signature.Models {
		p.models[model.Name] = model
		if model.Name == p.ctxName {
			p.ctxModel = model
			break
		}
	}
	if p.ctxModel == nil {
		return nil, fmt.Errorf("%w: missing context model", ErrInvalidSchema)
	}

	for _, enum := range p.signature.Enums {
		p.enums[enum.Name] = enum
	}

	return p, nil
}

func (p *Converter) Convert(data map[string]interface{}, encoder *polyglot.BufferEncoder) error {
	ctx, ok := data[p.ctxName]
	if !ok {
		return ErrInvalidData
	}

	ctxMap, ok := ctx.(map[string]interface{})
	if !ok {
		return ErrInvalidData
	}

	err := p.encodeModel(p.ctxModel, ctxMap, encoder)
	if err != nil {
		return fmt.Errorf("%w: error encoding context: %w", ErrInvalidData, err)
	}

	return nil
}

func (p *Converter) encodeModel(model *signature.ModelSchema, data map[string]interface{}, encoder *polyglot.BufferEncoder) (err error) {
	for _, m := range model.Models {
		modelData, ok := data[m.Name]
		if !ok {
			return fmt.Errorf("%w: missing model data", ErrInvalidData)
		}
		modelDataMap, ok := modelData.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid model data", ErrInvalidData)
		}

		schema, ok := p.models[m.Reference]
		if !ok {
			return fmt.Errorf("%w: missing model reference schema", ErrInvalidSchema)
		}

		err = p.encodeModel(schema, modelDataMap, encoder)
		if err != nil {
			return fmt.Errorf("%w: error encoding model %s: %w", ErrInvalidData, m.Name, err)
		}
	}

	for _, a := range model.ModelArrays {
		arrayData, ok := data[a.Name]
		if !ok {
			return fmt.Errorf("%w: missing model array data", ErrInvalidData)
		}

		arrayDataSlice, ok := arrayData.([]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid model array data", ErrInvalidData)
		}

		schema, ok := p.models[a.Reference]
		if !ok {
			return fmt.Errorf("%w: missing model array reference schema", ErrInvalidSchema)
		}

		encoder.Slice(uint32(len(arrayDataSlice)), polyglot.AnyKind)
		for _, ad := range arrayDataSlice {
			arrayDataMap, ok := ad.(map[string]interface{})
			if !ok {
				return fmt.Errorf("%w: invalid model array data", ErrInvalidData)
			}

			err = p.encodeModel(schema, arrayDataMap, encoder)
			if err != nil {
				return fmt.Errorf("%w: error encoding model array %s: %w", ErrInvalidData, a.Name, err)
			}
		}
	}

	for _, s := range model.Strings {
		stringData, ok := data[s.Name]
		if !ok {
			return fmt.Errorf("%w: missing string data", ErrInvalidData)
		}

		stringDataString, ok := stringData.(string)
		if !ok {
			return fmt.Errorf("%w: invalid string data", ErrInvalidData)
		}

		encoder.String(stringDataString)
	}

	for _, sa := range model.StringArrays {
		arrayData, ok := data[sa.Name]
		if !ok {
			return fmt.Errorf("%w: missing string array data", ErrInvalidData)
		}

		arrayDataSlice, ok := arrayData.([]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid string array data", ErrInvalidData)
		}

		encoder.Slice(uint32(len(arrayDataSlice)), polyglot.StringKind)
		for _, ad := range arrayDataSlice {
			j, ok := ad.(string)
			if !ok {
				return fmt.Errorf("%w: invalid string array data", ErrInvalidData)
			}
			encoder.String(j)
		}
	}

	for _, sm := range model.StringMaps {
		mapData, ok := data[sm.Name]
		if !ok {
			return fmt.Errorf("%w: missing string map data", ErrInvalidData)
		}

		mapDataMap, ok := mapData.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid string map data", ErrInvalidData)
		}

		err = encodeMap[string](p, polyglot.StringKind, sm.Value, mapDataMap, encoder.String, encoder)
		if err != nil {
			return err
		}
	}

	for _, i := range model.Int32s {
		int32Data, ok := data[i.Name]
		if !ok {
			return fmt.Errorf("%w: missing int32 data", ErrInvalidData)
		}

		int32DataInt, ok := int32Data.(float64)
		if !ok {
			return fmt.Errorf("%w: invalid int32 data", ErrInvalidData)
		}

		encoder.Int32(int32(int32DataInt))
	}

	for _, ia := range model.Int32Arrays {
		arrayData, ok := data[ia.Name]
		if !ok {
			return fmt.Errorf("%w: missing int32 array data", ErrInvalidData)
		}

		arrayDataSlice, ok := arrayData.([]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid int32 array data", ErrInvalidData)
		}

		encoder.Slice(uint32(len(arrayDataSlice)), polyglot.Int32Kind)
		for _, ad := range arrayDataSlice {
			j, ok := ad.(float64)
			if !ok {
				return fmt.Errorf("%w: invalid int32 array data", ErrInvalidData)
			}
			encoder.Int32(int32(j))
		}
	}

	for _, im := range model.Int32Maps {
		mapData, ok := data[im.Name]
		if !ok {
			return fmt.Errorf("%w: missing int32 map data", ErrInvalidData)
		}

		mapDataMap, ok := mapData.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid int32 map data", ErrInvalidData)
		}

		convertedMapDataMap := make(map[int32]interface{}, len(mapDataMap))
		for k, v := range mapDataMap {
			i, err := strconv.ParseInt(k, 10, 32)
			if err != nil {
				return fmt.Errorf("%w: invalid int32 map data", ErrInvalidData)
			}
			convertedMapDataMap[int32(i)] = v
		}

		err = encodeMap[int32](p, polyglot.Int32Kind, im.Value, convertedMapDataMap, encoder.Int32, encoder)
		if err != nil {
			return err
		}
	}

	for _, i := range model.Int64s {
		int64Data, ok := data[i.Name]
		if !ok {
			return fmt.Errorf("%w: missing int64 data", ErrInvalidData)
		}

		int64DataInt, ok := int64Data.(float64)
		if !ok {
			return fmt.Errorf("%w: invalid int64 data", ErrInvalidData)
		}

		encoder.Int64(int64(int64DataInt))
	}

	for _, ia := range model.Int64Arrays {
		arrayData, ok := data[ia.Name]
		if !ok {
			return fmt.Errorf("%w: missing int64 array data", ErrInvalidData)
		}

		arrayDataSlice, ok := arrayData.([]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid int64 array data", ErrInvalidData)
		}

		encoder.Slice(uint32(len(arrayDataSlice)), polyglot.Int64Kind)
		for _, ad := range arrayDataSlice {
			j, ok := ad.(float64)
			if !ok {
				return fmt.Errorf("%w: invalid int64 array data", ErrInvalidData)
			}
			encoder.Int64(int64(j))
		}
	}

	for _, im := range model.Int64Maps {
		mapData, ok := data[im.Name]
		if !ok {
			return fmt.Errorf("%w: missing int64 map data", ErrInvalidData)
		}

		mapDataMap, ok := mapData.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid int64 map data", ErrInvalidData)
		}

		convertedMapDataMap := make(map[int64]interface{}, len(mapDataMap))
		for k, v := range mapDataMap {
			i, err := strconv.ParseInt(k, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: invalid int64 map data", ErrInvalidData)
			}
			convertedMapDataMap[i] = v
		}

		err = encodeMap[int64](p, polyglot.Int64Kind, im.Value, convertedMapDataMap, encoder.Int64, encoder)
		if err != nil {
			return err
		}
	}

	for _, u := range model.Uint32s {
		uint32Data, ok := data[u.Name]
		if !ok {
			return fmt.Errorf("%w: missing uint32 data", ErrInvalidData)
		}

		uint32DataInt, ok := uint32Data.(float64)
		if !ok {
			return fmt.Errorf("%w: invalid uint32 data", ErrInvalidData)
		}

		encoder.Uint32(uint32(uint32DataInt))
	}

	for _, ua := range model.Uint32Arrays {
		arrayData, ok := data[ua.Name]
		if !ok {
			return fmt.Errorf("%w: missing uint32 array data", ErrInvalidData)
		}

		arrayDataSlice, ok := arrayData.([]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid uint32 array data", ErrInvalidData)
		}

		encoder.Slice(uint32(len(arrayDataSlice)), polyglot.Uint32Kind)
		for _, ad := range arrayDataSlice {
			j, ok := ad.(float64)
			if !ok {
				return fmt.Errorf("%w: invalid uint32 array data", ErrInvalidData)
			}
			encoder.Uint32(uint32(j))
		}
	}

	for _, um := range model.Uint32Maps {
		mapData, ok := data[um.Name]
		if !ok {
			return fmt.Errorf("%w: missing uint32 map data", ErrInvalidData)
		}

		mapDataMap, ok := mapData.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid uint32 map data", ErrInvalidData)
		}

		convertedMapDataMap := make(map[uint32]interface{}, len(mapDataMap))
		for k, v := range mapDataMap {
			i, err := strconv.ParseUint(k, 10, 32)
			if err != nil {
				return fmt.Errorf("%w: invalid uint32 map data", ErrInvalidData)
			}
			convertedMapDataMap[uint32(i)] = v
		}

		err = encodeMap[uint32](p, polyglot.Uint32Kind, um.Value, convertedMapDataMap, encoder.Uint32, encoder)
		if err != nil {
			return err
		}
	}

	for _, u := range model.Uint64s {
		uint64Data, ok := data[u.Name]
		if !ok {
			return fmt.Errorf("%w: missing uint64 data", ErrInvalidData)
		}

		uint64DataInt, ok := uint64Data.(float64)
		if !ok {
			return fmt.Errorf("%w: invalid uint64 data", ErrInvalidData)
		}

		encoder.Uint64(uint64(uint64DataInt))
	}

	for _, ua := range model.Uint64Arrays {
		arrayData, ok := data[ua.Name]
		if !ok {
			return fmt.Errorf("%w: missing uint64 array data", ErrInvalidData)
		}

		arrayDataSlice, ok := arrayData.([]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid uint64 array data", ErrInvalidData)
		}

		encoder.Slice(uint32(len(arrayDataSlice)), polyglot.Uint64Kind)
		for _, ad := range arrayDataSlice {
			j, ok := ad.(float64)
			if !ok {
				return fmt.Errorf("%w: invalid uint64 array data", ErrInvalidData)
			}
			encoder.Uint64(uint64(j))
		}
	}

	for _, um := range model.Uint64Maps {
		mapData, ok := data[um.Name]
		if !ok {
			return fmt.Errorf("%w: missing uint64 map data", ErrInvalidData)
		}

		mapDataMap, ok := mapData.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid uint64 map data", ErrInvalidData)
		}

		convertedMapDataMap := make(map[uint64]interface{}, len(mapDataMap))
		for k, v := range mapDataMap {
			i, err := strconv.ParseUint(k, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: invalid uint64 map data", ErrInvalidData)
			}
			convertedMapDataMap[uint64(i)] = v
		}

		err = encodeMap[uint64](p, polyglot.Uint64Kind, um.Value, convertedMapDataMap, encoder.Uint64, encoder)
		if err != nil {
			return err
		}
	}

	for _, f := range model.Float32s {
		float32Data, ok := data[f.Name]
		if !ok {
			return fmt.Errorf("%w: missing float32 data", ErrInvalidData)
		}

		float32DataFloat, ok := float32Data.(float64)
		if !ok {
			return fmt.Errorf("%w: invalid float32 data", ErrInvalidData)
		}

		encoder.Float32(float32(float32DataFloat))
	}

	for _, fa := range model.Float32Arrays {
		arrayData, ok := data[fa.Name]
		if !ok {
			return fmt.Errorf("%w: missing float32 array data", ErrInvalidData)
		}

		arrayDataSlice, ok := arrayData.([]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid float32 array data", ErrInvalidData)
		}

		encoder.Slice(uint32(len(arrayDataSlice)), polyglot.Float32Kind)
		for _, ad := range arrayDataSlice {
			j, ok := ad.(float64)
			if !ok {
				return fmt.Errorf("%w: invalid float32 array data", ErrInvalidData)
			}
			encoder.Float32(float32(j))
		}
	}

	for _, f := range model.Float64s {
		float64Data, ok := data[f.Name]
		if !ok {
			return fmt.Errorf("%w: missing float64 data", ErrInvalidData)
		}

		float64DataFloat, ok := float64Data.(float64)
		if !ok {
			return fmt.Errorf("%w: invalid float64 data", ErrInvalidData)
		}

		encoder.Float64(float64DataFloat)
	}

	for _, fa := range model.Float64Arrays {
		arrayData, ok := data[fa.Name]
		if !ok {
			return fmt.Errorf("%w: missing float64 array data", ErrInvalidData)
		}

		arrayDataSlice, ok := arrayData.([]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid float64 array data", ErrInvalidData)
		}

		encoder.Slice(uint32(len(arrayDataSlice)), polyglot.Float64Kind)
		for _, ad := range arrayDataSlice {
			j, ok := ad.(float64)
			if !ok {
				return fmt.Errorf("%w: invalid float64 array data", ErrInvalidData)
			}
			encoder.Float64(j)
		}
	}

	for _, e := range model.Enums {
		enumData, ok := data[e.Name]
		if !ok {
			return fmt.Errorf("%w: missing enum data", ErrInvalidData)
		}

		schema, ok := p.enums[e.Reference]
		if !ok {
			return fmt.Errorf("%w: missing enum reference schema", ErrInvalidSchema)
		}

		enumDataString, ok := enumData.(string)
		if !ok {
			return fmt.Errorf("%w: invalid enum data", ErrInvalidData)
		}

		err = p.encodeEnum(schema, enumDataString, encoder)
		if err != nil {
			return err
		}
	}

	for _, ea := range model.EnumArrays {
		arrayData, ok := data[ea.Name]
		if !ok {
			return fmt.Errorf("%w: missing enum array data", ErrInvalidData)
		}

		arrayDataSlice, ok := arrayData.([]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid enum array data", ErrInvalidData)
		}

		schema, ok := p.enums[ea.Reference]
		if !ok {
			return fmt.Errorf("%w: missing enum reference schema", ErrInvalidSchema)
		}

		encoder.Slice(uint32(len(arrayDataSlice)), polyglot.Uint32Kind)
		for _, ad := range arrayDataSlice {
			j, ok := ad.(string)
			if !ok {
				return fmt.Errorf("%w: invalid enum array data", ErrInvalidData)
			}
			err = p.encodeEnum(schema, j, encoder)
			if err != nil {
				return err
			}
		}
	}

	for _, em := range model.EnumMaps {
		mapData, ok := data[em.Name]
		if !ok {
			return fmt.Errorf("%w: missing enum map data", ErrInvalidData)
		}

		mapDataStringMap, ok := mapData.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid enum map data", ErrInvalidData)
		}

		schema, ok := p.enums[em.Reference]
		if !ok {
			return fmt.Errorf("%w: missing enum reference schema", ErrInvalidSchema)
		}

		mapDataMap := make(map[uint32]interface{})
		found := false
		for k, v := range mapDataStringMap {
			for i, ev := range schema.Values {
				if ev == k {
					mapDataMap[uint32(i)] = v
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("%w: invalid enum map data", ErrInvalidData)
			}
			found = false
		}

		err = encodeMap[uint32](p, polyglot.Uint32Kind, em.Value, mapDataMap, encoder.Uint32, encoder)
		if err != nil {
			return err
		}
	}

	for _, b := range model.Bytes {
		byteData, ok := data[b.Name]
		if !ok {
			return fmt.Errorf("%w: missing byte data", ErrInvalidData)
		}

		byteDataString, ok := byteData.(string)
		if !ok {
			return fmt.Errorf("%w: invalid byte data", ErrInvalidData)
		}

		d, err := base64.StdEncoding.DecodeString(byteDataString)
		if err != nil {
			return fmt.Errorf("%w: invalid byte data", ErrInvalidData)
		}

		encoder.Bytes(d)
	}

	for _, ba := range model.BytesArrays {
		arrayData, ok := data[ba.Name]
		if !ok {
			return fmt.Errorf("%w: missing byte array data", ErrInvalidData)
		}

		arrayDataSlice, ok := arrayData.([]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid byte array data", ErrInvalidData)
		}

		encoder.Slice(uint32(len(arrayDataSlice)), polyglot.BytesKind)
		for _, ad := range arrayDataSlice {
			j, ok := ad.(string)
			if !ok {
				return fmt.Errorf("%w: invalid byte array data", ErrInvalidData)
			}
			d, err := base64.StdEncoding.DecodeString(j)
			if err != nil {
				return fmt.Errorf("%w: invalid byte array data", ErrInvalidData)
			}
			encoder.Bytes(d)
		}
	}

	for _, b := range model.Bools {
		boolData, ok := data[b.Name]
		if !ok {
			return fmt.Errorf("%w: missing bool data", ErrInvalidData)
		}

		boolDataBool, ok := boolData.(bool)
		if !ok {
			return fmt.Errorf("%w: invalid bool data", ErrInvalidData)
		}

		encoder.Bool(boolDataBool)
	}

	for _, ba := range model.BoolArrays {
		arrayData, ok := data[ba.Name]
		if !ok {
			return fmt.Errorf("%w: missing bool array data", ErrInvalidData)
		}

		arrayDataSlice, ok := arrayData.([]interface{})
		if !ok {
			return fmt.Errorf("%w: invalid bool array data", ErrInvalidData)
		}

		encoder.Slice(uint32(len(arrayDataSlice)), polyglot.BoolKind)
		for _, ad := range arrayDataSlice {
			j, ok := ad.(bool)
			if !ok {
				return fmt.Errorf("%w: invalid bool array data", ErrInvalidData)
			}
			encoder.Bool(j)
		}
	}

	return nil
}

func (p *Converter) encodeEnum(enum *signature.EnumSchema, data string, encoder *polyglot.BufferEncoder) (err error) {
	for i, v := range enum.Values {
		if v == data {
			encoder.Uint32(uint32(i))
			return nil
		}
	}
	return fmt.Errorf("%w: invalid enum data", ErrInvalidData)
}

func encodeMap[T comparable](parser *Converter, keyKind polyglot.Kind, valueName string, mapData map[T]interface{}, keyEncoder func(T) *polyglot.BufferEncoder, encoder *polyglot.BufferEncoder) error {
	valueKind := polyglot.AnyKind
	isPrimitive := signature.ValidPrimitiveType(valueName)
	if isPrimitive {
		switch valueName {
		case "string":
			valueKind = polyglot.StringKind
		case "int32":
			valueKind = polyglot.Int32Kind
		case "int64":
			valueKind = polyglot.Int64Kind
		case "uint32":
			valueKind = polyglot.Uint32Kind
		case "uint64":
			valueKind = polyglot.Uint64Kind
		case "float32":
			valueKind = polyglot.Float32Kind
		case "float64":
			valueKind = polyglot.Float64Kind
		case "bool":
			valueKind = polyglot.BoolKind
		case "bytes":
			valueKind = polyglot.BytesKind
		default:
			return fmt.Errorf("%w: invalid primitive map data: %s", ErrInvalidData, valueName)
		}
	}
	encoder.Map(uint32(len(mapData)), keyKind, valueKind)
	for k, v := range mapData {
		keyEncoder(k)
		if isPrimitive {
			switch valueName {
			case "string":
				value, ok := v.(string)
				if !ok {
					return fmt.Errorf("%w: invalid primitive string map data", ErrInvalidData)
				}
				encoder.String(value)
			case "int32":
				value, ok := v.(float64)
				if !ok {
					return fmt.Errorf("%w: invalid primitive int32 map data", ErrInvalidData)
				}
				encoder.Int32(int32(value))
			case "int64":
				value, ok := v.(float64)
				if !ok {
					return fmt.Errorf("%w: invalid primitive int64 map data", ErrInvalidData)
				}
				encoder.Int64(int64(value))
			case "uint32":
				value, ok := v.(float64)
				if !ok {
					return fmt.Errorf("%w: invalid primitive uint32 map data", ErrInvalidData)
				}
				encoder.Uint32(uint32(value))
			case "uint64":
				value, ok := v.(float64)
				if !ok {
					return fmt.Errorf("%w: invalid primitive uint64 map data", ErrInvalidData)
				}
				encoder.Uint64(uint64(value))
			case "float32":
				value, ok := v.(float64)
				if !ok {
					return fmt.Errorf("%w: invalid primitive float32 map data", ErrInvalidData)
				}
				encoder.Float32(float32(value))
			case "float64":
				value, ok := v.(float64)
				if !ok {
					return fmt.Errorf("%w: invalid primitive float64 map data", ErrInvalidData)
				}
				encoder.Float64(value)
			case "bool":
				value, ok := v.(bool)
				if !ok {
					return fmt.Errorf("%w: invalid primitive bool map data", ErrInvalidData)
				}
				encoder.Bool(value)
			case "bytes":
				value, ok := v.(string)
				if !ok {
					return fmt.Errorf("%w: invalid primitive bytes map data", ErrInvalidData)
				}
				valueBytes, err := hex.DecodeString(value)
				if err != nil {
					return fmt.Errorf("%w: invalid primitive bytes map data", ErrInvalidData)
				}
				encoder.Bytes(valueBytes)
			default:
				return fmt.Errorf("%w: invalid primitive map data: %s", ErrInvalidData, valueName)
			}
		} else {
			model, ok := parser.models[valueName]
			if !ok {
				return fmt.Errorf("%w: invalid reference map schema", ErrInvalidData)
			}

			modelDataMap, ok := v.(map[string]interface{})
			if !ok {
				return fmt.Errorf("%w: invalid model data", ErrInvalidData)
			}

			err := parser.encodeModel(model, modelDataMap, encoder)
			if err != nil {
				return fmt.Errorf("%w: error encoding map %s: %w", ErrInvalidData, valueName, err)
			}
		}
	}
	return nil
}
