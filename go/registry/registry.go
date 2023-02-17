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
	"errors"
	"fmt"
	openapiClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/loopholelabs/auth"
	openapi "github.com/loopholelabs/auth/pkg/client/openapi"
	"github.com/loopholelabs/auth/pkg/client/session"
	"github.com/loopholelabs/scale/go/client"
	"github.com/loopholelabs/scale/go/client/models"
	"github.com/loopholelabs/scale/go/client/registry"
	"github.com/loopholelabs/scale/go/storage"
	"github.com/loopholelabs/scalefile/scalefunc"
	"io"
	"net/http"
	"net/url"
)

const (
	DefaultOrganization = "scale"
)

var (
	DefaultCookieURL = &url.URL{
		Scheme: "https",
		Host:   "scale.sh",
	}
)

type PullPolicy string

const (
	AlwaysPullPolicy       PullPolicy = "always"
	IfNotPresentPullPolicy            = "if-not-present"
	NeverPullPolicy                   = "never"
)

var (
	ErrHashMismatch   = errors.New("hash mismatch")
	ErrNoFunction     = errors.New("function does not exist locally and pull policy does not allow pulling from registry")
	ErrDownloadFailed = errors.New("scale function could not be pull from the registry")
)

type config struct {
	pullPolicy     PullPolicy
	cacheDirectory string
	apiKey         string
	cookieURL      *url.URL
	baseURL        string
	organization   string
}

type Option func(config *config)

func WithPullPolicy(pullPolicy PullPolicy) Option {
	return func(config *config) {
		config.pullPolicy = pullPolicy
	}
}

func WithCacheDirectory(cacheDirectory string) Option {
	return func(config *config) {
		config.cacheDirectory = cacheDirectory
	}
}

func WithAPIKey(apiKey string) Option {
	return func(config *config) {
		config.apiKey = apiKey
	}
}

func WithBaseURL(baseURL string) Option {
	return func(config *config) {
		config.baseURL = baseURL
	}
}

func WithOrganization(organization string) Option {
	return func(config *config) {
		config.organization = organization
	}
}

func WithCookieURL(cookieURL *url.URL) Option {
	return func(config *config) {
		config.cookieURL = cookieURL
	}
}

