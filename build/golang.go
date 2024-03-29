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
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/loopholelabs/scale/extension"

	"github.com/loopholelabs/scale/compile/golang"
	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/storage"
)

var (
	ErrNoGo     = errors.New("go not found in PATH. Please install go: https://golang.org/doc/install")
	ErrNoTinyGo = errors.New("tinygo not found in PATH. Please install tinygo: https://tinygo.org/getting-started/")

	ErrNoSignatureDependencyGoMod      = errors.New("signature dependency not found in go.mod file")
	ErrInvalidSignatureDependencyGoMod = errors.New("unable to parse signature dependency in go.mod file")
)

type LocalGolangOptions struct {
	// Stdout is the output writer for the various build commands
	Stdout io.Writer

	// Scalefile is the scalefile to be built
	Scalefile *scalefile.Schema

	// SignatureSchema is the schema of the signature
	//
	// Note: The SignatureSchema is only used to embed type information into the scale function,
	// and not as part of the build process
	SignatureSchema *signature.Schema

	// ExtensionSchemas are the schemas of the extensions. The array must be in the same order as the extensions
	// are defined in the scalefile
	//
	// Note: The ExtensionSchemas are only used to embed extension information into the scale function,
	// and not as part of the build process
	ExtensionSchemas []*extension.Schema

	// SourceDirectory is the directory where the source code is located
	SourceDirectory string

	// Storage is the storage handler to use for the build
	Storage *storage.BuildStorage

	// Release is whether to build in release mode
	Release bool

	// Target is the target to build for
	Target Target

	// GoBin is the optional path to the go binary
	GoBin string

	// TinyGoBin is the optional path to the tinygo binary
	TinyGoBin string

	// Args are the optional arguments to pass to the compiler
	Args []string

	Extensions []extension.Info
}

