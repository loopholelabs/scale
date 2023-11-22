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
	"context"
	"strings"
	"text/template"

	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/signature/generator/rust"

	interfacesVersion "github.com/loopholelabs/scale-extension-interfaces/version"

	polyglotVersion "github.com/loopholelabs/polyglot/version"

	scaleVersion "github.com/loopholelabs/scale/version"

	polyglotUtils "github.com/loopholelabs/polyglot/utils"

	"github.com/loopholelabs/scale/extension"
	"github.com/loopholelabs/scale/extension/generator/rust/templates"
	"github.com/loopholelabs/scale/signature/generator/rust/format"
	"github.com/loopholelabs/scale/signature/generator/utils"
)

const (
	defaultPackageName = "types"
)

var generator *Generator

// GenerateTypes generates the types for the extension
func GenerateTypes(extensionSchema *extension.Schema, packageName string) ([]byte, error) {
	return generator.GenerateTypes(extensionSchema, packageName)
}

// GenerateCargofile generates the cargo.toml file for the extension
func GenerateCargofile(packageName string, packageVersion string) ([]byte, error) {
	return generator.GenerateCargofile(packageName, packageVersion)
}

func GenerateGuest(extensionSchema *extension.Schema, extensionHash string, packageName string) ([]byte, error) {
	return generator.GenerateGuest(extensionSchema, extensionHash, packageName)
}

func init() {
	var err error
	generator, err = New()
	if err != nil {
		panic(err)
	}
}

// Generator is the rust generator
type Generator struct {
	templ     *template.Template
	signature *rust.Generator
	formatter *format.Formatter
}

// New creates a new rust generator
func New() (*Generator, error) {
	templ, err := template.New("").Funcs(templateFunctions()).ParseFS(templates.FS, "*.rs.templ")
	if err != nil {
		return nil, err
	}

	sig, err := rust.New()
	if err != nil {
		return nil, err
	}

	formatter, err := format.New()
	if err != nil {
		return nil, err
	}

	return &Generator{
		templ:     templ,
		signature: sig,
		formatter: formatter,
	}, nil
}

// GenerateTypes generates the types for the extension
func (g *Generator) GenerateTypes(extensionSchema *extension.Schema, packageName string) ([]byte, error) {
	signatureSchema := &signature.Schema{
		Version: extensionSchema.Version,
		Enums:   extensionSchema.Enums,
		Models:  extensionSchema.Models,
	}

	signatureSchema.SetHasLengthValidator(extensionSchema.HasLengthValidator())
	signatureSchema.SetHasCaseModifier(extensionSchema.HasCaseModifier())
	signatureSchema.SetHasLimitValidator(extensionSchema.HasLimitValidator())
	signatureSchema.SetHasRegexValidator(extensionSchema.HasRegexValidator())

	return g.signature.GenerateTypes(signatureSchema, packageName)
}

// GenerateCargofile generates the cargofile for the extension
func (g *Generator) GenerateCargofile(packageName string, packageVersion string) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "cargo.rs.templ", map[string]any{
		"polyglot_version":                   strings.TrimPrefix(polyglotVersion.Version(), "v"),
		"scale_signature_interfaces_version": strings.TrimPrefix(interfacesVersion.Version(), "v"),
		"package_name":                       packageName,
		"package_version":                    strings.TrimPrefix(packageVersion, "v"),
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GenerateGuest generates the guest bindings
func (g *Generator) GenerateGuest(extensionSchema *extension.Schema, extensionHash string, packageName string) ([]byte, error) {
	if packageName == "" {
		packageName = defaultPackageName
	}

	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "guest.rs.templ", map[string]any{
		"extension_schema": extensionSchema,
		"extension_hash":   extensionHash,
	})
	if err != nil {
		return nil, err
	}

	formatted, err := g.formatter.Format(context.Background(), buf.String())
	if err != nil {
		return nil, err
	}
	buf.Reset()
	err = g.templ.ExecuteTemplate(buf, "header.rs.templ", map[string]any{
		"generator_version": strings.TrimPrefix(scaleVersion.Version(), "v"),
		"package_name":      packageName,
	})
	if err != nil {
		return nil, err
	}
	return []byte(buf.String() + "\n\n" + formatted), nil
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
		"SnakeCase":               polyglotUtils.SnakeCase,
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
	case "string":
		return "String"
	case "int32":
		return "i32"
	case "int64":
		return "i64"
	case "uint32":
		return "u32"
	case "uint64":
		return "u64"
	case "float32":
		return "f32"
	case "float64":
		return "f64"
	case "bool":
		return "bool"
	case "bytes":
		return "Vec<u8>"
	default:
		return t
	}
}

func polyglotPrimitive(t string) string {
	switch t {
	case "string":
		return "Kind::String"
	case "int32":
		return "Kind::I32"
	case "int64":
		return "Kind::I64"
	case "uint32":
		return "Kind::U32"
	case "uint64":
		return "Kind::U64"
	case "float32":
		return "Kind::F32"
	case "float64":
		return "Kind::F64"
	case "bool":
		return "Kind::Bool"
	case "bytes":
		return "Kind::Bytes"
	default:
		return "Kind::Any"
	}
}

func polyglotPrimitiveEncode(t string) string {
	switch t {
	case "string":
		return "encode_string"
	case "int32":
		return "encode_i32"
	case "int64":
		return "encode_i64"
	case "uint32":
		return "encode_u32"
	case "uint64":
		return "encode_u64"
	case "float32":
		return "encode_f32"
	case "float64":
		return "encode_f64"
	case "bool":
		return "encode_bool"
	case "bytes":
		return "encode_bytes"
	default:
		return t
	}
}

func polyglotPrimitiveDecode(t string) string {
	switch t {
	case "string":
		return "decode_string"
	case "int32":
		return "decode_i32"
	case "int64":
		return "decode_i64"
	case "uint32":
		return "decode_u32"
	case "uint64":
		return "decode_u64"
	case "float32":
		return "decode_f32"
	case "float64":
		return "decode_f64"
	case "bool":
		return "decode_bool"
	case "bytes":
		return "decode_bytes"
	default:
		return ""
	}
}
