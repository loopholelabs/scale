//go:build integration && !generate

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
	"context"
	"encoding/hex"
	"fmt"
	"github.com/loopholelabs/scale"
	"github.com/loopholelabs/scale/build"
	hostSignature "github.com/loopholelabs/scale/integration/golang_tests/host_signature"
	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"testing"
)

func TestGolangE2E(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	s := new(signature.Schema)
	err = s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	hash, err := s.Hash()
	require.NoError(t, err)

	golangCompileDir := wd + "/golang_tests/compile"
	err = os.MkdirAll(golangCompileDir, 0755)
	require.NoError(t, err)

	t.Cleanup(func() {
		err = os.RemoveAll(golangCompileDir)
		require.NoError(t, err)
	})

	golangFunctionDir := wd + "/golang_tests/function"
	scf := &scalefile.Schema{
		Version:  scalefile.V1AlphaVersion,
		Name:     "example",
		Tag:      "latest",
		Language: string(scalefunc.Go),
		Signature: scalefile.SignatureSchema{
			Organization: "local",
			Name:         "example",
			Tag:          "latest",
		},
		Function: "Example",
	}

	schema, err := build.LocalGolang(&build.LocalGolangOptions{
		Scalefile:        scf,
		SignaturePath:    wd + "/golang_tests/signature",
		SignatureSchema:  s,
		SourceDirectory:  golangFunctionDir,
		StorageDirectory: golangCompileDir,
		Release:          false,
		Target:           build.WASITarget,
	})
	require.NoError(t, err)

	assert.Equal(t, scalefunc.V1Alpha, schema.Version)
	assert.Equal(t, scf.Name, schema.Name)
	assert.Equal(t, scf.Tag, schema.Tag)
	assert.Equal(t, fmt.Sprintf("%s/%s:%s", scf.Signature.Organization, scf.Signature.Name, scf.Signature.Tag), schema.SignatureName)
	assert.Equal(t, s, schema.SignatureSchema)
	assert.Equal(t, hex.EncodeToString(hash), schema.SignatureHash)
	assert.Equal(t, scalefunc.Go, schema.Language)
	assert.Equal(t, 0, len(schema.Dependencies))

	ctx := context.Background()
	runtime, err := scale.New(ctx, hostSignature.New, []*scalefunc.Schema{schema})
	require.NoError(t, err)

	instance, err := runtime.Instance()
	require.NoError(t, err)

	sig := hostSignature.New()

	err = instance.Run(ctx, sig)
	require.NoError(t, err)

	require.Equal(t, "This is a Golang Function", sig.Context.StringField)
}

func TestRustE2E(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	s := new(signature.Schema)
	err = s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	hash, err := s.Hash()
	require.NoError(t, err)

	rustCompileDir := wd + "/rust_tests/compile"
	err = os.MkdirAll(rustCompileDir, 0755)
	require.NoError(t, err)

	t.Cleanup(func() {
		err = os.RemoveAll(rustCompileDir)
		require.NoError(t, err)
	})

	rustFunctionDir := wd + "/rust_tests/function"
	scf := &scalefile.Schema{
		Version:  scalefile.V1AlphaVersion,
		Name:     "example",
		Tag:      "latest",
		Language: string(scalefunc.Rust),
		Signature: scalefile.SignatureSchema{
			Organization: "local",
			Name:         "example",
			Tag:          "latest",
		},
		Function: "example",
	}

	schema, err := build.LocalRust(&build.LocalRustOptions{
		Scalefile:        scf,
		SignaturePath:    wd + "/rust_tests/signature",
		SignatureSchema:  s,
		SourceDirectory:  rustFunctionDir,
		StorageDirectory: rustCompileDir,
		Release:          false,
		Target:           build.WASITarget,
	})
	require.NoError(t, err)

	assert.Equal(t, scalefunc.V1Alpha, schema.Version)
	assert.Equal(t, scf.Name, schema.Name)
	assert.Equal(t, scf.Tag, schema.Tag)
	assert.Equal(t, fmt.Sprintf("%s/%s:%s", scf.Signature.Organization, scf.Signature.Name, scf.Signature.Tag), schema.SignatureName)
	assert.Equal(t, s, schema.SignatureSchema)
	assert.Equal(t, hex.EncodeToString(hash), schema.SignatureHash)
	assert.Equal(t, scalefunc.Go, schema.Language)
	assert.Equal(t, 0, len(schema.Dependencies))

	ctx := context.Background()
	runtime, err := scale.New(ctx, hostSignature.New, []*scalefunc.Schema{schema})
	require.NoError(t, err)

	instance, err := runtime.Instance()
	require.NoError(t, err)

	sig := hostSignature.New()

	err = instance.Run(ctx, sig)
	require.NoError(t, err)

	require.Equal(t, "This is a Rust Function", sig.Context.StringField)
}

