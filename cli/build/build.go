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
	"errors"
	"fmt"
	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
)

var (
	ErrNoGo     = errors.New("go not found in PATH. Please install go: https://golang.org/doc/install")
	ErrNoTinyGo = errors.New("tinygo not found in PATH. Please install tinygo: https://tinygo.org/getting-started/")
	ErrNoCargo  = errors.New("cargo not found in PATH. Please install cargo: https://doc.rust-lang.org/cargo/getting-started/installation.html")
	ErrNoNpm    = errors.New("npm not found in PATH. Please install npm: https://docs.npmjs.com/downloading-and-installing-node-js-and-npm")
)

type Module struct {
	Source    string
	Name      string
	Signature string
}

type RustOptions struct {
	CargoBin string
	Args     []string
}

type GolangOptions struct {
	GoBin     string
	TinygoBin string
	Args      []string
}

type TypescriptOptions struct {
	NpmBin string
}

type Options struct {
	Scalefile  *scalefile.Schema
	Signature  *signature.Schema
	BaseDir    string
	Golang     *GolangOptions
	Rust       *RustOptions
	Typescript *TypescriptOptions
}

func LocalBuild(options *Options) (*scalefunc.Schema, error) {
	//scaleFunc := &scalefunc.Schema{
	//	Version:   scalefunc.V1Alpha,
	//	Name:      options.Scalefile.Name,
	//	Tag:       options.Scalefile.Tag,
	//	Signature: fmt.Sprintf("%s/%s:%s", options.Scalefile.Signature.Organization, options.Scalefile.Name, options.Scalefile.Tag),
	//}

	switch scalefunc.Language(options.Scalefile.Language) {
	//case scalefunc.Go:
	//	scaleFunc.Language = scalefunc.Go
	//	return GolangBuild(scaleFile, scaleFunc, goBin, tinygoBin, tinygoArgs, baseDir)
	//case scalefunc.Rust:
	//	scaleFunc.Language = scalefunc.Rust
	//	return RustBuild(scaleFile, scaleFunc, cargoBin, cargoArgs, baseDir)
	//case scalefunc.TypeScript:
	//	scaleFunc.Language = scalefunc.TypeScript
	//	return TypeScriptBuild(scaleFile, scaleFunc, npmBin, baseDir)
	default:
		return nil, fmt.Errorf("%s support not implemented", options.Scalefile.Language)
	}
}

