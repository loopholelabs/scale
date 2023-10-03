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
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/evanw/esbuild/pkg/api"
	polyglotVersion "github.com/loopholelabs/polyglot/version"

	interfacesVersion "github.com/loopholelabs/scale-extension-interfaces/version"
	scaleVersion "github.com/loopholelabs/scale/version"

	"github.com/loopholelabs/scale/extension"
	"github.com/loopholelabs/scale/extension/generator/typescript/templates"

	"github.com/loopholelabs/scale/signature/generator/utils"
)

const (
	defaultPackageName = "types"
	tsConfig           = `
{
  "compilerOptions": {
    "target": "es2020",
    "module": "commonjs",
    "esModuleInterop": true,
    "forceConsistentCasingInFileNames": true,
    "strict": true,
    "skipLibCheck": true,
    "resolveJsonModule": true,
    "sourceMap": true,
    "paths": {
      "signature": ["./"]
    },
    "types": ["node"]
  },
}`
)

var generator *Generator

type Transpiled struct {
	Typescript  []byte
	Javascript  []byte
	SourceMap   []byte
	Declaration []byte
}

// GenerateTypes generates the types for the extension
func GenerateTypes(extensionSchema *extension.Schema, packageName string) ([]byte, error) {
	return generator.GenerateTypes(extensionSchema, packageName)
}

// GenerateTypesTranspiled generates the types for the extension and transpiles it to javascript
func GenerateTypesTranspiled(extensionSchema *extension.Schema, packageName string, sourceName string) (*Transpiled, error) {
	typescriptSource, err := generator.GenerateTypes(extensionSchema, packageName)
	if err != nil {
		return nil, err
	}
	return generator.GenerateTypesTranspiled(extensionSchema, packageName, sourceName, string(typescriptSource))
}

// GeneratePackageJSON generates the package.json file for the extension
func GeneratePackageJSON(packageName string, packageVersion string) ([]byte, error) {
	return generator.GeneratePackageJSON(packageName, packageVersion)
}

// GenerateGuest generates the guest bindings for the extension
func GenerateGuest(extensionSchema *extension.Schema, extensionHash string, packageName string) ([]byte, error) {
	return generator.GenerateGuest(extensionSchema, extensionHash, packageName)
}

// GenerateGuestTranspiled generates the guest bindings and transpiles it to javascript
func GenerateGuestTranspiled(extensionSchema *extension.Schema, extensionHash string, packageName string, sourceName string) (*Transpiled, error) {
	typescriptSource, err := generator.GenerateGuest(extensionSchema, extensionHash, packageName)
	if err != nil {
		return nil, err
	}
	return generator.GenerateGuestTranspiled(extensionSchema, packageName, sourceName, string(typescriptSource))
}

// GenerateHost generates the host bindings for the extension
//
// Note: the given schema should already be normalized, validated, and modified to have its accessors and validators disabled
func GenerateHost(extensionSchema *extension.Schema, extensionHash string, packageName string) ([]byte, error) {
	return generator.GenerateHost(extensionSchema, extensionHash, packageName)
}

// GenerateHostTranspiled generates the host bindings and transpiles it to javascript
//
// Note: the given schema should already be normalized, validated, and modified to have its accessors and validators disabled
func GenerateHostTranspiled(extensionSchema *extension.Schema, extensionHash string, packageName string, sourceName string) (*Transpiled, error) {
	typescriptSource, err := generator.GenerateHost(extensionSchema, extensionHash, packageName)
	if err != nil {
		return nil, err
	}
	return generator.GenerateHostTranspiled(extensionSchema, packageName, sourceName, string(typescriptSource))
}

func init() {
	var err error
	generator, err = New()
	if err != nil {
		panic(err)
	}
}

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