func LocalGolang(options *LocalGolangOptions) (*scalefunc.V1BetaSchema, error) {
	var err error
	if options.GoBin != "" {
		stat, err := os.Stat(options.GoBin)
		if err != nil {
			return nil, fmt.Errorf("unable to find go binary %s: %w", options.GoBin, err)
		}
		if !(stat.Mode()&0111 != 0) {
			return nil, fmt.Errorf("go binary %s is not executable", options.GoBin)
		}
	} else {
		options.GoBin, err = exec.LookPath("go")
		if err != nil {
			return nil, ErrNoGo
		}
	}

	if options.TinyGoBin != "" {
		stat, err := os.Stat(options.TinyGoBin)
		if err != nil {
			return nil, fmt.Errorf("unable to find tinygo binary %s: %w", options.TinyGoBin, err)
		}
		if !(stat.Mode()&0111 != 0) {
			return nil, fmt.Errorf("tinygo binary %s is not executable", options.TinyGoBin)
		}
	} else {
		options.TinyGoBin, err = exec.LookPath("tinygo")
		if err != nil {
			return nil, ErrNoTinyGo
		}
	}

	if len(options.ExtensionSchemas) != len(options.Scalefile.Extensions) {
		return nil, fmt.Errorf("number of extension schemas does not match number of extensions in scalefile")
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

	modfileData, err := os.ReadFile(path.Join(options.SourceDirectory, "go.mod"))
	if err != nil {
		return nil, fmt.Errorf("unable to read go.mod file: %w", err)
	}

	manifest, err := golang.ParseManifest(modfileData)
	if err != nil {
		return nil, fmt.Errorf("unable to parse go.mod file: %w", err)
	}

	if !manifest.HasRequire("signature", "v0.1.0", true) {
		return nil, ErrNoSignatureDependencyGoMod
	}

	if !manifest.HasReplacement("signature", "", "", "", true) {
		return nil, ErrNoSignatureDependencyGoMod
	}

	signatureInfo := new(golang.SignatureInfo)

	signatureInfo.ImportVersion, signatureInfo.ImportPath = manifest.GetReplacement("signature")
	if signatureInfo.ImportVersion != "" && signatureInfo.ImportPath != "" {
		return nil, ErrInvalidSignatureDependencyGoMod
	}

	switch options.Scalefile.Signature.Organization {
	case "local":
		signatureInfo.Local = true
		if signatureInfo.ImportVersion != "" {
			return nil, fmt.Errorf("scalefile's signature block does not match go.mod: signature import version is %s for a local signature", signatureInfo.ImportVersion)
		}

		if !filepath.IsAbs(signatureInfo.ImportPath) {
			signatureInfo.ImportPath, err = filepath.Abs(path.Join(options.SourceDirectory, signatureInfo.ImportPath))
			if err != nil {
				return nil, fmt.Errorf("unable to parse signature dependency path: %w", err)
			}
		}
	default:
		signatureInfo.Local = false
		if signatureInfo.ImportVersion == "" {
			return nil, fmt.Errorf("scalefile's signature block does not match go.mod: signature import version is empty for a signature with organization %s", options.Scalefile.Signature.Organization)
		}
	}

	functionInfo := &golang.FunctionInfo{
		PackageName: strings.ToLower(options.Scalefile.Name),
		ImportPath:  options.SourceDirectory,
	}

	build, err := options.Storage.Mkdir()
	if err != nil {
		return nil, fmt.Errorf("unable to create build directory: %w", err)
	}
	defer func() {
		_ = options.Storage.Delete(build)
	}()

	// Copy over any replacements from the go.mod
	replacements := make([]golang.GoModReplacement, 0)

	r := manifest.GetReplacements()
	for _, resp := range r {
		// Check if the target is a local dir...
		newPath := resp.New.Path

		if !filepath.IsAbs(newPath) {
			newPath = filepath.Join(options.SourceDirectory, newPath)
		}
		replacements = append(replacements, golang.GoModReplacement{
			Name: fmt.Sprintf("%s %s", resp.Old.Path, resp.Old.Version),
			Path: fmt.Sprintf("%s %s", newPath, resp.New.Version),
		})
	}

	modfile, err := golang.GenerateGoModfile(signatureInfo, functionInfo, replacements)
	if err != nil {
		return nil, fmt.Errorf("unable to generate go.mod file: %w", err)
	}

	mainFile, err := golang.GenerateGoMain(options.Scalefile, options.SignatureSchema, functionInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to generate main.go file: %w", err)
	}

	compilePath := path.Join(build.Path, "compile")

	err = os.MkdirAll(compilePath, 0755)
	if err != nil {
		return nil, fmt.Errorf("unable to create compile directory: %w", err)
	}

	err = os.WriteFile(path.Join(compilePath, "go.mod"), modfile, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to create go.mod file: %w", err)
	}

	err = os.WriteFile(path.Join(compilePath, "main.go"), mainFile, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to create main.go file: %w", err)
	}

	cmd := exec.Command(options.GoBin, "mod", "tidy")
	cmd.Dir = compilePath
	cmd.Stdout = options.Stdout
	cmd.Stderr = options.Stdout
	cmd.Env = os.Environ()
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("unable to compile scale function: %w", err)
	}

	var target string
	switch options.Target {
	case WASITarget:
		target = "wasi"
	case WASMTarget:
		target = "wasm"
	default:
		return nil, fmt.Errorf("unknown build target %d", options.Target)
	}

	buildArgs := append([]string{"build", "-o", "scale.wasm", "-gc=conservative", "-opt=s"}, options.Args...)
	if options.Release {
		buildArgs = append(buildArgs, "-no-debug")
	}
	buildArgs = append(buildArgs, "-target", target, "main.go")

	cmd = exec.Command(options.TinyGoBin, buildArgs...)
	cmd.Dir = compilePath
	cmd.Stdout = options.Stdout
	cmd.Stderr = options.Stdout
	cmd.Env = os.Environ()
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("unable to compile scale function: %w", err)
	}

	data, err := os.ReadFile(path.Join(compilePath, "scale.wasm"))
	if err != nil {
		return nil, fmt.Errorf("unable to read compiled wasm file: %w", err)
	}

	signatureHash, err := options.SignatureSchema.Hash()
	if err != nil {
		return nil, fmt.Errorf("unable to hash signature: %w", err)
	}

	sig := scalefunc.V1BetaSignature{
		Name:         options.Scalefile.Signature.Name,
		Organization: options.Scalefile.Signature.Organization,
		Tag:          options.Scalefile.Signature.Tag,
		Schema:       options.SignatureSchema,
		Hash:         hex.EncodeToString(signatureHash),
	}

	exts := make([]scalefunc.V1BetaExtension, len(options.Scalefile.Extensions))
	for i, ext := range options.Scalefile.Extensions {
		extensionHash, err := options.ExtensionSchemas[i].Hash()
		if err != nil {
			return nil, fmt.Errorf("unable to hash extension %s: %w", ext.Name, err)
		}

		exts[i] = scalefunc.V1BetaExtension{
			Name:         ext.Name,
			Organization: ext.Organization,
			Tag:          ext.Tag,
			Schema:       options.ExtensionSchemas[i],
			Hash:         hex.EncodeToString(extensionHash),
		}
	}

	return &scalefunc.V1BetaSchema{
		Name:       options.Scalefile.Name,
		Tag:        options.Scalefile.Tag,
		Signature:  sig,
		Extensions: exts,
		Language:   scalefunc.Go,
		Manifest:   modfileData,
		Stateless:  options.Scalefile.Stateless,
		Function:   data,
	}, nil
}
