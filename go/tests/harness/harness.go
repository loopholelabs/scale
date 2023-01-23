/*
	Copyright 2022 Loophole Labs

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

package harness

import (
	"fmt"
	"github.com/loopholelabs/scale/go/compile"
	rustCompile "github.com/loopholelabs/scale/rust/compile"
	"github.com/loopholelabs/scalefile"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"os/exec"
	"path"
	"testing"
)

type Module struct {
	Path         string
	Name         string
	Signature    string
	Dependencies []*scalefile.Dependency
}

func GoSetup(t testing.TB, modules []*Module, importPath string) map[*Module]string {
	tinygo, err := exec.LookPath("tinygo")
	require.NoError(t, err, "tinygo not found in path")

	t.Cleanup(func() {
		for _, module := range modules {
			moduleDir := path.Dir(module.Path)
			err := os.RemoveAll(path.Join(moduleDir, fmt.Sprintf("%s-%s-build", module.Name, t.Name())))
			if !os.IsNotExist(err) {
				require.NoError(t, err, fmt.Sprintf("failed to remove module %s", module.Name))
			}
		}
	})

	g := compile.NewGenerator()

	generated := make(map[*Module]string)

	for _, module := range modules {
		_, err = os.Stat(module.Path)
		require.NoError(t, err, fmt.Sprintf("module %s not found", module.Name))

		moduleDir := path.Dir(module.Path)

		err = os.Mkdir(path.Join(moduleDir, fmt.Sprintf("%s-%s-build", module.Name, t.Name())), 0755)
		if !os.IsExist(err) {
			require.NoError(t, err, fmt.Sprintf("failed to create build directory for scale function %s", module.Name))
		}

		file, err := os.OpenFile(path.Join(moduleDir, fmt.Sprintf("%s-%s-build", module.Name, t.Name()), "main.go"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		require.NoError(t, err, fmt.Sprintf("failed to create main.go for scale function %s", module.Name))

		err = g.GenerateGoMain(file, module.Signature, fmt.Sprintf("%s/%s/%s-%s-build/scale", importPath, module.Name, module.Name, t.Name()))
		require.NoError(t, err, fmt.Sprintf("failed to generate main.go for scale function %s", module.Name))

		err = file.Close()
		require.NoError(t, err, fmt.Sprintf("failed to close main.go for scale function %s", module.Name))

		err = os.Mkdir(path.Join(moduleDir, fmt.Sprintf("%s-%s-build", module.Name, t.Name()), "scale"), 0755)
		if !os.IsExist(err) {
			require.NoError(t, err, fmt.Sprintf("failed to create scale directory for scale function %s", module.Name))
		}

		scale, err := os.OpenFile(path.Join(moduleDir, fmt.Sprintf("%s-%s-build", module.Name, t.Name()), "scale", "scale.go"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		require.NoError(t, err, fmt.Sprintf("failed to create scale.go for scale function %s", module.Name))

		file, err = os.Open(module.Path)
		require.NoError(t, err, fmt.Sprintf("failed to open scale function %s", module.Name))

		_, err = io.Copy(scale, file)
		require.NoError(t, err, fmt.Sprintf("failed to copy scale function %s", module.Name))

		err = scale.Close()
		require.NoError(t, err, fmt.Sprintf("failed to close scale.go for scale function %s", module.Name))

		err = file.Close()
		require.NoError(t, err, fmt.Sprintf("failed to close scale function %s", module.Name))

		wd, err := os.Getwd()
		require.NoError(t, err, fmt.Sprintf("failed to get working directory for scale function %s", module.Name))

		cmd := exec.Command(tinygo, "build", "-o", fmt.Sprintf("%s-%s.wasm", module.Name, t.Name()), "-scheduler=none", "-target=wasi", "--no-debug", "main.go")
		cmd.Dir = path.Join(wd, moduleDir, fmt.Sprintf("%s-%s-build", module.Name, t.Name()))

		err = cmd.Run()
		require.NoError(t, err, fmt.Sprintf("failed to build module %s", module.Name))

		generated[module] = path.Join(cmd.Dir, fmt.Sprintf("%s-%s.wasm", module.Name, t.Name()))
	}

	return generated
}

func RustSetup(t testing.TB, modules []*Module, importPath string) map[*Module]string {
	cargo, err := exec.LookPath("cargo")
	require.NoError(t, err, "cargo not found in path")

	t.Cleanup(func() {
		for _, module := range modules {
			moduleDir := path.Dir(module.Path)
			err := os.RemoveAll(path.Join(moduleDir, fmt.Sprintf("%s-%s-build", module.Name, t.Name())))
			if !os.IsNotExist(err) {
				require.NoError(t, err, fmt.Sprintf("failed to remove module %s", module.Name))
			}
		}
	})

	g := rustCompile.NewGenerator()

	generated := make(map[*Module]string)

	for _, module := range modules {
		_, err = os.Stat(module.Path)
		require.NoError(t, err, fmt.Sprintf("module %s not found", module.Name))

		moduleDir := path.Dir(module.Path)

		err = os.Mkdir(path.Join(moduleDir, fmt.Sprintf("%s-%s-build", module.Name, t.Name())), 0755)
		if !os.IsExist(err) {
			require.NoError(t, err, fmt.Sprintf("failed to create build directory for scale function %s", module.Name))
		}

		file, err := os.OpenFile(path.Join(moduleDir, fmt.Sprintf("%s-%s-build", module.Name, t.Name()), "lib.rs"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		require.NoError(t, err, fmt.Sprintf("failed to create lib.rs for scale function %s", module.Name))

		err = g.GenerateLibRs(file, module.Signature, importPath)

		cargoFile, err := os.OpenFile(path.Join(moduleDir, fmt.Sprintf("%s-%s-build", module.Name, t.Name()), "Cargo.toml"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		dependencies := []*scalefile.Dependency{}
		err = g.GenerateCargoTomlfile(cargoFile, dependencies)
		require.NoError(t, err, fmt.Sprintf("failed to generate lib.rs for scale function %s", module.Name))

		err = file.Close()
		require.NoError(t, err, fmt.Sprintf("failed to close lib.rs for scale function %s", module.Name))

		err = os.Mkdir(path.Join(moduleDir, fmt.Sprintf("%s-%s-build", module.Name, t.Name()), "scale"), 0755)
		if !os.IsExist(err) {
			require.NoError(t, err, fmt.Sprintf("failed to create scale directory for scale function %s", module.Name))
		}

		scale, err := os.OpenFile(path.Join(moduleDir, fmt.Sprintf("%s-%s-build", module.Name, t.Name()), "scale", "scale.rs"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		require.NoError(t, err, fmt.Sprintf("failed to create scale.go for scale function %s", module.Name))

		file, err = os.Open(module.Path)
		require.NoError(t, err, fmt.Sprintf("failed to open scale function %s", module.Name))

		_, err = io.Copy(scale, file)
		require.NoError(t, err, fmt.Sprintf("failed to copy scale function %s", module.Name))

		err = scale.Close()
		require.NoError(t, err, fmt.Sprintf("failed to close scale.go for scale function %s", module.Name))

		err = file.Close()
		require.NoError(t, err, fmt.Sprintf("failed to close scale function %s", module.Name))

		wd, err := os.Getwd()
		require.NoError(t, err, fmt.Sprintf("failed to get working directory for scale function %s", module.Name))

		cmd := exec.Command(cargo, "build", "--target", "wasm32-unknown-unknown", "--manifest-path", "Cargo.toml")

		cmd.Dir = path.Join(wd, moduleDir, fmt.Sprintf("%s-%s-build", module.Name, t.Name()))

		err = cmd.Run()
		require.NoError(t, err, fmt.Sprintf("wd:  %s", cmd.Dir))

		generated[module] = path.Join(cmd.Dir, "target/wasm32-unknown-unknown/debug/compile.wasm")
	}

	return generated
}
