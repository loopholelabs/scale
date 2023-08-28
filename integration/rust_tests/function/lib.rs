use signature::types;

pub fn example(ctx: Option<&mut types::ModelWithAllFieldTypes>) -> Result<Option<types::ModelWithAllFieldTypes>, Box<dyn std::error::Error>> {
    let c = ctx.unwrap();
    c.string_field = "TEST".to_string();
    return signature::next(Some(c));
}