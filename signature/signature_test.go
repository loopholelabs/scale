//go:build !integration

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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSchema(t *testing.T) {
	s := new(Schema)
	err := s.Decode([]byte(MasterTestingSchema))
	require.NoError(t, err)

	assert.Equal(t, V1AlphaVersion, s.Version)
	assert.Equal(t, "MasterSchema", s.Name)
	assert.Equal(t, "MasterSchemaTag", s.Tag)
	assert.Equal(t, "ModelWithAllFieldTypes", s.Context)

	assert.Equal(t, "EmptyModel", s.Models[0].Name)

	assert.Equal(t, "EmptyModelWithDescription", s.Models[1].Name)
	assert.Equal(t, "Test Description", s.Models[1].Description)

	assert.Equal(t, "ModelWithSingleStringField", s.Models[2].Name)
	assert.Equal(t, "StringField", s.Models[2].Strings[0].Name)
	assert.Equal(t, "DefaultValue", s.Models[2].Strings[0].Default)

	assert.Equal(t, "ModelWithSingleStringFieldAndDescription", s.Models[3].Name)
	assert.Equal(t, "Test Description", s.Models[3].Description)
	assert.Equal(t, "StringField", s.Models[3].Strings[0].Name)
	assert.Equal(t, "DefaultValue", s.Models[3].Strings[0].Default)

	assert.Equal(t, "ModelWithSingleInt32Field", s.Models[4].Name)
	assert.Equal(t, "Int32Field", s.Models[4].Int32s[0].Name)
	assert.Equal(t, int32(32), s.Models[4].Int32s[0].Default)

	assert.Equal(t, "ModelWithSingleInt32FieldAndDescription", s.Models[5].Name)
	assert.Equal(t, "Test Description", s.Models[5].Description)
	assert.Equal(t, "Int32Field", s.Models[5].Int32s[0].Name)
	assert.Equal(t, int32(32), s.Models[5].Int32s[0].Default)

	assert.Equal(t, "ModelWithMultipleFields", s.Models[6].Name)
	assert.Equal(t, "StringField", s.Models[6].Strings[0].Name)
	assert.Equal(t, "DefaultValue", s.Models[6].Strings[0].Default)
	assert.Equal(t, "Int32Field", s.Models[6].Int32s[0].Name)
	assert.Equal(t, int32(32), s.Models[6].Int32s[0].Default)

	assert.Equal(t, "ModelWithMultipleFieldsAndDescription", s.Models[7].Name)
	assert.Equal(t, "Test Description", s.Models[7].Description)
	assert.Equal(t, "StringField", s.Models[7].Strings[0].Name)
	assert.Equal(t, "DefaultValue", s.Models[7].Strings[0].Default)
	assert.Equal(t, "Int32Field", s.Models[7].Int32s[0].Name)
	assert.Equal(t, int32(32), s.Models[7].Int32s[0].Default)

	assert.Equal(t, "GenericEnum", s.Enums[0].Name)
	assert.Equal(t, []string{"FirstValue", "SecondValue", "DefaultValue"}, s.Enums[0].Values)

	assert.Equal(t, "ModelWithEnum", s.Models[8].Name)
	assert.Equal(t, "EnumField", s.Models[8].Enums[0].Name)
	assert.Equal(t, "GenericEnum", s.Models[8].Enums[0].Reference)

	assert.Equal(t, "ModelWithEnumAndDescription", s.Models[9].Name)
	assert.Equal(t, "Test Description", s.Models[9].Description)
	assert.Equal(t, "EnumField", s.Models[9].Enums[0].Name)
	assert.Equal(t, "GenericEnum", s.Models[9].Enums[0].Reference)

	assert.Equal(t, "ModelWithEnumAccessor", s.Models[10].Name)
	assert.Equal(t, "EnumField", s.Models[10].Enums[0].Name)
	assert.Equal(t, "GenericEnum", s.Models[10].Enums[0].Reference)
	assert.Equal(t, true, s.Models[10].Enums[0].Accessor)

	assert.Equal(t, "ModelWithEnumAccessorAndDescription", s.Models[11].Name)
	assert.Equal(t, "Test Description", s.Models[11].Description)
	assert.Equal(t, "EnumField", s.Models[11].Enums[0].Name)
	assert.Equal(t, "GenericEnum", s.Models[11].Enums[0].Reference)
	assert.Equal(t, true, s.Models[11].Enums[0].Accessor)

	assert.Equal(t, "ModelWithMultipleFieldsAccessor", s.Models[12].Name)
	assert.Equal(t, "StringField", s.Models[12].Strings[0].Name)
	assert.Equal(t, "DefaultValue", s.Models[12].Strings[0].Default)
	assert.Equal(t, true, *s.Models[12].Strings[0].Accessor)
	assert.Equal(t, "Int32Field", s.Models[12].Int32s[0].Name)
	assert.Equal(t, int32(32), s.Models[12].Int32s[0].Default)
	assert.Equal(t, true, *s.Models[12].Int32s[0].Accessor)

	assert.Equal(t, "ModelWithMultipleFieldsAccessorAndDescription", s.Models[13].Name)
	assert.Equal(t, "Test Description", s.Models[13].Description)
	assert.Equal(t, "StringField", s.Models[13].Strings[0].Name)
	assert.Equal(t, "DefaultValue", s.Models[13].Strings[0].Default)
	assert.Equal(t, true, *s.Models[13].Strings[0].Accessor)
	assert.Equal(t, "Int32Field", s.Models[13].Int32s[0].Name)
	assert.Equal(t, int32(32), s.Models[13].Int32s[0].Default)
	assert.Equal(t, true, *s.Models[13].Int32s[0].Accessor)

	assert.Equal(t, "ModelWithEmbeddedModels", s.Models[14].Name)
	assert.Equal(t, "EmbeddedEmptyModel", s.Models[14].Models[0].Name)
	assert.Equal(t, "EmptyModel", s.Models[14].Models[0].Reference)
	assert.Equal(t, "EmbeddedModelArrayWithMultipleFieldsAccessor", s.Models[14].ModelArrays[0].Name)
	assert.Equal(t, "ModelWithMultipleFieldsAccessor", s.Models[14].ModelArrays[0].Reference)
	assert.Equal(t, uint32(64), s.Models[14].ModelArrays[0].InitialSize)

	assert.Equal(t, "ModelWithEmbeddedModelsAndDescription", s.Models[15].Name)
	assert.Equal(t, "Test Description", s.Models[15].Description)
	assert.Equal(t, "EmbeddedEmptyModel", s.Models[15].Models[0].Name)
	assert.Equal(t, "EmptyModel", s.Models[15].Models[0].Reference)
	assert.Equal(t, "EmbeddedModelArrayWithMultipleFieldsAccessor", s.Models[15].ModelArrays[0].Name)
	assert.Equal(t, "ModelWithMultipleFieldsAccessor", s.Models[15].ModelArrays[0].Reference)
	assert.Equal(t, uint32(0), s.Models[15].ModelArrays[0].InitialSize)

	assert.Equal(t, "ModelWithEmbeddedModelsAccessor", s.Models[16].Name)
	assert.Equal(t, "EmbeddedEmptyModel", s.Models[16].Models[0].Name)
	assert.Equal(t, "EmptyModel", s.Models[16].Models[0].Reference)
	assert.Equal(t, true, s.Models[16].Models[0].Accessor)
	assert.Equal(t, "EmbeddedModelArrayWithMultipleFieldsAccessor", s.Models[16].ModelArrays[0].Name)
	assert.Equal(t, "ModelWithMultipleFieldsAccessor", s.Models[16].ModelArrays[0].Reference)
	assert.Equal(t, uint32(0), s.Models[16].ModelArrays[0].InitialSize)
	assert.Equal(t, true, s.Models[16].ModelArrays[0].Accessor)

	assert.Equal(t, "ModelWithEmbeddedModelsAccessorAndDescription", s.Models[17].Name)
	assert.Equal(t, "Test Description", s.Models[17].Description)
	assert.Equal(t, "EmbeddedEmptyModel", s.Models[17].Models[0].Name)
	assert.Equal(t, "EmptyModel", s.Models[17].Models[0].Reference)
	assert.Equal(t, true, s.Models[17].Models[0].Accessor)
	assert.Equal(t, "EmbeddedModelArrayWithMultipleFieldsAccessor", s.Models[17].ModelArrays[0].Name)
	assert.Equal(t, "ModelWithMultipleFieldsAccessor", s.Models[17].ModelArrays[0].Reference)
	assert.Equal(t, uint32(0), s.Models[17].ModelArrays[0].InitialSize)
	assert.Equal(t, true, s.Models[17].ModelArrays[0].Accessor)

	assert.Equal(t, "ModelWithAllFieldTypes", s.Models[18].Name)
	assert.Equal(t, "StringField", s.Models[18].Strings[0].Name)
	assert.Equal(t, "DefaultValue", s.Models[18].Strings[0].Default)
	assert.Equal(t, "StringArrayField", s.Models[18].StringArrays[0].Name)
	assert.Equal(t, uint32(0), s.Models[18].StringArrays[0].InitialSize)
	assert.Equal(t, "StringMapField", s.Models[18].StringMaps[0].Name)
	assert.Equal(t, "string", s.Models[18].StringMaps[0].Value)
	assert.Equal(t, "StringMapFieldEmbedded", s.Models[18].StringMaps[1].Name)
	assert.Equal(t, "EmptyModel", s.Models[18].StringMaps[1].Value)
	assert.Equal(t, "Int32Field", s.Models[18].Int32s[0].Name)
	assert.Equal(t, int32(32), s.Models[18].Int32s[0].Default)
	assert.Equal(t, "Int32ArrayField", s.Models[18].Int32Arrays[0].Name)
	assert.Equal(t, uint32(0), s.Models[18].Int32Arrays[0].InitialSize)
	assert.Equal(t, "Int32MapField", s.Models[18].Int32Maps[0].Name)
	assert.Equal(t, "int32", s.Models[18].Int32Maps[0].Value)
	assert.Equal(t, "Int32MapFieldEmbedded", s.Models[18].Int32Maps[1].Name)
	assert.Equal(t, "EmptyModel", s.Models[18].Int32Maps[1].Value)
	assert.Equal(t, "Int64Field", s.Models[18].Int64s[0].Name)
	assert.Equal(t, int64(64), s.Models[18].Int64s[0].Default)
	assert.Equal(t, "Int64ArrayField", s.Models[18].Int64Arrays[0].Name)
	assert.Equal(t, uint32(0), s.Models[18].Int64Arrays[0].InitialSize)
	assert.Equal(t, "Int64MapField", s.Models[18].Int64Maps[0].Name)
	assert.Equal(t, "int64", s.Models[18].Int64Maps[0].Value)
	assert.Equal(t, "Int64MapFieldEmbedded", s.Models[18].Int64Maps[1].Name)
	assert.Equal(t, "EmptyModel", s.Models[18].Int64Maps[1].Value)
	assert.Equal(t, "Float32Field", s.Models[18].Float32s[0].Name)
	assert.Equal(t, float32(32.32), s.Models[18].Float32s[0].Default)
	assert.Equal(t, "Float32ArrayField", s.Models[18].Float32Arrays[0].Name)
	assert.Equal(t, uint32(0), s.Models[18].Float32Arrays[0].InitialSize)
	assert.Equal(t, "Float64Field", s.Models[18].Float64s[0].Name)
	assert.Equal(t, float64(64.64), s.Models[18].Float64s[0].Default)
	assert.Equal(t, "Float64ArrayField", s.Models[18].Float64Arrays[0].Name)
	assert.Equal(t, uint32(0), s.Models[18].Float64Arrays[0].InitialSize)
	assert.Equal(t, "BoolField", s.Models[18].Bools[0].Name)
	assert.Equal(t, true, s.Models[18].Bools[0].Default)
	assert.Equal(t, "BoolArrayField", s.Models[18].BoolArrays[0].Name)
	assert.Equal(t, uint32(0), s.Models[18].BoolArrays[0].InitialSize)
	assert.Equal(t, "BytesField", s.Models[18].Bytes[0].Name)
	assert.Equal(t, uint32(512), s.Models[18].Bytes[0].InitialSize)
	assert.Equal(t, "BytesArrayField", s.Models[18].BytesArrays[0].Name)
	assert.Equal(t, uint32(0), s.Models[18].BytesArrays[0].InitialSize)
	assert.Equal(t, "EnumField", s.Models[18].Enums[0].Name)
	assert.Equal(t, "GenericEnum", s.Models[18].Enums[0].Reference)
	assert.Equal(t, "EnumArrayField", s.Models[18].EnumArrays[0].Name)
	assert.Equal(t, "GenericEnum", s.Models[18].EnumArrays[0].Reference)
	assert.Equal(t, uint32(0), s.Models[18].EnumArrays[0].InitialSize)
	assert.Equal(t, "EnumMapField", s.Models[18].EnumMaps[0].Name)
	assert.Equal(t, "GenericEnum", s.Models[18].EnumMaps[0].Reference)
	assert.Equal(t, "string", s.Models[18].EnumMaps[0].Value)
	assert.Equal(t, "EnumMapFieldEmbedded", s.Models[18].EnumMaps[1].Name)
	assert.Equal(t, "EmptyModel", s.Models[18].EnumMaps[1].Value)
	assert.Equal(t, "ModelField", s.Models[18].Models[0].Name)
	assert.Equal(t, "EmptyModel", s.Models[18].Models[0].Reference)
	assert.Equal(t, "ModelArrayField", s.Models[18].ModelArrays[0].Name)
	assert.Equal(t, "EmptyModel", s.Models[18].ModelArrays[0].Reference)
	assert.Equal(t, uint32(0), s.Models[18].ModelArrays[0].InitialSize)
}
