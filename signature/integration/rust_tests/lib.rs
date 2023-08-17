pub mod generated;

#[cfg(test)]
mod tests {
    use std::fs;
    use std::error::Error;
    use std::io::{Cursor, Seek, SeekFrom};
    use crate::generated::{Decode, Encode};

    use super::*;

    #[test]
    fn test_output() -> Result<(), Box<dyn Error>> {
        let mut buf = Cursor::new(Vec::new());

        let nil_model: Option<generated::EmptyModel> = None;
        generated::EmptyModel::encode(&nil_model, &mut buf)?;
        fs::write("../test_data/nil_model.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let empty_model = generated::EmptyModel::new();
        generated::EmptyModel::encode(&Some(empty_model.clone()), &mut buf)?;
        fs::write("../test_data/empty_model.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let empty_model_with_description = generated::EmptyModelWithDescription::new();
        generated::EmptyModelWithDescription::encode(&Some(empty_model_with_description), &mut buf)?;
        fs::write("../test_data/empty_model_with_description.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let mut model_with_single_string_field = generated::ModelWithSingleStringField::new();
        assert_eq!(model_with_single_string_field.string_field, String::from("DefaultValue"));
        model_with_single_string_field.string_field = String::from("hello world");
        generated::ModelWithSingleStringField::encode(&Some(model_with_single_string_field), &mut buf)?;
        fs::write("../test_data/model_with_single_string_field.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let mut model_with_single_string_field_and_description = generated::ModelWithSingleStringFieldAndDescription::new();
        assert_eq!(model_with_single_string_field_and_description.string_field, String::from("DefaultValue"));
        model_with_single_string_field_and_description.string_field = String::from("hello world");
        generated::ModelWithSingleStringFieldAndDescription::encode(&Some(model_with_single_string_field_and_description), &mut buf)?;
        fs::write("../test_data/model_with_single_string_field_and_description.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let mut model_with_single_int32_field = generated::ModelWithSingleInt32Field::new();
        assert_eq!(model_with_single_int32_field.int32_field, 32);
        model_with_single_int32_field.int32_field = 42;
        generated::ModelWithSingleInt32Field::encode(&Some(model_with_single_int32_field), &mut buf)?;
        fs::write("../test_data/model_with_single_int32_field.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let mut model_with_single_int32_field_and_description = generated::ModelWithSingleInt32FieldAndDescription::new();
        assert_eq!(model_with_single_int32_field_and_description.int32_field, 32);
        model_with_single_int32_field_and_description.int32_field = 42;
        generated::ModelWithSingleInt32FieldAndDescription::encode(&Some(model_with_single_int32_field_and_description), &mut buf)?;
        fs::write("../test_data/model_with_single_int32_field_and_description.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let mut model_with_multiple_fields = generated::ModelWithMultipleFields::new();
        assert_eq!(model_with_multiple_fields.string_field, String::from("DefaultValue"));
        assert_eq!(model_with_multiple_fields.int32_field, 32);
        model_with_multiple_fields.string_field = String::from("hello world");
        model_with_multiple_fields.int32_field = 42;
        generated::ModelWithMultipleFields::encode(&Some(model_with_multiple_fields), &mut buf)?;
        fs::write("../test_data/model_with_multiple_fields.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let mut model_with_multiple_fields_and_description = generated::ModelWithMultipleFieldsAndDescription::new();
        assert_eq!(model_with_multiple_fields_and_description.string_field, String::from("DefaultValue"));
        assert_eq!(model_with_multiple_fields_and_description.int32_field, 32);
        model_with_multiple_fields_and_description.string_field = String::from("hello world");
        model_with_multiple_fields_and_description.int32_field = 42;
        generated::ModelWithMultipleFieldsAndDescription::encode(&Some(model_with_multiple_fields_and_description), &mut buf)?;
        fs::write("../test_data/model_with_multiple_fields_and_description.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let mut model_with_enum = generated::ModelWithEnum::new();
        assert_eq!(model_with_enum.enum_field, generated::GenericEnum::DefaultValue);
        model_with_enum.enum_field = generated::GenericEnum::SecondValue;
        generated::ModelWithEnum::encode(&Some(model_with_enum), &mut buf)?;
        fs::write("../test_data/model_with_enum.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let mut model_with_enum_and_description = generated::ModelWithEnumAndDescription::new();
        assert_eq!(model_with_enum_and_description.enum_field, generated::GenericEnum::DefaultValue);
        model_with_enum_and_description.enum_field = generated::GenericEnum::SecondValue;
        generated::ModelWithEnumAndDescription::encode(&Some(model_with_enum_and_description), &mut buf)?;
        fs::write("../test_data/model_with_enum_and_description.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let mut model_with_enum_accessor = generated::ModelWithEnumAccessor::new();
        let default_enum_value = model_with_enum_accessor.get_enum_field();
        assert_eq!(*default_enum_value, generated::GenericEnum::DefaultValue);
        model_with_enum_accessor.set_enum_field(generated::GenericEnum::SecondValue);
        generated::ModelWithEnumAccessor::encode(&Some(model_with_enum_accessor), &mut buf)?;
        fs::write("../test_data/model_with_enum_accessor.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let mut model_with_enum_accessor_and_description = generated::ModelWithEnumAccessorAndDescription::new();
        let default_enum_value = model_with_enum_accessor_and_description.get_enum_field();
        assert_eq!(*default_enum_value, generated::GenericEnum::DefaultValue);
        model_with_enum_accessor_and_description.set_enum_field(generated::GenericEnum::SecondValue);
        generated::ModelWithEnumAccessorAndDescription::encode(&Some(model_with_enum_accessor_and_description), &mut buf)?;
        fs::write("../test_data/model_with_enum_accessor_and_description.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let mut model_with_multiple_fields_accessor = generated::ModelWithMultipleFieldsAccessor::new();
        let mut string_value = model_with_multiple_fields_accessor.get_string_field();
        let int32_value = model_with_multiple_fields_accessor.get_int32_field();
        assert_eq!(string_value, String::from("DefaultValue"));
        assert_eq!(int32_value, 32);
        assert!(model_with_multiple_fields_accessor.set_string_field(String::from("hello world")).is_err());
        assert!(model_with_multiple_fields_accessor.set_string_field(String::from("")).is_err());
        model_with_multiple_fields_accessor.set_string_field(String::from("hello"))?;
        string_value = model_with_multiple_fields_accessor.get_string_field();
        assert_eq!(string_value, String::from("HELLO"));
        assert!(model_with_multiple_fields_accessor.set_int32_field(-1).is_err());
        assert!(model_with_multiple_fields_accessor.set_int32_field(101).is_err());
        model_with_multiple_fields_accessor.set_int32_field(42)?;
        generated::ModelWithMultipleFieldsAccessor::encode(&Some(model_with_multiple_fields_accessor.clone()), &mut buf)?;
        fs::write("../test_data/model_with_multiple_fields_accessor.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let mut model_with_multiple_fields_accessor_and_description = generated::ModelWithMultipleFieldsAccessorAndDescription::new();
        let string_value = model_with_multiple_fields_accessor_and_description.get_string_field();
        let int32_value = model_with_multiple_fields_accessor_and_description.get_int32_field();
        assert_eq!(string_value, String::from("DefaultValue"));
        assert_eq!(int32_value, 32);
        model_with_multiple_fields_accessor_and_description.set_string_field(String::from("hello world"))?;
        model_with_multiple_fields_accessor_and_description.set_int32_field(42)?;
        generated::ModelWithMultipleFieldsAccessorAndDescription::encode(&Some(model_with_multiple_fields_accessor_and_description), &mut buf)?;
        fs::write("../test_data/model_with_multiple_fields_accessor_and_description.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let mut model_with_embedded_models = generated::ModelWithEmbeddedModels::new();
        assert!(model_with_embedded_models.embedded_empty_model.is_some());
        assert_eq!(model_with_embedded_models.embedded_model_array_with_multiple_fields_accessor.capacity(), 64);
        assert_eq!(model_with_embedded_models.embedded_model_array_with_multiple_fields_accessor.len(), 0);
        model_with_embedded_models.embedded_model_array_with_multiple_fields_accessor.push(model_with_multiple_fields_accessor.clone());
        generated::ModelWithEmbeddedModels::encode(&Some(model_with_embedded_models), &mut buf)?;
        fs::write("../test_data/model_with_embedded_models.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let mut model_with_embedded_models_and_description = generated::ModelWithEmbeddedModelsAndDescription::new();
        assert!(model_with_embedded_models_and_description.embedded_empty_model.is_some());
        assert_eq!(model_with_embedded_models_and_description.embedded_model_array_with_multiple_fields_accessor.capacity(), 0);
        assert_eq!(model_with_embedded_models_and_description.embedded_model_array_with_multiple_fields_accessor.len(), 0);
        model_with_embedded_models_and_description.embedded_model_array_with_multiple_fields_accessor.push(model_with_multiple_fields_accessor.clone());
        generated::ModelWithEmbeddedModelsAndDescription::encode(&Some(model_with_embedded_models_and_description), &mut buf)?;
        fs::write("../test_data/model_with_embedded_models_and_description.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let mut model_with_embedded_models_accessor = generated::ModelWithEmbeddedModelsAccessor::new();
        let embedded_model = model_with_embedded_models_accessor.get_embedded_empty_model();
        assert!(embedded_model.is_some());
        let embedded_model_array = model_with_embedded_models_accessor.get_embedded_model_array_with_multiple_fields_accessor().unwrap();
        assert_eq!(embedded_model_array.capacity(), 0);
        assert_eq!(embedded_model_array.len(), 0);
        model_with_embedded_models_accessor.set_embedded_model_array_with_multiple_fields_accessor(vec![model_with_multiple_fields_accessor.clone()]);
        generated::ModelWithEmbeddedModelsAccessor::encode(&Some(model_with_embedded_models_accessor), &mut buf)?;
        fs::write("../test_data/model_with_embedded_models_accessor.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let mut model_with_embedded_models_accessor_and_description = generated::ModelWithEmbeddedModelsAccessorAndDescription::new();
        let embedded_model = model_with_embedded_models_accessor_and_description.get_embedded_empty_model();
        assert!(embedded_model.is_some());
        let embedded_model_array = model_with_embedded_models_accessor_and_description.get_embedded_model_array_with_multiple_fields_accessor().unwrap();
        assert_eq!(embedded_model_array.capacity(), 0);
        assert_eq!(embedded_model_array.len(), 0);
        model_with_embedded_models_accessor_and_description.set_embedded_model_array_with_multiple_fields_accessor(vec![model_with_multiple_fields_accessor.clone()]);
        generated::ModelWithEmbeddedModelsAccessorAndDescription::encode(&Some(model_with_embedded_models_accessor_and_description), &mut buf)?;
        fs::write("../test_data/model_with_embedded_models_accessor_and_description.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        let mut model_with_all_field_types = generated::ModelWithAllFieldTypes::new();

        assert_eq!(model_with_all_field_types.string_field, String::from("DefaultValue"));
        model_with_all_field_types.string_field = "hello world".to_string();
        assert_eq!(model_with_all_field_types.string_array_field.capacity(), 0);
        assert_eq!(model_with_all_field_types.string_array_field.len(), 0);
        model_with_all_field_types.string_array_field.push("hello".to_string());
        model_with_all_field_types.string_array_field.push("world".to_string());
        assert_eq!(model_with_all_field_types.string_map_field.capacity(), 0);
        assert_eq!(model_with_all_field_types.string_map_field.len(), 0);
        model_with_all_field_types.string_map_field.insert("hello".to_string(), "world".to_string());
        model_with_all_field_types.string_map_field_embedded.insert("hello".to_string(), empty_model.clone());

        assert_eq!(model_with_all_field_types.int32_field, 32);
        model_with_all_field_types.int32_field = 42;
        assert_eq!(model_with_all_field_types.int32_array_field.capacity(), 0);
        assert_eq!(model_with_all_field_types.int32_array_field.len(), 0);
        model_with_all_field_types.int32_array_field.push(42);
        model_with_all_field_types.int32_array_field.push(84);
        assert_eq!(model_with_all_field_types.int32_map_field.capacity(), 0);
        assert_eq!(model_with_all_field_types.int32_map_field.len(), 0);
        model_with_all_field_types.int32_map_field.insert(42, 84);
        model_with_all_field_types.int32_map_field_embedded.insert(42, empty_model.clone());

        assert_eq!(model_with_all_field_types.int64_field, 64);
        model_with_all_field_types.int64_field = 100;
        assert_eq!(model_with_all_field_types.int64_array_field.capacity(), 0);
        assert_eq!(model_with_all_field_types.int64_array_field.len(), 0);
        model_with_all_field_types.int64_array_field.push(100);
        model_with_all_field_types.int64_array_field.push(200);
        assert_eq!(model_with_all_field_types.int64_map_field.capacity(), 0);
        assert_eq!(model_with_all_field_types.int64_map_field.len(), 0);
        model_with_all_field_types.int64_map_field.insert(100, 200);
        model_with_all_field_types.int64_map_field_embedded.insert(100, empty_model.clone());

        assert_eq!(model_with_all_field_types.uint32_field, 32);
        model_with_all_field_types.uint32_field = 42;
        assert_eq!(model_with_all_field_types.uint32_array_field.capacity(), 0);
        assert_eq!(model_with_all_field_types.uint32_array_field.len(), 0);
        model_with_all_field_types.uint32_array_field.push(42);
        model_with_all_field_types.uint32_array_field.push(84);
        assert_eq!(model_with_all_field_types.uint32_map_field.capacity(), 0);
        assert_eq!(model_with_all_field_types.uint32_map_field.len(), 0);
        model_with_all_field_types.uint32_map_field.insert(42, 84);
        model_with_all_field_types.uint32_map_field_embedded.insert(42, empty_model.clone());

        assert_eq!(model_with_all_field_types.uint64_field, 64);
        model_with_all_field_types.uint64_field = 100;
        assert_eq!(model_with_all_field_types.uint64_array_field.capacity(), 0);
        assert_eq!(model_with_all_field_types.uint64_array_field.len(), 0);
        model_with_all_field_types.uint64_array_field.push(100);
        model_with_all_field_types.uint64_array_field.push(200);
        assert_eq!(model_with_all_field_types.uint64_map_field.capacity(), 0);
        assert_eq!(model_with_all_field_types.uint64_map_field.len(), 0);
        model_with_all_field_types.uint64_map_field.insert(100, 200);
        model_with_all_field_types.uint64_map_field_embedded.insert(100, empty_model.clone());

        assert_eq!(model_with_all_field_types.float32_field, 32.32);
        model_with_all_field_types.float32_field = 42.0;
        assert_eq!(model_with_all_field_types.float32_array_field.capacity(), 0);
        assert_eq!(model_with_all_field_types.float32_array_field.len(), 0);
        model_with_all_field_types.float32_array_field.push(42.0);
        model_with_all_field_types.float32_array_field.push(84.0);

        assert_eq!(model_with_all_field_types.float64_field, 64.64);
        model_with_all_field_types.float64_field = 100.0;
        assert_eq!(model_with_all_field_types.float64_array_field.capacity(), 0);
        assert_eq!(model_with_all_field_types.float64_array_field.len(), 0);
        model_with_all_field_types.float64_array_field.push(100.0);
        model_with_all_field_types.float64_array_field.push(200.0);

        assert!(model_with_all_field_types.bool_field);
        model_with_all_field_types.bool_field = false;
        assert_eq!(model_with_all_field_types.bool_array_field.capacity(), 0);
        assert_eq!(model_with_all_field_types.bool_array_field.len(), 0);
        model_with_all_field_types.bool_array_field.push(true);
        model_with_all_field_types.bool_array_field.push(false);

        assert_eq!(model_with_all_field_types.bytes_field.capacity(), 512);
        assert_eq!(model_with_all_field_types.bytes_field.len(), 0);
        model_with_all_field_types.bytes_field.extend_from_slice(&[42, 84]);
        assert_eq!(model_with_all_field_types.bytes_array_field.capacity(), 0);
        assert_eq!(model_with_all_field_types.bytes_array_field.len(), 0);
        model_with_all_field_types.bytes_array_field.push(vec![42, 84]);
        model_with_all_field_types.bytes_array_field.push(vec![84, 42]);

        assert_eq!(model_with_all_field_types.enum_field, generated::GenericEnum::DefaultValue);
        model_with_all_field_types.enum_field = generated::GenericEnum::SecondValue;
        assert_eq!(model_with_all_field_types.enum_array_field.capacity(), 0);
        assert_eq!(model_with_all_field_types.enum_array_field.len(), 0);
        model_with_all_field_types.enum_array_field.push(generated::GenericEnum::FirstValue);
        model_with_all_field_types.enum_array_field.push(generated::GenericEnum::SecondValue);
        assert_eq!(model_with_all_field_types.enum_map_field.capacity(), 0);
        assert_eq!(model_with_all_field_types.enum_map_field.len(), 0);
        model_with_all_field_types.enum_map_field.insert(generated::GenericEnum::FirstValue, "hello world".to_string());
        model_with_all_field_types.enum_map_field_embedded.insert(generated::GenericEnum::FirstValue, empty_model.clone());

        assert!(model_with_all_field_types.model_field.is_some());
        assert_eq!(model_with_all_field_types.model_array_field.len(), 0);
        model_with_all_field_types.model_array_field.push(empty_model.clone());
        model_with_all_field_types.model_array_field.push(empty_model.clone());

        generated::ModelWithAllFieldTypes::encode(&Some(model_with_all_field_types), &mut buf)?;
        fs::write("../test_data/model_with_all_field_types.bin", buf.get_ref())?;
        buf.seek(SeekFrom::Start(0))?;
        buf.get_mut().clear();

        Ok(())
    }

    #[test]
    fn test_input() -> Result<(), Box<dyn Error>> {
        let mut null_model_data = fs::read("../test_data/nil_model.bin")?;
        let _null_model = generated::EmptyModel::decode(&mut Cursor::new(&mut null_model_data))?;
        assert!(_null_model.is_none());

        let mut empty_model_data = fs::read("../test_data/empty_model.bin")?;
        let _empty_model = generated::EmptyModel::decode(&mut Cursor::new(&mut empty_model_data))?;
        assert!(_empty_model.is_some());

        let mut empty_model_with_description_data = fs::read("../test_data/empty_model_with_description.bin")?;
        let _empty_model_with_description = generated::EmptyModelWithDescription::decode(&mut Cursor::new(&mut empty_model_with_description_data))?;
        assert!(_empty_model_with_description.is_some());

        let mut model_with_single_string_field_data = fs::read("../test_data/model_with_single_string_field.bin")?;
        let model_with_single_string_field = generated::ModelWithSingleStringField::decode(&mut Cursor::new(&mut model_with_single_string_field_data))?.unwrap();
        assert_eq!(model_with_single_string_field.string_field, "hello world");

        let mut model_with_single_string_field_and_description_data = fs::read("../test_data/model_with_single_string_field_and_description.bin")?;
        let model_with_single_string_field_and_description = generated::ModelWithSingleStringFieldAndDescription::decode(&mut Cursor::new(&mut model_with_single_string_field_and_description_data))?.unwrap();
        assert_eq!(model_with_single_string_field_and_description.string_field, "hello world");

        let mut model_with_single_int32_field_data = fs::read("../test_data/model_with_single_int32_field.bin")?;
        let model_with_single_int32_field = generated::ModelWithSingleInt32Field::decode(&mut Cursor::new(&mut model_with_single_int32_field_data))?.unwrap();
        assert_eq!(model_with_single_int32_field.int32_field, 42);

        let mut model_with_single_int32_field_and_description_data = fs::read("../test_data/model_with_single_int32_field_and_description.bin")?;
        let model_with_single_int32_field_and_description = generated::ModelWithSingleInt32FieldAndDescription::decode(&mut Cursor::new(&mut model_with_single_int32_field_and_description_data))?.unwrap();
        assert_eq!(model_with_single_int32_field_and_description.int32_field, 42);

        let mut model_with_multiple_fields_data = fs::read("../test_data/model_with_multiple_fields.bin")?;
        let model_with_multiple_fields = generated::ModelWithMultipleFields::decode(&mut Cursor::new(&mut model_with_multiple_fields_data))?.unwrap();
        assert_eq!(model_with_multiple_fields.string_field, "hello world");
        assert_eq!(model_with_multiple_fields.int32_field, 42);

        let mut model_with_multiple_fields_and_description_data = fs::read("../test_data/model_with_multiple_fields_and_description.bin")?;
        let model_with_multiple_fields_and_description = generated::ModelWithMultipleFieldsAndDescription::decode(&mut Cursor::new(&mut model_with_multiple_fields_and_description_data))?.unwrap();
        assert_eq!(model_with_multiple_fields_and_description.string_field, "hello world");
        assert_eq!(model_with_multiple_fields_and_description.int32_field, 42);

        let mut model_with_enum_data = fs::read("../test_data/model_with_enum.bin")?;
        let model_with_enum = generated::ModelWithEnum::decode(&mut Cursor::new(&mut model_with_enum_data))?.unwrap();
        assert_eq!(model_with_enum.enum_field, generated::GenericEnum::SecondValue);

        let mut model_with_enum_and_description_data = fs::read("../test_data/model_with_enum_and_description.bin")?;
        let model_with_enum_and_description = generated::ModelWithEnumAndDescription::decode(&mut Cursor::new(&mut model_with_enum_and_description_data))?.unwrap();
        assert_eq!(model_with_enum_and_description.enum_field, generated::GenericEnum::SecondValue);

        let mut model_with_enum_accessor_data = fs::read("../test_data/model_with_enum_accessor.bin")?;
        let model_with_enum_accessor = generated::ModelWithEnumAccessor::decode(&mut Cursor::new(&mut model_with_enum_accessor_data))?.unwrap();
        let enum_value = model_with_enum_accessor.get_enum_field();
        assert_eq!(*enum_value, generated::GenericEnum::SecondValue);

        let mut model_with_enum_accessor_and_description_data = fs::read("../test_data/model_with_enum_accessor_and_description.bin")?;
        let model_with_enum_accessor_and_description = generated::ModelWithEnumAccessorAndDescription::decode(&mut Cursor::new(&mut model_with_enum_accessor_and_description_data))?.unwrap();
        let enum_value = model_with_enum_accessor_and_description.get_enum_field();
        assert_eq!(*enum_value, generated::GenericEnum::SecondValue);

        let mut model_with_multiple_fields_accessor_data = fs::read("../test_data/model_with_multiple_fields_accessor.bin")?;
        let model_with_multiple_fields_accessor = generated::ModelWithMultipleFieldsAccessor::decode(&mut Cursor::new(&mut model_with_multiple_fields_accessor_data))?.unwrap();
        let string_field_value = model_with_multiple_fields_accessor.get_string_field();
        assert_eq!(string_field_value, "HELLO");
        let int32_field_value = model_with_multiple_fields_accessor.get_int32_field();
        assert_eq!(int32_field_value, 42);

        let mut model_with_multiple_fields_accessor_and_description_data = fs::read("../test_data/model_with_multiple_fields_accessor_and_description.bin")?;
        let model_with_multiple_fields_accessor_and_description = generated::ModelWithMultipleFieldsAccessorAndDescription::decode(&mut Cursor::new(&mut model_with_multiple_fields_accessor_and_description_data))?.unwrap();
        let string_field_value = model_with_multiple_fields_accessor_and_description.get_string_field();
        assert_eq!(string_field_value, "hello world");
        let int32_field_value = model_with_multiple_fields_accessor_and_description.get_int32_field();
        assert_eq!(int32_field_value, 42);

        let mut model_with_embedded_models_data = fs::read("../test_data/model_with_embedded_models.bin")?;
        let model_with_embedded_models = generated::ModelWithEmbeddedModels::decode(&mut Cursor::new(&mut model_with_embedded_models_data))?.unwrap();
        assert!(model_with_embedded_models.embedded_empty_model.is_some());
        assert_eq!(model_with_embedded_models.embedded_model_array_with_multiple_fields_accessor.len(), 1);
        assert_eq!(model_with_embedded_models.embedded_model_array_with_multiple_fields_accessor[0].get_int32_field(), 42);
        assert_eq!(model_with_embedded_models.embedded_model_array_with_multiple_fields_accessor[0].get_string_field(), "HELLO");

        let mut model_with_embedded_models_and_description_data = fs::read("../test_data/model_with_embedded_models_and_description.bin")?;
        let model_with_embedded_models_and_description = generated::ModelWithEmbeddedModelsAndDescription::decode(&mut Cursor::new(&mut model_with_embedded_models_and_description_data))?.unwrap();
        assert!(model_with_embedded_models_and_description.embedded_empty_model.is_some());
        assert_eq!(model_with_embedded_models_and_description.embedded_model_array_with_multiple_fields_accessor.len(), 1);
        assert_eq!(model_with_embedded_models_and_description.embedded_model_array_with_multiple_fields_accessor[0].get_int32_field(), 42);
        assert_eq!(model_with_embedded_models_and_description.embedded_model_array_with_multiple_fields_accessor[0].get_string_field(), "HELLO");

        let mut model_with_embedded_models_accessor_data = fs::read("../test_data/model_with_embedded_models_accessor.bin")?;
        let model_with_embedded_models_accessor = generated::ModelWithEmbeddedModelsAccessor::decode(&mut Cursor::new(&mut model_with_embedded_models_accessor_data))?.unwrap();
        let embedded_empty_model = model_with_embedded_models_accessor.get_embedded_empty_model();
        assert!(embedded_empty_model.is_some());
        let embedded_model_array_with_multiple_fields_accessor = model_with_embedded_models_accessor.get_embedded_model_array_with_multiple_fields_accessor().unwrap();
        assert_eq!(embedded_model_array_with_multiple_fields_accessor.len(), 1);
        assert_eq!(embedded_model_array_with_multiple_fields_accessor[0].get_int32_field(), 42);
        assert_eq!(embedded_model_array_with_multiple_fields_accessor[0].get_string_field(), "HELLO");

        let mut model_with_embedded_models_accessor_and_description_data = fs::read("../test_data/model_with_embedded_models_accessor_and_description.bin")?;
        let model_with_embedded_models_accessor_and_description = generated::ModelWithEmbeddedModelsAccessorAndDescription::decode(&mut Cursor::new(&mut model_with_embedded_models_accessor_and_description_data))?.unwrap();
        let embedded_empty_model = model_with_embedded_models_accessor_and_description.get_embedded_empty_model();
        assert!(embedded_empty_model.is_some());
        let embedded_model_array_with_multiple_fields_accessor = model_with_embedded_models_accessor_and_description.get_embedded_model_array_with_multiple_fields_accessor().unwrap();
        assert_eq!(embedded_model_array_with_multiple_fields_accessor[0].get_int32_field(), 42);
        assert_eq!(embedded_model_array_with_multiple_fields_accessor[0].get_string_field(), "HELLO");

        let mut model_with_all_field_types_data = fs::read("../test_data/model_with_all_field_types.bin")?;
        let model_with_all_field_types = generated::ModelWithAllFieldTypes::decode(&mut Cursor::new(&mut model_with_all_field_types_data))?.unwrap();
        assert_eq!(model_with_all_field_types.string_field, "hello world");
        assert_eq!(model_with_all_field_types.string_array_field.len(), 2);
        assert_eq!(model_with_all_field_types.string_array_field[0], "hello");
        assert_eq!(model_with_all_field_types.string_array_field[1], "world");
        assert_eq!(model_with_all_field_types.string_map_field.get("hello"), Some(&"world".to_string()));
        assert!(model_with_all_field_types.string_map_field_embedded.get("hello").is_some());

        assert_eq!(model_with_all_field_types.int32_field, 42);
        assert_eq!(model_with_all_field_types.int32_array_field.len(), 2);
        assert_eq!(model_with_all_field_types.int32_array_field[0], 42);
        assert_eq!(model_with_all_field_types.int32_array_field[1], 84);
        assert_eq!(model_with_all_field_types.int32_map_field.get(&42), Some(&84));
        assert!(model_with_all_field_types.int32_map_field_embedded.get(&42).is_some());

        assert_eq!(model_with_all_field_types.int64_field, 100);
        assert_eq!(model_with_all_field_types.int64_array_field.len(), 2);
        assert_eq!(model_with_all_field_types.int64_array_field[0], 100);
        assert_eq!(model_with_all_field_types.int64_array_field[1], 200);
        assert_eq!(model_with_all_field_types.int64_map_field.get(&100), Some(&200));
        assert!(model_with_all_field_types.int64_map_field_embedded.get(&100).is_some());

        assert_eq!(model_with_all_field_types.uint32_field, 42);
        assert_eq!(model_with_all_field_types.uint32_array_field.len(), 2);
        assert_eq!(model_with_all_field_types.uint32_array_field[0], 42);
        assert_eq!(model_with_all_field_types.uint32_array_field[1], 84);
        assert_eq!(model_with_all_field_types.uint32_map_field.get(&42), Some(&84));
        assert!(model_with_all_field_types.uint32_map_field_embedded.get(&42).is_some());

        assert_eq!(model_with_all_field_types.uint64_field, 100);
        assert_eq!(model_with_all_field_types.uint64_array_field.len(), 2);
        assert_eq!(model_with_all_field_types.uint64_array_field[0], 100);
        assert_eq!(model_with_all_field_types.uint64_array_field[1], 200);
        assert_eq!(model_with_all_field_types.uint64_map_field.get(&100), Some(&200));
        assert!(model_with_all_field_types.uint64_map_field_embedded.get(&100).is_some());

        assert_eq!(model_with_all_field_types.float32_field, 42.0);
        assert_eq!(model_with_all_field_types.float32_array_field.len(), 2);
        assert_eq!(model_with_all_field_types.float32_array_field[0], 42.0);
        assert_eq!(model_with_all_field_types.float32_array_field[1], 84.0);

        assert_eq!(model_with_all_field_types.float64_field, 100.0);
        assert_eq!(model_with_all_field_types.float64_array_field.len(), 2);
        assert_eq!(model_with_all_field_types.float64_array_field[0], 100.0);
        assert_eq!(model_with_all_field_types.float64_array_field[1], 200.0);

        assert_eq!(model_with_all_field_types.bool_field, false);
        assert_eq!(model_with_all_field_types.bool_array_field.len(), 2);
        assert_eq!(model_with_all_field_types.bool_array_field[0], true);
        assert_eq!(model_with_all_field_types.bool_array_field[1], false);

        assert_eq!(model_with_all_field_types.bytes_field, &[42, 84]);
        assert_eq!(model_with_all_field_types.bytes_array_field.len(), 2);
        assert_eq!(model_with_all_field_types.bytes_array_field[0], &[42, 84]);
        assert_eq!(model_with_all_field_types.bytes_array_field[1], &[84, 42]);

        assert_eq!(model_with_all_field_types.enum_field, generated::GenericEnum::SecondValue);
        assert_eq!(model_with_all_field_types.enum_array_field.len(), 2);
        assert_eq!(model_with_all_field_types.enum_array_field[0], generated::GenericEnum::FirstValue);
        assert_eq!(model_with_all_field_types.enum_array_field[1], generated::GenericEnum::SecondValue);
        assert_eq!(model_with_all_field_types.enum_map_field.get(&generated::GenericEnum::FirstValue), Some(&"hello world".to_string()));
        assert!(model_with_all_field_types.enum_map_field_embedded.get(&generated::GenericEnum::FirstValue).is_some());

        assert_eq!(model_with_all_field_types.model_array_field.len(), 2);

        Ok(())
    }
}