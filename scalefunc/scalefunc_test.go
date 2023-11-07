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
	"github.com/loopholelabs/scale/signature"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	t.Run("V1Alpha", func(t *testing.T) {
		v1Beta := &V1BetaSchema{
			Language: Go,
			Signature: V1BetaSignature{
				Name: "Test Signature",
				Schema: &signature.Schema{
					Version: signature.V1AlphaVersion,
					Context: "ctx",
					Models: []*signature.ModelSchema{
						{
							Name:        "ctx",
							Description: "test",
						},
					},
				},
				Hash: "Test Signature Hash",
			},
		}

		decoded := new(V1AlphaSchema)

		encoded := v1Beta.Encode()
		err := decoded.Decode(encoded)
		assert.ErrorIs(t, err, ErrVersion)

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

		v1Alpha.Language = "invalid"

		encoded = v1Alpha.Encode()
		err = decoded.Decode(encoded)
		assert.ErrorIs(t, err, ErrLanguage)

		masterTestingSchema := new(signature.Schema)
		err = masterTestingSchema.Decode([]byte(signature.MasterTestingSchema))
		assert.NoError(t, err)

		dependencies := make([]V1AlphaDependency, 3)
		dependencies[0] = V1AlphaDependency{
			Name:     "Test Dependency 1",
			Version:  "1.0.0",
			Metadata: make(map[string]string),
		}

		dependencies[0].Metadata["test0"] = "test1"

		dependencies[1] = V1AlphaDependency{
			Name:     "Test Dependency 2",
			Version:  "2.0.0",
			Metadata: make(map[string]string),
		}

		dependencies[1].Metadata["test2"] = "test3"

		dependencies[2] = V1AlphaDependency{
			Name:     "Test Dependency 3",
			Version:  "3.0.0",
			Metadata: make(map[string]string),
		}

		dependencies[2].Metadata["test4"] = "test5"

		v1Alpha = &V1AlphaSchema{
			Name:            "Test Name",
			Tag:             "Test Tag",
			SignatureName:   "Test Signature",
			SignatureSchema: masterTestingSchema,
			Dependencies:    dependencies,
			Language:        Go,
			Function:        []byte("Test Function Contents"),
		}

		encoded = v1Alpha.Encode()
		err = decoded.Decode(encoded)
		assert.NoError(t, err)

		assert.Equal(t, v1Alpha.Name, decoded.Name)
		assert.Equal(t, v1Alpha.Tag, decoded.Tag)
		assert.Equal(t, v1Alpha.Language, decoded.Language)
		assert.Equal(t, v1Alpha.Function, decoded.Function)
		assert.Equal(t, v1Alpha.SignatureName, decoded.SignatureName)

		encoded[decoded.Size+uint32(len(v1Alpha.Hash))-1] = 0
		err = decoded.Decode(encoded)
		assert.ErrorIs(t, err, ErrHash)
	})

	t.Run("V1Beta", func(t *testing.T) {
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

		decoded := new(V1BetaSchema)

		encoded := v1Alpha.Encode()
		err := decoded.Decode(encoded)
		assert.NoError(t, err)

		assert.Equal(t, v1Alpha.Language, decoded.Language)
		assert.Equal(t, v1Alpha.SignatureName, decoded.Signature.Name)
		assert.Equal(t, v1Alpha.SignatureHash, decoded.Signature.Hash)

		v1Beta := &V1BetaSchema{
			Language: "invalid",
			Signature: V1BetaSignature{
				Schema: &signature.Schema{
					Version: signature.V1AlphaVersion,
					Context: "ctx",
					Models: []*signature.ModelSchema{
						{
							Name:        "ctx",
							Description: "test",
						},
					},
				},
			},
		}

		encoded = v1Beta.Encode()
		err = decoded.Decode(encoded)
		assert.ErrorIs(t, err, ErrLanguage)
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
