//go:build integration

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

package integration

import (
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/signature/generator/golang"
	"github.com/loopholelabs/scale/signature/generator/rust"
	"github.com/loopholelabs/scale/signature/generator/typescript"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"testing"
)

const simpleSchema = `
version = "v1alpha"
name = "simpleSchema"
tag = "simpleSchematag"
context = "Context"
model Context {
	int32 A {
		default = 0
	}

	int32 B {
		default = 0
	}

	int32 C {
		default = 0
	}
}
`

func TestSimpleSchema(t *testing.T) {
	s := new(signature.Schema)
	err := s.Decode([]byte(simpleSchema))
	require.NoError(t, err)

	require.NoError(t, s.Validate())

	const golangDir = "./golang_tests"
	formatted, err := golang.Generate(s, "golang_tests", "v0.1.0")
	require.NoError(t, err)

	err = os.WriteFile(golangDir+"/generated.go", formatted, 0644)
	require.NoError(t, err)
}

func TestGolangToGolang(t *testing.T) {
	s := new(signature.Schema)
	err := s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	require.NoError(t, s.Validate())

	const golangDir = "./golang_tests"

	formatted, err := golang.Generate(s, "golang_tests", "v0.1.0")
	require.NoError(t, err)

	guest, err := golang.GenerateGuest(s, "golang_tests", "v0.1.0")
	require.NoError(t, err)

	err = os.WriteFile(golangDir+"/generated.go", formatted, 0644)
	require.NoError(t, err)

	err = os.WriteFile(golangDir+"/guest.go", guest, 0644)
	require.NoError(t, err)

	cmd := exec.Command("go", "test", "./...", "-v", "--tags=integration,golang", "-run", "TestOutput")
	cmd.Dir = golangDir
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("go", "test", "./...", "-v", "--tags=integration,golang", "-run", "TestInput")
	cmd.Dir = golangDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))
}

func TestRustToRust(t *testing.T) {
	s := new(signature.Schema)
	err := s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	require.NoError(t, s.Validate())

	const rustDir = "./rust_tests"

	formatted, err := rust.Generate(s, "rust_tests", "v0.1.0")
	require.NoError(t, err)

	err = os.WriteFile(rustDir+"/generated.rs", formatted, 0644)
	require.NoError(t, err)

	cmd := exec.Command("cargo", "test", "test_output")
	cmd.Dir = rustDir
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("cargo", "test", "test_input")
	cmd.Dir = rustDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))
}

func TestTypescriptToTypescript(t *testing.T) {
	s := new(signature.Schema)
	err := s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	require.NoError(t, s.Validate())

	const typescriptDir = "./typescript_tests"

	formatted, err := typescript.Generate(s, "typescript_tests", "v0.1.0")
	require.NoError(t, err)

	err = os.WriteFile(typescriptDir+"/generated.ts", formatted, 0644)
	require.NoError(t, err)

	cmd := exec.Command("npm", "run", "test", "--", "-t", "test-output")
	cmd.Dir = typescriptDir
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("npm", "run", "test", "--", "-t", "test-input")
	cmd.Dir = typescriptDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))
}

func TestGolangToRust(t *testing.T) {
	s := new(signature.Schema)
	err := s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	require.NoError(t, s.Validate())

	const golangDir = "./golang_tests"

	formatted, err := golang.Generate(s, "golang_tests", "v0.1.0")
	require.NoError(t, err)

	err = os.WriteFile(golangDir+"/generated.go", formatted, 0644)
	require.NoError(t, err)

	const rustDir = "./rust_tests"

	formatted, err = rust.Generate(s, "rust_tests", "v0.1.0")
	require.NoError(t, err)

	err = os.WriteFile(rustDir+"/generated.rs", formatted, 0644)
	require.NoError(t, err)

	cmd := exec.Command("go", "test", "./...", "-v", "--tags=integration,golang", "-run", "TestOutput")
	cmd.Dir = golangDir
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("cargo", "test", "test_input")
	cmd.Dir = rustDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))
}

