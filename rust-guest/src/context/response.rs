/*
	Copyright 2022 Loophole Labs

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		   http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

#![allow(unused_variables)]
use super::generated::{Response, StringList};
use std::collections::HashMap;

pub trait MutableResponse {
    fn status_code(&self) -> &i32;//R
    fn body(&self) -> &Vec<u8>;//R
    fn set_body(&mut self, body: String) -> &mut Self;
    fn set_body_bytes(&mut self, bytes: Vec<u8>) -> &mut Self;
    fn headers(&self) -> &HashMap<String, StringList>;//R
    fn get_headers(&self, key: &String) -> Option<&StringList>;
    fn set_headers(&mut self, key: String, value: Vec<String>);
}

impl MutableResponse for Response {
    fn status_code(&self) -> &i32 {
        &self.status_code
    }

    fn body(&self) -> &Vec<u8> {
        &self.body
    }

    fn set_body(&mut self, body: String) -> &mut Self {
        self.body = body.as_bytes().to_vec();
        self
    }

    fn set_body_bytes(&mut self, bytes: Vec<u8>) -> &mut Self {
        self.body = bytes;
        self
    }

    fn headers(&self) -> &HashMap<String, StringList> {
        &self.headers
    }

    fn get_headers(&self, key: &String) -> Option<&StringList>{
        self.headers.get(key)
    }

    fn set_headers(&mut self, key: String, value: Vec<String>) {
        self.headers.insert(key,  StringList{ value: value });
    }
}
