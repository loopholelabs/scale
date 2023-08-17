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
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/signature/generator/golang/templates"
	"go/format"
	"text/template"

	"github.com/loopholelabs/scale/signature/generator/utils"
)

const (
	defaultPackageName = "types"
)

var generator *Generator

// Generate generates the types for the signature
func Generate(schema *signature.Schema, packageName string, version string) ([]byte, error) {
	return generator.Generate(schema, packageName, version)
}

// GenerateModfile generates the go.mod file for the signature
func GenerateModfile(packageName string, polyglotVersion string) ([]byte, error) {
	return generator.GenerateModfile(packageName, polyglotVersion)
}

func GenerateGuest(schema *signature.Schema, hash string, packageName string, version string) ([]byte, error) {
	return generator.GenerateGuest(schema, hash, packageName, version)
}

func GenerateHost(schema *signature.Schema, hash string, packageName string, version string) ([]byte, error) {
	return generator.GenerateHost(schema, hash, packageName, version)
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

// Generate generates the types for the signature
func (g *Generator) Generate(schema *signature.Schema, packageName string, version string) ([]byte, error) {
	if packageName == "" {
		packageName = defaultPackageName
	}

	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "types.go.templ", map[string]any{
		"schema":  schema,
		"version": version,
		"package": packageName,
	})
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

// GenerateModfile generates the modfile for the signature
func (g *Generator) GenerateModfile(packageName string, polyglotVersion string) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "mod.go.templ", map[string]any{
		"polyglot_version": polyglotVersion,
		"package":          packageName,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GenerateGuest generates the guest bindings
func (g *Generator) GenerateGuest(schema *signature.Schema, hash string, packageName string, version string) ([]byte, error) {
	if packageName == "" {
		packageName = defaultPackageName
	}

	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "guest.go.templ", map[string]any{
		"schema":  schema,
		"hash":    hash,
		"version": version,
		"package": packageName,
	})
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

// GenerateHost generates the host bindings
//
// Note: the given schema should already be normalized, validated, and modified to have its accessors and validators disabled
func (g *Generator) GenerateHost(schema *signature.Schema, hash string, packageName string, version string) ([]byte, error) {
	if packageName == "" {
		packageName = defaultPackageName
	}

	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "host.go.templ", map[string]any{
		"schema":  schema,
		"hash":    hash,
		"version": version,
		"package": packageName,
	})
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

func templateFunctions() template.FuncMap {
	return template.FuncMap{
		"Primitive":               primitive,
		"IsPrimitive":             signature.ValidPrimitiveType,
		"PolyglotPrimitive":       polyglotPrimitive,
		"PolyglotPrimitiveEncode": polyglotPrimitiveEncode,
		"PolyglotPrimitiveDecode": polyglotPrimitiveDecode,
		"Deref":                   func(i *bool) bool { return *i },
		"LowerFirst":              func(s string) string { return string(s[0]+32) + s[1:] },
		"Params":                  utils.Params,
	}
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
