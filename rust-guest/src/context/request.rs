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
use super::generated::{Request, StringList};
use std::string::String;
use std::collections::HashMap;

pub trait MutableRequest {
    fn method(&mut self) -> &String;
    fn set_method(&mut self, method: String) -> &mut Self;
    fn remote_ip(&self) -> &String;
    fn body(&mut self) -> Vec<u8>;
    fn set_body(&mut self, body: String) -> &mut Self;
    fn set_body_bytes(&mut self, bytes: Vec<u8>) -> &mut Self;
    fn headers(&self) -> &HashMap<String, StringList>;
    fn get_headers(&self, key: &String) -> Option<&StringList>;
    fn set_headers(&mut self, key: String, value: Vec<String>);
}

impl MutableRequest for Request {
    fn method(&mut self) -> &String {
        &self.method
    }

    fn set_method(&mut self, method: String) -> &mut Self {
        self.method = method;
        self
    }

    fn remote_ip(&self) -> &String {
        &self.i_p
    }

    fn body(&mut self) -> Vec<u8> {
        self.body.clone()
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
