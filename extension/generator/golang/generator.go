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

	polyglotVersion "github.com/loopholelabs/polyglot/version"
	interfacesVersion "github.com/loopholelabs/scale-extension-interfaces/version"
	"github.com/loopholelabs/scale/extension"

	scaleVersion "github.com/loopholelabs/scale/version"

	"github.com/loopholelabs/scale/extension/generator/golang/templates"
	"github.com/loopholelabs/scale/signature/generator/utils"
)

const (
	defaultPackageName = "types"
)

var generator *Generator

func GenerateTypes(schema *extension.Schema, packageName string) ([]byte, error) {
	return generator.GenerateTypes(schema, packageName)
}

func GenerateInterfaces(schema *extension.Schema, packageName string, version string) ([]byte, error) {
	return generator.GenerateInterfaces(schema, packageName, version)
}

func GenerateGuest(schema *extension.Schema, extensionHash string, packageName string, version string) ([]byte, error) {
	return generator.GenerateGuest(schema, extensionHash, packageName, version)
}

func GenerateModfile(packageName string) ([]byte, error) {
	return generator.GenerateModfile(packageName)
}

func GenerateHost(schema *extension.Schema, extensionHash string, packageName string, version string) ([]byte, error) {
	return generator.GenerateHost(schema, extensionHash, packageName, version)
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
	templ, err := template.New("").Funcs(templateFunctions()).ParseFS(templates.FS, "*.go.templ")
	if err != nil {
		return nil, err
	}

	return &Generator{
		templ: templ,
	}, nil
}

// Generate generates the go code
func (g *Generator) GenerateTypes(schema *extension.Schema, packageName string) ([]byte, error) {
	if packageName == "" {
		packageName = defaultPackageName
	}

	ext, err := schema.CloneWithDisabledAccessorsValidatorsAndModifiers()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = g.templ.ExecuteTemplate(buf, "types.go.templ", map[string]any{
		"signature_schema":  ext,
		"generator_version": scaleVersion.Version(),
		"package_name":      packageName,
	})
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

func (g *Generator) GenerateInterfaces(schema *extension.Schema, packageName string, version string) ([]byte, error) {
	if packageName == "" {
		packageName = defaultPackageName
	}

	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "interfaces.go.templ", map[string]any{
		"schema":  schema,
		"version": version,
		"package": packageName,
	})
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

// GenerateGuest generates the guest bindings
func (g *Generator) GenerateGuest(schema *extension.Schema, schemaHash string, packageName string, version string) ([]byte, error) {
	if packageName == "" {
		packageName = defaultPackageName
	}

	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "guest.go.templ", map[string]any{
		"extension_schema": schema,
		"extension_hash":   schemaHash,
		"version":          version,
		"package":          packageName,
	})
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

// GenerateModfile generates the modfile for the signature
func (g *Generator) GenerateModfile(packageImportPath string) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "mod.go.templ", map[string]any{
		"polyglot_version":                   polyglotVersion.Version(),
		"scale_extension_interfaces_version": interfacesVersion.Version(),
		"package_import_path":                packageImportPath,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GenerateHost generates the host bindings
func (g *Generator) GenerateHost(schema *extension.Schema, schemaHash string, packageName string, version string) ([]byte, error) {
	if packageName == "" {
		packageName = defaultPackageName
	}

	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "host.go.templ", map[string]any{
		"extension_schema": schema,
		"extension_hash":   schemaHash,
		"version":          version,
		"package":          packageName,
	})
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

func templateFunctions() template.FuncMap {
	return template.FuncMap{
		"IsInterface":             isInterface,
		"Primitive":               primitive,
		"IsPrimitive":             extension.ValidPrimitiveType,
		"PolyglotPrimitive":       polyglotPrimitive,
		"PolyglotPrimitiveEncode": polyglotPrimitiveEncode,
		"PolyglotPrimitiveDecode": polyglotPrimitiveDecode,
		"Deref":                   func(i *bool) bool { return *i },
		"LowerFirst":              func(s string) string { return string(s[0]+32) + s[1:] },
		"Params":                  utils.Params,
	}
}

func isInterface(schema *extension.Schema, s string) bool {
	for _, i := range schema.Interfaces {
		if i.Name == s {
			return true
		}
	}
	return false
}

func primitive(t string) string {
	switch t {
	case "string", "int32", "int64", "uint32", "uint64", "float32", "float64", "bool":
		return t
	case "bytes":
		return "[]byte"
	default:
		return ""
	}
}

func polyglotPrimitive(t string) string {
	switch t {
	case "string":
		return "polyglot.StringKind"
	case "int32":
		return "polyglot.Int32Kind"
	case "int64":
		return "polyglot.Int64Kind"
	case "uint32":
		return "polyglot.Uint32Kind"
	case "uint64":
		return "polyglot.Uint64Kind"
	case "float32":
		return "polyglot.Float32Kind"
	case "float64":
		return "polyglot.Float64Kind"
	case "bool":
		return "polyglot.BoolKind"
	case "bytes":
		return "polyglot.BytesKind"
	default:
		return "polyglot.AnyKind"
	}
}

func polyglotPrimitiveEncode(t string) string {
	switch t {
	case "string":
		return "String"
	case "int32":
		return "Int32"
	case "int64":
		return "Int64"
	case "uint32":
		return "Uint32"
	case "uint64":
		return "Uint64"
	case "float32":
		return "Float32"
	case "float64":
		return "Float64"
	case "bool":
		return "Bool"
	case "bytes":
		return "Bytes"
	default:
		return ""
	}
}

func polyglotPrimitiveDecode(t string) string {
	switch t {
	case "string":
		return "String"
	case "int32":
		return "Int32"
	case "int64":
		return "Int64"
	case "uint32":
		return "Uint32"
	case "uint64":
		return "Uint64"
	case "float32":
		return "Float32"
	case "float64":
		return "Float64"
	case "bool":
		return "Bool"
	case "bytes":
		return "Bytes"
	default:
		return ""
	}
}
