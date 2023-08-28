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
	ErrNoNpm    = errors.New("npm not found in PATH. Please install npm: https://docs.npmjs.com/downloading-and-installing-node-js-and-npm")
)

type LocalGolangOptions struct {
	Version          string
	Scalefile        *scalefile.Schema
	SourceDirectory  string
	SignaturePath    string
	SignatureSchema  *signature.Schema
	StorageDirectory string
	Release          bool
	Target           Target
	GoBin            string
	TinyGoBin        string
	Args             []string
}

type LocalRustOptions struct {
	Version          string
	Scalefile        *scalefile.Schema
	SourceDirectory  string
	SignaturePath    string
	SignatureSchema  *signature.Schema
	StorageDirectory string
	Release          bool
	Target           Target
	Registry         string
	CargoBin         string
	Args             []string
}

type TypescriptOptions struct {
	NpmBin string
	Args   []string
}

func LocalGolang(options *LocalGolangOptions) (*scalefunc.Schema, error) {
	stb := storage.DefaultBuild
	sts := storage.DefaultSignature
	if options.StorageDirectory != "" {
		var err error
		stb, err = storage.NewBuild(options.StorageDirectory)
		if err != nil {
			return nil, fmt.Errorf("failed to instantiate builds storage for %s: %w", options.StorageDirectory, err)
		}

		sts, err = storage.NewSignature(options.StorageDirectory)
		if err != nil {
			return nil, fmt.Errorf("failed to instantiate signatures storage for %s: %w", options.StorageDirectory, err)
		}
	}

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

	build, err := stb.Mkdir()
	if err != nil {
		return nil, fmt.Errorf("unable to create build directory: %w", err)
	}
	defer func() {
		_ = stb.Delete(build)
	}()

	signatureImportPath := options.SignaturePath
	signatureImportVersion := ""
	var sig *signature.Schema
	if signatureImportPath == "" {
		switch options.Scalefile.Signature.Organization {
		case "local":
			storageSig, err := sts.Get(options.Scalefile.Signature.Name, options.Scalefile.Signature.Tag, options.Scalefile.Signature.Organization, "")
			if err != nil {
				return nil, fmt.Errorf("unable to get local signature: %w", err)
			}
			if storageSig == nil {
				return nil, fmt.Errorf("local signature %s:%s not found", options.Scalefile.Signature.Name, options.Scalefile.Signature.Tag)
			}

			signaturePath := path.Join(build.Path, "signature")
			err = os.MkdirAll(signaturePath, 0755)
			if err != nil {
				return nil, fmt.Errorf("unable to create local signature directory: %w", err)
			}
			err = storage.GenerateSignature(storageSig.Schema, storageSig.Name, storageSig.Tag, storageSig.Organization, signaturePath)
			if err != nil {
				return nil, fmt.Errorf("unable to generate local signature: %w", err)
			}
			sig = storageSig.Schema
			signatureImportPath = path.Join(signaturePath, "golang", "guest")
		default:
			return nil, fmt.Errorf("unknown signature organization %s", options.Scalefile.Signature.Organization)
		}
	} else {
		sig = options.SignatureSchema
	}

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

	modfile, err := golang.GenerateGoModfile(options.Scalefile, signatureImportPath, signatureImportVersion, options.SourceDirectory, dependencies, "compile")
	if err != nil {
		return nil, fmt.Errorf("unable to generate go.mod file: %w", err)
	}

	mainFile, err := golang.GenerateGoMain(sig, options.Scalefile, options.Version)
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
		return nil, fmt.Errorf("unknown build target %s", options.Target)
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

	hash, err := sig.Hash()
	if err != nil {
		return nil, fmt.Errorf("unable to hash signature: %w", err)
	}

	return &scalefunc.Schema{
		Version:         scalefunc.V1Alpha,
		Name:            options.Scalefile.Name,
		Tag:             options.Scalefile.Tag,
		SignatureName:   fmt.Sprintf("%s/%s:%s", options.Scalefile.Signature.Organization, options.Scalefile.Signature.Name, options.Scalefile.Signature.Tag),
		SignatureSchema: sig,
		SignatureHash:   hex.EncodeToString(hash),
		Language:        scalefunc.Go,
		Dependencies:    nil,
		Function:        data,
	}, nil
}

