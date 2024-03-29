// Code generated by scale-signature , DO NOT EDIT.
// output: local_rustfmt_latest_guest

#![allow(dead_code)]
#![allow(unused_imports)]
#![allow(unused_variables)]
#![allow(unused_mut)]

use num_enum::TryFromPrimitive;
use polyglot_rs::{Decoder, DecodingError, Encoder, Kind};
use regex::Regex;
use std::collections::HashMap;
use std::convert::TryFrom;
use std::io::Cursor;

pub trait Encode {
    fn encode<'a>(
        a: Option<&Self>,
        b: &'a mut Cursor<Vec<u8>>,
    ) -> Result<&'a mut Cursor<Vec<u8>>, Box<dyn std::error::Error>>
    where
        Self: Sized;
}

trait EncodeSelf {
    fn encode_self<'a, 'b>(
        &'b self,
        b: &'a mut Cursor<Vec<u8>>,
    ) -> Result<&'a mut Cursor<Vec<u8>>, Box<dyn std::error::Error>>;
}

pub trait Decode {
    fn decode(b: &mut Cursor<&mut Vec<u8>>) -> Result<Option<Self>, Box<dyn std::error::Error>>
    where
        Self: Sized;
}
#[derive(Clone, Debug, PartialEq)]
pub struct Context {
    pub data: String,
}

impl Context {
    pub fn new() -> Self {
        Self {
            data: "".to_string(),
        }
    }
}

impl Encode for Context {
    fn encode<'a>(
        a: Option<&Context>,
        e: &'a mut Cursor<Vec<u8>>,
    ) -> Result<&'a mut Cursor<Vec<u8>>, Box<dyn std::error::Error>> {
        a.encode_self(e)
    }
}

impl EncodeSelf for Context {
    fn encode_self<'a, 'b>(
        &'b self,
        e: &'a mut Cursor<Vec<u8>>,
    ) -> Result<&'a mut Cursor<Vec<u8>>, Box<dyn std::error::Error>> {
        e.encode_string(&self.data)?;

        Ok(e)
    }
}

impl EncodeSelf for Option<&Context> {
    fn encode_self<'a, 'b>(
        &'b self,
        e: &'a mut Cursor<Vec<u8>>,
    ) -> Result<&'a mut Cursor<Vec<u8>>, Box<dyn std::error::Error>> {
        if let Some(x) = self {
            x.encode_self(e)?;
        } else {
            e.encode_none()?;
        }
        Ok(e)
    }
}

impl EncodeSelf for Option<Context> {
    fn encode_self<'a, 'b>(
        &'b self,
        e: &'a mut Cursor<Vec<u8>>,
    ) -> Result<&'a mut Cursor<Vec<u8>>, Box<dyn std::error::Error>> {
        if let Some(x) = self {
            x.encode_self(e)?;
        } else {
            e.encode_none()?;
        }
        Ok(e)
    }
}

impl Decode for Context {
    fn decode(d: &mut Cursor<&mut Vec<u8>>) -> Result<Option<Context>, Box<dyn std::error::Error>> {
        if d.decode_none() {
            return Ok(None);
        }

        if let Ok(error) = d.decode_error() {
            return Err(error);
        }

        let mut x = Context::new();

        x.data = d.decode_string()?;

        Ok(Some(x))
    }
}
