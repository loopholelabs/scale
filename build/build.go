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
	"github.com/loopholelabs/scale/compile/golang"
	"github.com/loopholelabs/scale/compile/rust"
	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/storage"
	"os"
	"os/exec"
	"path"
)

type Target int

const (
	WASITarget Target = iota
)

var (
	ErrNoGo     = errors.New("go not found in PATH. Please install go: https://golang.org/doc/install")
	ErrNoTinyGo = errors.New("tinygo not found in PATH. Please install tinygo: https://tinygo.org/getting-started/")
	ErrNoCargo  = errors.New("cargo not found in PATH. Please install cargo: https://doc.rust-lang.org/cargo/getting-started/installation.html")
)

type LocalGolangOptions struct {
	// Version is the generator version
	Version string

	// Scalefile is the scalefile to be built
	Scalefile *scalefile.Schema

	// SourceDirectory is the directory where the source code is located
	SourceDirectory string

	// SignaturePath is the import path for the signature
	//
	// For local signatures this will be a path to the signature package
	// For remote signatures this will be a path to the signature repository
	SignaturePath string

	// SignatureVersion is the optional version of the signature
	//
	// This is required for remote signatures
	SignatureVersion string

	// SignatureSchema is the schema of the signature
	SignatureSchema *signature.Schema

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
}

type LocalRustOptions struct {
	// Version is the generator version
	Version string

	// Scalefile is the scalefile to be built
	Scalefile *scalefile.Schema

	// SourceDirectory is the directory where the source code is located
	SourceDirectory string

	// SignaturePackage is the package for the signature
	SignaturePackage string

	// SignaturePath is the optional import path for the signature
	//
	// This is required for local signatures
	SignaturePath string

	// SignatureVersion is the optional version of the signature
	//
	// This is required for remote signatures
	SignatureVersion string

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

func LocalGolang(options *LocalGolangOptions) (*scalefunc.Schema, error) {
	if options.GoBin != "" {
		stat, err := os.Stat(options.GoBin)
		if err != nil {
			return nil, fmt.Errorf("unable to find go binary %s: %w", options.GoBin, err)
		}
		if !(stat.Mode()&0111 != 0) {
			return nil, fmt.Errorf("go binary %s is not executable", options.GoBin)
		}
	} else {
		var err error
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
		var err error
		options.TinyGoBin, err = exec.LookPath("tinygo")
		if err != nil {
			return nil, ErrNoTinyGo
		}
	}

	_, err := os.Stat(options.SourceDirectory)
	if err != nil {
		return nil, fmt.Errorf("unable to find source directory %s: %w", options.SourceDirectory, err)
	}

	build, err := options.Storage.Mkdir()
	if err != nil {
		return nil, fmt.Errorf("unable to create build directory: %w", err)
	}
	defer func() {
		_ = options.Storage.Delete(build)
	}()

	dependencies := []*scalefunc.Dependency{
		{
			Name:     "signature",
			Version:  "v0.1.0",
			Metadata: nil,
		},
		{
			Name:     options.Scalefile.Name,
			Version:  "v0.1.0",
			Metadata: nil,
		},
	}

	modfile, err := golang.GenerateGoModfile(options.Scalefile, options.SignaturePath, options.SignatureVersion, options.SourceDirectory, dependencies, "compile")
	if err != nil {
		return nil, fmt.Errorf("unable to generate go.mod file: %w", err)
	}

	mainFile, err := golang.GenerateGoMain(options.SignatureSchema, options.Scalefile, options.Version)
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
	output, err := cmd.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("unable to compile scale function: %s", output)
		}
		return nil, fmt.Errorf("unable to compile scale function: %w", err)
	}

	var target string
	switch options.Target {
	case WASITarget:
		target = "wasi"
	default:
		return nil, fmt.Errorf("unknown build target %d", options.Target)
	}

	buildArgs := append([]string{"build", "-o", "scale.wasm"}, options.Args...)
	if options.Release {
		buildArgs = append(buildArgs, "-no-debug")
	}
	buildArgs = append(buildArgs, "-target", target, "main.go")

	cmd = exec.Command(options.TinyGoBin, buildArgs...)
	cmd.Dir = compilePath

	output, err = cmd.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("unable to compile scale function: %s", output)
		}
		return nil, fmt.Errorf("unable to compile scale function: %w", err)
	}

	data, err := os.ReadFile(path.Join(compilePath, "scale.wasm"))
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
		Language:        scalefunc.Go,
		Dependencies:    nil,
		Function:        data,
	}, nil
}

func LocalRust(options *LocalRustOptions) (*scalefunc.Schema, error) {
	if options.CargoBin != "" {
		stat, err := os.Stat(options.CargoBin)
		if err != nil {
			return nil, fmt.Errorf("unable to find cargo binary %s: %w", options.CargoBin, err)
		}
		if !(stat.Mode()&0111 != 0) {
			return nil, fmt.Errorf("cargo binary %s is not executable", options.CargoBin)
		}
	} else {
		var err error
		options.CargoBin, err = exec.LookPath("cargo")
		if err != nil {
			return nil, ErrNoCargo
		}
	}

	_, err := os.Stat(options.SourceDirectory)
	if err != nil {
		return nil, fmt.Errorf("unable to find source directory %s: %w", options.SourceDirectory, err)
	}

	build, err := options.Storage.Mkdir()
	if err != nil {
		return nil, fmt.Errorf("unable to create build directory: %w", err)
	}
	defer func() {
		_ = options.Storage.Delete(build)
	}()

	cargofile, err := rust.GenerateRustCargofile(options.Scalefile, options.SignaturePackage, options.SignatureVersion, options.SignaturePath, options.SourceDirectory, nil, "compile", "0.1.0")
	if err != nil {
		return nil, fmt.Errorf("unable to generate cargo.toml file: %w", err)
	}

	libFile, err := rust.GenerateRustLib(options.SignatureSchema, options.Scalefile, options.Version)
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
	output, err := cmd.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("unable to compile scale function: %s", output)
		}
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

	output, err = cmd.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("unable to compile scale function: %s", output)
		}
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
		Language:        scalefunc.Go,
		Dependencies:    nil,
		Function:        data,
	}, nil
}
