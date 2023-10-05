use signature::types;

pub fn scale(
    ctx: Option<types::Context>,
) -> Result<Option<types::Context>, Box<dyn std::error::Error>> {
    let mut unwrapped = ctx.unwrap();
    unwrapped.data = prettyplease::unparse(&syn::parse_str(unwrapped.data.as_str()).unwrap());
    Ok(Some(unwrapped))
}
