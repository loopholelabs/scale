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

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-openapi/strfmt"
	authClient "github.com/loopholelabs/auth/pkg/client"
	client "github.com/loopholelabs/auth/pkg/client/openapi"
	"github.com/loopholelabs/auth/pkg/client/session"
	"github.com/loopholelabs/auth/pkg/kind"
	"github.com/loopholelabs/cmdutils/pkg/config"
	apiClient "github.com/loopholelabs/scale/cli/client"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/url"
	"os"
	"path"
	"time"
)

var _ config.Config = (*Config)(nil)

var (
	ErrAPIEndpointRequired    = errors.New("api endpoint is required")
	ErrAuthEndpointRequired   = errors.New("auth endpoint is required")
	ErrSessionDomainRequired  = errors.New("session domain is required")
	ErrUpdateEndpointRequired = errors.New("update endpoint is required")
	ErrNoSession              = errors.New("no session found")
)

var (
	configFile string
	logFile    string
)

const (
	defaultConfigPath = "~/.config/scale"
	configName        = "scale.yml"
	logName           = "scale.log"

	DefaultAPIEndpoint    = "api.scale.sh"
	DefaultAuthEndpoint   = "auth.scale.sh"
	DefaultSessionDomain  = "scale.sh"
	DefaultUpdateEndpoint = "dl.scale.sh"

	sessionFileMode = 0600
)

// Config is dynamically sourced from various files and environment variables.
type Config struct {
	APIEndpoint       string                `mapstructure:"api_endpoint"`
	AuthEndpoint      string                `mapstructure:"auth_endpoint"`
	SessionDomain     string                `mapstructure:"session_domain"`
	UpdateEndpoint    string                `mapstructure:"update_endpoint"`
	DisableAutoUpdate bool                  `mapstructure:"disable_auto_update"`
	NoTelemetry       bool                  `mapstructure:"no_telemetry"`
	StorageDirectory  string                `mapstructure:"storage_directory"`
	Session           *session.Session      `mapstructure:"-"`
	authClient        *authClient.AuthAPIV1 `mapstructure:"-"`
	apiClient         *apiClient.ScaleAPIV1 `mapstructure:"-"`
	sessionCookieURL  *url.URL              `mapstructure:"-"`
}

func New() *Config {
	return &Config{
		APIEndpoint:    DefaultAPIEndpoint,
		AuthEndpoint:   DefaultAuthEndpoint,
		SessionDomain:  DefaultSessionDomain,
		UpdateEndpoint: DefaultUpdateEndpoint,
	}
}

func (c *Config) RootPersistentFlags(flags *pflag.FlagSet) {
	flags.StringVar(&c.APIEndpoint, "api-endpoint", DefaultAPIEndpoint, "The Scale API endpoint")
	flags.StringVar(&c.AuthEndpoint, "auth-endpoint", DefaultAuthEndpoint, "The Scale Authentication API endpoint")
	flags.StringVar(&c.SessionDomain, "session-domain", DefaultSessionDomain, "The Scale API session domain")
	flags.StringVar(&c.UpdateEndpoint, "update-endpoint", DefaultUpdateEndpoint, "The Scale Update API endpoint")
	flags.BoolVar(&c.DisableAutoUpdate, "disable-auto-update", false, "Disable automatic update checks")
	flags.BoolVar(&c.NoTelemetry, "no-telemetry", false, "Opt out of telemetry tracking")
	flags.StringVar(&c.StorageDirectory, "storage-directory", "", "The (optional) directory to store compiled scale functions and generated signatures")
}

func (c *Config) GlobalRequiredFlags(_ *cobra.Command) error {
	return nil
}

func (c *Config) Validate() error {
	err := viper.Unmarshal(c)
	if err != nil {
		return fmt.Errorf("unable to unmarshal config: %w", err)
	}

	if c.APIEndpoint == "" {
		return ErrAPIEndpointRequired
	}

	if c.AuthEndpoint == "" {
		return ErrAuthEndpointRequired
	}

	if c.SessionDomain == "" {
		return ErrSessionDomainRequired
	}

	if c.UpdateEndpoint == "" {
		return ErrUpdateEndpointRequired
	}

	c.sessionCookieURL = &url.URL{
		Scheme: "https",
		Host:   c.SessionDomain,
	}

	sessionPath, err := c.SessionPath()
	if err != nil {
		return fmt.Errorf("unable to get session path: %w", err)
	}

	if !c.IsAuthenticated() {
		stat, err := os.Stat(sessionPath)
		if err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("unable to stat %s: %w", sessionPath, err)
			}
		} else {
			if stat.Mode()&^sessionFileMode != 0 {
				err = os.Chmod(sessionPath, sessionFileMode)
				if err != nil {
					return fmt.Errorf("unable to chmod %s: %w", sessionPath, err)
				}
			}
			sessionData, err := os.ReadFile(sessionPath)
			if err != nil {
				return fmt.Errorf("unable to read %s: %w", sessionPath, err)
			}

			c.Session = new(session.Session)
			err = json.Unmarshal(sessionData, c.Session)
			if err != nil {
				return fmt.Errorf("unable to unmarshal %s: %w", sessionPath, err)
			}
		}
	}

	return nil
}

