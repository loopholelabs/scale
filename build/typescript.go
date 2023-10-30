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

package build

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"

	"github.com/loopholelabs/scale/compile/typescript"
	"github.com/loopholelabs/scale/compile/typescript/builder"

	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/storage"

	"github.com/loopholelabs/wasm-toolkit/pkg/customs"
	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/debug"
	"github.com/loopholelabs/wasm-toolkit/pkg/wasm/wasmfile"
)

var (
	ErrNoNPM = errors.New("npm not found in PATH. Please install npm: https://docs.npmjs.com/downloading-and-installing-node-js-and-npm")
)

const (
	tsConfig = `
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
    "types": ["node"]
  },
}`
)

type LocalTypescriptOptions struct {
	// Output is the output writer for the various build commands
	Output io.Writer

	// Scalefile is the scalefile to be built
	Scalefile *scalefile.Schema

	// SourceDirectory is the directory where the source code is located
	SourceDirectory string

	// SignatureSchema is the schema of the signature
	SignatureSchema *signature.Schema

	// Storage is the storage handler to use for the build
	Storage *storage.BuildStorage

	// Release is whether to build in release mode
	Release bool

	// Target is the target to build for
	Target Target

	// NPMBin is the optional path to the npm binary
	NPMBin string
}

func LocalTypescript(options *LocalTypescriptOptions) (*scalefunc.Schema, error) {
	var err error
	if options.NPMBin != "" {
		stat, err := os.Stat(options.NPMBin)
		if err != nil {
			return nil, fmt.Errorf("unable to find npm binary %s: %w", options.NPMBin, err)
		}
		if !(stat.Mode()&0111 != 0) {
			return nil, fmt.Errorf("npm binary %s is not executable", options.NPMBin)
		}
	} else {
		options.NPMBin, err = exec.LookPath("npm")
		if err != nil {
			return nil, ErrNoNPM
		}
	}

	if !filepath.IsAbs(options.SourceDirectory) {
		options.SourceDirectory, err = filepath.Abs(options.SourceDirectory)
		if err != nil {
			return nil, fmt.Errorf("unable to parse source directory: %w", err)
		}
	}

	_, err = os.Stat(options.SourceDirectory)
	if err != nil {
		return nil, fmt.Errorf("unable to find source directory %s: %w", options.SourceDirectory, err)
	}

	packageJSONData, err := os.ReadFile(path.Join(options.SourceDirectory, "package.json"))
	if err != nil {
		return nil, fmt.Errorf("unable to read Cargo.toml file: %w", err)
	}

	manifest, err := typescript.ParseManifest(packageJSONData)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Cargo.toml file: %w", err)
	}

	if !manifest.HasDependency("signature") {
		return nil, fmt.Errorf("signature dependency not found in package.json")
	}

	signaturePath := manifest.GetDependency("signature")
	if signaturePath == "" {
		return nil, fmt.Errorf("unable to parse signature dependency in package.json")
	}

	signatureImport := ""
	switch {
	case strings.HasPrefix(signaturePath, "http://") || strings.HasPrefix(signaturePath, "https://"):
		signatureImport = signaturePath
		signaturePath = ""
	case strings.HasPrefix(signaturePath, "file:"):
		signaturePath = strings.TrimPrefix(signaturePath, "file:")
	default:
		return nil, fmt.Errorf("unable to parse signature dependency path: %s", signaturePath)
	}

	if signaturePath == "" && options.Scalefile.Signature.Organization == "local" {
		return nil, fmt.Errorf("scalefile's signature block does not match package.json")
	}

	if signaturePath != "" && !filepath.IsAbs(signaturePath) {
		signaturePath, err = filepath.Abs(path.Join(options.SourceDirectory, signaturePath))
		if err != nil {
			return nil, fmt.Errorf("unable to parse signature dependency path: %w", err)
		}
	}

	build, err := options.Storage.Mkdir()
	if err != nil {
		return nil, fmt.Errorf("unable to create build directory: %w", err)
	}
	defer func() {
		_ = options.Storage.Delete(build)
	}()

	var target api.Platform
	switch options.Target {
	case WASITarget:
		target = api.PlatformNode
	case WASMTarget:
		target = api.PlatformBrowser
	default:
		return nil, fmt.Errorf("unknown build target %d", options.Target)
	}

	result := api.Build(api.BuildOptions{
		Bundle:      true,
		Platform:    target,
		Format:      api.FormatCommonJS,
		Define:      map[string]string{"global": "globalThis"},
		TsconfigRaw: tsConfig,
		EntryPoints: []string{path.Join(options.SourceDirectory, "index.ts")},
	})

	if len(result.Errors) > 0 {
		var errString strings.Builder
		for _, err := range result.Errors {
			errString.WriteString(err.Text)
			errString.WriteRune('\n')
		}
		return nil, fmt.Errorf("unable to compile scale function source using esbuild: %s", errString.String())
	}

	functionDir := path.Join(build.Path, "function")
	err = os.MkdirAll(functionDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("unable to create function dist directory: %w", err)
	}

	err = os.WriteFile(path.Join(functionDir, "index.js"), result.OutputFiles[0].Contents, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to write index.js file for function dist: %w", err)
	}

	packageJSONFile, err := typescript.GenerateTypescriptPackageJSON(options.Scalefile, signaturePath, signatureImport)
	if err != nil {
		return nil, fmt.Errorf("unable to generate package.json file: %w", err)
	}

	indexFile, err := typescript.GenerateTypescriptIndex(options.Scalefile, functionDir)
	if err != nil {
		return nil, fmt.Errorf("unable to generate index.ts file: %w", err)
	}

	compilePath := path.Join(build.Path, "compile")

	err = os.MkdirAll(compilePath, 0755)
	if err != nil {
		return nil, fmt.Errorf("unable to create compile directory: %w", err)
	}

	err = os.WriteFile(path.Join(compilePath, "package.json"), packageJSONFile, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to create package.json file: %w", err)
	}

	err = os.WriteFile(path.Join(compilePath, "index.ts"), indexFile, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to create index.ts file: %w", err)
	}

	cmd := exec.Command(options.NPMBin, "install")
	cmd.Dir = compilePath
	cmd.Stderr = options.Output
	cmd.Stdout = options.Output
	cmd.Env = os.Environ()
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("unable to compile scale function and update npm: %w", err)
	}

	result = api.Build(api.BuildOptions{
		Bundle:      true,
		Platform:    target,
		Format:      api.FormatCommonJS,
		Define:      map[string]string{"global": "globalThis"},
		TsconfigRaw: tsConfig,
		EntryPoints: []string{path.Join(compilePath, "index.ts")},
	})

	if len(result.Errors) > 0 {
		var errString strings.Builder
		for _, err := range result.Errors {
			errString.WriteString(err.Text)
			errString.WriteRune('\n')
		}
		return nil, fmt.Errorf("unable to compile scale function compiler using esbuild: %s", errString.String())
	}

	jsBuilderBinary := path.Join(build.Path, "js_builder")
	err = os.WriteFile(path.Join(build.Path, "js_builder"), builder.BuilderExecutable, 0770)
	if err != nil {
		return nil, fmt.Errorf("unable to write js_builder executable: %w", err)
	}

	cmd = exec.Command(jsBuilderBinary, "-o", path.Join(build.Path, "scale.wasm"))
	cmd.Stdin = strings.NewReader(string(result.OutputFiles[0].Contents))
	cmd.Stderr = options.Output
	cmd.Stdout = options.Output
	cmd.Env = os.Environ()
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("unable to compile scale function using js_builder: %w", err)
	}
	/*
		data, err := os.ReadFile(path.Join(build.Path, "scale.wasm"))
		if err != nil {
			return nil, fmt.Errorf("unable to read compiled wasm file: %w", err)
		}
	*/
	// Do the extension transform here...
	wfile, err := wasmfile.New(path.Join(build.Path, "scale.wasm"))
	if err != nil {
		return nil, fmt.Errorf("unable to read compiled wasm file: %w", err)
	}

	wfile.Debug = &debug.WasmDebug{}
	wfile.Debug.ParseNameSectionData(wfile.GetCustomSectionData("name"))

	// TODO: This config needs to come from generated extension bits...
	conf_imp := customs.RemapMuxImport{
		Source: customs.Import{
			Module: "env",
			Name:   "ext_mux",
		},
		Mapper: map[uint64]customs.Import{
			0: {
				Module: "env",
				Name:   "ext_5c7d22390f9101d459292d76c11b5e9f66c327b1766aae34b9cc75f9f40e8206_New",
			},
			1: {
				Module: "env",
				Name:   "ext_5c7d22390f9101d459292d76c11b5e9f66c327b1766aae34b9cc75f9f40e8206_HttpConnector_Fetch",
			},
		},
	}

	conf_exp := customs.RemapMuxExport{
		Source: "ext_resize",
		Mapper: map[uint64]string{
			0: "ext_5c7d22390f9101d459292d76c11b5e9f66c327b1766aae34b9cc75f9f40e8206_Resize",
		},
	}

	err = customs.MuxImport(wfile, conf_imp)
	if err != nil {
		return nil, fmt.Errorf("unable to parse extension remap")
	}

	err = customs.MuxExport(wfile, conf_exp)
	if err != nil {
		return nil, fmt.Errorf("unable to parse extension remap")
	}

	var wasm_bin bytes.Buffer
	err = wfile.EncodeBinary(&wasm_bin)
	if err != nil {
		return nil, fmt.Errorf("unable to parse extension remap")
	}

	hash, err := options.SignatureSchema.Hash()
	if err != nil {
		return nil, fmt.Errorf("unable to hash signature: %w", err)
	}

	return &scalefunc.Schema{
		Version:         scalefunc.V1Alpha,
		Name:            options.Scalefile.Name,
		Tag:             options.Scalefile.Tag,
		SignatureName:   fmt.Sprintf("%s/%s:%s", options.Scalefile.Signature.Organization, options.Scalefile.Signature.Name, options.Scalefile.Signature.Tag),
		SignatureSchema: options.SignatureSchema,
		SignatureHash:   hex.EncodeToString(hash),
		Language:        scalefunc.TypeScript,
		Stateless:       options.Scalefile.Stateless,
		Dependencies:    nil,
		Function:        wasm_bin.Bytes(),
	}, nil
}
