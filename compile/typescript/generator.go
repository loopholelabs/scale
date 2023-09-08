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
	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/version"
	"strings"
	"text/template"

	"github.com/loopholelabs/scale/compile/typescript/templates"
)

var generator *Generator

func GenerateTypescriptPackageJSON(packageSchema *scalefile.Schema, signaturePath string, signatureVersion string) ([]byte, error) {
	return generator.GenerateTypescriptPackageJSON(packageSchema, signaturePath, signatureVersion)
}

func GenerateTypescriptIndex(packageSchema *scalefile.Schema, functionPath string) ([]byte, error) {
	return generator.GenerateTypescriptIndex(packageSchema, functionPath)
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

func (g *Generator) GenerateTypescriptPackageJSON(packageSchema *scalefile.Schema, signaturePath string, signatureVersion string) ([]byte, error) {
	if signaturePath != "" && !strings.HasPrefix(signaturePath, "/") && !strings.HasPrefix(signaturePath, "./") && !strings.HasPrefix(signaturePath, "../") {
		signaturePath = "./" + signaturePath
	}

	buf := new(bytes.Buffer)
	err := g.template.ExecuteTemplate(buf, "packagejson.ts.templ", map[string]interface{}{
		"package_schema":    packageSchema,
		"signature_path":    signaturePath,
		"signature_version": signatureVersion,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (g *Generator) GenerateTypescriptIndex(packageSchema *scalefile.Schema, functionPath string) ([]byte, error) {
	if !strings.HasPrefix(functionPath, "/") && !strings.HasPrefix(functionPath, "./") && !strings.HasPrefix(functionPath, "../") {
		functionPath = "./" + functionPath
	}
	buf := new(bytes.Buffer)
	err := g.template.ExecuteTemplate(buf, "index.ts.templ", map[string]interface{}{
		"generator_version": strings.TrimPrefix(version.Version(), "v"),
		"package_schema":    packageSchema,
		"function_path":     functionPath,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