//func GolangBuild() (*scalefunc.Schema, error) {
//	module := &Module{
//		Name:      scaleFile.Name,
//		Source:    scaleFile.Source,
//		Signature: "github.com/loopholelabs/scale-signature-http",
//	}
//
//	if goBin != "" {
//		stat, err := os.Stat(goBin)
//		if err != nil {
//			return nil, fmt.Errorf("unable to find go binary %s: %w", goBin, err)
//		}
//		if !(stat.Mode()&0111 != 0) {
//			return nil, fmt.Errorf("go binary %s is not executable", goBin)
//		}
//	} else {
//		var err error
//		goBin, err = exec.LookPath("go")
//		if err != nil {
//			return nil, ErrNoGo
//		}
//	}
//
//	if tinygoBin != "" {
//		stat, err := os.Stat(tinygoBin)
//		if err != nil {
//			return nil, fmt.Errorf("unable to find tinygo binary %s: %w", tinygoBin, err)
//		}
//		if !(stat.Mode()&0111 != 0) {
//			return nil, fmt.Errorf("tinygo binary %s is not executable", tinygoBin)
//		}
//	} else {
//		var err error
//		tinygoBin, err = exec.LookPath("tinygo")
//		if err != nil {
//			return nil, ErrNoTinyGo
//		}
//	}
//
//	g := compile.NewGenerator()
//
//	moduleSourcePath := path.Join(baseDir, module.Source)
//	_, err := os.Stat(moduleSourcePath)
//	if err != nil {
//		return nil, fmt.Errorf("unable to find module %s: %w", moduleSourcePath, err)
//	}
//
//	buildDir := path.Join(baseDir, "build")
//	defer func() {
//		_ = os.RemoveAll(buildDir)
//	}()
//
//	err = os.Mkdir(buildDir, 0755)
//	if !os.IsExist(err) && err != nil {
//		return nil, fmt.Errorf("unable to create build %s directory: %w", buildDir, err)
//	}
//
//	file, err := os.OpenFile(path.Join(buildDir, "main.go"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
//	if err != nil {
//		return nil, fmt.Errorf("unable to create main.go file: %w", err)
//	}
//
//	err = g.GenerateGoMain(file, "compile/scale", module.Signature)
//	if err != nil {
//		return nil, fmt.Errorf("unable to generate main.go file: %w", err)
//	}
//
//	err = file.Close()
//	if err != nil {
//		return nil, fmt.Errorf("unable to close main.go file: %w", err)
//	}
//
//	file, err = os.OpenFile(path.Join(buildDir, "go.mod"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
//	if err != nil {
//		return nil, fmt.Errorf("unable to create go.mod file: %w", err)
//	}
//
//	deps := make([]*scalefile.Dependency, 0, len(scaleFile.Dependencies))
//	for _, dep := range scaleFile.Dependencies {
//		var d = dep
//		deps = append(deps, &d)
//	}
//	err = g.GenerateGoModfile(file, deps)
//	if err != nil {
//		return nil, fmt.Errorf("unable to generate go.mod file: %w", err)
//	}
//
//	err = file.Close()
//	if err != nil {
//		return nil, fmt.Errorf("unable to close go.mod file: %w", err)
//	}
//
//	scalePath := path.Join(buildDir, "scale")
//	err = os.Mkdir(scalePath, 0755)
//	if !os.IsExist(err) && err != nil {
//		return nil, fmt.Errorf("unable to create scale source directory %s: %w", scalePath, err)
//	}
//
//	scale, err := os.OpenFile(path.Join(scalePath, "scale.go"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
//	if err != nil {
//		return nil, fmt.Errorf("unable to create scale.go file: %w", err)
//	}
//
//	file, err = os.Open(moduleSourcePath)
//	if err != nil {
//		return nil, fmt.Errorf("unable to open scale source file: %w", err)
//	}
//
//	_, err = io.Copy(scale, file)
//	if err != nil {
//		return nil, fmt.Errorf("unable to copy scale source file: %w", err)
//	}
//
//	err = scale.Close()
//	if err != nil {
//		return nil, fmt.Errorf("unable to close scale.go file: %w", err)
//	}
//
//	err = file.Close()
//	if err != nil {
//		return nil, fmt.Errorf("unable to close scale source file: %w", err)
//	}
//
//	wd, err := os.Getwd()
//	if err != nil {
//		return nil, fmt.Errorf("unable to get working directory: %w", err)
//	}
//
//	cmd := exec.Command(goBin, "mod", "tidy")
//	cmd.Dir = path.Join(wd, buildDir)
//	output, err := cmd.CombinedOutput()
//	if err != nil {
//		if _, ok := err.(*exec.ExitError); ok {
//			return nil, fmt.Errorf("unable to compile scale function: %s", output)
//		}
//		return nil, fmt.Errorf("unable to compile scale function: %w", err)
//	}
//
//	tinygoArgs = append([]string{"build", "-o", "scale.wasm", "-target=wasi"}, tinygoArgs...)
//	tinygoArgs = append(tinygoArgs, "main.go")
//
//	cmd = exec.Command(tinygoBin, tinygoArgs...)
//	cmd.Dir = path.Join(wd, buildDir)
//
//	output, err = cmd.CombinedOutput()
//	if err != nil {
//		if _, ok := err.(*exec.ExitError); ok {
//			return nil, fmt.Errorf("unable to compile scale function: %s", output)
//		}
//		return nil, fmt.Errorf("unable to compile scale function: %w", err)
//	}
//
//	data, err := os.ReadFile(path.Join(cmd.Dir, "scale.wasm"))
//	if err != nil {
//		return nil, fmt.Errorf("unable to read compiled wasm file: %w", err)
//	}
//	scaleFunc.Function = data
//
//	return scaleFunc, nil
//}
