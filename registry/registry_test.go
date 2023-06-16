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

package registry

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/loopholelabs/scale/storage"
	"os"
	"testing"

	"github.com/loopholelabs/scale/scalefunc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testingAPIBaseURl = "https://api.dev.scale.sh/v1"

func TestPulldownCache(t *testing.T) {
	tempDir := t.TempDir()
	st, err := storage.New(tempDir)
	require.NoError(t, err)

	sf := &scalefunc.ScaleFunc{
		Version:   scalefunc.V1Alpha,
		Name:      "Test1",
		Signature: "signature1",
		Tag:       "1",
		Language:  scalefunc.Go,
		Function:  []byte("Hello world"),
	}

	function := "TestFunction"
	tag := "1"

	h := sha256.New()
	h.Write(sf.Function)

	bs := h.Sum(nil)

	hash := hex.EncodeToString(bs)
	err = st.Put(function, tag, DefaultOrganization, hash, sf)
	require.NoError(t, err)

	newsf, err := Download(function, tag, WithCacheDirectory(tempDir), WithPullPolicy(NeverPullPolicy))

	require.NoError(t, err)

	newsf.Size = 0
	newsf.Checksum = ""

	assert.EqualValues(t, sf, newsf)
	require.NoError(t, err)
}

func TestRegistryDownload(t *testing.T) {
	/* This test requires a valid API key for the scale dev api to run, the SCALE_API_KEY environment variable must be
	set in the testing environment. */
	apiKey := os.Getenv("SCALE_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping test, SCALE_API_KEY environment variable not set")
	}
	sf, err := Download("TestRegistryDownload", "1",
		WithAPIKey(apiKey),
		WithBaseURL(testingAPIBaseURl),
		WithOrganization("alex"),
	)
	require.NoError(t, err)
	require.Equal(t, "TestRegistryDownload", sf.Name)
	require.Equal(t, "1", sf.Tag)
	require.Equal(t, "signature1", sf.Signature)
	require.Equal(t, scalefunc.Go, sf.Language)
	require.Equal(t, scalefunc.V1Alpha, sf.Version)
}
