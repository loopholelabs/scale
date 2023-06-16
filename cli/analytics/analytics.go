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

package analytics

import (
	"github.com/denisbrodbeck/machineid"
	"github.com/posthog/posthog-go"
)

var (
	// APIKey is the PostHog API Key
	APIKey = ""

	// APIHost is the PostHog API Host
	APIHost = ""

	// AppID is the Loophole App ID
	AppID = ""
)

var (
	Client    posthog.Client
	MachineID string
)

var _ posthog.Logger = (*noopLogger)(nil)

type noopLogger struct{}

func (n *noopLogger) Logf(_ string, _ ...interface{})  {}
func (n noopLogger) Errorf(_ string, _ ...interface{}) {}

func init() {
	if APIKey == "" || APIHost == "" {
		return
	}
	var err error
	Client, err = posthog.NewWithConfig(APIKey, posthog.Config{
		Endpoint:  APIHost,
		BatchSize: 1,
		Logger:    new(noopLogger),
	})
	if err != nil {
		panic(err)
	}

	MachineID, err = machineid.ProtectedID(AppID)
	if err != nil {
		panic(err)
	}
}