// GenerateTypes generates the types for the extension
//
// This is not transpiled to javascript and does not include source maps or type definitions
func (g *Generator) GenerateTypes(extensionSchema *extension.Schema, packageName string) ([]byte, error) {
	if packageName == "" {
		packageName = defaultPackageName
	}

	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "types.ts.templ", map[string]any{
		"extension_schema":  extensionSchema,
		"generator_version": strings.TrimPrefix(scaleVersion.Version(), "v"),
		"package_name":      packageName,
	})
	if err != nil {
		return nil, err
	}

	return []byte(formatTS(buf.String())), nil
}

// GenerateTypesTranspiled takes the typescript source for the generated types and transpiles it to javascript
func (g *Generator) GenerateTypesTranspiled(extensionSchema *extension.Schema, packageName string, sourceName string, typescriptSource string) (*Transpiled, error) {
	result := api.Transform(typescriptSource, api.TransformOptions{
		Loader:      api.LoaderTS,
		Format:      api.FormatCommonJS,
		Sourcemap:   api.SourceMapExternal,
		SourceRoot:  sourceName,
		TsconfigRaw: tsConfig,
	})

	if len(result.Errors) > 0 {
		var errString strings.Builder
		for _, err := range result.Errors {
			errString.WriteString(err.Text)
			errString.WriteRune('\n')
		}
		return nil, errors.New(errString.String())
	}
	if packageName == "" {
		packageName = defaultPackageName
	}

	headerBuf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(headerBuf, "header.ts.templ", map[string]any{
		"generator_version": strings.Trim(scaleVersion.Version(), "v"),
		"package_name":      packageName,
	})
	if err != nil {
		return nil, err
	}

	declarationBuf := new(bytes.Buffer)
	err = g.templ.ExecuteTemplate(declarationBuf, "declaration.ts.templ", map[string]any{
		"extension_schema":  extensionSchema,
		"generator_version": strings.TrimPrefix(scaleVersion.Version(), "v"),
		"package_name":      packageName,
	})
	if err != nil {
		return nil, err
	}

	return &Transpiled{
		Typescript:  []byte(typescriptSource),
		Javascript:  append(append([]byte(headerBuf.String()+"\n\n"), result.Code...), []byte(fmt.Sprintf("//# sourceMappingURL=%s.map", sourceName))...),
		SourceMap:   result.Map,
		Declaration: []byte(formatTS(declarationBuf.String())),
	}, nil
}

// GeneratePackageJSON generates the package.json file for the extension
func (g *Generator) GeneratePackageJSON(packageName string, packageVersion string) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "package.ts.templ", map[string]any{
		"polyglot_version":                   strings.TrimPrefix(polyglotVersion.Version(), "v"),
		"scale_extension_interfaces_version": strings.TrimPrefix(interfacesVersion.Version(), "v"),
		"package_name":                       packageName,
		"package_version":                    strings.TrimPrefix(packageVersion, "v"),
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GenerateGuest generates the guest bindings for the extension
func (g *Generator) GenerateGuest(extensionSchema *extension.Schema, extensionHash string, packageName string) ([]byte, error) {
	if packageName == "" {
		packageName = defaultPackageName
	}

	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "guest.ts.templ", map[string]any{
		"signature_schema":  extensionSchema,
		"signature_hash":    extensionHash,
		"generator_version": strings.TrimPrefix(scaleVersion.Version(), "v"),
		"package_name":      packageName,
	})
	if err != nil {
		return nil, err
	}

	return []byte(formatTS(buf.String())), nil
}

