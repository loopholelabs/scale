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
	"github.com/loopholelabs/cmdutils"
	"github.com/loopholelabs/cmdutils/pkg/command"
	"github.com/loopholelabs/cmdutils/pkg/printer"
	"github.com/loopholelabs/scale/cli/analytics"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/loopholelabs/scale/client/access"
	"github.com/posthog/posthog-go"
	"github.com/spf13/cobra"
	"time"
)

// ListCmd encapsulates the commands for listing API Keys
func ListCmd() command.SetupCommand[*config.Config] {
	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
		listCmd := &cobra.Command{
			Use:   "list",
			Short: "list API Keys",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()
				client := ch.Config.APIClient()
				end := ch.Printer.PrintProgress("Retrieving API Keys...")
				res, err := client.Access.GetAccessApikey(access.NewGetAccessApikeyParamsWithContext(ctx))
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

				if len(res.Payload) == 0 && ch.Printer.Format() == printer.Human {
					ch.Printer.Println("No API Keys have been created yet.")
					return nil
				}

				keys := make([]apiKeyRedacted, 0, len(res.Payload))
				for _, key := range res.Payload {
					keys = append(keys, apiKeyRedacted{
						Created: key.CreatedAt,
						ID:      key.ID,
						Name:    key.Name,
					})
				}

				return ch.Printer.PrintResource(keys)
			},
		}

		cmd.AddCommand(listCmd)
	}
}