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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	s := &ScaleFunc{
		Version:  "invalid",
		Language: Go,
	}
	decoded := new(ScaleFunc)

	encoded := s.Encode()
	err := decoded.Decode(encoded)
	assert.ErrorIs(t, err, VersionErr)

	s.Version = V1Alpha
	s.Language = "invalid"

	encoded = s.Encode()
	err = decoded.Decode(encoded)
	assert.ErrorIs(t, err, LanguageErr)

	s = &ScaleFunc{
		Version:  V1Alpha,
		Name:     "Test Name",
		Language: Go,
		Function: []byte("Test Function Contents"),
	}

	encoded = s.Encode()
	err = decoded.Decode(encoded)
	assert.NoError(t, err)

	assert.Equal(t, s.Version, decoded.Version)
	assert.Equal(t, s.Name, decoded.Name)
	assert.Equal(t, s.Language, decoded.Language)
	assert.Equal(t, s.Function, decoded.Function)

	encoded[decoded.Size+uint32(len(s.Checksum))-1] = 0
	err = decoded.Decode(encoded)
	assert.ErrorIs(t, err, ChecksumErr)
}
