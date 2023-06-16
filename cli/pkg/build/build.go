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
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"

	//rustCompile "github.com/loopholelabs/scale/compile/rust"
	//tsCompile "github.com/loopholelabs/scale/compile/typescript"
	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/scalefunc"

	_ "embed"
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

func LocalBuild(scaleFile *scalefile.Schema, goBin string, tinygoBin string, cargoBin string, npmBin string, baseDir string, tinygoArgs []string, cargoArgs []string) (*scalefunc.ScaleFunc, error) {
	scaleFunc := &scalefunc.ScaleFunc{
		Version:   scalefunc.V1Alpha,
		Name:      scaleFile.Name,
		Tag:       scaleFile.Tag,
		Signature: scaleFile.Signature,
		Language:  scaleFile.Language,
	}

	switch scaleFunc.Language {
	case scalefunc.Go:
		return GolangBuild(scaleFile, scaleFunc, goBin, tinygoBin, tinygoArgs, baseDir)
	case scalefunc.Rust:
		return RustBuild(scaleFile, scaleFunc, cargoBin, cargoArgs, baseDir)
	case scalefunc.TypeScript:
		return TypeScriptBuild(scaleFile, scaleFunc, npmBin, baseDir)
	default:
		return nil, fmt.Errorf("%s support not implemented", scaleFile.Language)
	}
}

func GolangBuild(scaleFile *scalefile.ScaleFile, scaleFunc *scalefunc.ScaleFunc, goBin string, tinygoBin string, tinygoArgs []string, baseDir string) (*scalefunc.ScaleFunc, error) {
	module := &Module{
		Name:      scaleFile.Name,
		Source:    scaleFile.Source,
		Signature: "github.com/loopholelabs/scale-signature-http",
	}

	if goBin != "" {
		stat, err := os.Stat(goBin)
		if err != nil {
			return nil, fmt.Errorf("unable to find go binary %s: %w", goBin, err)
		}
		if !(stat.Mode()&0111 != 0) {
			return nil, fmt.Errorf("go binary %s is not executable", goBin)
		}
	} else {
		var err error
		goBin, err = exec.LookPath("go")
		if err != nil {
			return nil, ErrNoGo
		}
	}

	if tinygoBin != "" {
		stat, err := os.Stat(tinygoBin)
		if err != nil {
			return nil, fmt.Errorf("unable to find tinygo binary %s: %w", tinygoBin, err)
		}
		if !(stat.Mode()&0111 != 0) {
			return nil, fmt.Errorf("tinygo binary %s is not executable", tinygoBin)
		}
	} else {
		var err error
		tinygoBin, err = exec.LookPath("tinygo")
		if err != nil {
			return nil, ErrNoTinyGo
		}
	}

	g := compile.NewGenerator()

	moduleSourcePath := path.Join(baseDir, module.Source)
	_, err := os.Stat(moduleSourcePath)
	if err != nil {
		return nil, fmt.Errorf("unable to find module %s: %w", moduleSourcePath, err)
	}

	buildDir := path.Join(baseDir, "build")
	defer func() {
		_ = os.RemoveAll(buildDir)
	}()

	err = os.Mkdir(buildDir, 0755)
	if !os.IsExist(err) && err != nil {
		return nil, fmt.Errorf("unable to create build %s directory: %w", buildDir, err)
	}

	file, err := os.OpenFile(path.Join(buildDir, "main.go"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to create main.go file: %w", err)
	}

	err = g.GenerateGoMain(file, "compile/scale", module.Signature)
	if err != nil {
		return nil, fmt.Errorf("unable to generate main.go file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to close main.go file: %w", err)
	}

	file, err = os.OpenFile(path.Join(buildDir, "go.mod"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to create go.mod file: %w", err)
	}

	deps := make([]*scalefile.Dependency, 0, len(scaleFile.Dependencies))
	for _, dep := range scaleFile.Dependencies {
		var d = dep
		deps = append(deps, &d)
	}
	err = g.GenerateGoModfile(file, deps)
	if err != nil {
		return nil, fmt.Errorf("unable to generate go.mod file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to close go.mod file: %w", err)
	}

	scalePath := path.Join(buildDir, "scale")
	err = os.Mkdir(scalePath, 0755)
	if !os.IsExist(err) && err != nil {
		return nil, fmt.Errorf("unable to create scale source directory %s: %w", scalePath, err)
	}

	scale, err := os.OpenFile(path.Join(scalePath, "scale.go"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to create scale.go file: %w", err)
	}

	file, err = os.Open(moduleSourcePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open scale source file: %w", err)
	}

	_, err = io.Copy(scale, file)
	if err != nil {
		return nil, fmt.Errorf("unable to copy scale source file: %w", err)
	}

	err = scale.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to close scale.go file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to close scale source file: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to get working directory: %w", err)
	}

	cmd := exec.Command(goBin, "mod", "tidy")
	cmd.Dir = path.Join(wd, buildDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("unable to compile scale function: %s", output)
		}
		return nil, fmt.Errorf("unable to compile scale function: %w", err)
	}

	tinygoArgs = append([]string{"build", "-o", "scale.wasm", "-target=wasi"}, tinygoArgs...)
	tinygoArgs = append(tinygoArgs, "main.go")

	cmd = exec.Command(tinygoBin, tinygoArgs...)
	cmd.Dir = path.Join(wd, buildDir)

	output, err = cmd.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("unable to compile scale function: %s", output)
		}
		return nil, fmt.Errorf("unable to compile scale function: %w", err)
	}

	data, err := os.ReadFile(path.Join(cmd.Dir, "scale.wasm"))
	if err != nil {
		return nil, fmt.Errorf("unable to read compiled wasm file: %w", err)
	}
	scaleFunc.Function = data

	return scaleFunc, nil
}

