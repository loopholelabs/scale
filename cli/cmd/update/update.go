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

package update

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/loopholelabs/cmdutils"
	"github.com/loopholelabs/cmdutils/pkg/command"
	"github.com/loopholelabs/cmdutils/pkg/printer"
	"github.com/loopholelabs/releaser/pkg/client"
	"github.com/loopholelabs/scale/cli/analytics"
	"github.com/loopholelabs/scale/cli/cmd/utils"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/loopholelabs/scale/cli/internal/log"
	"github.com/loopholelabs/scale/cli/version"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Cmd encapsulates the commands for updating the CLI.
func Cmd() command.SetupCommand[*config.Config] {
	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
		updateCmd := &cobra.Command{
			Use:   "update",
			Short: "Update the Scale CLI to the latest version",
			Long:  "Update the Scale CLI to the latest version using the Scale Update service",
			PreRunE: func(cmd *cobra.Command, args []string) error {
				log.Init(ch.Config.GetLogFile())
				err := ch.Config.GlobalRequiredFlags(cmd)
				if err != nil {
					return err
				}
				return ch.Config.Validate()
			},
			PostRunE: utils.PostRunAnalytics(ch),
			RunE: func(cmd *cobra.Command, args []string) error {
				c := client.New(fmt.Sprintf("https://%s", ch.Config.UpdateEndpoint))
				latest, err := c.GetLatest()
				if err != nil {
					return fmt.Errorf("error getting latest version: %w", err)
				}
				if latest == version.Version {
					ch.Printer.Println("Scale CLI is already up to date")
					return nil
				}

				end := ch.Printer.PrintProgress(fmt.Sprintf("Updating Scale CLI to version %s...", printer.BoldGreen(latest)))
				executable, err := os.Executable()
				if err != nil {
					end()
					return fmt.Errorf("error getting executable path: %w", err)
				}

				stableExecutable, err := filepath.EvalSymlinks(executable)
				if err != nil {
					end()
					return fmt.Errorf("error getting stable executable path: %w", err)
				}

				executableName := filepath.Base(stableExecutable)
				executableDirectory := filepath.Dir(stableExecutable)
				tarExecutable := filepath.Join(executableDirectory, fmt.Sprintf("%s.tar.gz", executableName))
				previousExecutable := filepath.Join(executableDirectory, fmt.Sprintf("%s.previous", executableName))

				latestBinary, err := c.DownloadVersion(latest)
				if err != nil {
					end()
					return fmt.Errorf("error downloading latest version: %w", err)
				}

				latestTar, err := os.OpenFile(tarExecutable, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
				if err != nil {
					end()
					if errors.Is(err, os.ErrPermission) {
						return fmt.Errorf("scale update must be run as root, unable to open file %s", tarExecutable)
					}
					return fmt.Errorf("error opening file %s: %w", tarExecutable, err)
				}

				_, err = io.Copy(latestTar, bytes.NewReader(latestBinary))
				if err != nil {
					end()
					return fmt.Errorf("error writing to file %s: %w", tarExecutable, err)
				}
				_ = latestTar.Close()

				_ = os.Remove(previousExecutable)

				err = os.Rename(stableExecutable, previousExecutable)
				if err != nil {
					end()
					return fmt.Errorf("error renaming %s to %s: %w", stableExecutable, previousExecutable, err)
				}

				tarCmd := exec.Command("tar", "-xf", tarExecutable, "-C", executableDirectory)
				stdout, err := tarCmd.Output()
				if err != nil {
					end()
					ch.Printer.Printf("Error while updating Scale CLI, rolling back to previous state\n")
					rollbackErr := os.Rename(previousExecutable, stableExecutable)
					if rollbackErr != nil {
						return fmt.Errorf("error while rolling back Scale CLI to previous state: %w", rollbackErr)
					}
					return fmt.Errorf("error while updating Scale CLI: %w", err)
				}
				ch.Printer.Printf(string(stdout))
				err = os.Rename(filepath.Join(executableDirectory, "scale-cli"), stableExecutable)
				if err != nil {
					end()
					return fmt.Errorf("error renaming %s to %s: %w", filepath.Join(executableDirectory, "scale-cli"), stableExecutable, err)
				}
				_ = os.Remove(tarExecutable)
				_ = os.Remove(previousExecutable)
				end()

				if analytics.Client != nil {
					_ = analytics.Client.Enqueue(posthog.Capture{
						DistinctId: analytics.MachineID,
						Event:      "update",
						Timestamp:  time.Now(),
					})
				}

				ch.Printer.Printf("Scale CLI updated to version %s\n", printer.BoldGreen(latest))
				return nil
			},
		}

		cmd.AddCommand(updateCmd)
	}
}
