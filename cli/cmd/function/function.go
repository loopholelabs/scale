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

import (
	"github.com/loopholelabs/cmdutils"
	"github.com/loopholelabs/cmdutils/pkg/command"
	"github.com/loopholelabs/scale/cli/cmd/utils"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/spf13/cobra"
)

type scaleFunction struct {
	Name      string `header:"name" json:"name"`
	Tag       string `header:"tag" json:"tag"`
	Org       string `header:"org" json:"org"`
	Language  string `header:"language" json:"language"`
	Signature string `header:"signature" json:"signature"`
	Hash      string `header:"hash" json:"hash"`
	Version   string `header:"version" json:"version"`
}

// Cmd encapsulates the commands for functions.
func Cmd() command.SetupCommand[*config.Config] {
	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
		functionCmd := &cobra.Command{
			Use:                "function <command>",
			Aliases:            []string{"fn"},
			Short:              "Create, list, and manage local Scale Functions",
			PersistentPostRunE: utils.PostRunAnalytics(ch),
		}

		listSetup := ListCmd()
		listSetup(functionCmd, ch)

		newSetup := NewCmd(false)
		newSetup(functionCmd, ch)

		buildSetup := BuildCmd(false)
		buildSetup(functionCmd, ch)

		deleteSetup := DeleteCmd()
		deleteSetup(functionCmd, ch)

		exportSetup := ExportCmd()
		exportSetup(functionCmd, ch)

		runSetup := RunCmd(false)
		runSetup(functionCmd, ch)

		deploySetup := DeployCmd(false)
		deploySetup(functionCmd, ch)

		tagSetup := TagCmd()
		tagSetup(functionCmd, ch)

		newAliasSetup := NewCmd(true)
		newAliasSetup(cmd, ch)

		buildAliasSetup := BuildCmd(true)
		buildAliasSetup(cmd, ch)

		runAliasSetup := RunCmd(true)
		runAliasSetup(cmd, ch)

		deployAliasSetup := DeployCmd(true)
		deployAliasSetup(cmd, ch)

		cmd.AddCommand(functionCmd)
	}
}
