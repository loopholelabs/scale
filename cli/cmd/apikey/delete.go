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
	"github.com/loopholelabs/cmdutils/pkg/printer"
	"github.com/loopholelabs/scale/cli/analytics"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/loopholelabs/scale/client/access"
	"github.com/posthog/posthog-go"
	"github.com/spf13/cobra"
	"time"
)

// DeleteCmd encapsulates the commands for deleting API Keys
func DeleteCmd() command.SetupCommand[*config.Config] {
	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
		deleteCmd := &cobra.Command{
			Use:   "delete <id>",
			Args:  cobra.ExactArgs(1),
			Short: "delete an API Key with the given ID",
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()
				client := ch.Config.APIClient()
				id := args[0]

				end := ch.Printer.PrintProgress(fmt.Sprintf("Deleting API Key %s...", id))
				_, err := client.Access.DeleteAccessApikeyNameorid(access.NewDeleteAccessApikeyNameoridParamsWithContext(ctx).WithNameorid(id))
				end()
				if err != nil {
					return err
				}

				if analytics.Client != nil {
					_ = analytics.Client.Enqueue(posthog.Capture{
						DistinctId: analytics.MachineID,
						Event:      "delete-apikey",
						Timestamp:  time.Now(),
					})
				}

				if ch.Printer.Format() == printer.Human {
					ch.Printer.Printf("%s %s %s\n", printer.BoldRed("API Key"), printer.BoldGreen(id), printer.BoldRed("deleted"))
					return nil
				}

				return ch.Printer.PrintResource(map[string]string{
					"deleted": id,
				})
			},
		}

		cmd.AddCommand(deleteCmd)
	}
}

//
//func DeleteCmd(ch *cmdutil.Helper) *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "delete <id>",
//		Args:  cobra.ExactArgs(1),
//		Short: "delete an API Key with the given ID",
//		RunE: func(cmd *cobra.Command, args []string) error {
//			ctx := cmd.Context()
//			client, err := ch.Client()
//			if err != nil {
//				return err
//			}
//
//			id := args[0]
//
//			end := ch.Printer.PrintProgress(fmt.Sprintf("Deleting API Key %s...", id))
//			_, err = client.Access.DeleteAccessApikeyID(access.NewDeleteAccessApikeyIDParamsWithContext(ctx).WithID(id))
//			end()
//			if err != nil {
//				return err
//			}
//
//			if ch.Printer.Format() == printer.Human {
//				ch.Printer.Printf("API Key %s %s\n", printer.BoldGreen(id), printer.BoldRed("deleted"))
//				return nil
//			}
//
//			return ch.Printer.PrintResource(map[string]string{
//				"deleted": id,
//			})
//		},
//	}
//
//	return cmd
//}
