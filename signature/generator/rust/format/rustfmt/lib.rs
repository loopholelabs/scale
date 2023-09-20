use signature::types;

pub fn scale(
    ctx: Option<&mut types::Context>,
) -> Result<Option<types::Context>, Box<dyn std::error::Error>> {
    let unwrapped = ctx.unwrap();
    unwrapped.data = prettyplease::unparse(&syn::parse_str(unwrapped.data.as_str()).unwrap());
    signature::next(Some(unwrapped))
}
