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

package compile

import (
	"github.com/loopholelabs/scale/go/compile/templates"
	"github.com/loopholelabs/scalefile"
	"io"
	"text/template"
)

type Generator struct {
	template *template.Template
}

func NewGenerator() *Generator {
	templ := template.Must(template.New("main").ParseFS(templates.FS, "*"))
	return &Generator{
		template: templ,
	}
}

func (g *Generator) GenerateGoModfile(writer io.Writer, dependencies []*scalefile.Dependency) error {
	return g.template.ExecuteTemplate(writer, "go.mod.templ", map[string]interface{}{
		"dependencies": dependencies,
	})
}

func (g *Generator) GenerateGoMain(writer io.Writer, signature string, path string) error {
	return g.template.ExecuteTemplate(writer, "main.go.templ", map[string]interface{}{
		"path":      path,
		"signature": signature,
	})
}
