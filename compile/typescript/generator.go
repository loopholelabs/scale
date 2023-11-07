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

package typescript

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/version"

	"github.com/loopholelabs/scale/compile/typescript/templates"
)

var generator *Generator

func GenerateTypescriptPackageJSON(signatureInfo *SignatureInfo) ([]byte, error) {
	return generator.GenerateTypescriptPackageJSON(signatureInfo)
}

func GenerateTypescriptIndex(packageSchema *scalefile.Schema, function *FunctionInfo) ([]byte, error) {
	return generator.GenerateTypescriptIndex(packageSchema, function)
}

func init() {
	generator = New()
}

type Generator struct {
	template *template.Template
}

func New() *Generator {
	return &Generator{
		template: template.Must(template.New("main").ParseFS(templates.FS, "*ts.templ")),
	}
}

func (g *Generator) GenerateTypescriptPackageJSON(signatureInfo *SignatureInfo) ([]byte, error) {
	signatureInfo.normalize()

	buf := new(bytes.Buffer)
	err := g.template.ExecuteTemplate(buf, "packagejson.ts.templ", map[string]interface{}{
		"signature": signatureInfo,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (g *Generator) GenerateTypescriptIndex(scalefileSchema *scalefile.Schema, function *FunctionInfo) ([]byte, error) {
	function.normalize()

	buf := new(bytes.Buffer)
	err := g.template.ExecuteTemplate(buf, "index.ts.templ", map[string]interface{}{
		"generatorVersion": strings.TrimPrefix(version.Version(), "v"),
		"scalefileSchema":  scalefileSchema,
		"function":         function,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
