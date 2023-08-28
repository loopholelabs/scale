use signature::types;

pub fn example(ctx: Option<&mut types::ModelWithAllFieldTypes>) -> Result<Option<types::ModelWithAllFieldTypes>, Box<dyn std::error::Error>> {
    let unwrapped = ctx.unwrap();
    unwrapped.string_field = "This is a Rust Function".to_string();
    return signature::next(Some(unwrapped));
}