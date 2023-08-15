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

package generator

import (
	"golang.org/x/mod/zip"
)

type GolangModule struct {
	CanonicalVersion string
	ImportPath       string
	Files            []GolangFile
}

func (g *GolangModule) ModuleFiles() []zip.File {
	files := make([]zip.File, len(g.Files))
	for i, file := range g.Files {
		files[i] = file
	}
	return files
}
