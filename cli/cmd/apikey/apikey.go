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

package apikey

import (
	"github.com/loopholelabs/cmdutils"
	"github.com/loopholelabs/cmdutils/pkg/command"
	"github.com/loopholelabs/scale/cli/cmd/utils"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/spf13/cobra"
)

type apiKey struct {
	Name    string `header:"name" json:"name"`
	ID      string `header:"id" json:"id"`
	Value   string `header:"value" json:"value"`
	Created string `header:"created" json:"created"`
}

type apiKeyRedacted struct {
	Name    string `header:"name" json:"name"`
	ID      string `header:"id" json:"id"`
	Created string `header:"created" json:"created"`
}

// Cmd encapsulates the commands for authentication.
func Cmd() command.SetupCommand[*config.Config] {
	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
		apikeyCmd := &cobra.Command{
			Use:                "apikey",
			Short:              "Create, list, and manage API Keys",
			PersistentPreRunE:  utils.PreRunAuthenticatedAPI(ch),
			PersistentPostRunE: utils.PostRunAuthenticatedAPI(ch),
		}

		listSetup := ListCmd()
		listSetup(apikeyCmd, ch)

		getSetup := GetCmd()
		getSetup(apikeyCmd, ch)

		createSetup := CreateCmd()
		createSetup(apikeyCmd, ch)

		deleteSetup := DeleteCmd()
		deleteSetup(apikeyCmd, ch)

		cmd.AddCommand(apikeyCmd)
	}
}
