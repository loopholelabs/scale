//go:build !integration

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

package extension

import (
	"github.com/loopholelabs/scale/signature"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchema(t *testing.T) {
	s := new(Schema)
	err := s.Decode([]byte(MasterTestingSchema))
	require.NoError(t, err)

	assert.Equal(t, signature.V1AlphaVersion, s.Version)

	assert.Equal(t, 1, len(s.Functions))

	assert.Equal(t, "New", s.Functions[0].Name)
	assert.Equal(t, "HttpConfig", s.Functions[0].Params)
	assert.Equal(t, "HttpConnector", s.Functions[0].Return)
}
