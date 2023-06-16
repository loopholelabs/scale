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

package golang

import (
	"github.com/loopholelabs/scale/compile/templates"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature/schema"
	"io"
	"text/template"
)

type Generator struct {
	template *template.Template
}

func New() *Generator {
	return &Generator{
		template: template.Must(template.New("main").ParseFS(templates.FS, "*go.templ")),
	}
}

func (g *Generator) GenerateGoModfile(writer io.Writer, signatureImport string, signatureVersion string, dependencies []*scalefunc.Dependency, version string) error {
	return g.template.ExecuteTemplate(writer, "go.mod.templ", map[string]interface{}{
		"version":        version,
		"dependencies":   dependencies,
		"old_dependency": "signature",
		"old_version":    "v0.1.0",
		"new_dependency": signatureImport,
		"new_version":    signatureVersion,
	})
}

func (g *Generator) GenerateGoMain(writer io.Writer, signature *schema.Schema, schema *scalefunc.ScaleFunc, version string) error {
	return g.template.ExecuteTemplate(writer, "main.go.templ", map[string]interface{}{
		"version":   version,
		"signature": signature,
		"schema":    schema,
	})
}
