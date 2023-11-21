//go:build !integration && !generate

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

package scalefunc

import (
	"testing"

	"github.com/loopholelabs/scale/signature"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {

	t.Run("V1Alpha", func(t *testing.T) {
		v1Alpha := &V1AlphaSchema{
			Language:      Go,
			SignatureName: "Test Signature",
			SignatureSchema: &signature.Schema{
				Version: signature.V1AlphaVersion,
				Context: "ctx",
				Models: []*signature.ModelSchema{
					{
						Name:        "ctx",
						Description: "test",
					},
				},
			},
			SignatureHash: "Test Signature Hash",
		}

		decoded := new(V1AlphaSchema)

		encoded := v1Alpha.Encode()
		err := decoded.Decode(encoded)
		assert.NoError(t, err)

		assert.Equal(t, v1Alpha.Language, decoded.Language)

	})
}

func TestValidName(t *testing.T) {
	assert.True(t, ValidString("test"))
	assert.True(t, ValidString("test1"))
	assert.True(t, ValidString("test.1"))
	assert.True(t, ValidString("te---.-1"))
	assert.True(t, ValidString("test-1"))
	assert.False(t, ValidString("test_1"))
	assert.False(t, ValidString("test 1"))
	assert.False(t, ValidString("test1 "))
	assert.False(t, ValidString(" test1"))
	assert.False(t, ValidString("test1_"))
	assert.False(t, ValidString("test1?"))
	assert.False(t, ValidString("test1!"))
	assert.False(t, ValidString("test1@"))
	assert.False(t, ValidString("test1#"))
	assert.False(t, ValidString("test1$"))
	assert.False(t, ValidString("test1%"))
	assert.False(t, ValidString("test1^"))
	assert.False(t, ValidString("test1&"))
	assert.False(t, ValidString("test1*"))
	assert.False(t, ValidString("test1("))
	assert.False(t, ValidString("test1-1!"))
}
