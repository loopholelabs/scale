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
	"github.com/loopholelabs/scale/signature"
	"strings"
	"text/template"

	"github.com/loopholelabs/scale/signature/generator/templates"
	"github.com/loopholelabs/scale/signature/generator/utils"
)

const (
	defaultPackageName = "types"
)

// Generator is the typescript generator
type Generator struct {
	templ *template.Template
}

// New creates a new typescript generator
func New() (*Generator, error) {
	templ, err := template.New("").Funcs(templateFunctions()).ParseFS(templates.FS, "*.ts.templ")
	if err != nil {
		return nil, err
	}

	return &Generator{
		templ: templ,
	}, nil
}

// Generate generates the typescript code
func (g *Generator) Generate(schema *signature.Schema, packageName string, version string) ([]byte, error) {
	if packageName == "" {
		packageName = defaultPackageName
	}

	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "types.ts.templ", map[string]any{
		"schema":  schema,
		"version": version,
		"package": packageName,
	})
	if err != nil {
		return nil, err
	}

	return []byte(formatTS(buf.String())), nil
}

func templateFunctions() template.FuncMap {
	return template.FuncMap{
		"Primitive":               primitive,
		"IsPrimitive":             signature.ValidPrimitiveType,
		"PolyglotPrimitive":       polyglotPrimitive,
		"PolyglotPrimitiveEncode": polyglotPrimitiveEncode,
		"PolyglotPrimitiveDecode": polyglotPrimitiveDecode,
		"Deref":                   func(i *bool) bool { return *i },
		"CamelCase":               utils.CamelCase,
		"Params":                  utils.Params,
		"Constructor":             constructor,
	}
}

func primitive(t string) string {
	switch t {
	case "string":
		return "string"
	case "int32":
		return "number"
	case "int64":
		return "bigint"
	case "uint32":
		return "number"
	case "uint64":
		return "bigint"
	case "float32":
		return "number"
	case "float64":
		return "number"
	case "bool":
		return "boolean"
	case "bytes":
		return "Uint8Array"
	default:
		return t
	}
}

func constructor(t string) string {
	switch t {
	case "string":
		return "String"
	case "int32":
		return "Number"
	case "int64":
		return "BigInt"
	case "uint32":
		return "Number"
	case "uint64":
		return "BigInt"
	case "float32":
		return "Number"
	case "float64":
		return "Number"
	case "bool":
		return "Boolean"
	case "bytes":
		return "Uint8Array"
	default:
		return t
	}
}

func polyglotPrimitive(t string) string {
	switch t {
	case "string":
		return "Kind.String"
	case "int32":
		return "Kind.Int32"
	case "int64":
		return "Kind.Int64"
	case "uint32":
		return "Kind.Uint32"
	case "uint64":
		return "Kind.Uint64"
	case "float32":
		return "Kind.Float32"
	case "float64":
		return "Kind.Float64"
	case "bool":
		return "Kind.Boolean"
	case "bytes":
		return "Kind.Uint8Array"
	default:
		return "Kind.Any"
	}
}

func polyglotPrimitiveEncode(t string) string {
	switch t {
	case "string":
		return "string"
	case "int32":
		return "int32"
	case "int64":
		return "int64"
	case "uint32":
		return "uint32"
	case "uint64":
		return "uint64"
	case "float32":
		return "float32"
	case "float64":
		return "float64"
	case "bool":
		return "boolean"
	case "bytes":
		return "uint8Array"
	default:
		return t
	}
}

func polyglotPrimitiveDecode(t string) string {
	switch t {
	case "string":
		return "string"
	case "int32":
		return "int32"
	case "int64":
		return "int64"
	case "uint32":
		return "uint32"
	case "uint64":
		return "uint64"
	case "float32":
		return "float32"
	case "float64":
		return "float64"
	case "bool":
		return "boolean"
	case "bytes":
		return "uint8Array"
	default:
		return ""
	}
}

func formatTS(code string) string {
	var output strings.Builder
	indentLevel := 0
	lastLineEmpty := false
	lastLineOpenBrace := false
	for _, line := range strings.Split(code, "\n") {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			// Allow empty lines between classes and class members, but only 1 empty line not more.
			if indentLevel > 1 || lastLineEmpty || lastLineOpenBrace {
				continue
			} else {
				output.WriteRune('\n')
			}
			lastLineEmpty = true
		} else {
			if strings.HasPrefix(trimmedLine, "}") {
				indentLevel--
			}
			output.WriteString(strings.Repeat("  ", indentLevel))
			output.WriteString(trimmedLine)
			if strings.HasSuffix(trimmedLine, "{") {
				lastLineOpenBrace = true
				indentLevel++
			} else {
				lastLineOpenBrace = false
			}
			output.WriteRune('\n')
			lastLineEmpty = false
		}
	}
	return output.String()
}
