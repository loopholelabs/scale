mod opt;
mod options;

extern crate anyhow;
extern crate binaryen;
extern crate wizer;
extern crate structopt;

use crate::options::Options;
use anyhow::{bail, Context, Result};
use std::env;
use std::{fs, process::Command};
use structopt::StructOpt;

fn main() -> Result<()> {
    let opts = Options::from_args();
    let wizen = env::var("JS_WIZEN");

    if wizen.eq(&Ok("1".into())) {
        let wasm: &[u8] = include_bytes!(concat!(env!("OUT_DIR"), "/engine.wasm"));
        opt::Optimizer::new(wasm)
            .optimize(true)
            .write_optimized_wasm(opts.output)?;

        env::remove_var("JS_WIZEN");

        return Ok(());
    }

    let contents = fs::File::open(&opts.input)
        .with_context(|| format!("Failed to open input file {}", opts.input.display()))?;
    let self_cmd = env::args().next().unwrap();

    env::set_var("JS_WIZEN", "1");
    let status = Command::new(self_cmd)
        .arg(&opts.input)
        .arg("-o")
        .arg(&opts.output)
        .stdin(contents)
        .status()?;

    if !status.success() {
        bail!("Couldn't create wasm from input");
    }

    Ok(())
}
