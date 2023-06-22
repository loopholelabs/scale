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
	"github.com/loopholelabs/cmdutils"
	"github.com/loopholelabs/cmdutils/pkg/command"
	"github.com/loopholelabs/scale/cli/cmd/utils"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/spf13/cobra"
)

type signatureModel struct {
	Name    string `header:"name" json:"name"`
	Tag     string `header:"tag" json:"tag"`
	Org     string `header:"org" json:"org"`
	Hash    string `header:"hash" json:"hash"`
	Version string `header:"version" json:"version"`
}

// Cmd encapsulates the commands for signatures.
func Cmd() command.SetupCommand[*config.Config] {
	return func(cmd *cobra.Command, ch *cmdutils.Helper[*config.Config]) {
		signatureCmd := &cobra.Command{
			Use:                "signature <command>",
			Aliases:            []string{"sig"},
			Short:              "Create, list, and manage Scale Signatures",
			PersistentPostRunE: utils.PostRunAnalytics(ch),
		}

		newSetup := NewCmd(false)
		newSetup(signatureCmd, ch)

		generateSetup := GenerateCmd(false)
		generateSetup(signatureCmd, ch)

		listSetup := ListCmd(false)
		listSetup(signatureCmd, ch)

		deleteSetup := DeleteCmd(false)
		deleteSetup(signatureCmd, ch)

		useSetup := UseCmd(false)
		useSetup(signatureCmd, ch)

		cmd.AddCommand(signatureCmd)
	}
}
