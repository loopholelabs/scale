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

	"github.com/loopholelabs/scale/compile/golang/templates"
	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
)

const (
	DefaultVersion = "v0.1.0"
)

var generator *Generator

func GenerateGoModfile(schema *scalefile.Schema, signatureImport string, signatureVersion string, functionImport string, dependencies []*scalefunc.Dependency, packageName string) ([]byte, error) {
	return generator.GenerateGoModfile(schema, signatureImport, signatureVersion, functionImport, dependencies, packageName)
}

func GenerateGoMain(signature *signature.Schema, schema *scalefile.Schema, version string) ([]byte, error) {
	return generator.GenerateGoMain(signature, schema, version)
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

func (g *Generator) GenerateGoModfile(schema *scalefile.Schema, signatureImport string, signatureVersion string, functionImport string, dependencies []*scalefunc.Dependency, packageName string) ([]byte, error) {
	if !strings.HasPrefix(signatureImport, "/") && !strings.HasPrefix(signatureImport, "./") && !strings.HasPrefix(signatureImport, "../") {
		signatureImport = "./" + signatureImport
	}

	if !strings.HasPrefix(functionImport, "/") && !strings.HasPrefix(functionImport, "./") && !strings.HasPrefix(functionImport, "../") {
		functionImport = "./" + functionImport
	}

	buf := new(bytes.Buffer)
	err := g.template.ExecuteTemplate(buf, "mod.go.templ", map[string]interface{}{
		"package":                  packageName,
		"dependencies":             dependencies,
		"old_signature_dependency": "signature",
		"old_signature_version":    DefaultVersion,
		"old_function_dependency":  schema.Name,
		"old_function_version":     DefaultVersion,
		"new_signature_dependency": signatureImport,
		"new_signature_version":    signatureVersion,
		"new_function_dependency":  functionImport,
		"new_function_version":     "",
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (g *Generator) GenerateGoMain(signature *signature.Schema, schema *scalefile.Schema, version string) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := g.template.ExecuteTemplate(buf, "main.go.templ", map[string]interface{}{
		"version":   version,
		"signature": signature,
		"schema":    schema,
	})
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}
