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

//var (
//	InvalidDeployRegex = regexp.MustCompile(`[^A-Za-z0-9-]`)
//)
//
//// DeployCmd encapsulates the commands for deploying Functions
//func DeployCmd(hidden bool) command.SetupCommand[*config.Config] {
//	var name string
//	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
//		deployCmd := &cobra.Command{
//			Use:      "deploy [ ...[ <name>:<tag> ] | [ <org>/<name>:<tag> ] ]",
//			Args:     cobra.MinimumNArgs(1),
//			Short:    "deploy a Scale function",
//			Long:     "Deploy a Scale function",
//			Hidden:   hidden,
//			PreRunE:  utils.PreRunAuthenticatedAPI(ch),
//			PostRunE: utils.PostRunAuthenticatedAPI(ch),
//			RunE: func(cmd *cobra.Command, args []string) error {
//				st := storage.Default
//				if ch.Config.CacheDirectory != "" {
//					var err error
//					st, err = storage.New(ch.Config.CacheDirectory)
//					if err != nil {
//						return fmt.Errorf("failed to instantiate function storage for %s: %w", ch.Config.CacheDirectory, err)
//					}
//				}
//
//				if name != "" && InvalidDeployRegex.MatchString(name) {
//					return fmt.Errorf("invalid deploy name '%s', deploy names can only include letters, numbers, and dashes (`-`)", name)
//				}
//
//				fns := make([]*scalefunc.ScaleFunc, 0, len(args))
//				for _, f := range args {
//					parsed := utils.ParseFunction(f)
//					if parsed.Organization == "" {
//						parsed.Organization = utils.DefaultOrganization
//					}
//
//					if parsed.Organization != "" && !scalefunc.ValidString(parsed.Organization) {
//						return utils.InvalidStringError("organization name", parsed.Organization)
//					}
//
//					if parsed.Name == "" || !scalefunc.ValidString(parsed.Name) {
//						return utils.InvalidStringError("function name", parsed.Name)
//					}
//
//					if parsed.Tag == "" || !scalefunc.ValidString(parsed.Tag) {
//						return utils.InvalidStringError("function tag", parsed.Tag)
//					}
//
//					e, err := st.Get(parsed.Name, parsed.Tag, parsed.Organization, "")
//					if err != nil {
//						return fmt.Errorf("failed to get function %s: %w", f, err)
//					}
//
//					if e == nil {
//						end := ch.Printer.PrintProgress(fmt.Sprintf("Function %s was not found not found, pulling from the registry...", printer.BoldGreen(f)))
//						var opts []registry.Option
//						opts = append(opts, registry.WithClient(ch.Config.APIClient()), registry.WithStorage(st))
//						if parsed.Organization != "" && parsed.Organization != utils.DefaultOrganization {
//							opts = append(opts, registry.WithOrganization(parsed.Organization))
//						}
//						sf, err := registry.Download(parsed.Name, parsed.Tag, opts...)
//						end()
//						if err != nil {
//							if parsed.Organization == "" || parsed.Organization == utils.DefaultOrganization {
//								return fmt.Errorf("scale function %s:%s not found", parsed.Name, parsed.Tag)
//							} else {
//								return fmt.Errorf("scale function %s/%s:%s: not found", parsed.Organization, parsed.Name, parsed.Tag)
//							}
//						}
//
//						if analytics.Client != nil {
//							_ = analytics.Client.Enqueue(posthog.Capture{
//								DistinctId: analytics.MachineID,
//								Event:      "pull-registry",
//								Timestamp:  time.Now(),
//							})
//						}
//
//						if ch.Printer.Format() == printer.Human {
//							if parsed.Organization == "" || parsed.Organization == utils.DefaultOrganization {
//								ch.Printer.Printf("Pulled %s from the Scale Registry\n", printer.BoldGreen(fmt.Sprintf("%s:%s", sf.Name, sf.Tag)))
//							} else {
//								ch.Printer.Printf("Pulled %s from the Scale Registry\n", printer.BoldGreen(fmt.Sprintf("%s/%s:%s", parsed.Organization, sf.Name, sf.Tag)))
//							}
//						}
//						fns = append(fns, sf)
//					} else {
//						fns = append(fns, e.ScaleFunc)
//					}
//				}
//
//				if analytics.Client != nil {
//					_ = analytics.Client.Enqueue(posthog.Capture{
//						DistinctId: analytics.MachineID,
//						Event:      "deploy-function",
//						Timestamp:  time.Now(),
//						Properties: posthog.NewProperties().Set("chain-size", len(fns)),
//					})
//				}
//
//				buffer := polyglot.NewBuffer()
//				encoder := polyglot.Encoder(buffer)
//				encoder.Slice(uint32(len(fns)), polyglot.BytesKind)
//				for _, f := range fns {
//					encoder.Bytes(f.Encode())
//				}
//
//				ctx := cmd.Context()
//				client := ch.Config.APIClient()
//
//				reader := runtime.NamedReader("functions", bytes.NewReader(buffer.Bytes()))
//
//				var end = ch.Printer.PrintProgress("Deploying to Scale Cloud...")
//				params := deploy.NewPostDeployFunctionParamsWithContext(ctx).WithFunctions(reader).WithName(&name)
//				resp, err := client.Deploy.PostDeployFunction(params)
//				end()
//				if err != nil {
//					return fmt.Errorf("failed to deploy function: %w", err)
//				}
//
//				payload := resp.GetPayload()
//
//				// Wait for function to be deployed globally before printing the URL
//				time.Sleep(3 * time.Second)
//
//				ch.Printer.Printf("Functions deployed, available at %s\n", printer.BoldGreen(fmt.Sprintf("https://%s.%s", payload.Subdomain, payload.RootDomain)))
//				return nil
//			},
//		}
//
//		deployCmd.Flags().StringVarP(&name, "name", "n", "", "name of the function")
//
//		cmd.AddCommand(deployCmd)
//	}
//}
