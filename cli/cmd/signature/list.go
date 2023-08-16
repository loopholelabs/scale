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

package signature

import (
	"fmt"
	"github.com/loopholelabs/cmdutils"
	"github.com/loopholelabs/cmdutils/pkg/command"
	"github.com/loopholelabs/cmdutils/pkg/printer"
	"github.com/loopholelabs/scale/cli/analytics"
	"github.com/loopholelabs/scale/cli/cmd/utils"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/loopholelabs/scale/storage"
	"github.com/spf13/cobra"
)

// ListCmd encapsulates the commands for listing the available Signatures
func ListCmd(hidden bool) command.SetupCommand[*config.Config] {
	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
		listCmd := &cobra.Command{
			Use:     "list",
			Short:   "list locally available scale functions",
			Args:    cobra.NoArgs,
			PreRunE: utils.PreRunUpdateCheck(ch),
			RunE: func(cmd *cobra.Command, args []string) error {
				st := storage.DefaultSignature
				if ch.Config.StorageDirectory != "" {
					var err error
					st, err = storage.NewSignature(ch.Config.StorageDirectory)
					if err != nil {
						return fmt.Errorf("failed to instantiate signature storage for %s: %w", ch.Config.StorageDirectory, err)
					}
				}

				analytics.Event("list-signature")

				signatureEntries, err := st.List()
				if err != nil {
					return fmt.Errorf("failed to list scale signatures: %w", err)
				}

				if len(signatureEntries) == 0 && ch.Printer.Format() == printer.Human {
					ch.Printer.Println("No Scale Signatures available yet.")
					return nil
				}

				sigs := make([]signatureModel, len(signatureEntries))
				for i, entry := range signatureEntries {
					sigs[i] = signatureModel{
						Name:    entry.Schema.Name,
						Tag:     entry.Schema.Tag,
						Org:     entry.Organization,
						Hash:    entry.Hash,
						Version: entry.Schema.Version,
					}
				}

				return ch.Printer.PrintResource(sigs)
			},
		}

		cmd.AddCommand(listCmd)
	}
}
