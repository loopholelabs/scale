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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/loopholelabs/scalefile/scalefunc"
	"github.com/mitchellh/go-homedir"
	"io"
	"net/http"
	"os"
	"strings"
)

type PullPolicy string

const (
	AlwaysPullPolicy       PullPolicy = "always"
	IfNotPresentPullPolicy            = "if-not-present"
	NeverPullPolicy                   = "never"
)

const defaultCacheSubpath = ".cache/scale/functions"
const defaultAPIBaseURL = "https://api.scale.sh/v1"

var (
	ErrHashMismatch   = errors.New("hash mismatch")
	ErrNoFunction     = errors.New("function does not exist and pull policy is never")
	ErrDownloadFailed = errors.New("the scale func could not be retrieved from the server")
)

type Config struct {
	pullPolicy     PullPolicy
	cacheDirectory string
	apiKey         string
	apiBaseURL     string
	organization   string
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

func WithAPIKey(apiKey string) Option {
	return func(config *Config) {
		config.apiKey = apiKey
	}
}

func WithBaseURL(baseURL string) Option {
	return func(config *Config) {
		config.apiBaseURL = baseURL
	}
}

func WithOrganization(organization string) Option {
	return func(config *Config) {
		config.organization = organization
	}
}

// Create a new runtime for a specific scalefile
func New(function string, tag string, opts ...Option) (*scalefunc.ScaleFunc, error) {
	// Default config
	defaultPath, err := homedir.Expand(defaultCacheSubpath)
	if err != nil {
		return nil, err
	}
	conf := &Config{
		pullPolicy:     AlwaysPullPolicy,
		cacheDirectory: defaultPath,
		apiKey:         "",
		apiBaseURL:     defaultAPIBaseURL,
		organization:   "default",
	}
	for _, opt := range opts {
		opt(conf)
	}

	var sf *scalefunc.ScaleFunc

	// First check our local cache...
	if conf.pullPolicy != AlwaysPullPolicy {
		sf, err := getFromCache(function, tag, "", conf)

		if err == nil && conf.pullPolicy != AlwaysPullPolicy {
			return sf, err
		}

		if conf.pullPolicy == NeverPullPolicy {
			return nil, ErrNoFunction
		}
	}

	// Contact the API endpoint with the request
	response, err := apiRequest(function, tag, conf)
	if err != nil {
		return nil, err
	}

	sf, err = getFromCache(function, tag, response.Hash, conf)
	if err == nil {
		return sf, err
	}

	// Get the scalefunc from the URL
	httpResp, err := http.Get(response.PresignedURL)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode != 200 {
		return nil, ErrDownloadFailed
	}

	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	h := sha256.New()
	h.Write(data)

	bs := h.Sum(nil)

	hash := base64.URLEncoding.EncodeToString(bs)

	if hash != response.Hash {
		return nil, ErrHashMismatch
	}

	// Save to our local cache
	err = saveToCache(function, tag, hash, conf, data)
	if err != nil {
		return nil, err
	}

	// Decode to a scalefile
	err = sf.Decode(data)
	if err != nil {
		return nil, err
	}
	return sf, nil
}

// build a filename from the config
func buildFilename(function string, tag string, hash string, conf *Config) string {
	return fmt.Sprintf("%s.%s.%s.scale", function, tag, hash)
}

// Get a scalefile from the local cache
func getFromCache(function string, tag string, hash string, conf *Config) (*scalefunc.ScaleFunc, error) {
	if hash != "" {
		f := buildFilename(function, tag, hash, conf)

		// Try to read the scalefile
		path := fmt.Sprintf("%s%c%s", conf.cacheDirectory, os.PathSeparator, f)
		return scalefunc.Read(path)
	}

	filePrefix := fmt.Sprintf("%s.%s.", function, tag)
	files, err := os.ReadDir(conf.cacheDirectory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if strings.HasPrefix(name, filePrefix) {
			path := fmt.Sprintf("%s%c%s", conf.cacheDirectory, os.PathSeparator, name)
			return scalefunc.Read(path)
		}
	}
	return nil, ErrNoFunction
}

// Save a scalefile to our local cache
func saveToCache(function string, tag string, hash string, conf *Config, data []byte) error {
	err := os.MkdirAll(conf.cacheDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	// Overwrite the file
	f := buildFilename(function, tag, hash, conf)
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

// PulldownResponse API response when requesting a function
type GetFunctionResponse struct {
	Name         string `json:"name"`
	Tag          string `json:"tag"`
	Organization string `json:"organization"`
	Public       bool   `json:"public"`
	Hash         string `json:"hash"`
	PresignedURL string `json:"presigned_url"`
}

func apiRequest(function string, tag string, conf *Config) (*GetFunctionResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/registry/function/%s/%s/%s", conf.apiBaseURL, conf.organization, function, tag),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", conf.apiKey))
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		bodyString := string(bodyBytes)
		return nil, errors.New(bodyString)
	}

	response := &GetFunctionResponse{}
	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
