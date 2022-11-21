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

package signature

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseSignature(t *testing.T) {
	namespace, name, version := ParseSignature("")
	assert.Equal(t, "", namespace)
	assert.Equal(t, "", name)
	assert.Equal(t, "latest", version)

	namespace, name, version = ParseSignature("test")
	assert.Equal(t, "", namespace)
	assert.Equal(t, "test", name)
	assert.Equal(t, "latest", version)

	namespace, name, version = ParseSignature("test@1.0.0")
	assert.Equal(t, "", namespace)
	assert.Equal(t, "test", name)
	assert.Equal(t, "1.0.0", version)

	namespace, name, version = ParseSignature("nm/test")
	assert.Equal(t, "nm", namespace)
	assert.Equal(t, "test", name)
	assert.Equal(t, "latest", version)

	namespace, name, version = ParseSignature("nm/test@v0.0.1")
	assert.Equal(t, "nm", namespace)
	assert.Equal(t, "test", name)
	assert.Equal(t, "v0.0.1", version)
}
