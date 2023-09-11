//go:build integration && golang

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
	"github.com/loopholelabs/polyglot"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestOutput(t *testing.T) {
	buf := polyglot.NewBuffer()

	var nilModel *EmptyModel
	nilModel.Encode(buf)
	err := os.WriteFile("../../test_data/nil_model.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	emptyModel := NewEmptyModel()
	emptyModel.Encode(buf)
	err = os.WriteFile("../../test_data/empty_model.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	emptyModelWithDescription := NewEmptyModelWithDescription()
	emptyModelWithDescription.Encode(buf)
	err = os.WriteFile("../../test_data/empty_model_with_description.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	modelWithSingleStringField := NewModelWithSingleStringField()
	require.Equal(t, "DefaultValue", modelWithSingleStringField.StringField)
	modelWithSingleStringField.StringField = "hello world"
	modelWithSingleStringField.Encode(buf)
	err = os.WriteFile("../../test_data/model_with_single_string_field.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	modelWithSingleStringFieldAndDescription := NewModelWithSingleStringFieldAndDescription()
	require.Equal(t, "DefaultValue", modelWithSingleStringFieldAndDescription.StringField)
	modelWithSingleStringFieldAndDescription.StringField = "hello world"
	modelWithSingleStringFieldAndDescription.Encode(buf)
	err = os.WriteFile("../../test_data/model_with_single_string_field_and_description.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	modelWithSingleInt32Field := NewModelWithSingleInt32Field()
	require.Equal(t, int32(32), modelWithSingleInt32Field.Int32Field)
	modelWithSingleInt32Field.Int32Field = 42
	modelWithSingleInt32Field.Encode(buf)
	err = os.WriteFile("../../test_data/model_with_single_int32_field.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	modelWithSingleInt32FieldAndDescription := NewModelWithSingleInt32FieldAndDescription()
	require.Equal(t, int32(32), modelWithSingleInt32FieldAndDescription.Int32Field)
	modelWithSingleInt32FieldAndDescription.Int32Field = 42
	modelWithSingleInt32FieldAndDescription.Encode(buf)
	err = os.WriteFile("../../test_data/model_with_single_int32_field_and_description.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	modelWithMultipleFields := NewModelWithMultipleFields()
	require.Equal(t, "DefaultValue", modelWithMultipleFields.StringField)
	require.Equal(t, int32(32), modelWithMultipleFields.Int32Field)
	modelWithMultipleFields.StringField = "hello world"
	modelWithMultipleFields.Int32Field = 42
	modelWithMultipleFields.Encode(buf)
	err = os.WriteFile("../../test_data/model_with_multiple_fields.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	modelWithMultipleFieldsAndDescription := NewModelWithMultipleFieldsAndDescription()
	require.Equal(t, "DefaultValue", modelWithMultipleFieldsAndDescription.StringField)
	require.Equal(t, int32(32), modelWithMultipleFieldsAndDescription.Int32Field)
	modelWithMultipleFieldsAndDescription.StringField = "hello world"
	modelWithMultipleFieldsAndDescription.Int32Field = 42
	modelWithMultipleFieldsAndDescription.Encode(buf)
	err = os.WriteFile("../../test_data/model_with_multiple_fields_and_description.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	modelWithEnum := NewModelWithEnum()
	require.Equal(t, GenericEnumDefaultValue, modelWithEnum.EnumField)
	modelWithEnum.EnumField = GenericEnumSecondValue
	modelWithEnum.Encode(buf)
	err = os.WriteFile("../../test_data/model_with_enum.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	modelWithEnumAndDescription := NewModelWithEnumAndDescription()
	require.Equal(t, GenericEnumDefaultValue, modelWithEnumAndDescription.EnumField)
	modelWithEnumAndDescription.EnumField = GenericEnumSecondValue
	modelWithEnumAndDescription.Encode(buf)
	err = os.WriteFile("../../test_data/model_with_enum_and_description.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	modelWithEnumAccessor := NewModelWithEnumAccessor()
	defaultEnumValue, err := modelWithEnumAccessor.GetEnumField()
	require.NoError(t, err)
	require.Equal(t, GenericEnumDefaultValue, defaultEnumValue)
	err = modelWithEnumAccessor.SetEnumField(GenericEnumSecondValue)
	require.NoError(t, err)
	modelWithEnumAccessor.Encode(buf)
	err = os.WriteFile("../../test_data/model_with_enum_accessor.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	modelWithEnumAccessorAndDescription := NewModelWithEnumAccessorAndDescription()
	defaultEnumValue, err = modelWithEnumAccessorAndDescription.GetEnumField()
	require.NoError(t, err)
	require.Equal(t, GenericEnumDefaultValue, defaultEnumValue)
	err = modelWithEnumAccessorAndDescription.SetEnumField(GenericEnumSecondValue)
	require.NoError(t, err)
	modelWithEnumAccessorAndDescription.Encode(buf)
	err = os.WriteFile("../../test_data/model_with_enum_accessor_and_description.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	modelWithMultipleFieldsAccessor := NewModelWithMultipleFieldsAccessor()
	stringValue, err := modelWithMultipleFieldsAccessor.GetStringField()
	require.NoError(t, err)
	require.Equal(t, "DefaultValue", stringValue)
	err = modelWithMultipleFieldsAccessor.SetStringField("hello world")
	require.ErrorContains(t, err, "value must match ^[a-zA-Z0-9]*$")
	err = modelWithMultipleFieldsAccessor.SetStringField("")
	require.ErrorContains(t, err, "length must be between 1 and 20")
	err = modelWithMultipleFieldsAccessor.SetStringField("hello")
	require.NoError(t, err)
	stringValue, err = modelWithMultipleFieldsAccessor.GetStringField()
	require.NoError(t, err)
	require.Equal(t, strings.ToUpper("hello"), stringValue)
	int32Value, err := modelWithMultipleFieldsAccessor.GetInt32Field()
	require.NoError(t, err)
	require.Equal(t, int32(32), int32Value)
	err = modelWithMultipleFieldsAccessor.SetInt32Field(-1)
	require.ErrorContains(t, err, "value must be between 0 and 100")
	err = modelWithMultipleFieldsAccessor.SetInt32Field(101)
	require.ErrorContains(t, err, "value must be between 0 and 100")
	err = modelWithMultipleFieldsAccessor.SetInt32Field(42)
	require.NoError(t, err)
	modelWithMultipleFieldsAccessor.Encode(buf)
	err = os.WriteFile("../../test_data/model_with_multiple_fields_accessor.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	modelWithMultipleFieldsAccessorAndDescription := NewModelWithMultipleFieldsAccessorAndDescription()
	stringValue, err = modelWithMultipleFieldsAccessorAndDescription.GetStringField()
	require.NoError(t, err)
	require.Equal(t, "DefaultValue", stringValue)
	err = modelWithMultipleFieldsAccessorAndDescription.SetStringField("hello world")
	require.NoError(t, err)
	int32Value, err = modelWithMultipleFieldsAccessorAndDescription.GetInt32Field()
	require.NoError(t, err)
	require.Equal(t, int32(32), int32Value)
	err = modelWithMultipleFieldsAccessorAndDescription.SetInt32Field(42)
	require.NoError(t, err)
	modelWithMultipleFieldsAccessorAndDescription.Encode(buf)
	err = os.WriteFile("../../test_data/model_with_multiple_fields_accessor_and_description.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	modelWithEmbeddedModels := NewModelWithEmbeddedModels()
	require.NotNil(t, modelWithEmbeddedModels.EmbeddedEmptyModel)
	require.NotNil(t, modelWithEmbeddedModels.EmbeddedModelArrayWithMultipleFieldsAccessor)
	require.Equal(t, 64, cap(modelWithEmbeddedModels.EmbeddedModelArrayWithMultipleFieldsAccessor))
	require.Equal(t, 0, len(modelWithEmbeddedModels.EmbeddedModelArrayWithMultipleFieldsAccessor))
	require.IsType(t, []ModelWithMultipleFieldsAccessor{}, modelWithEmbeddedModels.EmbeddedModelArrayWithMultipleFieldsAccessor)
	modelWithEmbeddedModels.EmbeddedModelArrayWithMultipleFieldsAccessor = append(modelWithEmbeddedModels.EmbeddedModelArrayWithMultipleFieldsAccessor, *modelWithMultipleFieldsAccessor)
	modelWithEmbeddedModels.Encode(buf)
	err = os.WriteFile("../../test_data/model_with_embedded_models.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	modelWithEmbeddedModelsAndDescription := NewModelWithEmbeddedModelsAndDescription()
	require.NotNil(t, modelWithEmbeddedModelsAndDescription.EmbeddedEmptyModel)
	require.NotNil(t, modelWithEmbeddedModelsAndDescription.EmbeddedModelArrayWithMultipleFieldsAccessor)
	require.Equal(t, 0, cap(modelWithEmbeddedModelsAndDescription.EmbeddedModelArrayWithMultipleFieldsAccessor))
	require.Equal(t, 0, len(modelWithEmbeddedModelsAndDescription.EmbeddedModelArrayWithMultipleFieldsAccessor))
	require.IsType(t, []ModelWithMultipleFieldsAccessor{}, modelWithEmbeddedModelsAndDescription.EmbeddedModelArrayWithMultipleFieldsAccessor)
	modelWithEmbeddedModelsAndDescription.EmbeddedModelArrayWithMultipleFieldsAccessor = append(modelWithEmbeddedModelsAndDescription.EmbeddedModelArrayWithMultipleFieldsAccessor, *modelWithMultipleFieldsAccessor)
	modelWithEmbeddedModelsAndDescription.Encode(buf)
	err = os.WriteFile("../../test_data/model_with_embedded_models_and_description.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	modelWithEmbeddedModelsAccessor := NewModelWithEmbeddedModelsAccessor()
	embeddedModel, err := modelWithEmbeddedModelsAccessor.GetEmbeddedEmptyModel()
	require.NoError(t, err)
	require.NotNil(t, embeddedModel)
	embeddedModelArray, err := modelWithEmbeddedModelsAccessor.GetEmbeddedModelArrayWithMultipleFieldsAccessor()
	require.NoError(t, err)
	require.NotNil(t, embeddedModelArray)
	require.Equal(t, 0, cap(embeddedModelArray))
	require.Equal(t, 0, len(embeddedModelArray))
	require.IsType(t, []ModelWithMultipleFieldsAccessor{}, embeddedModelArray)
	err = modelWithEmbeddedModelsAccessor.SetEmbeddedModelArrayWithMultipleFieldsAccessor([]ModelWithMultipleFieldsAccessor{*modelWithMultipleFieldsAccessor})
	require.NoError(t, err)
	modelWithEmbeddedModelsAccessor.Encode(buf)
	err = os.WriteFile("../../test_data/model_with_embedded_models_accessor.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	modelWithEmbeddedModelsAccessorAndDescription := NewModelWithEmbeddedModelsAccessorAndDescription()
	embeddedModel, err = modelWithEmbeddedModelsAccessorAndDescription.GetEmbeddedEmptyModel()
	require.NoError(t, err)
	require.NotNil(t, embeddedModel)
	embeddedModelArray, err = modelWithEmbeddedModelsAccessorAndDescription.GetEmbeddedModelArrayWithMultipleFieldsAccessor()
	require.NoError(t, err)
	require.NotNil(t, embeddedModelArray)
	require.Equal(t, 0, cap(embeddedModelArray))
	require.Equal(t, 0, len(embeddedModelArray))
	require.IsType(t, []ModelWithMultipleFieldsAccessor{}, embeddedModelArray)
	err = modelWithEmbeddedModelsAccessorAndDescription.SetEmbeddedModelArrayWithMultipleFieldsAccessor([]ModelWithMultipleFieldsAccessor{*modelWithMultipleFieldsAccessor})
	require.NoError(t, err)
	modelWithEmbeddedModelsAccessorAndDescription.Encode(buf)
	err = os.WriteFile("../../test_data/model_with_embedded_models_accessor_and_description.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()

	modelWithAllFieldTypes := NewModelWithAllFieldTypes()

	require.Equal(t, "DefaultValue", modelWithAllFieldTypes.StringField)
	modelWithAllFieldTypes.StringField = "hello world"
	require.Equal(t, 0, cap(modelWithAllFieldTypes.StringArrayField))
	require.Equal(t, 0, len(modelWithAllFieldTypes.StringArrayField))
	require.IsType(t, []string{}, modelWithAllFieldTypes.StringArrayField)
	modelWithAllFieldTypes.StringArrayField = append(modelWithAllFieldTypes.StringArrayField, "hello", "world")
	require.Equal(t, 0, len(modelWithAllFieldTypes.StringMapField))
	require.IsType(t, map[string]string{}, modelWithAllFieldTypes.StringMapField)
	modelWithAllFieldTypes.StringMapField["hello"] = "world"
	require.Equal(t, 0, len(modelWithAllFieldTypes.StringMapFieldEmbedded))
	require.IsType(t, map[string]EmptyModel{}, modelWithAllFieldTypes.StringMapFieldEmbedded)
	modelWithAllFieldTypes.StringMapFieldEmbedded["hello"] = *emptyModel

	require.Equal(t, int32(32), modelWithAllFieldTypes.Int32Field)
	modelWithAllFieldTypes.Int32Field = 42
	require.Equal(t, 0, cap(modelWithAllFieldTypes.Int32ArrayField))
	require.Equal(t, 0, len(modelWithAllFieldTypes.Int32ArrayField))
	require.IsType(t, []int32{}, modelWithAllFieldTypes.Int32ArrayField)
	modelWithAllFieldTypes.Int32ArrayField = append(modelWithAllFieldTypes.Int32ArrayField, 42, 84)
	require.Equal(t, 0, len(modelWithAllFieldTypes.Int32MapField))
	require.IsType(t, map[int32]int32{}, modelWithAllFieldTypes.Int32MapField)
	modelWithAllFieldTypes.Int32MapField[42] = 84
	require.Equal(t, 0, len(modelWithAllFieldTypes.Int32MapFieldEmbedded))
	require.IsType(t, map[int32]EmptyModel{}, modelWithAllFieldTypes.Int32MapFieldEmbedded)
	modelWithAllFieldTypes.Int32MapFieldEmbedded[42] = *emptyModel

	require.Equal(t, int64(64), modelWithAllFieldTypes.Int64Field)
	modelWithAllFieldTypes.Int64Field = 100
	require.Equal(t, 0, cap(modelWithAllFieldTypes.Int64ArrayField))
	require.Equal(t, 0, len(modelWithAllFieldTypes.Int64ArrayField))
	require.IsType(t, []int64{}, modelWithAllFieldTypes.Int64ArrayField)
	modelWithAllFieldTypes.Int64ArrayField = append(modelWithAllFieldTypes.Int64ArrayField, 100, 200)
	require.Equal(t, 0, len(modelWithAllFieldTypes.Int64MapField))
	require.IsType(t, map[int64]int64{}, modelWithAllFieldTypes.Int64MapField)
	modelWithAllFieldTypes.Int64MapField[100] = 200
	require.Equal(t, 0, len(modelWithAllFieldTypes.Int64MapFieldEmbedded))
	require.IsType(t, map[int64]EmptyModel{}, modelWithAllFieldTypes.Int64MapFieldEmbedded)
	modelWithAllFieldTypes.Int64MapFieldEmbedded[100] = *emptyModel

	require.Equal(t, uint32(32), modelWithAllFieldTypes.Uint32Field)
	modelWithAllFieldTypes.Uint32Field = 42
	require.Equal(t, 0, cap(modelWithAllFieldTypes.Uint32ArrayField))
	require.Equal(t, 0, len(modelWithAllFieldTypes.Uint32ArrayField))
	require.IsType(t, []uint32{}, modelWithAllFieldTypes.Uint32ArrayField)
	modelWithAllFieldTypes.Uint32ArrayField = append(modelWithAllFieldTypes.Uint32ArrayField, 42, 84)
	require.Equal(t, 0, len(modelWithAllFieldTypes.Uint32MapField))
	require.IsType(t, map[uint32]uint32{}, modelWithAllFieldTypes.Uint32MapField)
	modelWithAllFieldTypes.Uint32MapField[42] = 84
	require.Equal(t, 0, len(modelWithAllFieldTypes.Uint32MapFieldEmbedded))
	require.IsType(t, map[uint32]EmptyModel{}, modelWithAllFieldTypes.Uint32MapFieldEmbedded)
	modelWithAllFieldTypes.Uint32MapFieldEmbedded[42] = *emptyModel

	require.Equal(t, uint64(64), modelWithAllFieldTypes.Uint64Field)
	modelWithAllFieldTypes.Uint64Field = 100
	require.Equal(t, 0, cap(modelWithAllFieldTypes.Uint64ArrayField))
	require.Equal(t, 0, len(modelWithAllFieldTypes.Uint64ArrayField))
	require.IsType(t, []uint64{}, modelWithAllFieldTypes.Uint64ArrayField)
	modelWithAllFieldTypes.Uint64ArrayField = append(modelWithAllFieldTypes.Uint64ArrayField, 100, 200)
	require.Equal(t, 0, len(modelWithAllFieldTypes.Uint64MapField))
	require.IsType(t, map[uint64]uint64{}, modelWithAllFieldTypes.Uint64MapField)
	modelWithAllFieldTypes.Uint64MapField[100] = 200
	require.Equal(t, 0, len(modelWithAllFieldTypes.Uint64MapFieldEmbedded))
	require.IsType(t, map[uint64]EmptyModel{}, modelWithAllFieldTypes.Uint64MapFieldEmbedded)
	modelWithAllFieldTypes.Uint64MapFieldEmbedded[100] = *emptyModel

	require.Equal(t, float32(32.32), modelWithAllFieldTypes.Float32Field)
	modelWithAllFieldTypes.Float32Field = 42.0
	require.Equal(t, 0, cap(modelWithAllFieldTypes.Float32ArrayField))
	require.Equal(t, 0, len(modelWithAllFieldTypes.Float32ArrayField))
	require.IsType(t, []float32{}, modelWithAllFieldTypes.Float32ArrayField)
	modelWithAllFieldTypes.Float32ArrayField = append(modelWithAllFieldTypes.Float32ArrayField, 42.0, 84.0)

	require.Equal(t, float64(64.64), modelWithAllFieldTypes.Float64Field)
	modelWithAllFieldTypes.Float64Field = 100.0
	require.Equal(t, 0, cap(modelWithAllFieldTypes.Float64ArrayField))
	require.Equal(t, 0, len(modelWithAllFieldTypes.Float64ArrayField))
	require.IsType(t, []float64{}, modelWithAllFieldTypes.Float64ArrayField)
	modelWithAllFieldTypes.Float64ArrayField = append(modelWithAllFieldTypes.Float64ArrayField, 100.0, 200.0)

	require.Equal(t, true, modelWithAllFieldTypes.BoolField)
	modelWithAllFieldTypes.BoolField = false
	require.Equal(t, 0, cap(modelWithAllFieldTypes.BoolArrayField))
	require.Equal(t, 0, len(modelWithAllFieldTypes.BoolArrayField))
	require.IsType(t, []bool{}, modelWithAllFieldTypes.BoolArrayField)
	modelWithAllFieldTypes.BoolArrayField = append(modelWithAllFieldTypes.BoolArrayField, true, false)

	require.Equal(t, 512, cap(modelWithAllFieldTypes.BytesField))
	require.Equal(t, 0, len(modelWithAllFieldTypes.BytesField))
	require.IsType(t, []byte{}, modelWithAllFieldTypes.BytesField)
	modelWithAllFieldTypes.BytesField = append(modelWithAllFieldTypes.BytesField, []byte{42, 84}...)
	require.Equal(t, 0, len(modelWithAllFieldTypes.BytesArrayField))
	require.IsType(t, [][]byte{}, modelWithAllFieldTypes.BytesArrayField)
	modelWithAllFieldTypes.BytesArrayField = append(modelWithAllFieldTypes.BytesArrayField, []byte{42, 84}, []byte{84, 42})

	require.Equal(t, GenericEnumDefaultValue, modelWithAllFieldTypes.EnumField)
	modelWithAllFieldTypes.EnumField = GenericEnumSecondValue
	require.Equal(t, 0, cap(modelWithAllFieldTypes.EnumArrayField))
	require.Equal(t, 0, len(modelWithAllFieldTypes.EnumArrayField))
	require.IsType(t, []GenericEnum{}, modelWithAllFieldTypes.EnumArrayField)
	modelWithAllFieldTypes.EnumArrayField = append(modelWithAllFieldTypes.EnumArrayField, GenericEnumFirstValue, GenericEnumSecondValue)
	require.Equal(t, 0, len(modelWithAllFieldTypes.EnumMapField))
	require.IsType(t, map[GenericEnum]string{}, modelWithAllFieldTypes.EnumMapField)
	modelWithAllFieldTypes.EnumMapField[GenericEnumFirstValue] = "hello world"
	require.Equal(t, 0, len(modelWithAllFieldTypes.EnumMapFieldEmbedded))
	require.IsType(t, map[GenericEnum]EmptyModel{}, modelWithAllFieldTypes.EnumMapFieldEmbedded)
	modelWithAllFieldTypes.EnumMapFieldEmbedded[GenericEnumFirstValue] = *emptyModel

	require.NotNil(t, modelWithAllFieldTypes.ModelField)
	require.Equal(t, 0, cap(modelWithAllFieldTypes.ModelArrayField))
	require.Equal(t, 0, len(modelWithAllFieldTypes.ModelArrayField))
	require.IsType(t, []EmptyModel{}, modelWithAllFieldTypes.ModelArrayField)
	modelWithAllFieldTypes.ModelArrayField = append(modelWithAllFieldTypes.ModelArrayField, *emptyModel, *emptyModel)

	modelWithAllFieldTypes.Encode(buf)
	err = os.WriteFile("../../test_data/model_with_all_field_types.bin", buf.Bytes(), 0644)
	require.NoError(t, err)
	buf.Reset()
}

func TestInput(t *testing.T) {
	nilModelData, err := os.ReadFile("../../test_data/nil_model.bin")
	require.NoError(t, err)
	nilModel, err := DecodeEmptyModel(nil, nilModelData)
	require.NoError(t, err)
	require.Nil(t, nilModel)

	emptyModelData, err := os.ReadFile("../../test_data/empty_model.bin")
	require.NoError(t, err)
	emptyModel, err := DecodeEmptyModel(nil, emptyModelData)
	require.NoError(t, err)
	require.NotNil(t, emptyModel)

	emptyModelWithDescriptionData, err := os.ReadFile("../../test_data/empty_model_with_description.bin")
	require.NoError(t, err)
	emptyModelWithDescription, err := DecodeEmptyModelWithDescription(nil, emptyModelWithDescriptionData)
	require.NoError(t, err)
	require.NotNil(t, emptyModelWithDescription)

	modelWithSingleStringFieldData, err := os.ReadFile("../../test_data/model_with_single_string_field.bin")
	require.NoError(t, err)
	modelWithSingleStringField, err := DecodeModelWithSingleStringField(nil, modelWithSingleStringFieldData)
	require.NoError(t, err)
	require.NotNil(t, modelWithSingleStringField)
	require.Equal(t, "hello world", modelWithSingleStringField.StringField)

	modelWithSingleStringFieldAndDescriptionData, err := os.ReadFile("../../test_data/model_with_single_string_field_and_description.bin")
	require.NoError(t, err)
	modelWithSingleStringFieldAndDescription, err := DecodeModelWithSingleStringFieldAndDescription(nil, modelWithSingleStringFieldAndDescriptionData)
	require.NoError(t, err)
	require.NotNil(t, modelWithSingleStringFieldAndDescription)
	require.Equal(t, "hello world", modelWithSingleStringFieldAndDescription.StringField)

	modelWithSingleInt32FieldData, err := os.ReadFile("../../test_data/model_with_single_int32_field.bin")
	require.NoError(t, err)
	modelWithSingleInt32Field, err := DecodeModelWithSingleInt32Field(nil, modelWithSingleInt32FieldData)
	require.NoError(t, err)
	require.NotNil(t, modelWithSingleInt32Field)
	require.Equal(t, int32(42), modelWithSingleInt32Field.Int32Field)

	modelWithSingleInt32FieldAndDescriptionData, err := os.ReadFile("../../test_data/model_with_single_int32_field_and_description.bin")
	require.NoError(t, err)
	modelWithSingleInt32FieldAndDescription, err := DecodeModelWithSingleInt32FieldAndDescription(nil, modelWithSingleInt32FieldAndDescriptionData)
	require.NoError(t, err)
	require.NotNil(t, modelWithSingleInt32FieldAndDescription)
	require.Equal(t, int32(42), modelWithSingleInt32FieldAndDescription.Int32Field)

	modelWithMultipleFieldsData, err := os.ReadFile("../../test_data/model_with_multiple_fields.bin")
	require.NoError(t, err)
	modelWithMultipleFields, err := DecodeModelWithMultipleFields(nil, modelWithMultipleFieldsData)
	require.NoError(t, err)
	require.NotNil(t, modelWithMultipleFields)
	require.Equal(t, "hello world", modelWithMultipleFields.StringField)
	require.Equal(t, int32(42), modelWithMultipleFields.Int32Field)

	modelWithMultipleFieldsAndDescriptionData, err := os.ReadFile("../../test_data/model_with_multiple_fields_and_description.bin")
	require.NoError(t, err)
	modelWithMultipleFieldsAndDescription, err := DecodeModelWithMultipleFieldsAndDescription(nil, modelWithMultipleFieldsAndDescriptionData)
	require.NoError(t, err)
	require.NotNil(t, modelWithMultipleFieldsAndDescription)
	require.Equal(t, "hello world", modelWithMultipleFieldsAndDescription.StringField)
	require.Equal(t, int32(42), modelWithMultipleFieldsAndDescription.Int32Field)

	modelWithEnumData, err := os.ReadFile("../../test_data/model_with_enum.bin")
	require.NoError(t, err)
	modelWithEnum, err := DecodeModelWithEnum(nil, modelWithEnumData)
	require.NoError(t, err)
	require.NotNil(t, modelWithEnum)
	require.Equal(t, GenericEnumSecondValue, modelWithEnum.EnumField)

	modelWithEnumAndDescriptionData, err := os.ReadFile("../../test_data/model_with_enum_and_description.bin")
	require.NoError(t, err)
	modelWithEnumAndDescription, err := DecodeModelWithEnumAndDescription(nil, modelWithEnumAndDescriptionData)
	require.NoError(t, err)
	require.NotNil(t, modelWithEnumAndDescription)
	require.Equal(t, GenericEnumSecondValue, modelWithEnumAndDescription.EnumField)

	modelWithEnumAccessorData, err := os.ReadFile("../../test_data/model_with_enum_accessor.bin")
	require.NoError(t, err)
	modelWithEnumAccessor, err := DecodeModelWithEnumAccessor(nil, modelWithEnumAccessorData)
	require.NoError(t, err)
	require.NotNil(t, modelWithEnumAccessor)
	enumValue, err := modelWithEnumAccessor.GetEnumField()
	require.NoError(t, err)
	require.Equal(t, GenericEnumSecondValue, enumValue)

	modelWithEnumAccessorAndDescriptionData, err := os.ReadFile("../../test_data/model_with_enum_accessor_and_description.bin")
	require.NoError(t, err)
	modelWithEnumAccessorAndDescription, err := DecodeModelWithEnumAccessorAndDescription(nil, modelWithEnumAccessorAndDescriptionData)
	require.NoError(t, err)
	require.NotNil(t, modelWithEnumAccessorAndDescription)
	enumValue, err = modelWithEnumAccessorAndDescription.GetEnumField()
	require.NoError(t, err)
	require.Equal(t, GenericEnumSecondValue, enumValue)

	modelWithMultipleFieldsAccessorData, err := os.ReadFile("../../test_data/model_with_multiple_fields_accessor.bin")
	require.NoError(t, err)
	modelWithMultipleFieldsAccessor, err := DecodeModelWithMultipleFieldsAccessor(nil, modelWithMultipleFieldsAccessorData)
	require.NoError(t, err)
	require.NotNil(t, modelWithMultipleFieldsAccessor)
	stringFieldValue, err := modelWithMultipleFieldsAccessor.GetStringField()
	require.NoError(t, err)
	require.Equal(t, "HELLO", stringFieldValue)
	int32FieldValue, err := modelWithMultipleFieldsAccessor.GetInt32Field()
	require.NoError(t, err)
	require.Equal(t, int32(42), int32FieldValue)

	modelWithMultipleFieldsAccessorAndDescriptionData, err := os.ReadFile("../../test_data/model_with_multiple_fields_accessor_and_description.bin")
	require.NoError(t, err)
	modelWithMultipleFieldsAccessorAndDescription, err := DecodeModelWithMultipleFieldsAccessorAndDescription(nil, modelWithMultipleFieldsAccessorAndDescriptionData)
	require.NoError(t, err)
	require.NotNil(t, modelWithMultipleFieldsAccessorAndDescription)
	stringFieldValue, err = modelWithMultipleFieldsAccessorAndDescription.GetStringField()
	require.NoError(t, err)
	require.Equal(t, "hello world", stringFieldValue)
	int32FieldValue, err = modelWithMultipleFieldsAccessorAndDescription.GetInt32Field()
	require.NoError(t, err)
	require.Equal(t, int32(42), int32FieldValue)

	modelWithEmbeddedModelsData, err := os.ReadFile("../../test_data/model_with_embedded_models.bin")
	require.NoError(t, err)
	modelWithEmbeddedModels, err := DecodeModelWithEmbeddedModels(nil, modelWithEmbeddedModelsData)
	require.NoError(t, err)
	require.NotNil(t, modelWithEmbeddedModels)
	require.NotNil(t, modelWithEmbeddedModels.EmbeddedEmptyModel)
	require.NotNil(t, modelWithEmbeddedModels.EmbeddedModelArrayWithMultipleFieldsAccessor)
	require.Equal(t, 1, cap(modelWithEmbeddedModels.EmbeddedModelArrayWithMultipleFieldsAccessor))
	require.Equal(t, 1, len(modelWithEmbeddedModels.EmbeddedModelArrayWithMultipleFieldsAccessor))
	require.IsType(t, []ModelWithMultipleFieldsAccessor{}, modelWithEmbeddedModels.EmbeddedModelArrayWithMultipleFieldsAccessor)
	require.Equal(t, *modelWithMultipleFieldsAccessor, modelWithEmbeddedModels.EmbeddedModelArrayWithMultipleFieldsAccessor[0])

	modelWithEmbeddedModelsAndDescriptionData, err := os.ReadFile("../../test_data/model_with_embedded_models_and_description.bin")
	require.NoError(t, err)
	modelWithEmbeddedModelsAndDescription, err := DecodeModelWithEmbeddedModelsAndDescription(nil, modelWithEmbeddedModelsAndDescriptionData)
	require.NoError(t, err)
	require.NotNil(t, modelWithEmbeddedModelsAndDescription)
	require.NotNil(t, modelWithEmbeddedModelsAndDescription.EmbeddedEmptyModel)
	require.NotNil(t, modelWithEmbeddedModelsAndDescription.EmbeddedModelArrayWithMultipleFieldsAccessor)
	require.Equal(t, 1, cap(modelWithEmbeddedModelsAndDescription.EmbeddedModelArrayWithMultipleFieldsAccessor))
	require.Equal(t, 1, len(modelWithEmbeddedModelsAndDescription.EmbeddedModelArrayWithMultipleFieldsAccessor))
	require.IsType(t, []ModelWithMultipleFieldsAccessor{}, modelWithEmbeddedModelsAndDescription.EmbeddedModelArrayWithMultipleFieldsAccessor)
	require.Equal(t, *modelWithMultipleFieldsAccessor, modelWithEmbeddedModelsAndDescription.EmbeddedModelArrayWithMultipleFieldsAccessor[0])

	modelWithEmbeddedModelsAccessorData, err := os.ReadFile("../../test_data/model_with_embedded_models_accessor.bin")
	require.NoError(t, err)
	modelWithEmbeddedModelsAccessor, err := DecodeModelWithEmbeddedModelsAccessor(nil, modelWithEmbeddedModelsAccessorData)
	require.NoError(t, err)
	require.NotNil(t, modelWithEmbeddedModelsAccessor)
	embeddedEmptyModel, err := modelWithEmbeddedModelsAccessor.GetEmbeddedEmptyModel()
	require.NoError(t, err)
	require.NotNil(t, embeddedEmptyModel)
	embeddedModelArrayWithMultipleFieldsAccessor, err := modelWithEmbeddedModelsAccessor.GetEmbeddedModelArrayWithMultipleFieldsAccessor()
	require.NoError(t, err)
	require.Equal(t, 1, cap(embeddedModelArrayWithMultipleFieldsAccessor))
	require.Equal(t, 1, len(embeddedModelArrayWithMultipleFieldsAccessor))
	require.IsType(t, []ModelWithMultipleFieldsAccessor{}, embeddedModelArrayWithMultipleFieldsAccessor)
	require.Equal(t, *modelWithMultipleFieldsAccessor, embeddedModelArrayWithMultipleFieldsAccessor[0])

	modelWithEmbeddedModelsAccessorAndDescriptionData, err := os.ReadFile("../../test_data/model_with_embedded_models_accessor_and_description.bin")
	require.NoError(t, err)
	modelWithEmbeddedModelsAccessorAndDescription, err := DecodeModelWithEmbeddedModelsAccessorAndDescription(nil, modelWithEmbeddedModelsAccessorAndDescriptionData)
	require.NoError(t, err)
	require.NotNil(t, modelWithEmbeddedModelsAccessorAndDescription)
	embeddedEmptyModel, err = modelWithEmbeddedModelsAccessorAndDescription.GetEmbeddedEmptyModel()
	require.NoError(t, err)
	require.NotNil(t, embeddedEmptyModel)
	embeddedModelArrayWithMultipleFieldsAccessor, err = modelWithEmbeddedModelsAccessorAndDescription.GetEmbeddedModelArrayWithMultipleFieldsAccessor()
	require.NoError(t, err)
	require.Equal(t, 1, cap(embeddedModelArrayWithMultipleFieldsAccessor))
	require.Equal(t, 1, len(embeddedModelArrayWithMultipleFieldsAccessor))
	require.IsType(t, []ModelWithMultipleFieldsAccessor{}, embeddedModelArrayWithMultipleFieldsAccessor)
	require.Equal(t, *modelWithMultipleFieldsAccessor, embeddedModelArrayWithMultipleFieldsAccessor[0])

	modelWithAllFieldTypesData, err := os.ReadFile("../../test_data/model_with_all_field_types.bin")
	require.NoError(t, err)
	modelWithAllFieldTypes, err := DecodeModelWithAllFieldTypes(nil, modelWithAllFieldTypesData)
	require.NoError(t, err)
	require.NotNil(t, modelWithAllFieldTypes)

	require.Equal(t, "hello world", modelWithAllFieldTypes.StringField)
	require.Equal(t, 2, len(modelWithAllFieldTypes.StringArrayField))
	require.Equal(t, "hello", modelWithAllFieldTypes.StringArrayField[0])
	require.Equal(t, "world", modelWithAllFieldTypes.StringArrayField[1])
	require.Equal(t, "world", modelWithAllFieldTypes.StringMapField["hello"])
	require.Equal(t, *emptyModel, modelWithAllFieldTypes.StringMapFieldEmbedded["hello"])

	require.Equal(t, int32(42), modelWithAllFieldTypes.Int32Field)
	require.Equal(t, 2, len(modelWithAllFieldTypes.Int32ArrayField))
	require.Equal(t, int32(42), modelWithAllFieldTypes.Int32ArrayField[0])
	require.Equal(t, int32(84), modelWithAllFieldTypes.Int32ArrayField[1])
	require.Equal(t, int32(84), modelWithAllFieldTypes.Int32MapField[42])
	require.Equal(t, *emptyModel, modelWithAllFieldTypes.Int32MapFieldEmbedded[42])

	require.Equal(t, int64(100), modelWithAllFieldTypes.Int64Field)
	require.Equal(t, 2, len(modelWithAllFieldTypes.Int64ArrayField))
	require.Equal(t, int64(100), modelWithAllFieldTypes.Int64ArrayField[0])
	require.Equal(t, int64(200), modelWithAllFieldTypes.Int64ArrayField[1])
	require.Equal(t, int64(200), modelWithAllFieldTypes.Int64MapField[100])
	require.Equal(t, *emptyModel, modelWithAllFieldTypes.Int64MapFieldEmbedded[100])

	require.Equal(t, uint32(42), modelWithAllFieldTypes.Uint32Field)
	require.Equal(t, 2, len(modelWithAllFieldTypes.Uint32ArrayField))
	require.Equal(t, uint32(42), modelWithAllFieldTypes.Uint32ArrayField[0])
	require.Equal(t, uint32(84), modelWithAllFieldTypes.Uint32ArrayField[1])
	require.Equal(t, uint32(84), modelWithAllFieldTypes.Uint32MapField[42])
	require.Equal(t, *emptyModel, modelWithAllFieldTypes.Uint32MapFieldEmbedded[42])

	require.Equal(t, uint64(100), modelWithAllFieldTypes.Uint64Field)
	require.Equal(t, 2, len(modelWithAllFieldTypes.Uint64ArrayField))
	require.Equal(t, uint64(100), modelWithAllFieldTypes.Uint64ArrayField[0])
	require.Equal(t, uint64(200), modelWithAllFieldTypes.Uint64ArrayField[1])
	require.Equal(t, uint64(200), modelWithAllFieldTypes.Uint64MapField[100])
	require.Equal(t, *emptyModel, modelWithAllFieldTypes.Uint64MapFieldEmbedded[100])

	require.Equal(t, float32(42.0), modelWithAllFieldTypes.Float32Field)
	require.Equal(t, 2, len(modelWithAllFieldTypes.Float32ArrayField))
	require.Equal(t, float32(42.0), modelWithAllFieldTypes.Float32ArrayField[0])
	require.Equal(t, float32(84.0), modelWithAllFieldTypes.Float32ArrayField[1])

	require.Equal(t, float64(100.0), modelWithAllFieldTypes.Float64Field)
	require.Equal(t, 2, len(modelWithAllFieldTypes.Float64ArrayField))
	require.Equal(t, float64(100.0), modelWithAllFieldTypes.Float64ArrayField[0])
	require.Equal(t, float64(200.0), modelWithAllFieldTypes.Float64ArrayField[1])

	require.Equal(t, false, modelWithAllFieldTypes.BoolField)
	require.Equal(t, 2, len(modelWithAllFieldTypes.BoolArrayField))
	require.Equal(t, true, modelWithAllFieldTypes.BoolArrayField[0])
	require.Equal(t, false, modelWithAllFieldTypes.BoolArrayField[1])

	require.Equal(t, []byte{42, 84}, modelWithAllFieldTypes.BytesField)
	require.Equal(t, 2, len(modelWithAllFieldTypes.BytesArrayField))
	require.Equal(t, []byte{42, 84}, modelWithAllFieldTypes.BytesArrayField[0])
	require.Equal(t, []byte{84, 42}, modelWithAllFieldTypes.BytesArrayField[1])

	require.Equal(t, GenericEnumSecondValue, modelWithAllFieldTypes.EnumField)
	require.Equal(t, 2, len(modelWithAllFieldTypes.EnumArrayField))
	require.Equal(t, GenericEnumFirstValue, modelWithAllFieldTypes.EnumArrayField[0])
	require.Equal(t, GenericEnumSecondValue, modelWithAllFieldTypes.EnumArrayField[1])
	require.Equal(t, "hello world", modelWithAllFieldTypes.EnumMapField[GenericEnumFirstValue])
	require.Equal(t, *emptyModel, modelWithAllFieldTypes.EnumMapFieldEmbedded[GenericEnumFirstValue])

	require.Equal(t, 2, len(modelWithAllFieldTypes.ModelArrayField))
	require.Equal(t, *emptyModel, modelWithAllFieldTypes.ModelArrayField[0])
	require.Equal(t, *emptyModel, modelWithAllFieldTypes.ModelArrayField[1])
}
