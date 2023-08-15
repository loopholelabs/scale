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
	"github.com/loopholelabs/cmdutils"
	"github.com/loopholelabs/cmdutils/pkg/command"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/spf13/cobra"
)

// Cmd encapsulates the commands for authentication.
func Cmd() command.SetupCommand[*config.Config] {
	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
		authCmd := &cobra.Command{
			Use:   "auth",
			Short: "Login and Logout using the Scale Authentication API",
			Long:  "Manage access to the Scale API using the Scale Authentication API",
		}

		loginSetup := LoginCmd(false)
		loginSetup(authCmd, ch)

		logoutSetup := LogoutCmd(false)
		logoutSetup(authCmd, ch)

		statusSetup := StatusCmd()
		statusSetup(authCmd, ch)

		loginAliasSetup := LoginCmd(true)
		loginAliasSetup(cmd, ch)

		logoutAliasSetup := LogoutCmd(true)
		logoutAliasSetup(cmd, ch)

		cmd.AddCommand(authCmd)
	}
}
