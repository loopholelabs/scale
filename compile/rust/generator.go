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

package rust

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/loopholelabs/scale/version"

	"github.com/loopholelabs/scale/compile/rust/templates"
	"github.com/loopholelabs/scale/scalefile"
)

var generator *Generator

func GenerateRustCargofile(signatureInfo *SignatureInfo, functionInfo *FunctionInfo) ([]byte, error) {
	return generator.GenerateRustCargofile(signatureInfo, functionInfo)
}

func GenerateRustLib(scalefileSchema *scalefile.Schema, functionInfo *FunctionInfo) ([]byte, error) {
	return generator.GenerateRustLib(scalefileSchema, functionInfo)
}

func init() {
	generator = New()
}

type Generator struct {
	template *template.Template
}

func New() *Generator {
	return &Generator{
		template: template.Must(template.New("main").ParseFS(templates.FS, "*rs.templ")),
	}
}

func (g *Generator) GenerateRustCargofile(signatureInfo *SignatureInfo, functionInfo *FunctionInfo) ([]byte, error) {
	signatureInfo.normalize()
	functionInfo.normalize()

	buf := new(bytes.Buffer)
	err := g.template.ExecuteTemplate(buf, "cargo.rs.templ", map[string]interface{}{
		"signature": signatureInfo,
		"function":  functionInfo,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (g *Generator) GenerateRustLib(scalefileSchema *scalefile.Schema, functionInfo *FunctionInfo) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := g.template.ExecuteTemplate(buf, "lib.rs.templ", map[string]interface{}{
		"generatorVersion": strings.TrimPrefix(version.Version(), "v"),
		"scalefileSchema":  scalefileSchema,
		"function":         functionInfo,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