func (c *Config) DefaultConfigDir() (string, error) {
	dir, err := homedir.Expand(defaultConfigPath)
	if err != nil {
		return "", fmt.Errorf("can't expand path %q: %s", defaultConfigPath, err)
	}

	return dir, nil
}

func (c *Config) DefaultConfigFile() string {
	return configName
}

func (c *Config) DefaultLogFile() string {
	return logName
}

func (c *Config) DefaultConfigPath() (string, error) {
	configDir, err := c.DefaultConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(configDir, c.DefaultConfigFile()), nil
}

func (c *Config) DefaultLogPath() (string, error) {
	configDir, err := c.DefaultConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(configDir, c.DefaultLogFile()), nil
}

func (c *Config) GetConfigFile() string {
	return configFile
}

func (c *Config) GetLogFile() string {
	return logFile
}

func (c *Config) SetLogFile(file string) {
	logFile = file
}

func (c *Config) SetConfigFile(file string) {
	configFile = file
}

func (c *Config) IsAuthenticated() bool {
	if c.Session == nil || c.Session.Value == "" {
		return false
	}
	switch c.Session.Kind {
	case kind.Session:
		if c.Session.Expiry.After(time.Now()) {
			return true
		}
		return false
	case kind.ServiceSession, kind.APIKey:
		return true
	default:
		return false
	}
}

// SessionPath is the path for the session file
func (c *Config) SessionPath() (string, error) {
	dir, err := c.DefaultConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(dir, "session"), nil
}

// NewAuthenticatedAPIClient creates an Authenticated Scale API client from our configuration
func (c *Config) NewAuthenticatedAPIClient() (*apiClient.ScaleAPIV1, error) {
	if !c.IsAuthenticated() {
		return nil, ErrNoSession
	}

	cl, err := client.AuthenticatedClient(c.SessionCookieURL(), c.APIEndpoint, apiClient.DefaultBasePath, apiClient.DefaultSchemes, nil, c.Session)
	if err != nil {
		return nil, err
	}

	return apiClient.New(cl, strfmt.Default), nil
}

// NewUnauthenticatedAPIClient creates an Unauthenticated Scale API client from our configuration
func (c *Config) NewUnauthenticatedAPIClient() *apiClient.ScaleAPIV1 {
	cl := client.UnauthenticatedClient(c.APIEndpoint, apiClient.DefaultBasePath, apiClient.DefaultSchemes, nil)
	return apiClient.New(cl, strfmt.Default)
}

func (c *Config) SetAPIClient(apiClient *apiClient.ScaleAPIV1) {
	c.apiClient = apiClient
}

func (c *Config) APIClient() *apiClient.ScaleAPIV1 {
	return c.apiClient
}

// NewUnauthenticatedAuthClient creates an Unauthenticated Scale Auth API client from our configuration
func (c *Config) NewUnauthenticatedAuthClient() *authClient.AuthAPIV1 {
	cl := client.UnauthenticatedClient(c.AuthEndpoint, authClient.DefaultBasePath, authClient.DefaultSchemes, nil)
	return authClient.New(cl, strfmt.Default)
}

// NewAuthenticatedAuthClient creates an Authenticated Scale Auth API client from our configuration
func (c *Config) NewAuthenticatedAuthClient() (*authClient.AuthAPIV1, error) {
	if !c.IsAuthenticated() {
		return nil, ErrNoSession
	}
	cl, err := client.AuthenticatedClient(c.SessionCookieURL(), c.AuthEndpoint, authClient.DefaultBasePath, authClient.DefaultSchemes, nil, c.Session)
	if err != nil {
		return nil, err
	}

	return authClient.New(cl, strfmt.Default), nil
}

func (c *Config) SetAuthClient(authClient *authClient.AuthAPIV1) {
	c.authClient = authClient
}

func (c *Config) AuthClient() *authClient.AuthAPIV1 {
	return c.authClient
}

func (c *Config) SessionCookieURL() *url.URL {
	return c.sessionCookieURL
}

func (c *Config) WriteSession() error {
	sessionPath, err := c.SessionPath()
	if err != nil {
		return err
	}

	_, err = os.Stat(sessionPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path.Dir(sessionPath), 0771)
		if err != nil {
			return fmt.Errorf("unable to create directory %s: %w", path.Dir(sessionPath), err)
		}
	} else if err != nil {
		return fmt.Errorf("unable to stat %s: %w", sessionPath, err)
	}

	sessionData, err := json.Marshal(c.Session)
	if err != nil {
		return err
	}

	err = os.WriteFile(sessionPath, sessionData, sessionFileMode)
	if err != nil {
		return fmt.Errorf("unable to write %s: %w", sessionPath, err)
	}

	return nil
}
