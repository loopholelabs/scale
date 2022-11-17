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

package main

import (
	"github.com/loopholelabs/scale/signature/generator"
	"io"
	"os"
)

func main() {
	gen := generator.New()

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	req, err := gen.UnmarshalRequest(data)
	if err != nil {
		panic(err)
	}

	res, err := gen.Generate(req)
	if err != nil {
		panic(err)
	}

	data, err = gen.MarshalResponse(res)
	if err != nil {
		panic(err)
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		panic(err)
	}
}
