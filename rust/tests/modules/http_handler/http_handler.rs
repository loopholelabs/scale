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


use scale_signature_http::context::Context;

pub fn scale(ctx: &mut Context) -> Result<&mut Context, Box<dyn std::error::Error>> {
    let response = ctx.response();
    let mut body = response.body();
    body.append(&mut "-modified".as_bytes().to_vec());
    response.set_body_bytes(body);
    Ok(ctx)
}
