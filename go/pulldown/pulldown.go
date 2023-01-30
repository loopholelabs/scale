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

type PullPolicy string

const (
	AlwaysPullPolicy       PullPolicy = "always"
	IfNotPresentPullPolicy            = "if-not-present"
	NeverPullPolicy                   = "never"
)

var (
	ErrHashMismatch = errors.New("hash mismatch")
	ErrNoFunction   = errors.New("function does not exist and pull policy is never")
)

type Config struct {
	pullPolicy     PullPolicy
	cacheDirectory string
	apiKey         string
}

type Option func(config *Config)

func WithPullPolicy(pullPolicy PullPolicy) Option {
	return func(config *Config) {
		config.pullPolicy = pullPolicy
	}
}

func WithCacheDirectory(cacheDirectory string) Option {
	return func(config *Config) {
		config.cacheDirectory = cacheDirectory
	}
}

func WithApiKey(apiKey string) Option {
	return func(config *Config) {
		config.apiKey = apiKey
	}
}

// Create a new runtime for a specific scalefile
func New(function string, opts ...Option) (*scalefile.ScaleFile, error) {
	// Default config
	conf := &Config{
		pullPolicy:     AlwaysPullPolicy,
		cacheDirectory: "~/.cache/scale/functions",
		apiKey:         "",
	}
	for _, opt := range opts {
		opt(conf)
	}

	// First check our local cache...
	sf, err := getFromCache(function, conf)

	if err == nil && conf.pullPolicy != AlwaysPullPolicy {
		return sf, err
	}

	if conf.pullPolicy == NeverPullPolicy {
		return nil, ErrNoFunction
	}

	// Contact the API endpoint with the request
	response, err := apiRequest(function, conf)
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
		return nil, ErrHashMismatch
	}

	// Save to our local cache
	err = saveToCache(function, conf, data)
	if err != nil {
		return nil, err
	}

	// Decode to a scalefile
	reader := bytes.NewReader(data)
	return scalefile.Decode(reader)
}

// build a filename from the config
// TODO: We should use the hash in the filename for optimization
func buildFilename(function string, conf *Config) string {
	return fmt.Sprintf("%s.scale", function)
}

// Get a scalefile from the local cache
func getFromCache(function string, conf *Config) (*scalefile.ScaleFile, error) {
	f := buildFilename(function, conf)

	// Try to read the scalefile
	path := fmt.Sprintf("%s%c%s", conf.cacheDirectory, os.PathSeparator, f)
	return scalefile.Read(path)
}

// Save a scalefile to our local cache
func saveToCache(function string, conf *Config, data []byte) error {
	err := os.MkdirAll(conf.cacheDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	// Overwrite the file
	f := buildFilename(function, conf)
	path := fmt.Sprintf("%s%c%s", conf.cacheDirectory, os.PathSeparator, f)

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

func removeCache(conf *Config) error {
	return os.RemoveAll(conf.cacheDirectory)
}

// What we get back from the API call
type PulldownResponse struct {
	URL  string
	Hash string
}

// Perform the scale api request to find the correct URL and hash
func apiRequest(function string, conf *Config) (PulldownResponse, error) {
	// TODO

	return PulldownResponse{
		URL:  "http://google.com",
		Hash: "1234",
	}, nil
}
