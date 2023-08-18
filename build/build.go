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
	"github.com/loopholelabs/scale/cli/version"
	"github.com/loopholelabs/scale/compile/golang"
	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/storage"
	"os"
	"os/exec"
	"path"
)

var (
	ErrNoGo     = errors.New("go not found in PATH. Please install go: https://golang.org/doc/install")
	ErrNoTinyGo = errors.New("tinygo not found in PATH. Please install tinygo: https://tinygo.org/getting-started/")
	ErrNoCargo  = errors.New("cargo not found in PATH. Please install cargo: https://doc.rust-lang.org/cargo/getting-started/installation.html")
	ErrNoNpm    = errors.New("npm not found in PATH. Please install npm: https://docs.npmjs.com/downloading-and-installing-node-js-and-npm")
)

type LocalGolangOptions struct {
	Scalefile        *scalefile.Schema
	Signature        *signature.Schema
	SourceDirectory  string
	StorageDirectory string
	GoBin            string
	TinyGoBin        string
	Args             []string
}

type RustOptions struct {
	CargoBin string
	Args     []string
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

	var signatureImportPath string
	var signatureImportVersion string
	var sig *storage.Signature
	if options.Signature != nil {
		hash, err := options.Signature.Hash()
		if err != nil {
			return nil, fmt.Errorf("unable to hash signature: %w", err)
		}
		sig = &storage.Signature{
			Schema:       options.Signature,
			Hash:         hex.EncodeToString(hash),
			Organization: "local",
		}

		signaturePath := path.Join(build.Path, "signature")
		err = os.MkdirAll(signaturePath, 0755)
		if err != nil {
			return nil, fmt.Errorf("unable to create local signature directory: %w", err)
		}
		err = storage.GenerateSignature(sig.Schema, signaturePath)
		if err != nil {
			return nil, fmt.Errorf("unable to generate local signature: %w", err)
		}

		signatureImportPath = path.Join(signaturePath, "golang", "guest")
	} else {
		switch options.Scalefile.Signature.Organization {
		case "local":
			sig, err = sts.Get(options.Scalefile.Signature.Name, options.Scalefile.Signature.Tag, options.Scalefile.Signature.Organization, "")
			if err != nil {
				return nil, fmt.Errorf("unable to get local signature: %w", err)
			}
			if sig == nil {
				return nil, fmt.Errorf("local signature %s:%s not found", options.Scalefile.Signature.Name, options.Scalefile.Signature.Tag)
			}

			signaturePath := path.Join(build.Path, "signature")
			err = os.MkdirAll(signaturePath, 0755)
			if err != nil {
				return nil, fmt.Errorf("unable to create local signature directory: %w", err)
			}
			err = storage.GenerateSignature(sig.Schema, signaturePath)
			if err != nil {
				return nil, fmt.Errorf("unable to generate local signature: %w", err)
			}

			signatureImportPath = path.Join(signaturePath, "golang", "guest")
		default:
			return nil, fmt.Errorf("unknown signature organization %s", options.Scalefile.Signature.Organization)
		}
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

	modfile, err := golang.GenerateGoModfile(options.Scalefile, signatureImportPath, signatureImportVersion, options.SourceDirectory, dependencies, "compile", version.Version)
	if err != nil {
		return nil, fmt.Errorf("unable to generate go.mod file: %w", err)
	}

	mainFile, err := golang.GenerateGoMain(sig.Schema, options.Scalefile, version.Version)
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

	tinygoArgs := append([]string{"build", "-o", "scale.wasm", "-target=wasi"}, options.Args...)
	tinygoArgs = append(tinygoArgs, "main.go")

	cmd = exec.Command(options.TinyGoBin, tinygoArgs...)
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

	return &scalefunc.Schema{
		Version:         scalefunc.V1Alpha,
		Name:            options.Scalefile.Name,
		Tag:             options.Scalefile.Tag,
		SignatureName:   fmt.Sprintf("%s/%s:%s", options.Scalefile.Signature.Organization, options.Scalefile.Signature.Name, options.Scalefile.Signature.Tag),
		SignatureSchema: sig.Schema,
		SignatureHash:   sig.Hash,
		Language:        scalefunc.Go,
		Dependencies:    nil,
		Function:        data,
	}, nil
}
