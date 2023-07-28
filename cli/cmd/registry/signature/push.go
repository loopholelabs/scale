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
	"bytes"
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/loopholelabs/cmdutils"
	"github.com/loopholelabs/cmdutils/pkg/command"
	"github.com/loopholelabs/cmdutils/pkg/printer"
	"github.com/loopholelabs/scale/cli/analytics"
	"github.com/loopholelabs/scale/cli/cmd/utils"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/loopholelabs/scale/client/registry"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/storage"
	"github.com/posthog/posthog-go"
	"github.com/spf13/cobra"
	"time"
)

// PushCmd encapsulates the commands for pushing Signatures
func PushCmd() command.SetupCommand[*config.Config] {
	var public bool
	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
		pushCmd := &cobra.Command{
			Use:      "push <org>/<name>:<tag> [flags]",
			Args:     cobra.ExactArgs(1),
			PreRunE:  utils.PreRunAuthenticatedAPI(ch),
			PostRunE: utils.PostRunAuthenticatedAPI(ch),
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
					return utils.InvalidStringError("function name", parsed.Name)
				}

				if parsed.Tag == "" || !scalefunc.ValidString(parsed.Tag) {
					return utils.InvalidStringError("function tag", parsed.Tag)
				}

				ctx := cmd.Context()
				client := ch.Config.APIClient()

				end := ch.Printer.PrintProgress(fmt.Sprintf("Pushing signature %s/%s:%s to the Scale Registry...", parsed.Organization, parsed.Name, parsed.Tag))

				e, err := st.Get(parsed.Name, parsed.Tag, parsed.Organization, "")
				if err != nil {
					return fmt.Errorf("failed to find signature %s/%s:%s: %w", parsed.Organization, parsed.Name, parsed.Tag, err)
				}
				if e == nil {
					return fmt.Errorf("signature %s/%s:%s does not exist", parsed.Organization, parsed.Name, parsed.Tag)
				}

				if analytics.Client != nil {
					_ = analytics.Client.Enqueue(posthog.Capture{
						DistinctId: analytics.MachineID,
						Event:      "push-signature",
						Timestamp:  time.Now(),
					})
				}

				encodedSchemaReader, err := e.Schema.Encode()
				if err != nil {
					return fmt.Errorf("failed to encode schema: %w", err)
				}

				res, err := client.Registry.PostRegistrySignature(registry.NewPostRegistrySignatureParamsWithContext(ctx).WithSignature(runtime.NamedReader("signature", bytes.NewReader(encodedSchemaReader))).WithPublic(&public))
				end()
				if err != nil {
					return err
				}

				p := "publicly"
				if !res.GetPayload().Public {
					p = "privately"
				}
				if ch.Printer.Format() == printer.Human {
					ch.Printer.Printf("Pushed signature %s to the Scale Registry\n", printer.BoldGreen(fmt.Sprintf("%s/%s:%s", res.GetPayload().Organization, res.GetPayload().Name, res.GetPayload().Tag)))
					ch.Printer.Printf("The following packages are now %s available:\n", printer.BoldBlue(p))
					ch.Printer.Printf("  %s: %s\n", printer.BoldBlue("Golang (Guest)"), printer.BoldGreen(res.GetPayload().GolangImportPathGuest))
					ch.Printer.Printf("  %s: %s\n\n", printer.BoldBlue("Golang (Host)"), printer.BoldGreen(res.GetPayload().GolangImportPathHost))
					return nil
				}

				return ch.Printer.PrintResource(map[string]string{
					"name":   res.GetPayload().Name,
					"tag":    res.GetPayload().Tag,
					"org":    res.GetPayload().Organization,
					"public": fmt.Sprintf("%t", res.GetPayload().Public),
				})

			},
		}

		pushCmd.Flags().BoolVar(&public, "public", false, "make the signature public")

		cmd.AddCommand(pushCmd)
	}
}
