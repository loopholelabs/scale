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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/loopholelabs/cmdutils"
	"github.com/loopholelabs/cmdutils/pkg/command"
	"github.com/loopholelabs/cmdutils/pkg/printer"
	"github.com/loopholelabs/scale/cli/analytics"
	"github.com/loopholelabs/scale/cli/cmd/utils"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/storage"
	"github.com/posthog/posthog-go"
	"github.com/spf13/cobra"
	"path"
	"time"
)

// GenerateCmd encapsulates the commands for generating a Signature from a Signature File
func GenerateCmd(hidden bool) command.SetupCommand[*config.Config] {
	var name string
	var tag string
	var org string
	var directory string

	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
		generateCmd := &cobra.Command{
			Use:     "generate [flags]",
			Args:    cobra.ExactArgs(0),
			Short:   "generate a scale signature from a signature file",
			Hidden:  hidden,
			PreRunE: utils.PreRunUpdateCheck(ch),
			RunE: func(cmd *cobra.Command, args []string) error {
				signaturePath := path.Join(directory, "scale.signature")
				signatureFile, err := signature.ReadSchema(signaturePath)
				if err != nil {
					return fmt.Errorf("failed to read signature file at %s: %w", signaturePath, err)
				}

				err = signatureFile.Validate()
				if err != nil {
					return fmt.Errorf("failed to validate signature file: %w", err)
				}

				if org == "" {
					org = utils.DefaultOrganization
				}

				if name == "" {
					name = signatureFile.Name
				} else {
					signatureFile.Name = name
				}

				if tag == "" {
					tag = signatureFile.Tag
				} else {
					signatureFile.Tag = tag
				}

				if !scalefunc.ValidString(org) {
					return utils.InvalidStringError("organization", org)
				}

				if signatureFile.Name == "" || !scalefunc.ValidString(signatureFile.Name) {
					return utils.InvalidStringError("name", signatureFile.Name)
				}

				if signatureFile.Tag == "" || !scalefunc.ValidString(signatureFile.Tag) {
					return utils.InvalidStringError("tag", signatureFile.Tag)
				}

				if analytics.Client != nil {
					_ = analytics.Client.Enqueue(posthog.Capture{
						DistinctId: analytics.MachineID,
						Event:      "generate-signature",
						Timestamp:  time.Now(),
					})
				}

				end := ch.Printer.PrintProgress(fmt.Sprintf("Generating scale signature %s/%s:%s...", org, signatureFile.Name, signatureFile.Tag))
				hash := sha256.New()
				hashData, err := signatureFile.Encode()
				if err != nil {
					end()
					return fmt.Errorf("failed to encode signature file: %w", err)
				}
				hash.Write(hashData)
				checksum := hex.EncodeToString(hash.Sum(nil))

				st := storage.DefaultSignature
				if ch.Config.StorageDirectory != "" {
					st, err = storage.NewSignature(ch.Config.StorageDirectory)
					if err != nil {
						end()
						return fmt.Errorf("failed to instantiate function storage for %s: %w", ch.Config.StorageDirectory, err)
					}
				}

				oldEntry, err := st.Get(signatureFile.Name, signatureFile.Tag, org, "")
				if err != nil {
					end()
					return fmt.Errorf("failed to check if scale signature already exists: %w", err)
				}

				if oldEntry != nil {
					err = st.Delete(name, tag, oldEntry.Organization, oldEntry.Hash)
					if err != nil {
						end()
						return fmt.Errorf("failed to delete existing scale signature %s:%s: %w", name, tag, err)
					}
				}

				err = st.Put(signatureFile.Name, signatureFile.Tag, org, checksum, signatureFile)
				if err != nil {
					end()
					return fmt.Errorf("failed to store scale signature: %w", err)
				}

				end()

				if ch.Printer.Format() == printer.Human {
					ch.Printer.Printf("Successfully generated scale signature %s\n", printer.BoldGreen(fmt.Sprintf("%s/%s:%s", org, signatureFile.Name, signatureFile.Tag)))
					return nil
				}

				return ch.Printer.PrintResource(map[string]string{
					"name":      name,
					"tag":       tag,
					"org":       org,
					"directory": directory,
				})
			},
		}

		generateCmd.Flags().StringVarP(&directory, "directory", "d", ".", "the directory containing the signature file")
		generateCmd.Flags().StringVarP(&name, "name", "n", "", "the (optional) name of this scale signature")
		generateCmd.Flags().StringVarP(&tag, "tag", "t", "", "the (optional) tag of this scale signature")
		generateCmd.Flags().StringVarP(&org, "org", "o", utils.DefaultOrganization, "the (optional) organization of this scale signature")

		cmd.AddCommand(generateCmd)
	}
}
