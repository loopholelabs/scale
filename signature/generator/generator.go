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

package generator

import (
	"fmt"
	"github.com/loopholelabs/polyglot-go/pkg/generator"
	"github.com/loopholelabs/polyglot-go/pkg/utils"
	polyglotTemplates "github.com/loopholelabs/polyglot-go/templates"
	"github.com/loopholelabs/scale/signature/templates"
	"github.com/loopholelabs/scale/signature/templates/override"
	"github.com/yalue/merged_fs"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
	"io"
	"text/template"
)

const version = "v0.0.1"

type Generator struct {
	options           *protogen.Options
	polyglotTemplate  *template.Template
	generatorTemplate *template.Template
	CustomFields      func() string
	CustomEncode      func() string
	CustomDecode      func() string
}

func New() (g *Generator) {
	fs := merged_fs.NewMergedFS(override.FS, polyglotTemplates.FS)
	polyglotTemplate := template.Must(template.New("main").Funcs(template.FuncMap{
		"CamelCase":          utils.CamelCaseFullName,
		"CamelCaseName":      utils.CamelCaseName,
		"MakeIterable":       utils.MakeIterable,
		"Counter":            utils.Counter,
		"FirstLowerCase":     utils.FirstLowerCase,
		"FirstLowerCaseName": utils.FirstLowerCaseName,
		"FindValue":          generator.FindValue,
		"GetKind":            generator.GetKind,
		"GetLUTEncoder":      generator.GetLUTEncoder,
		"GetLUTDecoder":      generator.GetLUTDecoder,
		"GetEncodingFields":  generator.GetEncodingFields,
		"GetDecodingFields":  generator.GetDecodingFields,
		"GetKindLUT":         generator.GetKindLUT,
		"CustomFields": func() string {
			return g.CustomFields()
		},
		"CustomEncode": func() string {
			return g.CustomEncode()
		},
		"CustomDecode": func() string {
			return g.CustomDecode()
		},
	}).ParseFS(fs, "*"))

	generatorTemplate := template.Must(template.New("main").Funcs(template.FuncMap{
		"CamelCase": utils.CamelCase,
	}).ParseFS(templates.FS, "*"))

	g = &Generator{
		options: &protogen.Options{
			ParamFunc:         func(name string, value string) error { return nil },
			ImportRewriteFunc: func(path protogen.GoImportPath) protogen.GoImportPath { return path },
		},
		polyglotTemplate:  polyglotTemplate,
		generatorTemplate: generatorTemplate,
		CustomEncode:      func() string { return "" },
		CustomDecode:      func() string { return "" },
		CustomFields:      func() string { return "" },
	}
	return g
}

func (g *Generator) UnmarshalRequest(buf []byte) (*pluginpb.CodeGeneratorRequest, error) {
	req := new(pluginpb.CodeGeneratorRequest)
	return req, proto.Unmarshal(buf, req)
}

func (g *Generator) MarshalResponse(res *pluginpb.CodeGeneratorResponse) ([]byte, error) {
	return proto.Marshal(res)
}

func (g *Generator) Generate(req *pluginpb.CodeGeneratorRequest) (res *pluginpb.CodeGeneratorResponse, err error) {
	plugin, err := g.options.New(req)
	if err != nil {
		return nil, err
	}

	for _, f := range plugin.Files {
		if !f.Generate {
			continue
		}
		genFile := plugin.NewGeneratedFile(fmt.Sprintf("%s.signature.go", f.GeneratedFilenamePrefix), f.GoImportPath)
		hostFile := plugin.NewGeneratedFile("host.go", f.GoImportPath)
		packageFile := plugin.NewGeneratedFile(fmt.Sprintf("%s.go", f.GeneratedFilenamePrefix), f.GoImportPath)

		packageName := string(f.Desc.Package().Name())
		if packageName == "" {
			packageName = string(f.GoPackageName)
		}

		err = g.ExecuteHostGeneratorTemplate(hostFile, packageName, f.Desc.Path())
		if err != nil {
			return nil, err
		}

		err = g.ExecutePackageGeneratorTemplate(packageFile, packageName, f.Desc.Path())
		if err != nil {
			return nil, err
		}

		err = g.ExecutePolyglotTemplate(genFile, f, packageName, true)
		if err != nil {
			return nil, err
		}
	}

	return plugin.Response(), nil
}

func (g *Generator) ExecutePolyglotTemplate(writer io.Writer, protoFile *protogen.File, packageName string, header bool) error {
	return g.polyglotTemplate.ExecuteTemplate(writer, "base.templ", map[string]interface{}{
		"pluginVersion":   version,
		"sourcePath":      protoFile.Desc.Path(),
		"package":         packageName,
		"requiredImports": generator.RequiredImports,
		"enums":           protoFile.Desc.Enums(),
		"messages":        protoFile.Desc.Messages(),
		"header":          header,
	})
}

func (g *Generator) ExecuteProtoGeneratorTemplate(writer io.Writer, packageName string) error {
	return g.generatorTemplate.ExecuteTemplate(writer, "package.proto.templ", map[string]interface{}{
		"package": packageName,
	})
}

func (g *Generator) ExecuteHostGeneratorTemplate(writer io.Writer, packageName string, source string) error {
	return g.generatorTemplate.ExecuteTemplate(writer, "host.go.templ", map[string]interface{}{
		"package":       packageName,
		"pluginVersion": version,
		"sourcePath":    source,
	})
}

func (g *Generator) ExecutePackageGeneratorTemplate(writer io.Writer, packageName string, source string) error {
	return g.generatorTemplate.ExecuteTemplate(writer, "package.go.templ", map[string]interface{}{
		"package":       packageName,
		"pluginVersion": version,
		"sourcePath":    source,
	})
}

func (g *Generator) ExecuteSignatureGeneratorTemplate(writer io.Writer, packageName string, packagePath string) error {
	return g.generatorTemplate.ExecuteTemplate(writer, "signature.go.templ", map[string]interface{}{
		"package":       packageName,
		"path":          packagePath,
		"pluginVersion": version,
	})
}
