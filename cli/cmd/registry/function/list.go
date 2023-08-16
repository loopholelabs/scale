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

package function

import (
	"fmt"
	"github.com/loopholelabs/cmdutils"
	"github.com/loopholelabs/cmdutils/pkg/command"
	"github.com/loopholelabs/cmdutils/pkg/printer"
	"github.com/loopholelabs/scale/cli/analytics"
	"github.com/loopholelabs/scale/cli/client/registry"
	"github.com/loopholelabs/scale/cli/client/userinfo"
	"github.com/loopholelabs/scale/cli/cmd/utils"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/spf13/cobra"
)

// ListCmd encapsulates the commands for listing Functions
func ListCmd() command.SetupCommand[*config.Config] {
	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
		listCmd := &cobra.Command{
			Use:      "list [<org>] [flags]",
			Short:    "list all the scale function for an organization from the registry",
			Long:     "List all the scale functions available in an organization from the registry. If no org is specified, it will default to the user's organization.",
			Args:     cobra.RangeArgs(0, 1),
			PreRunE:  utils.PreRunAuthenticatedAPI(ch),
			PostRunE: utils.PostRunAuthenticatedAPI(ch),
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()
				client := ch.Config.APIClient()

				infoRes, err := client.Userinfo.PostUserinfo(userinfo.NewPostUserinfoParamsWithContext(ctx))
				if err != nil {
					return err
				}

				org := infoRes.GetPayload().Organization
				if len(args) > 0 {
					org = args[0]
				}
				if org == "" || !scalefunc.ValidString(org) {
					return utils.InvalidStringError("organization name", org)
				}

				end := ch.Printer.PrintProgress(fmt.Sprintf("Retrieving Scale Functions for the '%s' organization from the Registry...", org))
				res, err := client.Registry.GetRegistryFunctionOrg(registry.NewGetRegistryFunctionOrgParamsWithContext(ctx).WithOrg(org))
				end()
				if err != nil {
					return err
				}

				analytics.Event("list-registry")

				if len(res.GetPayload()) == 0 && ch.Printer.Format() == printer.Human {
					ch.Printer.Println("No functions available in this organization yet.")
					return nil
				}

				ret := make([]functionModel, 0, len(res.GetPayload()))
				for _, fn := range res.GetPayload() {
					ret = append(ret, functionModel{
						Name:   fn.Name,
						Tag:    fn.Tag,
						Hash:   fn.Hash,
						Org:    fn.Organization,
						Public: fmt.Sprintf("%t", fn.Public),
					})
				}

				return ch.Printer.PrintResource(ret)

			},
		}

		cmd.AddCommand(listCmd)
	}
}
