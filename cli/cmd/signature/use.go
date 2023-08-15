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
	"github.com/loopholelabs/scale/cli/cmd/utils"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/loopholelabs/scale/compile/golang"
	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/storage"
	"github.com/spf13/cobra"
	"os"
	"path"
)

// UseCmd encapsulates the commands for using a Signature
func UseCmd(hidden bool) command.SetupCommand[*config.Config] {
	var directory string
	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
		useCmd := &cobra.Command{
			Use:     "use <org>/<name>:<tag> [flags]",
			Args:    cobra.ExactArgs(1),
			Short:   "create a new scale signature with the given name and tag",
			Hidden:  hidden,
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

				parsed := utils.ParseFunction(args[0])
				if parsed.Organization != "" && !scalefunc.ValidString(parsed.Organization) {
					return utils.InvalidStringError("organization name", parsed.Organization)
				}

				if parsed.Name == "" || !scalefunc.ValidString(parsed.Name) {
					return utils.InvalidStringError("signature name", parsed.Name)
				}

				if parsed.Tag == "" || !scalefunc.ValidString(parsed.Tag) {
					return utils.InvalidStringError("signature tag", parsed.Tag)
				}

				signaturePath, err := st.Path(parsed.Name, parsed.Tag, parsed.Organization, "")
				if err != nil || signaturePath == "" {
					return fmt.Errorf("failed to use signature %s/%s:%s: %w", parsed.Organization, parsed.Name, parsed.Tag, err)
				}

				sf, err := scalefile.ReadSchema(path.Join(directory, "scalefile"))
				if err != nil {
					return fmt.Errorf("failed to use signature %s/%s:%s: %w", parsed.Organization, parsed.Name, parsed.Tag, err)
				}

				switch scalefunc.Language(sf.Language) {
				case scalefunc.Go:
					modfileData, err := os.ReadFile(path.Join(directory, "go.mod"))
					if err != nil {
						return fmt.Errorf("failed to use signature %s/%s:%s: %w", parsed.Organization, parsed.Name, parsed.Tag, err)
					}

					m, err := golang.ParseManifest(modfileData)
					if err != nil {
						return fmt.Errorf("failed to use signature %s/%s:%s: %w", parsed.Organization, parsed.Name, parsed.Tag, err)
					}

					err = m.RemoveReplacement("signature", golang.DefaultVersion)
					if err != nil {
						return fmt.Errorf("failed to use signature %s/%s:%s: %w", parsed.Organization, parsed.Name, parsed.Tag, err)
					}

					err = m.RemoveReplacement("signature", "")
					if err != nil {
						return fmt.Errorf("failed to use signature %s/%s:%s: %w", parsed.Organization, parsed.Name, parsed.Tag, err)
					}

					err = m.AddReplacement("signature", golang.DefaultVersion, path.Join(signaturePath, "golang"), "")
					if err != nil {
						return fmt.Errorf("failed to use signature %s/%s:%s: %w", parsed.Organization, parsed.Name, parsed.Tag, err)
					}
					modfileData, err = m.Write()
					if err != nil {
						return fmt.Errorf("failed to use signature %s/%s:%s: %w", parsed.Organization, parsed.Name, parsed.Tag, err)
					}

					err = os.WriteFile(path.Join(directory, "go.mod"), modfileData, 0644)
					if err != nil {
						return fmt.Errorf("failed to use signature %s/%s:%s: %w", parsed.Organization, parsed.Name, parsed.Tag, err)
					}
				default:
					return fmt.Errorf("failed to use signature %s/%s:%s: unknown or unsupported language", parsed.Organization, parsed.Name, parsed.Tag)
				}

				sf.Signature.Name = parsed.Name
				sf.Signature.Tag = parsed.Tag
				sf.Signature.Organization = parsed.Organization

				sfData, err := sf.Encode()
				if err != nil {
					return fmt.Errorf("failed to use signature %s/%s:%s: %w", parsed.Organization, parsed.Name, parsed.Tag, err)
				}

				err = os.WriteFile(path.Join(directory, "scalefile"), sfData, 0644)
				if err != nil {
					return fmt.Errorf("failed to use signature %s/%s:%s: %w", parsed.Organization, parsed.Name, parsed.Tag, err)
				}

				if ch.Printer.Format() == printer.Human {
					ch.Printer.Printf("Successfully using scale signature %s\n", printer.BoldGreen(fmt.Sprintf("%s/%s:%s", parsed.Organization, parsed.Name, parsed.Tag)))
					return nil
				}

				return ch.Printer.PrintResource(map[string]string{
					"path": signaturePath,
					"name": parsed.Name,
					"org":  parsed.Organization,
					"tag":  parsed.Tag,
				})
			},
		}

		useCmd.Flags().StringVarP(&directory, "directory", "d", ".", "the directory that contains the scalefile and the function source")

		cmd.AddCommand(useCmd)
	}
}
