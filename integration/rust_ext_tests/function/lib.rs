/*
    Copyright 2023 Loophole Labs

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

use signature::types;

use extension::types::Stringval;
use extension::New;
use extension::World;
use extension::Example;

pub fn example(
    ctx: Option<types::ModelWithAllFieldTypes>,
) -> Result<Option<types::ModelWithAllFieldTypes>, Box<dyn std::error::Error>> {
    println!("This is a Rust Function, which calls an extension.");

    let ex_op = New(Stringval{ value: "".to_string() });

    if let Err(e) = ex_op {
      return Err(e)
    }

    let ex = ex_op.unwrap();

    let hello_op = ex.unwrap().Hello(Stringval{ value: "".to_string() });

    if let Err(e) = hello_op {
      return Err(e)
    }

    let hello = hello_op.unwrap().unwrap().value;

    let world_op = World(Stringval{ value: "".to_string() });

    if let Err(e) = world_op {
      return Err(e)
    }

    let world = world_op.unwrap().unwrap().value;

    let mut unwrapped = ctx.unwrap();
    unwrapped.string_field = format!("This is a Rust Function. Extension New().Hello()={hello} World()={world}");
    return signature::next(Some(unwrapped));
    // return Ok(Some(unwrapped))
}
