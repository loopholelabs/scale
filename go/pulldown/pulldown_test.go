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

package pulldown

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/loopholelabs/scalefile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPulldownCache(t *testing.T) {
	sc := StoreConfig{
		CacheDirectory: "testCache",
		PullPolicy:     POLICY_PULL_IF_NOT_PRESENT,
	}

	pc := PulldownConfig{
		APIKey:       "123",
		Organization: "loophole",
		Function:     "Test1",
		Tag:          "0.1.0",
	}

	sf := &scalefile.ScaleFile{
		Version:   scalefile.V1Alpha,
		Name:      "Test1",
		Signature: "signature1",
		Language:  scalefile.Go,
		Source:    "Hello world",
	}

	var b bytes.Buffer
	bwriter := bufio.NewWriter(&b)

	err := scalefile.Encode(bwriter, sf)
	require.NoError(t, err)
	err = bwriter.Flush()
	require.NoError(t, err)

	// Write it to the cache...
	err = saveToCache(pc, sc, b.Bytes())
	require.NoError(t, err)

	// Try reading a scalefile from the cache...
	newsf, err := New(pc, sc)
	require.NoError(t, err)

	// Now assert that the newsf is same as our original sf.
	var newb bytes.Buffer
	newbwriter := bufio.NewWriter(&newb)

	err = scalefile.Encode(newbwriter, newsf)
	require.NoError(t, err)
	err = newbwriter.Flush()
	require.NoError(t, err)

	assert.Equal(t, b.Bytes(), newb.Bytes())

	// cleanup - remove cache dir
	err = removeCache(sc)
	require.NoError(t, err)
}
