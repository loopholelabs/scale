//go:build !integration && !generate

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

package golang

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/loopholelabs/scale/signature"
)

func TestGenerator(t *testing.T) {
	s := new(signature.Schema)
	err := s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	formatted, err := Generate(s, "types")
	require.NoError(t, err)

	// os.WriteFile("./generated.txt", formatted, 0644)

	master, err := os.ReadFile("./generated.txt")
	require.NoError(t, err)
	require.Equal(t, string(master), string(formatted))

	t.Log(string(formatted))
}