func TestGolangToTypescript(t *testing.T) {
	s := new(signature.Schema)
	err := s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	require.NoError(t, s.Validate())

	const golangDir = "./golang_tests"

	formatted, err := golang.Generate(s, "golang_tests", "v0.1.0")
	require.NoError(t, err)

	err = os.WriteFile(golangDir+"/generated.go", formatted, 0644)
	require.NoError(t, err)

	const typescriptDir = "./typescript_tests"

	formatted, err = typescript.Generate(s, "typescript_tests", "v0.1.0")
	require.NoError(t, err)

	err = os.WriteFile(typescriptDir+"/generated.ts", formatted, 0644)
	require.NoError(t, err)

	cmd := exec.Command("go", "test", "./...", "-v", "--tags=integration,golang", "-run", "TestOutput")
	cmd.Dir = golangDir
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("npm", "run", "test", "--", "-t", "test-input")
	cmd.Dir = typescriptDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))
}

func TestRustToGolang(t *testing.T) {
	s := new(signature.Schema)
	err := s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	require.NoError(t, s.Validate())

	const golangDir = "./golang_tests"

	formatted, err := golang.Generate(s, "golang_tests", "v0.1.0")
	require.NoError(t, err)

	err = os.WriteFile(golangDir+"/generated.go", formatted, 0644)
	require.NoError(t, err)

	const rustDir = "./rust_tests"

	formatted, err = rust.Generate(s, "rust_tests", "v0.1.0")
	require.NoError(t, err)

	err = os.WriteFile(rustDir+"/generated.rs", formatted, 0644)
	require.NoError(t, err)

	cmd := exec.Command("cargo", "test", "test_output")
	cmd.Dir = rustDir
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("go", "test", "./...", "-v", "--tags=integration,golang", "-run", "TestInput")
	cmd.Dir = golangDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))
}

func TestRustToTypescript(t *testing.T) {
	s := new(signature.Schema)
	err := s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	require.NoError(t, s.Validate())

	const typescriptDir = "./typescript_tests"

	formatted, err := typescript.Generate(s, "typescript_tests", "v0.1.0")
	require.NoError(t, err)

	err = os.WriteFile(typescriptDir+"/generated.ts", formatted, 0644)
	require.NoError(t, err)

	const rustDir = "./rust_tests"

	formatted, err = rust.Generate(s, "rust_tests", "v0.1.0")
	require.NoError(t, err)

	err = os.WriteFile(rustDir+"/generated.rs", formatted, 0644)
	require.NoError(t, err)

	cmd := exec.Command("cargo", "test", "test_output")
	cmd.Dir = rustDir
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("npm", "run", "test", "--", "-t", "test-input")
	cmd.Dir = typescriptDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))
}

func TestTypescriptToGolang(t *testing.T) {
	s := new(signature.Schema)
	err := s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	require.NoError(t, s.Validate())

	const typescriptDir = "./typescript_tests"

	formatted, err := typescript.Generate(s, "typescript_tests", "v0.1.0")
	require.NoError(t, err)

	err = os.WriteFile(typescriptDir+"/generated.ts", formatted, 0644)
	require.NoError(t, err)

	const golangDir = "./golang_tests"

	formatted, err = golang.Generate(s, "golang_tests", "v0.1.0")
	require.NoError(t, err)

	err = os.WriteFile(golangDir+"/generated.go", formatted, 0644)
	require.NoError(t, err)

	cmd := exec.Command("npm", "run", "test", "--", "-t", "test-output")
	cmd.Dir = typescriptDir
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("go", "test", "./...", "-v", "--tags=integration,golang", "-run", "TestInput")
	cmd.Dir = golangDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))
}

func TestTypescriptToRust(t *testing.T) {
	s := new(signature.Schema)
	err := s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	require.NoError(t, s.Validate())

	const typescriptDir = "./typescript_tests"

	formatted, err := typescript.Generate(s, "typescript_tests", "v0.1.0")
	require.NoError(t, err)

	err = os.WriteFile(typescriptDir+"/generated.ts", formatted, 0644)
	require.NoError(t, err)

	const rustDir = "./rust_tests"

	formatted, err = rust.Generate(s, "rust_tests", "v0.1.0")
	require.NoError(t, err)

	err = os.WriteFile(rustDir+"/generated.rs", formatted, 0644)
	require.NoError(t, err)

	cmd := exec.Command("npm", "run", "test", "--", "-t", "test-output")
	cmd.Dir = typescriptDir
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("cargo", "test", "test_input")
	cmd.Dir = rustDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))
}
