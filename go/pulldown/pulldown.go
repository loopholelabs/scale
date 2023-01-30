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
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/loopholelabs/scalefile"
)

const POLICY_PULL_ALWAYS = "always"
const POLICY_PULL_IF_NOT_PRESENT = "ifNotPresent"
const POLICY_PULL_NEVER = "never"

var ERROR_HASH_MISMATCH = errors.New("Hash mismatch")
var ERROR_NO_FUNCTION = errors.New("Function does not exist and pull policy was never")

// Specifies a pulldown request
type PulldownConfig struct {
	APIKey       string
	Organization string
	Function     string
	Tag          string
}

// StoreConfig
type StoreConfig struct {
	CacheDirectory string
	PullPolicy     string
}

// Default store config
var DefaultStoreConfig = StoreConfig{
	CacheDirectory: "~/.cache/scale/functions",
	PullPolicy:     "allways",
}

// What we get back from the API call
type PulldownResponse struct {
	URL  string
	Hash string
}

// Create a new runtime for a specific scalefile
func New(pc PulldownConfig, sc StoreConfig) (*scalefile.ScaleFile, error) {
	// First check our local cache...

	sf, err := getFromCache(pc, sc)

	if err == nil && sc.PullPolicy != POLICY_PULL_ALWAYS {
		return sf, err
	}

	if sc.PullPolicy == POLICY_PULL_NEVER {
		return nil, ERROR_NO_FUNCTION
	}

	// Contact the API endpoint with the request
	response, err := apiRequest(pc)
	if err != nil {
		return sf, err
	}

	// Get the scalefile from the URL
	httpResp, err := http.Get(response.URL)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	// TODO: Verify the hash with response.Hash
	h := sha256.New()
	h.Write(data)

	bs := h.Sum(nil)
	// Convert it to hexstring

	s := hex.EncodeToString(bs)

	if s != response.Hash {
		return nil, ERROR_HASH_MISMATCH
	}

	// Save to our local cache
	err = saveToCache(pc, sc, data)
	if err != nil {
		return nil, err
	}

	// Decode to a scalefile
	reader := bytes.NewReader(data)
	return scalefile.Decode(reader)
}

// build a filename from the config
// TODO: We should use the hash in the filename for optimization
func buildFilename(pc PulldownConfig) string {
	return fmt.Sprintf("%s-%s-%s.scale", pc.Organization, pc.Function, pc.Tag)
}

// Get a scalefile from the local cache
func getFromCache(pc PulldownConfig, sc StoreConfig) (*scalefile.ScaleFile, error) {
	f := buildFilename(pc)

	// Try to read the scalefile
	path := fmt.Sprintf("%s%c%s", sc.CacheDirectory, os.PathSeparator, f)
	return scalefile.Read(path)
}

// Save a scalefile to our local cache
func saveToCache(pc PulldownConfig, sc StoreConfig, data []byte) error {
	err := os.MkdirAll(sc.CacheDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	// Overwrite the file
	f := buildFilename(pc)
	path := fmt.Sprintf("%s%c%s", sc.CacheDirectory, os.PathSeparator, f)

	fh, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	_, err = fh.Write(data)
	if err != nil {
		return err
	}

	err = fh.Close()
	if err != nil {
		return err
	}

	return nil
}

func removeCache(sc StoreConfig) error {
	return os.RemoveAll(sc.CacheDirectory)
}

// Perform the scale api request to find the correct URL and hash
func apiRequest(pc PulldownConfig) (PulldownResponse, error) {
	// TODO

	return PulldownResponse{
		URL:  "http://google.com",
		Hash: "1234",
	}, nil
}
