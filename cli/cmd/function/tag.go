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
	"github.com/loopholelabs/scale/cli/cmd/utils"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/storage"
	"github.com/spf13/cobra"
)

// TagCmd encapsulates the commands for tagging and renaming Functions
func TagCmd() command.SetupCommand[*config.Config] {
	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
		tagCmd := &cobra.Command{
			Use:     "tag [<current_name>:<current_tag>] | [<current_org>/<current_name>:<current_tag>] [<new_name>:<new_tag>] | [<new_org>/<new_name>:<new_tag>]",
			Short:   "tag a locally available Scale Function",
			Args:    cobra.ExactArgs(2),
			PreRunE: utils.PreRunUpdateCheck(ch),
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
					return fmt.Errorf("failed to tag function %s/%s:%s: %w", parsed.Organization, parsed.Name, parsed.Tag, err)
				}
				if e == nil {
					return fmt.Errorf("function %s/%s:%s does not exist", parsed.Organization, parsed.Name, parsed.Tag)
				}

				newParsed := utils.ParseFunction(args[1])
				if newParsed.Organization == "" {
					newParsed.Organization = utils.DefaultOrganization
				}

				if newParsed.Organization != "" && !scalefunc.ValidString(newParsed.Organization) {
					return utils.InvalidStringError("organization name", newParsed.Organization)
				}

				if newParsed.Name == "" || !scalefunc.ValidString(newParsed.Name) {
					return utils.InvalidStringError("function name", newParsed.Name)
				}

				if newParsed.Tag == "" || !scalefunc.ValidString(newParsed.Tag) {
					return utils.InvalidStringError("function tag", newParsed.Tag)
				}

				e.ScaleFunc.Name = newParsed.Name
				e.ScaleFunc.Tag = newParsed.Tag
				e.Organization = newParsed.Organization

				err = st.Put(e.ScaleFunc.Name, e.ScaleFunc.Tag, e.Organization, e.Hash, e.ScaleFunc)
				if err != nil {
					return fmt.Errorf("failed to tag function %s/%s:%s: %w", parsed.Organization, parsed.Name, parsed.Tag, err)
				}

				if ch.Printer.Format() == printer.Human {
					if parsed.Organization != "" {
						if newParsed.Organization != "" {
							ch.Printer.Printf("Tagged scale function %s to %s\n", printer.BoldGreen(fmt.Sprintf("%s/%s:%s", parsed.Organization, parsed.Name, parsed.Tag)), printer.BoldBlue(fmt.Sprintf("%s/%s:%s", newParsed.Organization, newParsed.Name, newParsed.Tag)))
						} else {
							ch.Printer.Printf("Tagged scale function %s to %s\n", printer.BoldGreen(fmt.Sprintf("%s/%s:%s", parsed.Organization, parsed.Name, parsed.Tag)), printer.BoldBlue(fmt.Sprintf("%s:%s", newParsed.Name, newParsed.Tag)))
						}
					} else {
						if newParsed.Organization != "" {
							ch.Printer.Printf("Tagged scale function %s to %s\n", printer.BoldGreen(fmt.Sprintf("%s:%s", parsed.Name, parsed.Tag)), printer.BoldBlue(fmt.Sprintf("%s/%s:%s", newParsed.Organization, newParsed.Name, newParsed.Tag)))
						} else {
							ch.Printer.Printf("Tagged scale function %s to %s\n", printer.BoldGreen(fmt.Sprintf("%s:%s", parsed.Name, parsed.Tag)), printer.BoldBlue(fmt.Sprintf("%s:%s", newParsed.Name, newParsed.Tag)))
						}
					}
					return nil
				}

				return ch.Printer.PrintResource(map[string]string{
					"org":  newParsed.Organization,
					"name": newParsed.Name,
					"tag":  newParsed.Tag,
					"hash": e.Hash,
				})
			},
		}

		cmd.AddCommand(tagCmd)
	}
}
