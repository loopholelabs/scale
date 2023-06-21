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

package registry

import (
	"fmt"
	"github.com/loopholelabs/cmdutils"
	"github.com/loopholelabs/cmdutils/pkg/command"
	"github.com/loopholelabs/cmdutils/pkg/printer"
	"github.com/loopholelabs/scale/cli/analytics"
	"github.com/loopholelabs/scale/cli/cmd/utils"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/loopholelabs/scale/client/registry"
	"github.com/loopholelabs/scale/client/userinfo"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/storage"
	"github.com/posthog/posthog-go"
	"github.com/spf13/cobra"
	"time"
)

// PushCmd encapsulates the commands for pushing Functions
func PushCmd(hidden bool) command.SetupCommand[*config.Config] {
	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
		var name string
		var tag string
		var org string
		var public bool
		pushCmd := &cobra.Command{
			Use:      "push [<name>:<tag> | [<org>/<name>:<tag>]",
			Short:    "push a locally available scale function to the registry",
			Long:     "Push a locally available scale function to the registry. The function must be available in the local storage directory. If no storage directory is specified, the default storage directory will be used. If the org is not specified, it will default to the user's default organization.",
			Hidden:   hidden,
			Args:     cobra.ExactArgs(1),
			PreRunE:  utils.PreRunAuthenticatedAPI(ch),
			PostRunE: utils.PostRunAuthenticatedAPI(ch),
			RunE: func(cmd *cobra.Command, args []string) error {
				st := storage.DefaultFunction
				if ch.Config.StorageDirectory != "" {
					var err error
					st, err = storage.NewFunction(ch.Config.StorageDirectory)
					if err != nil {
						return fmt.Errorf("failed to instantiate function storage for %s: %w", ch.Config.StorageDirectory, err)
					}
				}

				parsed := utils.ParseFunction(args[0])
				if parsed.Organization == "" {
					parsed.Organization = utils.DefaultOrganization
				}

				if parsed.Organization != "" && !scalefunc.ValidString(parsed.Organization) {
					return utils.InvalidStringError("organization name", parsed.Organization)
				}

				if parsed.Name == "" || !scalefunc.ValidString(parsed.Name) {
					return utils.InvalidStringError("function name", parsed.Name)
				}

				if parsed.Tag == "" || !scalefunc.ValidString(parsed.Tag) {
					return utils.InvalidStringError("function tag", parsed.Tag)
				}

				e, err := st.Get(parsed.Name, parsed.Tag, parsed.Organization, "")
				if err != nil {
					return fmt.Errorf("failed to push function %s/%s:%s: %w", parsed.Organization, parsed.Name, parsed.Tag, err)
				}
				if e == nil {
					return fmt.Errorf("function %s/%s:%s does not exist", parsed.Organization, parsed.Name, parsed.Tag)
				}

				if org != "" {
					parsed.Organization = org
				}

				if name != "" {
					parsed.Name = name
				}

				if tag != "" {
					parsed.Tag = tag
				}

				if parsed.Organization != "" && !scalefunc.ValidString(parsed.Organization) {
					return utils.InvalidStringError("organization name", parsed.Organization)
				}

				if parsed.Name == "" || !scalefunc.ValidString(parsed.Name) {
					return utils.InvalidStringError("function name", parsed.Name)
				}

				if parsed.Tag == "" || !scalefunc.ValidString(parsed.Tag) {
					return utils.InvalidStringError("function tag", parsed.Tag)
				}

				e.Schema.Name = parsed.Name
				e.Schema.Tag = parsed.Tag

				ctx := cmd.Context()
				client := ch.Config.APIClient()

				var end func()
				if parsed.Organization == utils.DefaultOrganization {
					userInfoRes, err := client.Userinfo.PostUserinfo(userinfo.NewPostUserinfoParamsWithContext(ctx))
					if err != nil {
						return err
					}
					ch.Printer.Printf("No organization specified, using user's default organization %s\n", printer.BoldGreen(userInfoRes.GetPayload().Organization))
					end = ch.Printer.PrintProgress(fmt.Sprintf("Pushing %s/%s:%s to Scale Registry...", userInfoRes.GetPayload().Organization, parsed.Name, parsed.Tag))
				} else {
					end = ch.Printer.PrintProgress(fmt.Sprintf("Pushing %s/%s:%s to Scale Registry...", parsed.Organization, parsed.Name, parsed.Tag))
				}

				params := registry.NewPostRegistryFunctionParamsWithContext(ctx).WithFunction(utils.NewScaleFunctionNamedReadCloser(e.Schema)).WithPublic(&public)
				if parsed.Organization != utils.DefaultOrganization {
					params = params.WithOrganization(&parsed.Organization)
				}

				res, err := client.Registry.PostRegistryFunction(params)
				end()
				if err != nil {
					return err
				}

				if analytics.Client != nil {
					_ = analytics.Client.Enqueue(posthog.Capture{
						DistinctId: analytics.MachineID,
						Event:      "push-registry",
						Timestamp:  time.Now(),
					})
				}

				if ch.Printer.Format() == printer.Human {
					ch.Printer.Printf("Pushed %s to the Scale Registry\n", printer.BoldGreen(fmt.Sprintf("%s/%s:%s", res.GetPayload().Organization, res.GetPayload().Name, res.GetPayload().Tag)))
					return nil
				}

				return ch.Printer.PrintResource(map[string]string{
					"name":   res.GetPayload().Name,
					"tag":    res.GetPayload().Tag,
					"org":    res.GetPayload().Organization,
					"public": fmt.Sprintf("%t", res.GetPayload().Public),
					"hash":   res.GetPayload().Hash,
				})
			},
		}

		pushCmd.Flags().StringVarP(&name, "name", "n", "", "the name of the function")
		pushCmd.Flags().StringVarP(&tag, "tag", "t", "", "the tag of the function")
		pushCmd.Flags().StringVarP(&org, "org", "o", "", "the organization of the function")
		pushCmd.Flags().BoolVarP(&public, "public", "p", false, "make the function public")

		cmd.AddCommand(pushCmd)
	}
}
