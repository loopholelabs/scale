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

package boilerplate

import (
	"bytes"
	"text/template"

	"github.com/loopholelabs/scale/signature/generator/templates"
)

var generator *Generator

func Generate(name string, tag string) ([]byte, error) {
	return generator.Generate(name, tag)
}

func init() {
	var err error
	generator, err = New()
	if err != nil {
		panic(err)
	}
}

// Generator is the go generator
type Generator struct {
	templ *template.Template
}

// New creates a new go generator
func New() (*Generator, error) {
	templ, err := template.New("").ParseFS(templates.FS, "boilerplate.templ")
	if err != nil {
		return nil, err
	}

	return &Generator{
		templ: templ,
	}, nil
}

// Generate generates the go code
func (g *Generator) Generate(name string, tag string) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "boilerplate.templ", map[string]any{
		"Name": name,
		"Tag":  tag,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
