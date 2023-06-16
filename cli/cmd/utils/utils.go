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

package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-openapi/runtime"
	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/loopholelabs/auth"
	"github.com/loopholelabs/auth/pkg/client/session"
	"github.com/loopholelabs/cmdutils"
	"github.com/loopholelabs/releaser/pkg/client"
	"github.com/loopholelabs/scale/cli/analytics"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/loopholelabs/scale/cli/internal/log"
	"github.com/loopholelabs/scale/cli/version"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/spf13/cobra"
	"io"
	"strings"
)

const (
	DefaultOrganization = "local"
	DefaultTag          = "latest"
)

var (
	ErrNotAuthenticated = errors.New("You must be authenticated to use this command. Please run 'scale auth login' to authenticate.")
)

var _ runtime.NamedReadCloser = (*ScaleFunctionNamedReadCloser)(nil)

type ScaleFunctionNamedReadCloser struct {
	reader io.ReadCloser
	name   string
}

func NewScaleFunctionNamedReadCloser(sf *scalefunc.ScaleFunc) *ScaleFunctionNamedReadCloser {
	return &ScaleFunctionNamedReadCloser{
		reader: io.NopCloser(bytes.NewReader(sf.Encode())),
		name:   sf.Name,
	}
}

func (s *ScaleFunctionNamedReadCloser) Read(p []byte) (n int, err error) {
	return s.reader.Read(p)
}

func (s *ScaleFunctionNamedReadCloser) Close() error {
	return s.reader.Close()
}

func (s *ScaleFunctionNamedReadCloser) Name() string {
	return s.name
}

func PreRunUpdateCheck(ch *cmdutils.Helper[*config.Config]) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		log.Init(ch.Config.GetLogFile())
		err := ch.Config.GlobalRequiredFlags(cmd)
		if err != nil {
			return err
		}
		err = ch.Config.Validate()
		if err != nil {
			return err
		}

		if !ch.Config.DisableAutoUpdate {
			updateClient := client.New(fmt.Sprintf("https://%s", ch.Config.UpdateEndpoint))
			latest, err := updateClient.GetLatest()
			if err == nil {
				if latest != version.Version {
					ch.Printer.Printf("A new version of the Scale CLI is available: %s. Please run 'scale update' to update.\n\n", latest)
				}
			}
		}

		return nil
	}
}

func PreRunAuthenticatedAPI(ch *cmdutils.Helper[*config.Config]) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		log.Init(ch.Config.GetLogFile())
		err := ch.Config.GlobalRequiredFlags(cmd)
		if err != nil {
			return err
		}
		err = ch.Config.Validate()
		if err != nil {
			return err
		}

		if !ch.Config.IsAuthenticated() {
			return ErrNotAuthenticated
		}

		c, err := ch.Config.NewAuthenticatedAPIClient()
		if err != nil {
			return err
		}

		ch.Config.SetAPIClient(c)

		if !ch.Config.DisableAutoUpdate {
			updateClient := client.New(fmt.Sprintf("https://%s", ch.Config.UpdateEndpoint))
			latest, err := updateClient.GetLatest()
			if err == nil {
				if latest != version.Version {
					ch.Printer.Printf("A new version of the Scale CLI is available: %s. Please run 'scale update' to update.\n\n", latest)
				}
			}
		}

		return nil
	}
}

func PreRunOptionalAuthenticatedAPI(ch *cmdutils.Helper[*config.Config]) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		log.Init(ch.Config.GetLogFile())
		err := ch.Config.GlobalRequiredFlags(cmd)
		if err != nil {
			return err
		}

		err = ch.Config.Validate()
		if err != nil {
			return err
		}

		if ch.Config.IsAuthenticated() {
			c, err := ch.Config.NewAuthenticatedAPIClient()
			if err == nil {
				ch.Config.SetAPIClient(c)
			}
		}

		if !ch.Config.DisableAutoUpdate {
			updateClient := client.New(fmt.Sprintf("https://%s", ch.Config.UpdateEndpoint))
			latest, err := updateClient.GetLatest()
			if err == nil {
				if latest != version.Version {
					ch.Printer.Printf("A new version of the Scale CLI is available: %s. Please run 'scale update' to update.\n\n", latest)
				}
			}
		}

		return nil
	}
}

func PostRunAuthenticatedAPI(ch *cmdutils.Helper[*config.Config]) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		c := ch.Config.APIClient()
		if c != nil && c.Transport != nil {
			cookies := c.Transport.(*runtimeClient.Runtime).Jar.Cookies(config.DefaultCookieURL)
			if len(cookies) == 0 {
				return nil
			}
			ch.Config.Session = session.New(auth.KindSession, cookies[0].Value, cookies[0].Expires)

			err := ch.Config.WriteSession()
			if err != nil {
				return fmt.Errorf("error updating session: %w", err)
			}

		}
		if analytics.Client != nil {
			_ = analytics.Client.Close()
			analytics.Client = nil
		}
		return nil
	}
}

func PostRunAnalytics(_ *cmdutils.Helper[*config.Config]) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if analytics.Client != nil {
			_ = analytics.Client.Close()
			analytics.Client = nil
		}
		return nil
	}
}

type ParsedFunction struct {
	Organization string
	Name         string
	Tag          string
}

func ParseFunction(fn string) *ParsedFunction {
	orgSplit := strings.Split(fn, "/")
	if len(orgSplit) == 1 {
		orgSplit = []string{"", fn}
	}
	tagSplit := strings.Split(orgSplit[1], ":")
	if len(tagSplit) == 1 {
		tagSplit = []string{tagSplit[0], ""}
	}
	return &ParsedFunction{
		Organization: orgSplit[0],
		Name:         tagSplit[0],
		Tag:          tagSplit[1],
	}
}

func InvalidStringError(kind string, str string) error {
	return fmt.Errorf("invalid %s '%s', %ss can only include letters, numbers, periods (`.`), and dashes (`-`)", kind, str, kind)
}
