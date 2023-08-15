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

package auth

import (
	"bufio"
	"fmt"
	"github.com/loopholelabs/auth/pkg/client/logout"
	"github.com/loopholelabs/cmdutils"
	"github.com/loopholelabs/cmdutils/pkg/command"
	"github.com/loopholelabs/cmdutils/pkg/printer"
	"github.com/loopholelabs/scale/cli/analytics"
	"github.com/loopholelabs/scale/cli/cmd/utils"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// LogoutCmd encapsulates the commands for logging out
func LogoutCmd(hidden bool) command.SetupCommand[*config.Config] {
	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
		logoutCmd := &cobra.Command{
			Use:      "logout",
			Args:     cobra.NoArgs,
			Short:    "Log out of the Scale API",
			Hidden:   hidden,
			PreRunE:  utils.PreRunUpdateCheck(ch),
			PostRunE: utils.PostRunAnalytics(ch),
			RunE: func(cmd *cobra.Command, args []string) error {
				if !ch.Config.IsAuthenticated() {
					ch.Printer.Println("Already logged out. Exiting...")
					return nil
				}

				if printer.IsTTY && ch.Printer.Format() == printer.Human {
					ch.Printer.Println("Press Enter to log out of the Scale API.")
					_ = waitForEnter(cmd.InOrStdin())
				}

				ctx := cmd.Context()

				end := ch.Printer.PrintProgress("Logging out...")
				defer end()

				c, err := ch.Config.NewAuthenticatedAuthClient()
				if err != nil {
					return fmt.Errorf("error creating authenticated auth client: %w", err)
				}

				_, err = c.Logout.PostLogout(logout.NewPostLogoutParamsWithContext(ctx))
				if err != nil {
					return fmt.Errorf("error logging out: %w", err)
				}

				err = deleteSession(ch.Config)
				if err != nil {
					return err
				}

				analytics.Event("logout")

				end()
				ch.Printer.Println("Successfully logged out.")

				return nil
			},
		}

		cmd.AddCommand(logoutCmd)
	}
}

func deleteSession(c *config.Config) error {
	sessionPath, err := c.SessionPath()
	if err != nil {
		return err
	}

	err = os.Remove(sessionPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrap(err, "error removing session file")
		}
	}

	return nil
}
func waitForEnter(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	return scanner.Err()
}
