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
	"text/template"

	"github.com/loopholelabs/scale/signature"

	"github.com/loopholelabs/scale/version"

	"github.com/loopholelabs/scale/compile/golang/templates"
	"github.com/loopholelabs/scale/scalefile"
)

var generator *Generator

type GoModReplacement struct {
	Name string
	Path string
}

func GenerateGoModfile(signatureInfo *SignatureInfo, functionInfo *FunctionInfo, replacements []GoModReplacement) ([]byte, error) {
	return generator.GenerateGoModfile(signatureInfo, functionInfo, replacements)
}

func GenerateGoMain(scalefileSchema *scalefile.Schema, signatureSchema *signature.Schema, functionInfo *FunctionInfo) ([]byte, error) {
	return generator.GenerateGoMain(scalefileSchema, signatureSchema, functionInfo)
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

func (g *Generator) GenerateGoModfile(signatureInfo *SignatureInfo, functionInfo *FunctionInfo, replacements []GoModReplacement) ([]byte, error) {
	signatureInfo.normalize()
	functionInfo.normalize()

	buf := new(bytes.Buffer)
	err := g.template.ExecuteTemplate(buf, "mod.go.templ", map[string]interface{}{
		"function":     functionInfo,
		"signature":    signatureInfo,
		"replacements": replacements,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (g *Generator) GenerateGoMain(scalefileSchema *scalefile.Schema, signatureSchema *signature.Schema, functionInfo *FunctionInfo) ([]byte, error) {
	functionInfo.normalize()

	buf := new(bytes.Buffer)
	err := g.template.ExecuteTemplate(buf, "main.go.templ", map[string]interface{}{
		"generatorVersion": version.Version(),
		"scalefileSchema":  scalefileSchema,
		"signatureSchema":  signatureSchema,
		"function":         functionInfo,
	})
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}