func New(name string, tag string, opts ...Option) (*scalefunc.ScaleFunc, error) {
	conf := &config{
		pullPolicy:   IfNotPresentPullPolicy,
		baseURL:      client.DefaultHost,
		cookieURL:    DefaultCookieURL,
		organization: DefaultOrganization,
	}
	for _, opt := range opts {
		opt(conf)
	}

	if conf.baseURL == "" {
		conf.baseURL = client.DefaultHost
	}

	if conf.cookieURL == nil {
		conf.cookieURL = DefaultCookieURL
	}

	if conf.organization == "" {
		conf.organization = DefaultOrganization
	}

	var err error
	st := storage.Default
	if conf.cacheDirectory != "" {
		st, err = storage.New(conf.cacheDirectory)
		if err != nil {
			return nil, fmt.Errorf("failed to create storage for directory %s: %w", conf.cacheDirectory, err)
		}
	}

	var o *openapiClient.Runtime
	if conf.apiKey != "" {
		o, err = openapi.AuthenticatedClient(conf.cookieURL, conf.baseURL, client.DefaultBasePath, client.DefaultSchemes, nil, &session.Session{
			Kind:  auth.KindAPIKey,
			Value: conf.apiKey,
		})
		if err != nil {
			return nil, err
		}
	} else {
		o = openapi.UnauthenticatedClient(conf.baseURL, client.DefaultBasePath, client.DefaultSchemes, nil)
	}

	c := client.New(o, strfmt.Default)

	switch conf.pullPolicy {
	case NeverPullPolicy:
		entry, err := st.Get(name, tag, conf.organization, "")
		if err != nil {
			return nil, err
		}
		if entry != nil {
			return entry.ScaleFunc, nil
		}
		return nil, ErrNoFunction
	case IfNotPresentPullPolicy:
		entry, err := st.Get(name, tag, conf.organization, "")
		if err != nil {
			return nil, err
		}
		if entry != nil {
			return entry.ScaleFunc, nil
		}
		var fn *models.ModelsGetFunctionResponse
		if conf.organization == DefaultOrganization {
			res, err := c.Registry.GetRegistryFunctionNameTag(registry.NewGetRegistryFunctionNameTagParams().WithName(name).WithTag(tag))
			if err != nil {
				return nil, fmt.Errorf("failed to get function %s/%s:%s from registry %s: %w", conf.organization, name, tag, conf.baseURL, err)
			}
			fn = res.GetPayload()
		} else {
			res, err := c.Registry.GetRegistryFunctionOrganizationNameTag(registry.NewGetRegistryFunctionOrganizationNameTagParams().WithName(name).WithTag(tag).WithOrganization(conf.organization))
			if err != nil {
				return nil, fmt.Errorf("failed to get function %s/%s:%s from registry %s: %w", conf.organization, name, tag, conf.baseURL, err)
			}
			fn = res.GetPayload()
		}
		res, err := http.Get(fn.PresignedURL)
		if err != nil {
			return nil, err
		}
		defer func() {
			_ = res.Body.Close()
		}()

		if res.StatusCode != 200 {
			return nil, ErrDownloadFailed
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		h := sha256.New()
		h.Write(data)
		bs := h.Sum(nil)
		hash := base64.URLEncoding.EncodeToString(bs)
		if hash != fn.Hash {
			return nil, ErrHashMismatch
		}

		sf := new(scalefunc.ScaleFunc)
		err = sf.Decode(data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode retrieved scale function %s/%s:%s: %w", conf.organization, name, tag, err)
		}

		err = st.Put(name, tag, conf.organization, hash, sf)
		if err != nil {
			return nil, fmt.Errorf("failed to store retrieved scale function %s/%s:%s: %w", conf.organization, name, tag, err)
		}

		return sf, nil
	case AlwaysPullPolicy:
		entry, err := st.Get(name, tag, conf.organization, "")
		if err != nil {
			return nil, err
		}
		var fn *models.ModelsGetFunctionResponse
		if conf.organization == DefaultOrganization {
			res, err := c.Registry.GetRegistryFunctionNameTag(registry.NewGetRegistryFunctionNameTagParams().WithName(name).WithTag(tag))
			if err != nil {
				return nil, fmt.Errorf("failed to get function %s/%s:%s from registry %s: %w", conf.organization, name, tag, conf.baseURL, err)
			}
			fn = res.GetPayload()
		} else {
			res, err := c.Registry.GetRegistryFunctionOrganizationNameTag(registry.NewGetRegistryFunctionOrganizationNameTagParams().WithName(name).WithTag(tag).WithOrganization(conf.organization))
			if err != nil {
				return nil, fmt.Errorf("failed to get function %s/%s:%s from registry %s: %w", conf.organization, name, tag, conf.baseURL, err)
			}
			fn = res.GetPayload()
		}
		if entry != nil {
			if fn.Hash == entry.Hash {
				return entry.ScaleFunc, nil
			}
		}

		res, err := http.Get(fn.PresignedURL)
		if err != nil {
			return nil, err
		}
		defer func() {
			_ = res.Body.Close()
		}()

		if res.StatusCode != 200 {
			return nil, ErrDownloadFailed
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		h := sha256.New()
		h.Write(data)
		bs := h.Sum(nil)
		computedHash := base64.URLEncoding.EncodeToString(bs)
		if computedHash != fn.Hash {
			return nil, ErrHashMismatch
		}

		sf := new(scalefunc.ScaleFunc)
		err = sf.Decode(data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode retrieved scale function %s/%s:%s: %w", conf.organization, name, tag, err)
		}

		if entry != nil {
			err = st.Delete(name, tag, entry.Organization, entry.Hash)
			if err != nil {
				return nil, fmt.Errorf("failed to delete existing scale function %s/%s:%s: %w", entry.Organization, name, tag, err)
			}
		}

		err = st.Put(name, tag, conf.organization, computedHash, sf)
		if err != nil {
			return nil, fmt.Errorf("failed to store retrieved scale function %s/%s:%s: %w", conf.organization, name, tag, err)
		}

		return sf, nil
	default:
		return nil, fmt.Errorf("unknown pull policy %s", conf.pullPolicy)
	}
}
