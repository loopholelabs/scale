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

	"github.com/stretchr/testify/assert"

	"github.com/loopholelabs/scale/signature"
)

func TestEncodeDecode(t *testing.T) {
	s := &Schema{
		Version:  "invalid",
		Language: Go,
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
	}

	decoded := new(Schema)

	encoded := s.Encode()
	err := decoded.Decode(encoded)
	assert.ErrorIs(t, err, ErrVersion)

	s.Version = V1Alpha
	s.Language = "invalid"

	encoded = s.Encode()
	err = decoded.Decode(encoded)
	assert.ErrorIs(t, err, ErrLanguage)

	masterTestingSchema := new(signature.Schema)
	err = masterTestingSchema.Decode([]byte(signature.MasterTestingSchema))
	assert.NoError(t, err)

	dependencies := make([]Dependency, 3)
	dependencies[0] = Dependency{
		Name:     "Test Dependency 1",
		Version:  "1.0.0",
		Metadata: make(map[string]string),
	}

	dependencies[0].Metadata["test0"] = "test1"

	dependencies[1] = Dependency{
		Name:     "Test Dependency 2",
		Version:  "2.0.0",
		Metadata: make(map[string]string),
	}

	dependencies[1].Metadata["test2"] = "test3"

	dependencies[2] = Dependency{
		Name:     "Test Dependency 3",
		Version:  "3.0.0",
		Metadata: make(map[string]string),
	}

	dependencies[2].Metadata["test4"] = "test5"

	s = &Schema{
		Version:         V1Alpha,
		Name:            "Test Name",
		Tag:             "Test Tag",
		SignatureName:   "Test Signature",
		SignatureSchema: masterTestingSchema,
		Dependencies:    dependencies,
		Language:        Go,
		Function:        []byte("Test Function Contents"),
	}

	encoded = s.Encode()
	err = decoded.Decode(encoded)
	assert.NoError(t, err)

	assert.Equal(t, s.Version, decoded.Version)
	assert.Equal(t, s.Name, decoded.Name)
	assert.Equal(t, s.Tag, decoded.Tag)
	assert.Equal(t, s.Language, decoded.Language)
	assert.Equal(t, s.Function, decoded.Function)
	assert.Equal(t, s.SignatureName, decoded.SignatureName)

	encoded[decoded.Size+uint32(len(s.Hash))-1] = 0
	err = decoded.Decode(encoded)
	assert.ErrorIs(t, err, ErrHash)
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
