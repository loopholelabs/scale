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

	"github.com/loopholelabs/scale/compile/rust"
	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/storage"
)

var (
	ErrNoCargo = errors.New("cargo not found in PATH. Please install cargo: https://doc.rust-lang.org/cargo/getting-started/installation.html")

	ErrNoSignatureDependencyCargoToml      = errors.New("signature dependency not found in Cargo.toml file")
	ErrInvalidSignatureDependencyCargoToml = errors.New("unable to parse signature dependency in Cargo.toml file")
)

type LocalRustOptions struct {
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

	// CargoBin is the optional path to the cargo binary
	CargoBin string

	// Args are the optional arguments to pass to the compiler
	Args []string
}

func LocalRust(options *LocalRustOptions) (*scalefunc.V1BetaSchema, error) {
	var err error
	if options.CargoBin != "" {
		stat, err := os.Stat(options.CargoBin)
		if err != nil {
			return nil, fmt.Errorf("unable to find cargo binary %s: %w", options.CargoBin, err)
		}
		if !(stat.Mode()&0111 != 0) {
			return nil, fmt.Errorf("cargo binary %s is not executable", options.CargoBin)
		}
	} else {
		options.CargoBin, err = exec.LookPath("cargo")
		if err != nil {
			return nil, ErrNoCargo
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

	cargoFileData, err := os.ReadFile(path.Join(options.SourceDirectory, "Cargo.toml"))
	if err != nil {
		return nil, fmt.Errorf("unable to read Cargo.toml file: %w", err)
	}

	manifest, err := rust.ParseManifest(cargoFileData)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Cargo.toml file: %w", err)
	}

	if !manifest.HasDependency("signature") {
		return nil, ErrNoSignatureDependencyCargoToml
	}

	parsedSignatureDependency := manifest.GetDependency("signature")
	if parsedSignatureDependency == nil {
		return nil, ErrInvalidSignatureDependencyCargoToml
	}

	signatureInfo := &rust.SignatureInfo{
		PackageName: parsedSignatureDependency.Package,
	}

	switch options.Scalefile.Signature.Organization {
	case "local":
		signatureInfo.Local = true
		if parsedSignatureDependency.Registry != "" {
			return nil, fmt.Errorf("scalefile's signature block does not match Cargo.toml: signature import registry is %s for a local signature", parsedSignatureDependency.Registry)
		}
		if !filepath.IsAbs(parsedSignatureDependency.Path) {
			parsedSignatureDependency.Path, err = filepath.Abs(path.Join(options.SourceDirectory, parsedSignatureDependency.Path))
			if err != nil {
				return nil, fmt.Errorf("unable to parse signature dependency path: %w", err)
			}
		}
		signatureInfo.ImportPath = parsedSignatureDependency.Path
	default:
		signatureInfo.Local = false
		if parsedSignatureDependency.Registry == "" {
			return nil, fmt.Errorf("scalefile's signature block does not match Cargo.toml: signature import registry is empty for a signature with organization %s", options.Scalefile.Signature.Organization)
		}

		if parsedSignatureDependency.Registry != "scale" {
			return nil, fmt.Errorf("scalefile's signature block does not match Cargo.toml: signature import registry is %s for a signature with organization %s", parsedSignatureDependency.Registry, options.Scalefile.Signature.Organization)
		}

		signatureInfo.ImportVersion = parsedSignatureDependency.Version
	}

	functionInfo := &rust.FunctionInfo{
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

	cargofile, err := rust.GenerateRustCargofile(signatureInfo, functionInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to generate cargo.toml file: %w", err)
	}

	libFile, err := rust.GenerateRustLib(options.Scalefile)
	if err != nil {
		return nil, fmt.Errorf("unable to generate lib.rs file: %w", err)
	}

	compilePath := path.Join(build.Path, "compile")

	err = os.MkdirAll(compilePath, 0755)
	if err != nil {
		return nil, fmt.Errorf("unable to create compile directory: %w", err)
	}

	err = os.WriteFile(path.Join(compilePath, "Cargo.toml"), cargofile, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to create cargo.toml file: %w", err)
	}

	err = os.WriteFile(path.Join(compilePath, "lib.rs"), libFile, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to create lib.rs file: %w", err)
	}

	var target string
	switch options.Target {
	case WASITarget:
		target = "wasm32-wasi"
	case WASMTarget:
		target = "wasm32-unknown-unknown"
	default:
		return nil, fmt.Errorf("unknown build target %d", options.Target)
	}

	cmd := exec.Command(options.CargoBin, "check", "--target", target)
	cmd.Dir = compilePath
	cmd.Stdout = options.Stdout
	cmd.Stderr = options.Stdout
	cmd.Env = os.Environ()
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("unable to compile scale function: %w", err)
	}

	buildArgs := append([]string{"build"}, options.Args...)
	if options.Release {
		buildArgs = append(buildArgs, "--release")
	}
	buildArgs = append(buildArgs, "--target", target)

	cmd = exec.Command(options.CargoBin, buildArgs...)
	cmd.Dir = compilePath
	cmd.Stdout = options.Stdout
	cmd.Stderr = options.Stdout
	cmd.Env = os.Environ()
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("unable to compile scale function: %w", err)
	}

	targetFolder := "debug"
	if options.Release {
		targetFolder = "release"
	}
	data, err := os.ReadFile(path.Join(compilePath, "target", target, targetFolder, "compile.wasm"))
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
		Language:   scalefunc.Rust,
		Manifest:   cargoFileData,
		Stateless:  options.Scalefile.Stateless,
		Function:   data,
	}, nil
}
