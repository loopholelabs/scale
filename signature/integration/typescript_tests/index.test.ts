import * as generated from "./generated";
import * as polyglot from "@loopholelabs/polyglot";
import * as fs from 'fs';

test('test-output', () => {
    const nilModelEncoder = new polyglot.Encoder();
   generated.EmptyModel.encode_undefined(nilModelEncoder);
    fs.writeFileSync('../binaries/nil_model.bin', nilModelEncoder.bytes, 'binary');

    const emptyModelEncoder = new polyglot.Encoder();
    const emptyModel = new generated.EmptyModel();
    emptyModel.encode(emptyModelEncoder);
    fs.writeFileSync('../binaries/empty_model.bin', emptyModelEncoder.bytes, 'binary');

    const emptyModelWithDescriptionEncoder = new polyglot.Encoder();
    const emptyModelWithDescription = new generated.EmptyModelWithDescription();
    emptyModelWithDescription.encode(emptyModelWithDescriptionEncoder);
    fs.writeFileSync('../binaries/empty_model_with_description.bin', emptyModelWithDescriptionEncoder.bytes, 'binary');

    const modelWithSingleStringFieldEncoder = new polyglot.Encoder();
    const modelWithSingleStringField = new generated.ModelWithSingleStringField();
    expect(modelWithSingleStringField.stringField).toEqual('DefaultValue');
    modelWithSingleStringField.stringField = 'hello world';
    modelWithSingleStringField.encode(modelWithSingleStringFieldEncoder);
    fs.writeFileSync('../binaries/model_with_single_string_field.bin', modelWithSingleStringFieldEncoder.bytes, 'binary');

    const modelWithSingleStringFieldAndDescriptionEncoder = new polyglot.Encoder();
    const modelWithSingleStringFieldAndDescription = new generated.ModelWithSingleStringFieldAndDescription();
    expect(modelWithSingleStringFieldAndDescription.stringField).toEqual('DefaultValue');
    modelWithSingleStringFieldAndDescription.stringField = 'hello world';
    modelWithSingleStringFieldAndDescription.encode(modelWithSingleStringFieldAndDescriptionEncoder);
    fs.writeFileSync(
        '../binaries/model_with_single_string_field_and_description.bin',
        modelWithSingleStringFieldAndDescriptionEncoder.bytes,
        'binary'
    );

    const modelWithSingleInt32FieldEncoder = new polyglot.Encoder();
    const modelWithSingleInt32Field = new generated.ModelWithSingleInt32Field();
    expect(modelWithSingleInt32Field.int32Field).toEqual(32);
    modelWithSingleInt32Field.int32Field = 42;
    modelWithSingleInt32Field.encode(modelWithSingleInt32FieldEncoder);
    fs.writeFileSync('../binaries/model_with_single_int32_field.bin', modelWithSingleInt32FieldEncoder.bytes, 'binary');

    const modelWithSingleInt32FieldAndDescriptionEncoder = new polyglot.Encoder();
    const modelWithSingleInt32FieldAndDescription = new generated.ModelWithSingleInt32FieldAndDescription();
    expect(modelWithSingleInt32FieldAndDescription.int32Field).toEqual(32);
    modelWithSingleInt32FieldAndDescription.int32Field = 42;
    modelWithSingleInt32FieldAndDescription.encode(modelWithSingleInt32FieldAndDescriptionEncoder);
    fs.writeFileSync('../binaries/model_with_single_int32_field_and_description.bin', modelWithSingleInt32FieldAndDescriptionEncoder.bytes, 'binary');

    const modelWithMultipleFieldsEncoder = new polyglot.Encoder();
    const modelWithMultipleFields = new generated.ModelWithMultipleFields();
    expect(modelWithMultipleFields.stringField).toEqual('DefaultValue');
    expect(modelWithMultipleFields.int32Field).toEqual(32);
    modelWithMultipleFields.stringField = 'hello world';
    modelWithMultipleFields.int32Field = 42;
    modelWithMultipleFields.encode(modelWithMultipleFieldsEncoder);
    fs.writeFileSync('../binaries/model_with_multiple_fields.bin', modelWithMultipleFieldsEncoder.bytes, 'binary');

    const modelWithMultipleFieldsAndDescriptionEncoder = new polyglot.Encoder();
    const modelWithMultipleFieldsAndDescription = new generated.ModelWithMultipleFieldsAndDescription();
    expect(modelWithMultipleFieldsAndDescription.stringField).toEqual('DefaultValue');
    expect(modelWithMultipleFieldsAndDescription.int32Field).toEqual(32);
    modelWithMultipleFieldsAndDescription.stringField = 'hello world';
    modelWithMultipleFieldsAndDescription.int32Field = 42;
    modelWithMultipleFieldsAndDescription.encode(modelWithMultipleFieldsAndDescriptionEncoder);
    fs.writeFileSync('../binaries/model_with_multiple_fields_and_description.bin', modelWithMultipleFieldsAndDescriptionEncoder.bytes, 'binary');

    const modelWithEnumEncoder = new polyglot.Encoder();
    const modelWithEnum = new generated.ModelWithEnum();
    expect(modelWithEnum.enumField).toEqual(generated.GenericEnum.DefaultValue);
    modelWithEnum.enumField = generated.GenericEnum.SecondValue;
    modelWithEnum.encode(modelWithEnumEncoder);
    fs.writeFileSync('../binaries/model_with_enum.bin', modelWithEnumEncoder.bytes, 'binary');

    const modelWithEnumAndDescriptionEncoder = new polyglot.Encoder();
    const modelWithEnumAndDescription = new generated.ModelWithEnumAndDescription();
    expect(modelWithEnumAndDescription.enumField).toEqual(generated.GenericEnum.DefaultValue);
    modelWithEnumAndDescription.enumField = generated.GenericEnum.SecondValue;
    modelWithEnumAndDescription.encode(modelWithEnumAndDescriptionEncoder);
    fs.writeFileSync('../binaries/model_with_enum_and_description.bin', modelWithEnumAndDescriptionEncoder.bytes, 'binary');

    const modelWithEnumAccessorEncoder = new polyglot.Encoder();
    const modelWithEnumAccessor = new generated.ModelWithEnumAccessor();
    let enumValue = modelWithEnumAccessor.enumField
    expect(enumValue).toEqual(generated.GenericEnum.DefaultValue);
    modelWithEnumAccessor.enumField = generated.GenericEnum.SecondValue;
    modelWithEnumAccessor.encode(modelWithEnumAccessorEncoder);
    fs.writeFileSync('../binaries/model_with_enum_accessor.bin', modelWithEnumAccessorEncoder.bytes, 'binary');

    const modelWithEnumAccessorAndDescriptionEncoder = new polyglot.Encoder();
    const modelWithEnumAccessorAndDescription = new generated.ModelWithEnumAccessorAndDescription();
    enumValue = modelWithEnumAccessorAndDescription.enumField
    expect(enumValue).toEqual(generated.GenericEnum.DefaultValue);
    modelWithEnumAccessorAndDescription.enumField = generated.GenericEnum.SecondValue;
    modelWithEnumAccessorAndDescription.encode(modelWithEnumAccessorAndDescriptionEncoder);
    fs.writeFileSync('../binaries/model_with_enum_accessor_and_description.bin', modelWithEnumAccessorAndDescriptionEncoder.bytes, 'binary');

    const modelWithMultipleFieldsAccessorEncoder = new polyglot.Encoder();
    const modelWithMultipleFieldsAccessor = new generated.ModelWithMultipleFieldsAccessor();
    let stringFieldValue = modelWithMultipleFieldsAccessor.stringField;
    expect(stringFieldValue).toEqual('DefaultValue');
    try {
        modelWithMultipleFieldsAccessor.stringField = 'hello world';
        fail('Expected error to be thrown');
    } catch (e) {
        // @ts-ignore
        expect(e.message).toEqual('value must match ^[a-zA-Z0-9]*$');
    }
    try {
        modelWithMultipleFieldsAccessor.stringField = "";
        fail('Expected error to be thrown');
    } catch (e) {
        // @ts-ignore
        expect(e.message).toEqual('length must be between 1 and 20');
    }
    modelWithMultipleFieldsAccessor.stringField = 'hello';
    stringFieldValue = modelWithMultipleFieldsAccessor.stringField;
    expect(stringFieldValue).toEqual('HELLO');
    let int32FieldValue = modelWithMultipleFieldsAccessor.int32Field;
    expect(int32FieldValue).toEqual(32);
    try {
        modelWithMultipleFieldsAccessor.int32Field = -1;
        fail('Expected error to be thrown');
    } catch (e) {
        // @ts-ignore
        expect(e.message).toEqual('value must be between 0 and 100');
    }
    try {
        modelWithMultipleFieldsAccessor.int32Field = 101;
        fail('Expected error to be thrown');
    } catch (e) {
        // @ts-ignore
        expect(e.message).toEqual('value must be between 0 and 100');
    }
    modelWithMultipleFieldsAccessor.int32Field = 42;
    modelWithMultipleFieldsAccessor.encode(modelWithMultipleFieldsAccessorEncoder);
    fs.writeFileSync('../binaries/model_with_multiple_fields_accessor.bin', modelWithMultipleFieldsAccessorEncoder.bytes, 'binary');

    const modelWithMultipleFieldsAccessorAndDescriptionEncoder = new polyglot.Encoder();
    const modelWithMultipleFieldsAccessorAndDescription = new generated.ModelWithMultipleFieldsAccessorAndDescription();
    stringFieldValue = modelWithMultipleFieldsAccessorAndDescription.stringField;
    expect(stringFieldValue).toEqual('DefaultValue');
    modelWithMultipleFieldsAccessorAndDescription.stringField = 'hello world';
    int32FieldValue = modelWithMultipleFieldsAccessorAndDescription.int32Field;
    expect(int32FieldValue).toEqual(32);
    modelWithMultipleFieldsAccessorAndDescription.int32Field = 42;
    modelWithMultipleFieldsAccessorAndDescription.encode(modelWithMultipleFieldsAccessorAndDescriptionEncoder);
    fs.writeFileSync('../binaries/model_with_multiple_fields_accessor_and_description.bin', modelWithMultipleFieldsAccessorAndDescriptionEncoder.bytes, 'binary');

    const modelWithEmbeddedModelsEncoder = new polyglot.Encoder();
    const modelWithEmbeddedModels = new generated.ModelWithEmbeddedModels();
    expect(modelWithEmbeddedModels.embeddedEmptyModel).not.toBeNull();
    expect(modelWithEmbeddedModels.embeddedEmptyModel).not.toBeUndefined();
    expect(modelWithEmbeddedModels.embeddedModelArrayWithMultipleFieldsAccessor).not.toBeNull();
    expect(modelWithEmbeddedModels.embeddedModelArrayWithMultipleFieldsAccessor).not.toBeUndefined();
    expect(modelWithEmbeddedModels.embeddedModelArrayWithMultipleFieldsAccessor.length).toEqual(0);
    modelWithEmbeddedModels.embeddedModelArrayWithMultipleFieldsAccessor.push(modelWithMultipleFieldsAccessor);
    modelWithEmbeddedModels.encode(modelWithEmbeddedModelsEncoder);
    fs.writeFileSync('../binaries/model_with_embedded_models.bin', modelWithEmbeddedModelsEncoder.bytes, 'binary');

    const modelWithEmbeddedModelsAndDescriptionEncoder = new polyglot.Encoder();
    const modelWithEmbeddedModelsAndDescription = new generated.ModelWithEmbeddedModelsAndDescription();
    expect(modelWithEmbeddedModelsAndDescription.embeddedEmptyModel).not.toBeNull();
    expect(modelWithEmbeddedModelsAndDescription.embeddedEmptyModel).not.toBeUndefined();
    expect(modelWithEmbeddedModelsAndDescription.embeddedModelArrayWithMultipleFieldsAccessor).not.toBeNull();
    expect(modelWithEmbeddedModelsAndDescription.embeddedModelArrayWithMultipleFieldsAccessor).not.toBeUndefined();
    expect(modelWithEmbeddedModelsAndDescription.embeddedModelArrayWithMultipleFieldsAccessor.length).toEqual(0);
    modelWithEmbeddedModelsAndDescription.embeddedModelArrayWithMultipleFieldsAccessor.push(modelWithMultipleFieldsAccessor);
    modelWithEmbeddedModelsAndDescription.encode(modelWithEmbeddedModelsAndDescriptionEncoder);
    fs.writeFileSync('../binaries/model_with_embedded_models_and_description.bin', modelWithEmbeddedModelsAndDescriptionEncoder.bytes, 'binary');

    const modelWithEmbeddedModelsAccessorEncoder = new polyglot.Encoder();
    const modelWithEmbeddedModelsAccessor = new generated.ModelWithEmbeddedModelsAccessor();
    let embeddedModel = modelWithEmbeddedModelsAccessor.embeddedEmptyModel;
    expect(embeddedModel).not.toBeNull();
    expect(embeddedModel).not.toBeUndefined();
    let embeddedModelArray = modelWithEmbeddedModelsAccessor.embeddedModelArrayWithMultipleFieldsAccessor;
    expect(embeddedModelArray).not.toBeNull();
    expect(embeddedModelArray).not.toBeUndefined();
    expect(embeddedModelArray.length).toEqual(0);
    modelWithEmbeddedModelsAccessor.embeddedModelArrayWithMultipleFieldsAccessor.push(modelWithMultipleFieldsAccessor);
    modelWithEmbeddedModelsAccessor.encode(modelWithEmbeddedModelsAccessorEncoder);
    fs.writeFileSync('../binaries/model_with_embedded_models_accessor.bin', modelWithEmbeddedModelsAccessorEncoder.bytes, 'binary');

    const modelWithEmbeddedModelsAccessorAndDescriptionEncoder = new polyglot.Encoder();
    const modelWithEmbeddedModelsAccessorAndDescription = new generated.ModelWithEmbeddedModelsAccessorAndDescription();
    embeddedModel = modelWithEmbeddedModelsAccessorAndDescription.embeddedEmptyModel;
    expect(embeddedModel).not.toBeNull();
    expect(embeddedModel).not.toBeUndefined();
    embeddedModelArray = modelWithEmbeddedModelsAccessorAndDescription.embeddedModelArrayWithMultipleFieldsAccessor;
    expect(embeddedModelArray).not.toBeNull();
    expect(embeddedModelArray).not.toBeUndefined();
    expect(embeddedModelArray.length).toEqual(0);
    modelWithEmbeddedModelsAccessorAndDescription.embeddedModelArrayWithMultipleFieldsAccessor.push(modelWithMultipleFieldsAccessor);
    modelWithEmbeddedModelsAccessorAndDescription.encode(modelWithEmbeddedModelsAccessorAndDescriptionEncoder);
    fs.writeFileSync('../binaries/model_with_embedded_models_accessor_and_description.bin', modelWithEmbeddedModelsAccessorAndDescriptionEncoder.bytes, 'binary');

    const modelWithAllFieldTypesEncoder = new polyglot.Encoder();
    const modelWithAllFieldTypes = new generated.ModelWithAllFieldTypes();
    expect(modelWithAllFieldTypes.stringField).toEqual('DefaultValue');
    modelWithAllFieldTypes.stringField = 'hello world';
    expect(modelWithAllFieldTypes.stringArrayField.length).toEqual(0);
    expect(modelWithAllFieldTypes.stringMapField).toEqual(new Map<string, string>());
    expect(modelWithAllFieldTypes.stringMapFieldEmbedded).toEqual(new Map<string, generated.EmptyModel>());
    modelWithAllFieldTypes.stringArrayField.push('hello', 'world');
    modelWithAllFieldTypes.stringMapField.set('hello', 'world');
    modelWithAllFieldTypes.stringMapFieldEmbedded.set('hello', emptyModel);

    expect(modelWithAllFieldTypes.int32Field).toEqual(32);
    modelWithAllFieldTypes.int32Field = 42;
    expect(modelWithAllFieldTypes.int32ArrayField.length).toEqual(0);
    expect(modelWithAllFieldTypes.int32MapField).toEqual(new Map<number, number>());
    expect(modelWithAllFieldTypes.int32MapFieldEmbedded).toEqual(new Map<number, generated.EmptyModel>());
    modelWithAllFieldTypes.int32ArrayField.push(42, 84);
    modelWithAllFieldTypes.int32MapField.set(42, 84);
    modelWithAllFieldTypes.int32MapFieldEmbedded.set(42, emptyModel);

    expect(modelWithAllFieldTypes.int64Field).toEqual(BigInt(64));
    modelWithAllFieldTypes.int64Field = BigInt(100);
    expect(modelWithAllFieldTypes.int64ArrayField.length).toEqual(0);
    expect(modelWithAllFieldTypes.int64MapField).toEqual(new Map<bigint, bigint>());
    expect(modelWithAllFieldTypes.int64MapFieldEmbedded).toEqual(new Map<bigint, generated.EmptyModel>());
    modelWithAllFieldTypes.int64ArrayField.push(BigInt(100), BigInt(200));
    modelWithAllFieldTypes.int64MapField.set(BigInt(100), BigInt(200));
    modelWithAllFieldTypes.int64MapFieldEmbedded.set(BigInt(100), emptyModel);

    expect(modelWithAllFieldTypes.uint32Field).toEqual(32);
    modelWithAllFieldTypes.uint32Field = 42;
    expect(modelWithAllFieldTypes.uint32ArrayField.length).toEqual(0);
    expect(modelWithAllFieldTypes.uint32MapField).toEqual(new Map<number, number>());
    expect(modelWithAllFieldTypes.uint32MapFieldEmbedded).toEqual(new Map<number, generated.EmptyModel>());
    modelWithAllFieldTypes.uint32ArrayField.push(42, 84);
    modelWithAllFieldTypes.uint32MapField.set(42, 84);
    modelWithAllFieldTypes.uint32MapFieldEmbedded.set(42, emptyModel);

    expect(modelWithAllFieldTypes.uint64Field).toEqual(BigInt(64));
    modelWithAllFieldTypes.uint64Field = BigInt(100);
    expect(modelWithAllFieldTypes.uint64ArrayField.length).toEqual(0);
    expect(modelWithAllFieldTypes.uint64MapField).toEqual(new Map<bigint, bigint>());
    expect(modelWithAllFieldTypes.uint64MapFieldEmbedded).toEqual(new Map<bigint, generated.EmptyModel>());
    modelWithAllFieldTypes.uint64ArrayField.push(BigInt(100), BigInt(200));
    modelWithAllFieldTypes.uint64MapField.set(BigInt(100), BigInt(200));
    modelWithAllFieldTypes.uint64MapFieldEmbedded.set(BigInt(100), emptyModel);

    expect(modelWithAllFieldTypes.float32Field).toEqual(32.32);
    modelWithAllFieldTypes.float32Field = 42.0;
    expect(modelWithAllFieldTypes.float32ArrayField.length).toEqual(0);
    modelWithAllFieldTypes.float32ArrayField.push(42.0, 84.0);

    expect(modelWithAllFieldTypes.float64Field).toEqual(64.64);
    modelWithAllFieldTypes.float64Field = 100.0;
    expect(modelWithAllFieldTypes.float64ArrayField.length).toEqual(0);
    modelWithAllFieldTypes.float64ArrayField.push(100.0, 200.0);

    expect(modelWithAllFieldTypes.boolField).toEqual(true);
    modelWithAllFieldTypes.boolField = false;
    expect(modelWithAllFieldTypes.boolArrayField.length).toEqual(0);
    modelWithAllFieldTypes.boolArrayField.push(true, false);

    expect(modelWithAllFieldTypes.bytesField.length).toEqual(512);
    modelWithAllFieldTypes.bytesField = Uint8Array.from([42, 84]);
    expect(modelWithAllFieldTypes.bytesArrayField.length).toEqual(0);
    modelWithAllFieldTypes.bytesArrayField.push(Uint8Array.from([42, 84]), Uint8Array.from([84, 42]));

    expect(modelWithAllFieldTypes.enumField).toEqual(generated.GenericEnum.DefaultValue);
    modelWithAllFieldTypes.enumField = generated.GenericEnum.SecondValue;
    expect(modelWithAllFieldTypes.enumArrayField.length).toEqual(0);
    modelWithAllFieldTypes.enumArrayField.push(generated.GenericEnum.FirstValue, generated.GenericEnum.SecondValue);
    expect(modelWithAllFieldTypes.enumMapField).toEqual(new Map<generated.GenericEnum, string>());
    expect(modelWithAllFieldTypes.enumMapFieldEmbedded).toEqual(new Map<generated.GenericEnum, generated.EmptyModel>());
    modelWithAllFieldTypes.enumMapField.set(generated.GenericEnum.FirstValue, 'hello world');
    modelWithAllFieldTypes.enumMapFieldEmbedded.set(generated.GenericEnum.FirstValue, emptyModel);

    expect(modelWithAllFieldTypes.modelField).not.toBeNull();
    expect(modelWithAllFieldTypes.modelField).not.toBeUndefined();
    expect(modelWithAllFieldTypes.modelArrayField.length).toEqual(0);
    modelWithAllFieldTypes.modelArrayField.push(emptyModel, emptyModel);

    modelWithAllFieldTypes.encode(modelWithAllFieldTypesEncoder);
    fs.writeFileSync('../binaries/model_with_all_field_types.bin', modelWithAllFieldTypesEncoder.bytes, 'binary');
});

