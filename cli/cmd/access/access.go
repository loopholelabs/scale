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

package access

import (
	"github.com/loopholelabs/cmdutils"
	"github.com/loopholelabs/cmdutils/pkg/command"
	"github.com/loopholelabs/scale/cli/cmd/access/apikey"
	"github.com/loopholelabs/scale/cli/cmd/utils"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/spf13/cobra"
)

// Cmd encapsulates the commands for access.
func Cmd() command.SetupCommand[*config.Config] {
	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
		accessCmd := &cobra.Command{
			Use:                "access",
			Short:              "Commands for managing access to the Scale API",
			PersistentPreRunE:  utils.PreRunAuthenticatedAPI(ch),
			PersistentPostRunE: utils.PostRunAuthenticatedAPI(ch),
		}

		apikeySetup := apikey.Cmd()
		apikeySetup(accessCmd, ch)

		cmd.AddCommand(accessCmd)
	}
}
