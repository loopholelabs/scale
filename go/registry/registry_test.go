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
	"encoding/base64"
	"os"
	"testing"

	"github.com/loopholelabs/scalefile/scalefunc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testingAPIBaseURl = "https://api.dev.scale.sh/v1"

func TestPulldownCache(t *testing.T) {
	conf := &Config{
		cacheDirectory: "testCache",
		pullPolicy:     NeverPullPolicy,
		apiKey:         "123",
	}

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

	hash := base64.URLEncoding.EncodeToString(bs)

	b := sf.Encode()

	// Write it to the cache...
	err := saveToCache(function, tag, hash, conf, b)
	require.NoError(t, err)

	// Try reading a scalefile from the cache...
	// Also tests that the With* works
	newsf, err := New(function, tag,
		WithCacheDirectory(conf.cacheDirectory),
		WithPullPolicy(conf.pullPolicy),
		WithAPIKey(conf.apiKey))

	require.NoError(t, err)

	// Now assert that the newsf is same as our original sf.

	newb := newsf.Encode()

	assert.Equal(t, b, newb)
	// cleanup - remove cache dir
	err = removeCache(conf)
	require.NoError(t, err)
}

func TestRegistryDownload(t *testing.T) {
	/* This test requires a valid API key for the scale dev api to run, the SCALE_API_KEY environment variable must be
	set in the testing environment. */
	apiKey := os.Getenv("SCALE_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping test, SCALE_API_KEY environment variable not set")
	}
	sf, err := New("TestRegistryDownload", "1",
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
