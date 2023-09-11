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

	"github.com/loopholelabs/scale/compile/rust"
	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/storage"
)

var (
	ErrNoCargo = errors.New("cargo not found in PATH. Please install cargo: https://doc.rust-lang.org/cargo/getting-started/installation.html")
)

type LocalRustOptions struct {
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

	// CargoBin is the optional path to the cargo binary
	CargoBin string

	// Args are the optional arguments to pass to the compiler
	Args []string
}

func LocalRust(options *LocalRustOptions) (*scalefunc.Schema, error) {
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
		return nil, fmt.Errorf("signature dependency not found in Cargo.toml")
	}

	signatureDependency := manifest.GetDependency("signature")
	if signatureDependency == nil {
		return nil, fmt.Errorf("unable to parse signature dependency in Cargo.toml")
	}

	if (signatureDependency.Registry == "scale" && options.Scalefile.Signature.Organization == "local") || (signatureDependency.Registry == "" && options.Scalefile.Signature.Organization != "local") {
		return nil, fmt.Errorf("scalefile's signature block does not match Cargo.toml")
	}

	if signatureDependency.Registry == "" && !filepath.IsAbs(signatureDependency.Path) {
		signatureDependency.Path, err = filepath.Abs(path.Join(options.SourceDirectory, signatureDependency.Path))
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

	cargofile, err := rust.GenerateRustCargofile(options.Scalefile, signatureDependency, options.SourceDirectory)
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

	cmd := exec.Command(options.CargoBin, "check")
	cmd.Dir = compilePath
	cmd.Stdout = options.Output
	cmd.Stderr = options.Output
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("unable to compile scale function: %w", err)
	}

	var target string
	switch options.Target {
	case WASITarget:
		target = "wasm32-wasi"
	default:
		return nil, fmt.Errorf("unknown build target %d", options.Target)
	}

	buildArgs := append([]string{"build"}, options.Args...)
	if options.Release {
		buildArgs = append(buildArgs, "--release")
	}
	buildArgs = append(buildArgs, "--target", target)

	cmd = exec.Command(options.CargoBin, buildArgs...)
	cmd.Dir = compilePath
	cmd.Stdout = options.Output
	cmd.Stderr = options.Output
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
		Language:        scalefunc.Rust,
		Stateless:       options.Scalefile.Stateless,
		Dependencies:    nil,
		Function:        data,
	}, nil
}
