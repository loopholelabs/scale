use polyglot_rs::{Decoder, Encoder};
use std::io::Cursor;

pub trait Encode {
    fn encode(self, b: &mut Cursor<Vec<u8>>) -> Result<&mut Cursor<Vec<u8>>, Box<dyn std::error::Error>>;
    fn internal_error(self, b: &mut Cursor<Vec<u8>>, error: Box<dyn std::error::Error>);
}

pub trait Decode {
    fn decode(b: &mut Cursor<&mut Vec<u8>>) -> Result<Option<Self>, Box<dyn std::error::Error>>
    where
        Self: Sized;
}

#[derive(Clone)]
pub struct BadContext {
    pub data: u32
}

impl Encode for BadContext {
    fn encode(self, b: &mut Cursor<Vec<u8>>) -> Result<&mut Cursor<Vec<u8>>, Box<dyn std::error::Error>> {
        b.encode_u32(self.data)?;
        Ok(b)
    }

    fn internal_error(self, b: &mut Cursor<Vec<u8>>, error: Box<dyn std::error::Error>) {
        b.encode_error(error).unwrap();
    }
}

impl Decode for BadContext {
    fn decode(b: &mut Cursor<&mut Vec<u8>>) -> Result<Option<BadContext>, Box<dyn std::error::Error>> {
        if b.decode_none() {
            return Ok(None);
        }

        if let Ok(error) = b.decode_error() {
            return Err(error);
        }

        Ok(Some(BadContext {
            data: b.decode_u32()?,
        }))
    }
}