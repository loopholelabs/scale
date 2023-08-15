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

package main

import (
	"github.com/loopholelabs/cmdutils/pkg/command"
	"github.com/loopholelabs/scale/cli/cmd/access/apikey"
	"github.com/loopholelabs/scale/cli/cmd/auth"
	"github.com/loopholelabs/scale/cli/cmd/function"
	"github.com/loopholelabs/scale/cli/cmd/registry"
	"github.com/loopholelabs/scale/cli/cmd/signature"
	"github.com/loopholelabs/scale/cli/cmd/update"
	"github.com/loopholelabs/scale/cli/internal/config"
	"github.com/loopholelabs/scale/cli/version"
)

var Cmd = command.New[*config.Config](
	"scale",
	"A CLI for working with Scale Functions and communicating with the Scale API.",
	"scale is a CLI Library for working with Scale Functions and communicating with the Scale API.",
	true,
	version.V,
	config.New,
	[]command.SetupCommand[*config.Config]{signature.Cmd(), function.Cmd(), auth.Cmd(), apikey.Cmd(), registry.Cmd(), update.Cmd()},
)
