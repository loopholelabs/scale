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

package apikey

import (
	"fmt"
	"github.com/loopholelabs/cmdutils"
	"github.com/loopholelabs/cmdutils/pkg/command"
	"github.com/loopholelabs/scale/cli/analytics"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/loopholelabs/scale/client/access"
	"github.com/posthog/posthog-go"
	"github.com/spf13/cobra"
	"time"
)

// GetCmd encapsulates the commands for getting API Keys
func GetCmd() command.SetupCommand[*config.Config] {
	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
		getCmd := &cobra.Command{
			Use:   "get <name>",
			Args:  cobra.ExactArgs(1),
			Short: "get information about an API Key with the given name",
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()
				client := ch.Config.APIClient()
				name := args[0]

				end := ch.Printer.PrintProgress(fmt.Sprintf("Retrieving API Key %s...", name))
				res, err := client.Access.GetAccessApikeyNameorid(access.NewGetAccessApikeyNameoridParamsWithContext(ctx).WithNameorid(name))
				end()
				if err != nil {
					return err
				}

				if analytics.Client != nil {
					_ = analytics.Client.Enqueue(posthog.Capture{
						DistinctId: analytics.MachineID,
						Event:      "get-apikey",
						Timestamp:  time.Now(),
					})
				}

				return ch.Printer.PrintResource(apiKeyRedacted{
					Name:    res.GetPayload().Name,
					Created: res.GetPayload().CreatedAt,
					ID:      res.GetPayload().ID,
				})
			},
		}

		cmd.AddCommand(getCmd)
	}
}
