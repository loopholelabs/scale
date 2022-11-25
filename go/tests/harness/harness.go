package harness

import (
	"fmt"
	"github.com/loopholelabs/scale/go/compile"
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

func Setup(t testing.TB, modules []*Module) map[*Module]string {
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

		err = g.GenerateGoMain(file, module.Signature, fmt.Sprintf("github.com/loopholelabs/scale/go/tests/modules/%s/%s-%s-build/scale", module.Name, module.Name, t.Name()))
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