func TestGolangToGolang(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	s := new(signature.Schema)
	err = s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	golangSignatureDir := wd + "/golang_tests/signature"
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = golangSignatureDir
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("go", "test", "./...", "-v", "--tags=integration,golang", "-run", "TestOutput")
	cmd.Dir = golangSignatureDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("go", "test", "./...", "-v", "--tags=integration,golang", "-run", "TestInput")
	cmd.Dir = golangSignatureDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))
}

func TestRustToRust(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	s := new(signature.Schema)
	err = s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	rustSignatureDir := wd + "/rust_tests/signature"
	cmd := exec.Command("cargo", "test", "test_output")
	cmd.Dir = rustSignatureDir
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("cargo", "test", "test_input")
	cmd.Dir = rustSignatureDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))
}

//func TestTypescriptToTypescript(t *testing.T) {
//	s := new(signature.Schema)
//	err := s.Decode([]byte(signature.MasterTestingSchema))
//	require.NoError(t, err)
//
//	const typescriptDir = "./typescript_tests"
//
//	formatted, err := typescript.Generate(s, "typescript_tests", "v0.1.0")
//	require.NoError(t, err)
//
//	err = os.WriteFile(typescriptDir+"/generated.ts", formatted, 0644)
//	require.NoError(t, err)
//
//	cmd := exec.Command("npm", "run", "test", "--", "-t", "test-output")
//	cmd.Dir = typescriptDir
//	out, err := cmd.CombinedOutput()
//	assert.NoError(t, err)
//	t.Log(string(out))
//
//	cmd = exec.Command("npm", "run", "test", "--", "-t", "test-input")
//	cmd.Dir = typescriptDir
//	out, err = cmd.CombinedOutput()
//	assert.NoError(t, err)
//	t.Log(string(out))
//}

func TestGolangToRust(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	s := new(signature.Schema)
	err = s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	golangSignatureDir := wd + "/golang_tests/signature"
	rustSignatureDir := wd + "/rust_tests/signature"

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = golangSignatureDir
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("go", "test", "./...", "-v", "--tags=integration,golang", "-run", "TestOutput")
	cmd.Dir = golangSignatureDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("cargo", "check")
	cmd.Dir = rustSignatureDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("cargo", "test", "test_input")
	cmd.Dir = rustSignatureDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))
}

//func TestGolangToTypescript(t *testing.T) {
//	wd, err := os.Getwd()
//	require.NoError(t, err)
//
//	s := new(signature.Schema)
//	err = s.Decode([]byte(signature.MasterTestingSchema))
//	require.NoError(t, err)
//
//	golangSignatureDir := wd + "/golang_tests/signature"
//
//	const typescriptDir = "./typescript_tests"
//
//	formatted, err := typescript.Generate(s, "typescript_tests", "v0.1.0")
//	require.NoError(t, err)
//
//	err = os.WriteFile(typescriptDir+"/generated.ts", formatted, 0644)
//	require.NoError(t, err)
//
//	cmd := exec.Command("go", "mod", "tidy")
//	cmd.Dir = golangSignatureDir
//	out, err := cmd.CombinedOutput()
//	assert.NoError(t, err)
//	t.Log(string(out))
//
//	cmd = exec.Command("go", "test", "./...", "-v", "--tags=integration,golang", "-run", "TestOutput")
//	cmd.Dir = golangSignatureDir
//	out, err = cmd.CombinedOutput()
//	assert.NoError(t, err)
//	t.Log(string(out))
//
//	cmd = exec.Command("npm", "run", "test", "--", "-t", "test-input")
//	cmd.Dir = typescriptDir
//	out, err = cmd.CombinedOutput()
//	assert.NoError(t, err)
//	t.Log(string(out))
//}

func TestRustToGolang(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	s := new(signature.Schema)
	err = s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	golangSignatureDir := wd + "/golang_tests/signature"
	rustSignatureDir := wd + "/rust_tests/signature"

	cmd := exec.Command("cargo", "check")
	cmd.Dir = rustSignatureDir
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("cargo", "test", "test_output")
	cmd.Dir = rustSignatureDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = golangSignatureDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("go", "test", "./...", "-v", "--tags=integration,golang", "-run", "TestInput")
	cmd.Dir = golangSignatureDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))
}

