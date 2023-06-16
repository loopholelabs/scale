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
	"github.com/loopholelabs/scale/registry"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/storage"
	"github.com/spf13/cobra"
	"time"
)

// PullCmd encapsulates the commands for pulling Functions
func PullCmd(hidden bool) command.SetupCommand[*config.Config] {
	var force bool
	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
		pullCmd := &cobra.Command{
			Use:      "pull [<name>:<tag> | [<org>/<name>:<tag>]",
			Short:    "pull a scale function from the registry",
			Long:     "Pull a scale function from the registry. If the org is not specified, it will default to the official `scale` organization.",
			Hidden:   hidden,
			Args:     cobra.ExactArgs(1),
			PreRunE:  utils.PreRunOptionalAuthenticatedAPI(ch),
			PostRunE: utils.PostRunAuthenticatedAPI(ch),
			RunE: func(cmd *cobra.Command, args []string) error {
				st := storage.Default
				if ch.Config.CacheDirectory != "" {
					var err error
					st, err = storage.New(ch.Config.CacheDirectory)
					if err != nil {
						return fmt.Errorf("failed to instantiate function storage for %s: %w", ch.Config.CacheDirectory, err)
					}
				}

				parsed := utils.ParseFunction(args[0])
				if parsed.Organization != "" && !scalefunc.ValidString(parsed.Organization) {
					return utils.InvalidStringError("organization name", parsed.Organization)
				}

				if parsed.Name == "" || !scalefunc.ValidString(parsed.Name) {
					return utils.InvalidStringError("function name", parsed.Name)
				}

				if parsed.Tag == "" || !scalefunc.ValidString(parsed.Tag) {
					return utils.InvalidStringError("function tag", parsed.Tag)
				}

				var end func()
				if parsed.Organization == "" {
					end = ch.Printer.PrintProgress(fmt.Sprintf("Pulling %s:%s from Scale Registry...", parsed.Name, parsed.Tag))
				} else {
					end = ch.Printer.PrintProgress(fmt.Sprintf("Pulling %s/%s:%s from Scale Registry...", parsed.Organization, parsed.Name, parsed.Tag))
				}

				var opts []registry.Option
				opts = append(opts, registry.WithClient(ch.Config.APIClient()), registry.WithStorage(st))
				if parsed.Organization != "" {
					opts = append(opts, registry.WithOrganization(parsed.Organization))
				}

				if force {
					opts = append(opts, registry.WithPullPolicy(registry.AlwaysPullPolicy))
				}

				sf, err := registry.Download(parsed.Name, parsed.Tag, opts...)
				end()
				if err != nil {
					if parsed.Organization == "" {
						return fmt.Errorf("failed to pull function %s:%s: %w", parsed.Name, parsed.Tag, err)
					} else {
						return fmt.Errorf("failed to pull function %s/%s:%s: %w", parsed.Organization, parsed.Name, parsed.Tag, err)
					}
				}

				if analytics.Client != nil {
					_ = analytics.Client.Enqueue(posthog.Capture{
						DistinctId: analytics.MachineID,
						Event:      "pull-registry",
						Timestamp:  time.Now(),
					})
				}

				if ch.Printer.Format() == printer.Human {
					if parsed.Organization == "" {
						ch.Printer.Printf("Pulled %s from the Scale Registry\n", printer.BoldGreen(fmt.Sprintf("%s:%s", sf.Name, sf.Tag)))
					} else {
						ch.Printer.Printf("Pulled %s from the Scale Registry\n", printer.BoldGreen(fmt.Sprintf("%s/%s:%s", parsed.Organization, sf.Name, sf.Tag)))
					}
					return nil
				}

				if parsed.Organization == "" {
					parsed.Organization = "scale"
				}

				return ch.Printer.PrintResource(map[string]string{
					"name": sf.Name,
					"tag":  sf.Tag,
					"org":  parsed.Organization,
				})
			},
		}

		pullCmd.Flags().BoolVar(&force, "force", false, "force overwrite of existing function")

		cmd.AddCommand(pullCmd)
	}
}