// GenerateGuestTranspiled takes the typescript source for the generated guest bindings and transpiles it to javascript
func (g *Generator) GenerateGuestTranspiled(extensionSchema *extension.Schema, packageName string, sourceName string, typescriptSource string) (*Transpiled, error) {
	result := api.Transform(typescriptSource, api.TransformOptions{
		Loader:      api.LoaderTS,
		Format:      api.FormatCommonJS,
		Sourcemap:   api.SourceMapExternal,
		SourceRoot:  sourceName,
		TsconfigRaw: tsConfig,
	})

	if len(result.Errors) > 0 {
		var errString strings.Builder
		for _, err := range result.Errors {
			errString.WriteString(err.Text)
			errString.WriteRune('\n')
		}
		return nil, errors.New(errString.String())
	}
	if packageName == "" {
		packageName = defaultPackageName
	}

	headerBuf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(headerBuf, "header.ts.templ", map[string]any{
		"generator_version": strings.Trim(scaleVersion.Version(), "v"),
		"package_name":      packageName,
	})
	if err != nil {
		return nil, err
	}

	declarationBuf := new(bytes.Buffer)
	err = g.templ.ExecuteTemplate(declarationBuf, "declaration-guest.ts.templ", map[string]any{
		"extension_schema":  extensionSchema,
		"generator_version": strings.TrimPrefix(scaleVersion.Version(), "v"),
		"package_name":      packageName,
	})
	if err != nil {
		return nil, err
	}

	return &Transpiled{
		Typescript:  []byte(typescriptSource),
		Javascript:  append(append([]byte(headerBuf.String()+"\n\n"), result.Code...), []byte(fmt.Sprintf("//# sourceMappingURL=%s.map", sourceName))...),
		SourceMap:   result.Map,
		Declaration: []byte(formatTS(declarationBuf.String())),
	}, nil
}

// GenerateHost generates the host bindings for the extension
//
// Note: the given schema should already be normalized, validated, and modified to have its accessors and validators disabled
func (g *Generator) GenerateHost(extensionSchema *extension.Schema, extensionHash string, packageName string) ([]byte, error) {
	if packageName == "" {
		packageName = defaultPackageName
	}

	buf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(buf, "host.ts.templ", map[string]any{
		"extension_schema":  extensionSchema,
		"extension_hash":    extensionHash,
		"generator_version": strings.TrimPrefix(scaleVersion.Version(), "v"),
		"package_name":      packageName,
	})
	if err != nil {
		return nil, err
	}

	return []byte(formatTS(buf.String())), nil
}

// GenerateHostTranspiled takes the typescript source for the generated host bindings and transpiles it to javascript
//
// Note: the given schema should already be normalized, validated, and modified to have its accessors and validators disabled
func (g *Generator) GenerateHostTranspiled(extensionSchema *extension.Schema, packageName string, sourceName string, typescriptSource string) (*Transpiled, error) {
	result := api.Transform(typescriptSource, api.TransformOptions{
		Loader:      api.LoaderTS,
		Format:      api.FormatCommonJS,
		Sourcemap:   api.SourceMapExternal,
		SourceRoot:  sourceName,
		TsconfigRaw: tsConfig,
	})

	if len(result.Errors) > 0 {
		var errString strings.Builder
		for _, err := range result.Errors {
			errString.WriteString(err.Text)
			errString.WriteRune('\n')
		}
		return nil, errors.New(errString.String())
	}
	if packageName == "" {
		packageName = defaultPackageName
	}

	headerBuf := new(bytes.Buffer)
	err := g.templ.ExecuteTemplate(headerBuf, "header.ts.templ", map[string]any{
		"generator_version": strings.Trim(scaleVersion.Version(), "v"),
		"package_name":      packageName,
	})
	if err != nil {
		return nil, err
	}

	declarationBuf := new(bytes.Buffer)
	err = g.templ.ExecuteTemplate(declarationBuf, "declaration-host.ts.templ", map[string]any{
		"extension_schema":  extensionSchema,
		"generator_version": strings.TrimPrefix(scaleVersion.Version(), "v"),
		"package_name":      packageName,
	})
	if err != nil {
		return nil, err
	}

	return &Transpiled{
		Typescript:  []byte(typescriptSource),
		Javascript:  append(append([]byte(headerBuf.String()+"\n\n"), result.Code...), []byte(fmt.Sprintf("//# sourceMappingURL=%s.map", sourceName))...),
		SourceMap:   result.Map,
		Declaration: []byte(formatTS(declarationBuf.String())),
	}, nil
}

func templateFunctions() template.FuncMap {
	return template.FuncMap{
		"Primitive":               primitive,
		"IsPrimitive":             extension.ValidPrimitiveType,
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

//nolint:revive
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