func RustBuild(scaleFile *scalefile.ScaleFile, scaleFunc *scalefunc.ScaleFunc, cargoBin string, cargoArgs []string, baseDir string) (*scalefunc.ScaleFunc, error) {
	module := &Module{
		Name:      scaleFile.Name,
		Source:    scaleFile.Source,
		Signature: "scale_signature_http",
	}

	if cargoBin != "" {
		stat, err := os.Stat(cargoBin)
		if err != nil {
			return nil, fmt.Errorf("unable to find cargo binary %s: %w", cargoBin, err)
		}
		if !(stat.Mode()&0111 != 0) {
			return nil, fmt.Errorf("cargo binary %s is not executable", cargoBin)
		}
	} else {
		var err error
		cargoBin, err = exec.LookPath("cargo")
		if err != nil {
			return nil, ErrNoCargo
		}
	}

	g := rustCompile.NewGenerator()

	moduleSourcePath := path.Join(baseDir, module.Source)
	_, err := os.Stat(moduleSourcePath)
	if err != nil {
		return nil, fmt.Errorf("unable to find module %s: %w", moduleSourcePath, err)
	}

	buildDir := path.Join(baseDir, "build")
	defer func() {
		_ = os.RemoveAll(buildDir)
	}()

	err = os.Mkdir(buildDir, 0755)
	if !os.IsExist(err) && err != nil {
		return nil, fmt.Errorf("unable to create build %s directory: %w", buildDir, err)
	}

	file, err := os.OpenFile(path.Join(buildDir, "lib.rs"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to create lib.rs file: %w", err)
	}

	err = g.GenerateRsLib(file, "scale/scale.rs", module.Signature)
	if err != nil {
		return nil, fmt.Errorf("unable to generate lib.rs file: %w", err)
	}

	cargoFile, err := os.OpenFile(path.Join(buildDir, "Cargo.toml"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to create Cargo.toml file: %w", err)
	}

	deps := make([]*scalefile.Dependency, 0, len(scaleFile.Dependencies))
	for _, dep := range scaleFile.Dependencies {
		var d = dep
		deps = append(deps, &d)
	}
	err = g.GenerateRsCargo(cargoFile, deps, module.Signature, "")
	if err != nil {
		return nil, fmt.Errorf("unable to generate Cargo.toml file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to close lib.rs file: %w", err)
	}

	scalePath := path.Join(buildDir, "scale")
	err = os.Mkdir(scalePath, 0755)
	if !os.IsExist(err) && err != nil {
		return nil, fmt.Errorf("unable to create scale source directory %s: %w", scalePath, err)
	}

	scale, err := os.OpenFile(path.Join(scalePath, "scale.rs"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to create scale.rs file: %w", err)
	}

	file, err = os.Open(moduleSourcePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open scale source file: %w", err)
	}

	_, err = io.Copy(scale, file)
	if err != nil {
		return nil, fmt.Errorf("unable to copy scale source file: %w", err)
	}

	err = scale.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to close scale.rs file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to close scale source file: %w", err)
	}

	cargoArgs = append([]string{"build", "--target", "wasm32-wasi", "--manifest-path", "Cargo.toml"}, cargoArgs...)

	cmd := exec.Command(cargoBin, cargoArgs...)
	cmd.Dir = buildDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("unable to compile scale function: %s", output)
		}
		return nil, fmt.Errorf("unable to compile scale function: %w", err)
	}

	outputPath := "target/wasm32-wasi/debug/compile.wasm"
	for _, arg := range cargoArgs {
		if arg == "--release" {
			outputPath = "target/wasm32-wasi/release/compile.wasm"
			break
		}
	}

	data, err := os.ReadFile(path.Join(cmd.Dir, outputPath))
	if err != nil {
		return nil, fmt.Errorf("unable to read compiled wasm file: %w", err)
	}
	scaleFunc.Function = data

	return scaleFunc, nil
}

