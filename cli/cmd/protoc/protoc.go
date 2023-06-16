/*
	Copyright 2022 Loophole Labs

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

package protoc

import (
	"github.com/loopholelabs/scale-signature/generator"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// Cmd returns the base command for the protoc version of scale.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "protoc",
		Short: "Start in Protoc Plugin Mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			gen := generator.New()

			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}

			req, err := gen.UnmarshalRequest(data)
			if err != nil {
				return err
			}

			res, err := gen.Generate(req)
			if err != nil {
				return err
			}

			data, err = gen.MarshalResponse(res)
			if err != nil {
				return err
			}

			_, err = os.Stdout.Write(data)
			if err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
