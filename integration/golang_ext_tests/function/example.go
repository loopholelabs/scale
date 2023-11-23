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

package example

import (
	extension "extension"
	"fmt"
	"signature"
)

func Example(ctx *signature.ModelWithAllFieldTypes) (*signature.ModelWithAllFieldTypes, error) {
	fmt.Printf("This is a Golang Function, which calls an extension.")

	v, err := extension.New(&extension.Stringval{Value: ""})
	if err != nil {
		return nil, err
	}
	hello, err := v.Hello(&extension.Stringval{Value: ""})
	if err != nil {
		return nil, err
	}
	world, err := extension.World(&extension.Stringval{Value: ""})
	if err != nil {
		return nil, err
	}

	if ctx != nil {
		ctx.StringField = fmt.Sprintf("This is a Golang Function. Extension New().Hello()=%s World()=%s", hello.Value, world.Value)
	}
	return signature.Next(ctx)
}
