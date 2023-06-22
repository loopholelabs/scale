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

//// BuildCmd encapsulates the commands for building Functions
//func BuildCmd(hidden bool) command.SetupCommand[*config.Config] {
//	var name string
//	var tag string
//	var org string
//	var directory string
//
//	var goBin string
//	var tinygoBin string
//	var cargoBin string
//	var npmBin string
//
//	var tinygoArgs []string
//	var cargoArgs []string
//
//	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
//		buildCmd := &cobra.Command{
//			Use:     "build [flags]",
//			Args:    cobra.ExactArgs(0),
//			Short:   "build a scale function locally and store it in the cache",
//			Long:    "Build a scale function locally and store it in the cache. The scalefile must be in the current directory or specified with the --directory flag.",
//			Hidden:  hidden,
//			PreRunE: utils.PreRunUpdateCheck(ch),
//			RunE: func(cmd *cobra.Command, args []string) error {
//				scaleFilePath := path.Join(directory, "scalefile")
//				scaleFile, err := scalefile.Read(scaleFilePath)
//				if err != nil {
//					return fmt.Errorf("failed to read scalefile at %s: %w", scaleFilePath, err)
//				}
//
//				if org == "" {
//					org = utils.DefaultOrganization
//				}
//
//				if name == "" {
//					name = scaleFile.Name
//				} else {
//					scaleFile.Name = name
//				}
//
//				if tag == "" {
//					tag = scaleFile.Tag
//				} else {
//					scaleFile.Tag = tag
//				}
//
//				if !scalefunc.ValidString(org) {
//					return utils.InvalidStringError("organization", org)
//				}
//
//				if scaleFile.Name == "" || !scalefunc.ValidString(scaleFile.Name) {
//					return utils.InvalidStringError("name", scaleFile.Name)
//				}
//
//				if scaleFile.Tag == "" || !scalefunc.ValidString(scaleFile.Tag) {
//					return utils.InvalidStringError("tag", scaleFile.Tag)
//				}
//
//				if analytics.Client != nil {
//					_ = analytics.Client.Enqueue(posthog.Capture{
//						DistinctId: analytics.MachineID,
//						Event:      "build-function",
//						Timestamp:  time.Now(),
//						Properties: posthog.NewProperties().Set("language", scaleFile.Language),
//					})
//				}
//
//				end := ch.Printer.PrintProgress(fmt.Sprintf("Building scale function %s:%s...", scaleFile.Name, scaleFile.Tag))
//				scaleFunc, err := build.LocalBuild(scaleFile, goBin, tinygoBin, cargoBin, npmBin, directory, tinygoArgs, cargoArgs)
//				end()
//				if err != nil {
//					return fmt.Errorf("failed to build scale function: %w", err)
//				}
//
//				hash := sha256.New()
//				hash.Write(scaleFunc.Encode())
//				checksum := hex.EncodeToString(hash.Sum(nil))
//
//				st := storage.Default
//				if ch.Config.CacheDirectory != "" {
//					st, err = storage.New(ch.Config.CacheDirectory)
//					if err != nil {
//						return fmt.Errorf("failed to instantiate function storage for %s: %w", ch.Config.CacheDirectory, err)
//					}
//				}
//
//				oldEntry, err := st.Get(scaleFunc.Name, scaleFunc.Tag, org, "")
//				if err != nil {
//					return fmt.Errorf("failed to check if scale function already exists: %w", err)
//				}
//
//				if oldEntry != nil {
//					err = st.Delete(name, tag, oldEntry.Organization, oldEntry.Hash)
//					if err != nil {
//						return fmt.Errorf("failed to delete existing scale function %s:%s: %w", name, tag, err)
//					}
//				}
//
//				err = st.Put(scaleFunc.Name, scaleFunc.Tag, org, checksum, scaleFunc)
//				if err != nil {
//					return fmt.Errorf("failed to store scale function: %w", err)
//				}
//
//				if org == utils.DefaultOrganization {
//					org = ""
//				}
//
//				if ch.Printer.Format() == printer.Human {
//					if org != "" {
//						ch.Printer.Printf("Successfully built scale function %s\n", printer.BoldGreen(fmt.Sprintf("%s/%s:%s", org, scaleFunc.Name, scaleFunc.Tag)))
//					} else {
//						ch.Printer.Printf("Successfully built scale function %s\n", printer.BoldGreen(fmt.Sprintf("%s:%s", scaleFunc.Name, scaleFunc.Tag)))
//					}
//					return nil
//				}
//
//				return ch.Printer.PrintResource(map[string]string{
//					"name":      name,
//					"tag":       tag,
//					"org":       org,
//					"directory": directory,
//				})
//			},
//		}
//
//		buildCmd.Flags().StringVarP(&directory, "directory", "d", ".", "the directory containing the scalefile")
//		buildCmd.Flags().StringVarP(&name, "name", "n", "", "the (optional) name of this scale function")
//		buildCmd.Flags().StringVarP(&tag, "tag", "t", "", "the (optional) tag of this scale function")
//		buildCmd.Flags().StringVarP(&org, "org", "o", "", "the (optional) organization of this scale function")
//
//		buildCmd.Flags().StringVar(&tinygoBin, "tinygo", "", "the (optional) path to the tinygo binary")
//		buildCmd.Flags().StringVar(&goBin, "go", "", "the (optional) path to the go binary")
//		buildCmd.Flags().StringVar(&cargoBin, "cargo", "", "the (optional) path to the cargo binary")
//		buildCmd.Flags().StringVar(&npmBin, "npm", "", "the (optional) path to the npm binary")
//
//		buildCmd.Flags().StringSliceVar(&tinygoArgs, "tinygo-args", []string{"-scheduler=none", "--no-debug"}, "list of (optional) tinygo build arguments")
//		buildCmd.Flags().StringSliceVar(&cargoArgs, "cargo-args", []string{"--release"}, "list of (optional) cargo build arguments")
//
//		cmd.AddCommand(buildCmd)
//	}
//}