func TypeScriptBuild(scaleFile *scalefile.ScaleFile, scaleFunc *scalefunc.ScaleFunc, npmBin string, baseDir string) (*scalefunc.ScaleFunc, error) {
	module := &Module{
		Name:      scaleFile.Name,
		Source:    scaleFile.Source,
		Signature: "scale-signature-http",
	}

	if npmBin != "" {
		stat, err := os.Stat(npmBin)
		if err != nil {
			return nil, fmt.Errorf("unable to find npm binary %s: %w", npmBin, err)
		}
		if !(stat.Mode()&0111 != 0) {
			return nil, fmt.Errorf("npm binary %s is not executable", npmBin)
		}
	} else {
		var err error
		npmBin, err = exec.LookPath("npm")
		if err != nil {
			return nil, ErrNoNpm
		}
	}

	g := tsCompile.NewGenerator()

	moduleSourcePath := path.Join(baseDir, module.Source)
	_, err := os.Stat(moduleSourcePath)
	if err != nil {
		return nil, fmt.Errorf("unable to find module %s: %w", moduleSourcePath, err)
	}

	buildDir := path.Join(baseDir, "build")
	defer func() {
		_ = os.RemoveAll(buildDir)
	}()

	err = os.Mkdir(buildDir, 0755)
	if !os.IsExist(err) && err != nil {
		return nil, fmt.Errorf("unable to create build %s directory: %w", buildDir, err)
	}

	file, err := os.OpenFile(path.Join(buildDir, "runner.ts"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to create runner.ts file: %w", err)
	}

	err = g.GenerateRunner(file, "", module.Signature)
	if err != nil {
		return nil, fmt.Errorf("unable to generate runner.ts file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to close runner.ts file: %w", err)
	}

	scale, err := os.OpenFile(path.Join(buildDir, "scale.ts"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to create scale.ts file: %w", err)
	}

	file, err = os.Open(moduleSourcePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open scale source file: %w", err)
	}

	_, err = io.Copy(scale, file)
	if err != nil {
		return nil, fmt.Errorf("unable to copy scale source file: %w", err)
	}

	err = scale.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to close scale.ts file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to close scale source file: %w", err)
	}

	if len(jsbuilderBin) == 0 {
		return nil, fmt.Errorf("No jsbuilder was included for this architecture")
	}

	reader := bytes.NewReader(jsbuilderBin)
	gzr, err := gzip.NewReader(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to create gzip reader for jsbuilder binary: %w", err)
	}

	jsf, err := os.Create(path.Join(buildDir, "jsbuilder"))

	_, err = io.Copy(jsf, gzr)
	if err != nil {
		return nil, fmt.Errorf("unable to extract jsbuilder binary: %w", err)
	}

	err = jsf.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to close jsbuilder binary: %w", err)
	}

	err = os.Chmod(path.Join(buildDir, "jsbuilder"), 0770)
	if err != nil {
		return nil, fmt.Errorf("unable to set permissions on jsbuilder binary: %w", err)
	}

	packageFile, err := os.OpenFile(path.Join(buildDir, "package.json"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to create package.json file: %w", err)
	}

	deps := make([]*scalefile.Dependency, 0, len(scaleFile.Dependencies))
	for _, dep := range scaleFile.Dependencies {
		var d = dep
		deps = append(deps, &d)
	}
	err = g.GeneratePackage(packageFile, deps, module.Signature, "")
	if err != nil {
		return nil, fmt.Errorf("unable to generate package.json file: %w", err)
	}

	cmdInstall := exec.Command(npmBin, "install")
	cmdInstall.Dir = buildDir

	outputInstall, err := cmdInstall.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("unable to compile scale function: %s", outputInstall)
		}
		return nil, fmt.Errorf("unable to compile scale function: %w", err)
	}

	cmdBuild := exec.Command(npmBin, "run", "build")
	cmdBuild.Dir = buildDir

	outputBuild, err := cmdBuild.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("unable to compile scale function: %s", outputBuild)
		}
		return nil, fmt.Errorf("unable to compile scale function: %w", err)
	}

	cmdJSBuilder := exec.Command("./jsbuilder", "dist/runner.js")
	cmdJSBuilder.Dir = buildDir

	outputJSBuilder, err := cmdJSBuilder.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("unable to compile scale function: %s", outputJSBuilder)
		}
		return nil, fmt.Errorf("unable to compile scale function: %w", err)
	}

	data, err := os.ReadFile(path.Join(buildDir, "index.wasm"))
	if err != nil {
		return nil, fmt.Errorf("unable to read compiled wasm file: %w", err)
	}

	scaleFunc.Function = data

	return scaleFunc, nil
}