test('test-input', () => {
    const nilModelData = fs.readFileSync("../binaries/nil_model.bin")
    const nilModel = generated.EmptyModel.decode(new polyglot.Decoder(nilModelData));
    expect(nilModel).toBeUndefined();

    const emptyModelData = fs.readFileSync("../binaries/empty_model.bin")
    const emptyModel = generated.EmptyModel.decode(new polyglot.Decoder(emptyModelData));
    expect(emptyModel).not.toBeNull();
    expect(emptyModel).not.toBeUndefined();

    const emptyModelWithDescriptionData = fs.readFileSync("../binaries/empty_model_with_description.bin")
    const emptyModelWithDescription = generated.EmptyModelWithDescription.decode(new polyglot.Decoder(emptyModelWithDescriptionData));
    expect(emptyModelWithDescription).not.toBeNull();
    expect(emptyModelWithDescription).not.toBeUndefined();

    const modelWithSingleStringFieldData = fs.readFileSync("../binaries/model_with_single_string_field.bin")
    const modelWithSingleStringField = generated.ModelWithSingleStringField.decode(new polyglot.Decoder(modelWithSingleStringFieldData));
    expect(modelWithSingleStringField).not.toBeNull();
    expect(modelWithSingleStringField).not.toBeUndefined();
    expect(modelWithSingleStringField?.stringField).toEqual("hello world");

    const modelWithSingleStringFieldAndDescriptionData = fs.readFileSync("../binaries/model_with_single_string_field_and_description.bin")
    const modelWithSingleStringFieldAndDescription = generated.ModelWithSingleStringFieldAndDescription.decode(new polyglot.Decoder(modelWithSingleStringFieldAndDescriptionData));
    expect(modelWithSingleStringFieldAndDescription).not.toBeNull();
    expect(modelWithSingleStringFieldAndDescription).not.toBeUndefined();
    expect(modelWithSingleStringFieldAndDescription?.stringField).toEqual("hello world");

    const modelWithSingleInt32FieldData = fs.readFileSync("../binaries/model_with_single_int32_field.bin")
    const modelWithSingleInt32Field = generated.ModelWithSingleInt32Field.decode(new polyglot.Decoder(modelWithSingleInt32FieldData));
    expect(modelWithSingleInt32Field).not.toBeNull();
    expect(modelWithSingleInt32Field).not.toBeUndefined();
    expect(modelWithSingleInt32Field?.int32Field).toEqual(42);

    const modelWithSingleInt32FieldAndDescriptionData = fs.readFileSync("../binaries/model_with_single_int32_field_and_description.bin")
    const modelWithSingleInt32FieldAndDescription = generated.ModelWithSingleInt32FieldAndDescription.decode(new polyglot.Decoder(modelWithSingleInt32FieldAndDescriptionData));
    expect(modelWithSingleInt32FieldAndDescription).not.toBeNull();
    expect(modelWithSingleInt32FieldAndDescription).not.toBeUndefined();
    expect(modelWithSingleInt32FieldAndDescription?.int32Field).toEqual(42);

    const modelWithMultipleFieldsData = fs.readFileSync("../binaries/model_with_multiple_fields.bin")
    const modelWithMultipleFields = generated.ModelWithMultipleFields.decode(new polyglot.Decoder(modelWithMultipleFieldsData));
    expect(modelWithMultipleFields).not.toBeNull();
    expect(modelWithMultipleFields).not.toBeUndefined();
    expect(modelWithMultipleFields?.stringField).toEqual("hello world");
    expect(modelWithMultipleFields?.int32Field).toEqual(42);

    const modelWithMultipleFieldsAndDescriptionData = fs.readFileSync("../binaries/model_with_multiple_fields_and_description.bin")
    const modelWithMultipleFieldsAndDescription = generated.ModelWithMultipleFieldsAndDescription.decode(new polyglot.Decoder(modelWithMultipleFieldsAndDescriptionData));
    expect(modelWithMultipleFieldsAndDescription).not.toBeNull();
    expect(modelWithMultipleFieldsAndDescription).not.toBeUndefined();
    expect(modelWithMultipleFieldsAndDescription?.stringField).toEqual("hello world");
    expect(modelWithMultipleFieldsAndDescription?.int32Field).toEqual(42);

    const modelWithEnumData = fs.readFileSync("../binaries/model_with_enum.bin")
    const modelWithEnum = generated.ModelWithEnum.decode(new polyglot.Decoder(modelWithEnumData));
    expect(modelWithEnum).not.toBeNull();
    expect(modelWithEnum).not.toBeUndefined();
    expect(modelWithEnum?.enumField).toEqual(generated.GenericEnum.SecondValue);

    const modelWithEnumAndDescriptionData = fs.readFileSync("../binaries/model_with_enum_and_description.bin")
    const modelWithEnumAndDescription = generated.ModelWithEnumAndDescription.decode(new polyglot.Decoder(modelWithEnumAndDescriptionData));
    expect(modelWithEnumAndDescription).not.toBeNull();
    expect(modelWithEnumAndDescription).not.toBeUndefined();
    expect(modelWithEnumAndDescription?.enumField).toEqual(generated.GenericEnum.SecondValue);

    const modelWithEnumAccessorData = fs.readFileSync("../binaries/model_with_enum_accessor.bin")
    const modelWithEnumAccessor = generated.ModelWithEnumAccessor.decode(new polyglot.Decoder(modelWithEnumAccessorData));
    expect(modelWithEnumAccessor).not.toBeNull();
    expect(modelWithEnumAccessor).not.toBeUndefined();
    expect(modelWithEnumAccessor?.enumField).toEqual(generated.GenericEnum.SecondValue);

    const modelWithEnumAccessorAndDescriptionData = fs.readFileSync("../binaries/model_with_enum_accessor_and_description.bin")
    const modelWithEnumAccessorAndDescription = generated.ModelWithEnumAccessorAndDescription.decode(new polyglot.Decoder(modelWithEnumAccessorAndDescriptionData));
    expect(modelWithEnumAccessorAndDescription).not.toBeNull();
    expect(modelWithEnumAccessorAndDescription).not.toBeUndefined();
    expect(modelWithEnumAccessorAndDescription?.enumField).toEqual(generated.GenericEnum.SecondValue);

    const modelWithMultipleFieldsAccessorData = fs.readFileSync("../binaries/model_with_multiple_fields_accessor.bin")
    const modelWithMultipleFieldsAccessor = generated.ModelWithMultipleFieldsAccessor.decode(new polyglot.Decoder(modelWithMultipleFieldsAccessorData));
    expect(modelWithMultipleFieldsAccessor).not.toBeNull();
    expect(modelWithMultipleFieldsAccessor).not.toBeUndefined();
    expect(modelWithMultipleFieldsAccessor?.stringField).toEqual("HELLO");
    expect(modelWithMultipleFieldsAccessor?.int32Field).toEqual(42);

    const modelWithMultipleFieldsAccessorAndDescriptionData = fs.readFileSync("../binaries/model_with_multiple_fields_accessor_and_description.bin")
    const modelWithMultipleFieldsAccessorAndDescription = generated.ModelWithMultipleFieldsAccessorAndDescription.decode(new polyglot.Decoder(modelWithMultipleFieldsAccessorAndDescriptionData));
    expect(modelWithMultipleFieldsAccessorAndDescription).not.toBeNull();
    expect(modelWithMultipleFieldsAccessorAndDescription).not.toBeUndefined();
    expect(modelWithMultipleFieldsAccessorAndDescription?.stringField).toEqual("hello world");
    expect(modelWithMultipleFieldsAccessorAndDescription?.int32Field).toEqual(42);

    const modelWithEmbeddedModelsData = fs.readFileSync("../binaries/model_with_embedded_models.bin")
    const modelWithEmbeddedModels = generated.ModelWithEmbeddedModels.decode(new polyglot.Decoder(modelWithEmbeddedModelsData));
    expect(modelWithEmbeddedModels).not.toBeNull();
    expect(modelWithEmbeddedModels).not.toBeUndefined();
    expect(modelWithEmbeddedModels?.embeddedEmptyModel).not.toBeNull();
    expect(modelWithEmbeddedModels?.embeddedEmptyModel).not.toBeUndefined();
    expect(modelWithEmbeddedModels?.embeddedModelArrayWithMultipleFieldsAccessor).toHaveLength(1);
    expect(modelWithEmbeddedModels?.embeddedModelArrayWithMultipleFieldsAccessor[0]?.int32Field).toEqual(42);
    expect(modelWithEmbeddedModels?.embeddedModelArrayWithMultipleFieldsAccessor[0]?.stringField).toEqual("HELLO");

    const modelWithEmbeddedModelsAndDescriptionData = fs.readFileSync("../binaries/model_with_embedded_models_and_description.bin")
    const modelWithEmbeddedModelsAndDescription = generated.ModelWithEmbeddedModelsAndDescription.decode(new polyglot.Decoder(modelWithEmbeddedModelsAndDescriptionData));
    expect(modelWithEmbeddedModelsAndDescription).not.toBeNull();
    expect(modelWithEmbeddedModelsAndDescription).not.toBeUndefined();
    expect(modelWithEmbeddedModelsAndDescription?.embeddedEmptyModel).not.toBeNull();
    expect(modelWithEmbeddedModelsAndDescription?.embeddedEmptyModel).not.toBeUndefined();
    expect(modelWithEmbeddedModelsAndDescription?.embeddedModelArrayWithMultipleFieldsAccessor).toHaveLength(1);
    expect(modelWithEmbeddedModelsAndDescription?.embeddedModelArrayWithMultipleFieldsAccessor[0]?.int32Field).toEqual(42);
    expect(modelWithEmbeddedModelsAndDescription?.embeddedModelArrayWithMultipleFieldsAccessor[0]?.stringField).toEqual("HELLO");

    const modelWithEmbeddedModelsAccessorData = fs.readFileSync("../binaries/model_with_embedded_models_accessor.bin")
    const modelWithEmbeddedModelsAccessor = generated.ModelWithEmbeddedModelsAccessor.decode(new polyglot.Decoder(modelWithEmbeddedModelsAccessorData));
    expect(modelWithEmbeddedModelsAccessor).not.toBeNull();
    expect(modelWithEmbeddedModelsAccessor).not.toBeUndefined();
    expect(modelWithEmbeddedModelsAccessor?.embeddedEmptyModel).not.toBeNull();
    expect(modelWithEmbeddedModelsAccessor?.embeddedEmptyModel).not.toBeUndefined();
    expect(modelWithEmbeddedModelsAccessor?.embeddedModelArrayWithMultipleFieldsAccessor).toHaveLength(1);
    expect(modelWithEmbeddedModelsAccessor?.embeddedModelArrayWithMultipleFieldsAccessor[0]?.int32Field).toEqual(42);
    expect(modelWithEmbeddedModelsAccessor?.embeddedModelArrayWithMultipleFieldsAccessor[0]?.stringField).toEqual("HELLO");

    const modelWithEmbeddedModelsAccessorAndDescriptionData = fs.readFileSync("../binaries/model_with_embedded_models_accessor_and_description.bin")
    const modelWithEmbeddedModelsAccessorAndDescription = generated.ModelWithEmbeddedModelsAccessorAndDescription.decode(new polyglot.Decoder(modelWithEmbeddedModelsAccessorAndDescriptionData));
    expect(modelWithEmbeddedModelsAccessorAndDescription).not.toBeNull();
    expect(modelWithEmbeddedModelsAccessorAndDescription).not.toBeUndefined();
    expect(modelWithEmbeddedModelsAccessorAndDescription?.embeddedEmptyModel).not.toBeNull();
    expect(modelWithEmbeddedModelsAccessorAndDescription?.embeddedEmptyModel).not.toBeUndefined();
    expect(modelWithEmbeddedModelsAccessorAndDescription?.embeddedModelArrayWithMultipleFieldsAccessor).toHaveLength(1);
    expect(modelWithEmbeddedModelsAccessorAndDescription?.embeddedModelArrayWithMultipleFieldsAccessor[0]?.int32Field).toEqual(42);
    expect(modelWithEmbeddedModelsAccessorAndDescription?.embeddedModelArrayWithMultipleFieldsAccessor[0]?.stringField).toEqual("HELLO");

    const modelWithAllFieldTypesData = fs.readFileSync("../binaries/model_with_all_field_types.bin")
    const modelWithAllFieldTypes = generated.ModelWithAllFieldTypes.decode(new polyglot.Decoder(modelWithAllFieldTypesData));
    expect(modelWithAllFieldTypes).not.toBeNull();
    expect(modelWithAllFieldTypes).not.toBeUndefined();
    expect(modelWithAllFieldTypes?.stringField).toEqual("hello world");
    expect(modelWithAllFieldTypes?.stringArrayField).toHaveLength(2);
    expect(modelWithAllFieldTypes?.stringArrayField[0]).toEqual("hello");
    expect(modelWithAllFieldTypes?.stringArrayField[1]).toEqual("world");
    expect(modelWithAllFieldTypes?.stringMapField.get("hello")).toEqual("world");
    expect(modelWithAllFieldTypes?.stringMapFieldEmbedded.get("hello")).toEqual(emptyModel);

    expect(modelWithAllFieldTypes?.int32Field).toEqual(42);
    expect(modelWithAllFieldTypes?.int32ArrayField).toHaveLength(2);
    expect(modelWithAllFieldTypes?.int32ArrayField[0]).toEqual(42);
    expect(modelWithAllFieldTypes?.int32ArrayField[1]).toEqual(84);
    expect(modelWithAllFieldTypes?.int32MapField.get(42)).toEqual(84);
    expect(modelWithAllFieldTypes?.int32MapFieldEmbedded.get(42)).toEqual(emptyModel);

    expect(modelWithAllFieldTypes?.int64Field).toEqual(BigInt(100));
    expect(modelWithAllFieldTypes?.int64ArrayField).toHaveLength(2);
    expect(modelWithAllFieldTypes?.int64ArrayField[0]).toEqual(BigInt(100));
    expect(modelWithAllFieldTypes?.int64ArrayField[1]).toEqual(BigInt(200));
    expect(modelWithAllFieldTypes?.int64MapField.get(BigInt(100))).toEqual(BigInt(200));
    expect(modelWithAllFieldTypes?.int64MapFieldEmbedded.get(BigInt(100))).toEqual(emptyModel);

    expect(modelWithAllFieldTypes?.uint32Field).toEqual(42);
    expect(modelWithAllFieldTypes?.uint32ArrayField).toHaveLength(2);
    expect(modelWithAllFieldTypes?.uint32ArrayField[0]).toEqual(42);
    expect(modelWithAllFieldTypes?.uint32ArrayField[1]).toEqual(84);
    expect(modelWithAllFieldTypes?.uint32MapField.get(42)).toEqual(84);
    expect(modelWithAllFieldTypes?.uint32MapFieldEmbedded.get(42)).toEqual(emptyModel);

    expect(modelWithAllFieldTypes?.uint64Field).toEqual(BigInt(100));
    expect(modelWithAllFieldTypes?.uint64ArrayField).toHaveLength(2);
    expect(modelWithAllFieldTypes?.uint64ArrayField[0]).toEqual(BigInt(100));
    expect(modelWithAllFieldTypes?.uint64ArrayField[1]).toEqual(BigInt(200));
    expect(modelWithAllFieldTypes?.uint64MapField.get(BigInt(100))).toEqual(BigInt(200));
    expect(modelWithAllFieldTypes?.uint64MapFieldEmbedded.get(BigInt(100))).toEqual(emptyModel);

    expect(modelWithAllFieldTypes?.float32Field).toEqual(42.0);
    expect(modelWithAllFieldTypes?.float32ArrayField).toHaveLength(2);
    expect(modelWithAllFieldTypes?.float32ArrayField[0]).toEqual(42.0);
    expect(modelWithAllFieldTypes?.float32ArrayField[1]).toEqual(84.0);

    expect(modelWithAllFieldTypes?.float64Field).toEqual(100.0);
    expect(modelWithAllFieldTypes?.float64ArrayField).toHaveLength(2);
    expect(modelWithAllFieldTypes?.float64ArrayField[0]).toEqual(100.0);
    expect(modelWithAllFieldTypes?.float64ArrayField[1]).toEqual(200.0);

    expect(modelWithAllFieldTypes?.boolField).toEqual(false);
    expect(modelWithAllFieldTypes?.boolArrayField).toHaveLength(2);
    expect(modelWithAllFieldTypes?.boolArrayField[0]).toEqual(true);
    expect(modelWithAllFieldTypes?.boolArrayField[1]).toEqual(false);

    expect(modelWithAllFieldTypes?.bytesField).toEqual(Buffer.from([42, 84]));
    expect(modelWithAllFieldTypes?.bytesArrayField).toHaveLength(2);
    expect(modelWithAllFieldTypes?.bytesArrayField[0]).toEqual(Buffer.from([42, 84]));
    expect(modelWithAllFieldTypes?.bytesArrayField[1]).toEqual(Buffer.from([84, 42]));

    expect(modelWithAllFieldTypes?.enumField).toEqual(generated.GenericEnum.SecondValue);
    expect(modelWithAllFieldTypes?.enumArrayField).toHaveLength(2);
    expect(modelWithAllFieldTypes?.enumArrayField[0]).toEqual(generated.GenericEnum.FirstValue);
    expect(modelWithAllFieldTypes?.enumArrayField[1]).toEqual(generated.GenericEnum.SecondValue);
    expect(modelWithAllFieldTypes?.enumMapField.get(generated.GenericEnum.FirstValue)).toEqual("hello world");
    expect(modelWithAllFieldTypes?.enumMapFieldEmbedded.get(generated.GenericEnum.FirstValue)).toEqual(emptyModel);

    expect(modelWithAllFieldTypes?.modelArrayField).toHaveLength(2);
});