func LocalRust(options *LocalRustOptions) (*scalefunc.Schema, error) {
	stb := storage.DefaultBuild
	sts := storage.DefaultSignature
	if options.StorageDirectory != "" {
		var err error
		stb, err = storage.NewBuild(options.StorageDirectory)
		if err != nil {
			return nil, fmt.Errorf("failed to instantiate builds storage for %s: %w", options.StorageDirectory, err)
		}

		sts, err = storage.NewSignature(options.StorageDirectory)
		if err != nil {
			return nil, fmt.Errorf("failed to instantiate signatures storage for %s: %w", options.StorageDirectory, err)
		}
	}

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

	build, err := stb.Mkdir()
	if err != nil {
		return nil, fmt.Errorf("unable to create build directory: %w", err)
	}
	defer func() {
		_ = stb.Delete(build)
	}()

	signatureImportPath := options.SignaturePath
	signatureImportVersion := ""
	var sig *signature.Schema
	if signatureImportPath == "" {
		switch options.Scalefile.Signature.Organization {
		case "local":
			storageSig, err := sts.Get(options.Scalefile.Signature.Name, options.Scalefile.Signature.Tag, options.Scalefile.Signature.Organization, "")
			if err != nil {
				return nil, fmt.Errorf("unable to get local signature: %w", err)
			}
			if storageSig == nil {
				return nil, fmt.Errorf("local signature %s:%s not found", options.Scalefile.Signature.Name, options.Scalefile.Signature.Tag)
			}

			signaturePath := path.Join(build.Path, "signature")
			err = os.MkdirAll(signaturePath, 0755)
			if err != nil {
				return nil, fmt.Errorf("unable to create local signature directory: %w", err)
			}
			err = storage.GenerateSignature(storageSig.Schema, storageSig.Name, storageSig.Tag, storageSig.Organization, signaturePath)
			if err != nil {
				return nil, fmt.Errorf("unable to generate local signature: %w", err)
			}
			sig = storageSig.Schema
			signatureImportPath = path.Join(signaturePath, "rust", "guest")
		default:
			return nil, fmt.Errorf("unknown signature organization %s", options.Scalefile.Signature.Organization)
		}
	} else {
		sig = options.SignatureSchema
	}

	cargofile, err := rust.GenerateRustCargofile(options.Scalefile, options.Registry, fmt.Sprintf("%s_%s_%s_guest", options.Scalefile.Signature.Organization, options.Scalefile.Signature.Name, options.Scalefile.Signature.Tag), signatureImportVersion, signatureImportPath, options.SourceDirectory, nil, "compile", "0.1.0")
	if err != nil {
		return nil, fmt.Errorf("unable to generate cargo.toml file: %w", err)
	}

	libFile, err := rust.GenerateRustLib(sig, options.Scalefile, options.Version)
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
		return nil, fmt.Errorf("unknown build target %s", options.Target)
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

	hash, err := sig.Hash()
	if err != nil {
		return nil, fmt.Errorf("unable to hash signature: %w", err)
	}

	return &scalefunc.Schema{
		Version:         scalefunc.V1Alpha,
		Name:            options.Scalefile.Name,
		Tag:             options.Scalefile.Tag,
		SignatureName:   fmt.Sprintf("%s/%s:%s", options.Scalefile.Signature.Organization, options.Scalefile.Signature.Name, options.Scalefile.Signature.Tag),
		SignatureSchema: sig,
		SignatureHash:   hex.EncodeToString(hash),
		Language:        scalefunc.Go,
		Dependencies:    nil,
		Function:        data,
	}, nil
}
