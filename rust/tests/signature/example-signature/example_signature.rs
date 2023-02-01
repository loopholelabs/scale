use polyglot_rs::{Decoder, DecodingError, Encoder};
use std::io;
use std::io::Cursor;

pub trait Encode {
    fn encode(self, b: &mut Cursor<Vec<u8>>) -> Result<&mut Cursor<Vec<u8>>, io::Error>;
    fn internal_error(self, b: &mut Cursor<Vec<u8>>, error: &str);
}

pub trait Decode {
    fn decode(b: &mut Cursor<&mut Vec<u8>>) -> Result<Option<Self>, DecodingError>
    where
        Self: Sized;
}

#[derive(Clone)]
pub struct ExampleContext {
    pub data: String
}

impl Encode for ExampleContext {
    fn encode(self, b: &mut Cursor<Vec<u8>>) -> Result<&mut Cursor<Vec<u8>>, io::Error> {
        b.encode_string(&*self.data)?;
        Ok(b)
    }

    fn internal_error(self, b: &mut Cursor<Vec<u8>>, error: &str) {
        b.encode_error(error).unwrap();
    }
}

impl Decode for ExampleContext {
    fn decode(b: &mut Cursor<&mut Vec<u8>>) -> Result<Option<ExampleContext>, DecodingError> {
        if b.decode_none() {
            return Ok(None);
        }

        Ok(Some(ExampleContext {
            data: b.decode_string()?,
        }))
    }
}