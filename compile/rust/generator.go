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
	"github.com/loopholelabs/scale/compile/rust/templates"
	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
	"strings"
	"text/template"
)

var generator *Generator

func GenerateRustCargofile(schema *scalefile.Schema, registry string, signatureVersion string, signaturePath string, functionPath string, dependencies []*scalefunc.Dependency, packageName string, packageVersion string) ([]byte, error) {
	return generator.GenerateRustCargofile(schema, registry, signatureVersion, signaturePath, functionPath, dependencies, packageName, packageVersion)
}

func GenerateRustLib(signature *signature.Schema, schema *scalefile.Schema, version string) ([]byte, error) {
	return generator.GenerateRustLib(signature, schema, version)
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

func (g *Generator) GenerateRustCargofile(schema *scalefile.Schema, registry string, signatureVersion string, signaturePath string, functionPath string, dependencies []*scalefunc.Dependency, packageName string, packageVersion string) ([]byte, error) {
	if !strings.HasPrefix(signaturePath, "/") && !strings.HasPrefix(signaturePath, "./") && !strings.HasPrefix(signaturePath, "../") {
		signaturePath = "./" + signaturePath
	}

	if !strings.HasPrefix(functionPath, "/") && !strings.HasPrefix(functionPath, "./") && !strings.HasPrefix(functionPath, "../") {
		functionPath = "./" + functionPath
	}

	buf := new(bytes.Buffer)
	err := g.template.ExecuteTemplate(buf, "cargo.rs.templ", map[string]interface{}{
		"version":              packageVersion,
		"package":              packageName,
		"dependencies":         dependencies,
		"registry":             registry,
		"signature_dependency": "signature",
		"signature_path":       signaturePath,
		"signature_version":    signatureVersion,
		"function_dependency":  schema.Name,
		"function_path":        functionPath,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (g *Generator) GenerateRustLib(signature *signature.Schema, schema *scalefile.Schema, scaleVersion string) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := g.template.ExecuteTemplate(buf, "lib.rs.templ", map[string]interface{}{
		"version":   scaleVersion,
		"signature": signature,
		"schema":    schema,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