//func TestRustToTypescript(t *testing.T) {
//	wd, err := os.Getwd()
//	require.NoError(t, err)
//
//	s := new(signature.Schema)
//	err = s.Decode([]byte(signature.MasterTestingSchema))
//	require.NoError(t, err)
//
//	rustSignatureDir := wd + "/rust_tests/signature"
//
//	const typescriptDir = "./typescript_tests"
//
//	formatted, err := typescript.Generate(s, "typescript_tests", "v0.1.0")
//	require.NoError(t, err)
//
//	err = os.WriteFile(typescriptDir+"/generated.ts", formatted, 0644)
//	require.NoError(t, err)
//
//	cmd := exec.Command("cargo", "check")
//	cmd.Dir = rustSignatureDir
//	out, err := cmd.CombinedOutput()
//	assert.NoError(t, err)
//	t.Log(string(out))
//
//	cmd = exec.Command("cargo", "test", "test_output")
//	cmd.Dir = rustSignatureDir
//	out, err = cmd.CombinedOutput()
//	assert.NoError(t, err)
//	t.Log(string(out))
//
//	cmd = exec.Command("npm", "run", "test", "--", "-t", "test-input")
//	cmd.Dir = typescriptDir
//	out, err = cmd.CombinedOutput()
//	assert.NoError(t, err)
//	t.Log(string(out))
//}

//func TestTypescriptToGolang(t *testing.T) {
//	wd, err := os.Getwd()
//	require.NoError(t, err)
//
//	s := new(signature.Schema)
//	err = s.Decode([]byte(signature.MasterTestingSchema))
//	require.NoError(t, err)
//
//	golangSignatureDir := wd + "/golang_tests/signature"
//	typescriptSignatureDir := wd + "/typescript_tests/signature"
//
//	formatted, err := typescript.Generate(s, "typescript_tests", "v0.1.0")
//	require.NoError(t, err)
//
//	err = os.WriteFile(typescriptSignatureDir+"/generated.ts", formatted, 0644)
//	require.NoError(t, err)
//
//	cmd := exec.Command("npm", "run", "test", "--", "-t", "test-output")
//	cmd.Dir = typescriptSignatureDir
//	out, err := cmd.CombinedOutput()
//	assert.NoError(t, err)
//	t.Log(string(out))
//
//	cmd = exec.Command("go", "mod", "tidy")
//	cmd.Dir = golangSignatureDir
//	out, err = cmd.CombinedOutput()
//	assert.NoError(t, err)
//	t.Log(string(out))
//
//	cmd = exec.Command("go", "test", "./...", "-v", "--tags=integration,golang", "-run", "TestInput")
//	cmd.Dir = golangSignatureDir
//	out, err = cmd.CombinedOutput()
//	assert.NoError(t, err)
//	t.Log(string(out))
//}

//func TestTypescriptToRust(t *testing.T) {
//	wd, err := os.Getwd()
//	require.NoError(t, err)
//
//	s := new(signature.Schema)
//	err = s.Decode([]byte(signature.MasterTestingSchema))
//	require.NoError(t, err)
//
//	rustSignatureDir := wd + "/rust_tests/signature"
//	typescriptSignatureDir := wd + "/typescript_tests/signature"
//
//	formatted, err := typescript.Generate(s, "typescript_tests", "v0.1.0")
//	require.NoError(t, err)
//
//	err = os.WriteFile(typescriptSignatureDir+"/generated.ts", formatted, 0644)
//	require.NoError(t, err)
//
//	cmd := exec.Command("npm", "run", "test", "--", "-t", "test-output")
//	cmd.Dir = typescriptSignatureDir
//	out, err := cmd.CombinedOutput()
//	assert.NoError(t, err)
//	t.Log(string(out))
//
//	cmd = exec.Command("cargo", "check")
//	cmd.Dir = rustSignatureDir
//	out, err = cmd.CombinedOutput()
//	assert.NoError(t, err)
//	t.Log(string(out))
//
//	cmd = exec.Command("cargo", "test", "test_input")
//	cmd.Dir = rustSignatureDir
//	out, err = cmd.CombinedOutput()
//	assert.NoError(t, err)
//	t.Log(string(out))
//}
