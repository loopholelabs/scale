use std::path::PathBuf;
use structopt::StructOpt;

#[derive(Debug, StructOpt)]
#[structopt(name = "jsbuilder", about = "JavaScript to WebAssembly")]
pub struct Options {
    #[structopt(parse(from_os_str))]
    pub input: PathBuf,

    #[structopt(short = "o", parse(from_os_str), default_value = "index.wasm")]
    pub output: PathBuf,
}
