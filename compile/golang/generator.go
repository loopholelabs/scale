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
	"bytes"
	"go/format"
	"strings"
	"text/template"

	"github.com/loopholelabs/scale/version"

	"github.com/loopholelabs/scale/compile/golang/templates"
	"github.com/loopholelabs/scale/extension"
	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/signature"
)

var generator *Generator

func GenerateGoModfile(packageSchema *scalefile.Schema, signatureImport string, signatureVersion string, functionImport string, extensions []extension.ExtensionInfo) ([]byte, error) {
	return generator.GenerateGoModfile(packageSchema, signatureImport, signatureVersion, functionImport, extensions)
}

func GenerateGoMain(packageSchema *scalefile.Schema, signatureSchema *signature.Schema) ([]byte, error) {
	return generator.GenerateGoMain(packageSchema, signatureSchema)
}

func init() {
	generator = New()
}

type Generator struct {
	template *template.Template
}

func New() *Generator {
	return &Generator{
		template: template.Must(template.New("main").ParseFS(templates.FS, "*go.templ")),
	}
}

func (g *Generator) GenerateGoModfile(packageSchema *scalefile.Schema, signatureImport string, signatureVersion string, functionImport string, extensions []extension.ExtensionInfo) ([]byte, error) {
	if signatureVersion == "" && !strings.HasPrefix(signatureImport, "/") && !strings.HasPrefix(signatureImport, "./") && !strings.HasPrefix(signatureImport, "../") {
		signatureImport = "./" + signatureImport
	}

	if !strings.HasPrefix(functionImport, "/") && !strings.HasPrefix(functionImport, "./") && !strings.HasPrefix(functionImport, "../") {
		functionImport = "./" + functionImport
	}

	buf := new(bytes.Buffer)
	err := g.template.ExecuteTemplate(buf, "mod.go.templ", map[string]interface{}{
		"package_schema":    packageSchema,
		"extensions":        extensions,
		"signature_import":  signatureImport,
		"signature_version": signatureVersion,
		"function_import":   functionImport,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (g *Generator) GenerateGoMain(packageSchema *scalefile.Schema, signatureSchema *signature.Schema) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := g.template.ExecuteTemplate(buf, "main.go.templ", map[string]interface{}{
		"generator_version": version.Version(),
		"signature_schema":  signatureSchema,
		"package_schema":    packageSchema,
	})
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}
