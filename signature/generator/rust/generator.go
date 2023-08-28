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
	polyglotUtils "github.com/loopholelabs/polyglot/utils"
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/signature/generator/rust/templates"
	"github.com/loopholelabs/scale/signature/generator/utils"
	"strings"
	"text/template"
)

const (
	defaultPackageName = "types"
)

var generator *Generator

// Generate generates the types for the signature
func Generate(schema *signature.Schema, packageName string, scaleVersion string) ([]byte, error) {
	return generator.Generate(schema, packageName, scaleVersion)
}

// GenerateCargofile generates the cargo.toml file for the signature
func GenerateCargofile(packageName string, packageVersion string, scaleVersion string, polyglotVersion string) ([]byte, error) {
	return generator.GenerateCargofile(packageName, packageVersion, scaleVersion, polyglotVersion)
}

func GenerateGuest(schema *signature.Schema, hash string, packageName string, scaleVersion string) ([]byte, error) {
	return generator.GenerateGuest(schema, hash, packageName, scaleVersion)
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
	templ *template.Template
}

// New creates a new rust generator
func New() (*Generator, error) {
	templ, err := template.New("").Funcs(templateFunctions()).ParseFS(templates.FS, "*.rs.templ")
	if err != nil {
		return nil, err
	}

	return &Generator{
		templ: templ,
	}, nil
}

// Generate generates the rust code for the given schema
func (g *Generator) Generate(schema *signature.Schema, packageName string, scaleVersion string) ([]byte, error) {
	if packageName == "" {
		packageName = defaultPackageName
	}

	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "types.rs.templ", map[string]any{
		"schema":  schema,
		"version": scaleVersion,
		"package": packageName,
	})
	if err != nil {
		return nil, err
	}

	// TODO: Compile `rustfmt` to a scale fn
	//cmd := exec.Command("rustfmt")
	//cmd.Stdin = bytes.NewReader(buf.Bytes())
	//output, err := cmd.CombinedOutput()
	//if err != nil {
	//	return nil, err
	//}

	return buf.Bytes(), nil
}

// GenerateCargofile generates the cargofile for the signature
func (g *Generator) GenerateCargofile(packageName string, packageVersion string, scaleVersion string, polyglotVersion string) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "cargo.rs.templ", map[string]any{
		"scale_version":    strings.TrimPrefix(scaleVersion, "v"),
		"polyglot_version": strings.TrimPrefix(polyglotVersion, "v"),
		"package":          packageName,
		"version":          packageVersion,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GenerateGuest generates the guest bindings
func (g *Generator) GenerateGuest(schema *signature.Schema, hash string, packageName string, scaleVersion string) ([]byte, error) {
	if packageName == "" {
		packageName = defaultPackageName
	}

	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "guest.rs.templ", map[string]any{
		"schema":  schema,
		"hash":    hash,
		"version": strings.TrimPrefix(scaleVersion, "v"),
		"package": packageName,
	})
	if err != nil {
		return nil, err
	}

	// TODO: Compile `rustfmt` to a scale fn
	//cmd := exec.Command("rustfmt")
	//cmd.Stdin = bytes.NewReader(buf.Bytes())
	//output, err := cmd.CombinedOutput()
	//if err != nil {
	//	return nil, err
	//}

	return buf.Bytes(), nil
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
		"SnakeCase":               polyglotUtils.SnakeCase,
		"Params":                  utils.Params,
	}
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
