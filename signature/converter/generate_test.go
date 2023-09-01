//go:build generate

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

package converter

import (
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/signature/generator/golang"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGenerateConverterSchema(t *testing.T) {
	s := new(signature.Schema)
	err := s.Decode([]byte(testSchema))
	require.NoError(t, err)

	formatted, err := golang.GenerateTypes(s, "generated")
	require.NoError(t, err)

	err = os.WriteFile("./converter_tests/generated.go", formatted, 0644)
	require.NoError(t, err)
}
