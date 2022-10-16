package tests

import (
	"github.com/loopholelabs/scale-go/scalefunc"
	"os"
	"os/exec"
	"testing"
)

type TestCase struct {
	Name   string
	Module string
	Run    func(scalefunc.ScaleFunc, *testing.T)
}

func TestMain(m *testing.M) {
	err := exec.Command("sh", "compile.sh").Run()
	if err != nil {
		panic(err)
	}
	code := m.Run()

	err = exec.Command("sh", "cleanup.sh").Run()
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}
