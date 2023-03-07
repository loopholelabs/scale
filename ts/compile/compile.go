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
	"io"
	"text/template"

	"github.com/loopholelabs/scale/ts/compile/templates"
)

type Generator struct {
	template *template.Template
}

func NewGenerator() *Generator {
	templ := template.Must(template.New("typescript").ParseFS(templates.FS, "*"))
	return &Generator{
		template: templ,
	}
}

func (g *Generator) GenerateRunner(writer io.Writer, path string, signature string) error {
	return g.template.ExecuteTemplate(writer, "runner.ts.templ", map[string]interface{}{
		"path":      path,
		"signature": signature,
	})
}
