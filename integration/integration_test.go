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
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/loopholelabs/scale"
	"github.com/loopholelabs/scale/build"
	hostSignature "github.com/loopholelabs/scale/integration/golang_tests/host_signature"
	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/storage"
)

func compileGolangGuest(t *testing.T) *scalefunc.V1BetaSchema {
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

	stb, err := storage.NewBuild(golangCompileDir)
	require.NoError(t, err)

	schema, err := build.LocalGolang(&build.LocalGolangOptions{
		Stdout:          os.Stdout,
		Scalefile:       scf,
		SourceDirectory: golangFunctionDir,
		SignatureSchema: s,
		Storage:         stb,
		Release:         false,
		Target:          build.WASITarget,
	})
	require.NoError(t, err)

	assert.Equal(t, scf.Name, schema.Name)
	assert.Equal(t, scf.Tag, schema.Tag)
	assert.Equal(t, scf.Signature.Name, schema.Signature.Name)
	assert.Equal(t, scf.Signature.Organization, schema.Signature.Organization)
	assert.Equal(t, scf.Signature.Tag, schema.Signature.Tag)
	assert.Equal(t, s, schema.Signature.Schema)
	assert.Equal(t, hex.EncodeToString(hash), schema.Signature.Hash)
	assert.Equal(t, scalefunc.Go, schema.Language)
	assert.Equal(t, 0, len(schema.Extensions))

	return schema
}

func compileRustGuest(t *testing.T) *scalefunc.V1BetaSchema {
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

	stb, err := storage.NewBuild(rustCompileDir)
	require.NoError(t, err)

	schema, err := build.LocalRust(&build.LocalRustOptions{
		Stdout:          os.Stdout,
		Scalefile:       scf,
		SourceDirectory: rustFunctionDir,
		SignatureSchema: s,
		Storage:         stb,
		Release:         false,
		Target:          build.WASITarget,
	})
	require.NoError(t, err)

	assert.Equal(t, scf.Name, schema.Name)
	assert.Equal(t, scf.Tag, schema.Tag)
	assert.Equal(t, scf.Signature.Name, schema.Signature.Name)
	assert.Equal(t, scf.Signature.Organization, schema.Signature.Organization)
	assert.Equal(t, scf.Signature.Tag, schema.Signature.Tag)
	assert.Equal(t, s, schema.Signature.Schema)
	assert.Equal(t, hex.EncodeToString(hash), schema.Signature.Hash)
	assert.Equal(t, scalefunc.Rust, schema.Language)
	assert.Equal(t, 0, len(schema.Extensions))

	return schema
}

func compileTypescriptGuest(t *testing.T) *scalefunc.V1BetaSchema {
	return compileTypescriptGuestFunction(t, "function")
}

func compileTypescriptGuestTimers(t *testing.T) *scalefunc.V1BetaSchema {
	return compileTypescriptGuestFunction(t, "function_timers")
}

func compileTypescriptGuestFunction(t *testing.T, f string) *scalefunc.V1BetaSchema {
	wd, err := os.Getwd()
	require.NoError(t, err)

	s := new(signature.Schema)
	err = s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	hash, err := s.Hash()
	require.NoError(t, err)

	typescriptCompileDir := wd + "/typescript_tests/compile"
	err = os.MkdirAll(typescriptCompileDir, 0755)
	require.NoError(t, err)

	t.Cleanup(func() {
		err = os.RemoveAll(typescriptCompileDir)
		require.NoError(t, err)
	})

	typescriptFunctionDir := wd + "/typescript_tests/" + f
	scf := &scalefile.Schema{
		Version:  scalefile.V1AlphaVersion,
		Name:     "example",
		Tag:      "latest",
		Language: string(scalefunc.TypeScript),
		Signature: scalefile.SignatureSchema{
			Organization: "local",
			Name:         "example",
			Tag:          "latest",
		},
		Function: "example",
	}

	stb, err := storage.NewBuild(typescriptCompileDir)
	require.NoError(t, err)

	schema, err := build.LocalTypescript(&build.LocalTypescriptOptions{
		Stdout:          os.Stdout,
		Scalefile:       scf,
		SourceDirectory: typescriptFunctionDir,
		SignatureSchema: s,
		Storage:         stb,
		Release:         false,
		Target:          build.WASITarget,
	})
	require.NoError(t, err)

	assert.Equal(t, scf.Name, schema.Name)
	assert.Equal(t, scf.Tag, schema.Tag)
	assert.Equal(t, scf.Signature.Name, schema.Signature.Name)
	assert.Equal(t, scf.Signature.Organization, schema.Signature.Organization)
	assert.Equal(t, scf.Signature.Tag, schema.Signature.Tag)
	assert.Equal(t, s, schema.Signature.Schema)
	assert.Equal(t, hex.EncodeToString(hash), schema.Signature.Hash)
	assert.Equal(t, scalefunc.TypeScript, schema.Language)
	assert.Equal(t, 0, len(schema.Extensions))

	return schema
}

func TestGolangHostGolangGuest(t *testing.T) {
	t.Log("Starting TestGolangHostGolangGuest")
	schema := compileGolangGuest(t)
	cfg := scale.NewConfig(hostSignature.New).WithFunction(schema).WithStdout(os.Stdout).WithStderr(os.Stderr)
	runtime, err := scale.New(cfg)
	require.NoError(t, err)

	instance, err := runtime.Instance()
	require.NoError(t, err)

	sig := hostSignature.New()

	ctx := context.Background()
	err = instance.Run(ctx, sig)
	require.NoError(t, err)

	require.Equal(t, "This is a Golang Function", sig.Context.StringField)
}

func TestGolangHostRustGuest(t *testing.T) {
	t.Log("Starting TestGolangHostRustGuest")
	schema := compileRustGuest(t)
	cfg := scale.NewConfig(hostSignature.New).WithFunction(schema).WithStdout(os.Stdout).WithStderr(os.Stderr)
	runtime, err := scale.New(cfg)
	require.NoError(t, err)

	instance, err := runtime.Instance()
	require.NoError(t, err)

	sig := hostSignature.New()

	ctx := context.Background()
	err = instance.Run(ctx, sig)
	require.NoError(t, err)

	require.Equal(t, "This is a Rust Function", sig.Context.StringField)
}

func TestGolangHostTypescriptGuest(t *testing.T) {
	t.Log("Starting TestGolangHostTypescriptGuest")
	schema := compileTypescriptGuest(t)
	cfg := scale.NewConfig(hostSignature.New).WithFunction(schema).WithStdout(os.Stdout).WithStderr(os.Stderr).WithRawOutput(true)
	runtime, err := scale.New(cfg)
	require.NoError(t, err)

	instance, err := runtime.Instance()
	require.NoError(t, err)

	sig := hostSignature.New()

	ctx := context.Background()
	err = instance.Run(ctx, sig)
	require.NoError(t, err)
	require.NotNil(t, sig)
	require.NotNil(t, sig.Context)

	require.Equal(t, "This is a Typescript Function", sig.Context.StringField)
}

func TestGolangHostTypescriptGuestTimers(t *testing.T) {
	t.Log("Starting TestGolangHostTypescriptGuestTimers")
	schema := compileTypescriptGuestTimers(t)
	cfg := scale.NewConfig(hostSignature.New).WithFunction(schema).WithStdout(os.Stdout).WithStderr(os.Stderr).WithRawOutput(true)
	runtime, err := scale.New(cfg)
	require.NoError(t, err)

	instance, err := runtime.Instance()
	require.NoError(t, err)

	sig := hostSignature.New()

	ctx := context.Background()
	err = instance.Run(ctx, sig)
	require.NoError(t, err)
	require.NotNil(t, sig)
	require.NotNil(t, sig.Context)

	require.Equal(t, "This is a Typescript Function.  INTERVAL 18 INTERVAL 19 BLAH INTERVAL 18 INTERVAL 19 TIMEOUT 50 INTERVAL 18 INTERVAL 19 INTERVAL 18 INTERVAL 19 INTERVAL 18 INTERVAL 19 DELAY 100 TIMEOUT 151 DELAY 200", sig.Context.StringField)
}

func TestTypescriptHostTypescriptGuest(t *testing.T) {
	t.Log("Starting TestTypescriptHostTypescriptGuest")
	wd, err := os.Getwd()
	require.NoError(t, err)

	schema := compileTypescriptGuest(t)
	err = os.WriteFile(wd+"/typescript.scale", schema.Encode(), 0644)
	require.NoError(t, err)
	t.Cleanup(func() {
		err = os.Remove(wd + "/typescript.scale")
		require.NoError(t, err)
	})

	cmd := exec.Command("npm", "run", "test", "--", "-t", "test-typescript-host-typescript-guest")
	cmd.Dir = wd
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))
}

func TestTypescriptHostGolangGuest(t *testing.T) {
	t.Log("Starting TestTypescriptHostGolangGuest")
	wd, err := os.Getwd()
	require.NoError(t, err)

	schema := compileGolangGuest(t)
	err = os.WriteFile(wd+"/golang.scale", schema.Encode(), 0644)
	require.NoError(t, err)
	t.Cleanup(func() {
		err = os.Remove(wd + "/golang.scale")
		require.NoError(t, err)
	})

	cmd := exec.Command("npm", "run", "test", "--", "-t", "test-typescript-host-golang-guest")
	cmd.Dir = wd
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))
}

func TestTypescriptHostRustGuest(t *testing.T) {
	t.Log("Starting TestTypescriptHostRustGuest")
	wd, err := os.Getwd()
	require.NoError(t, err)

	schema := compileRustGuest(t)
	err = os.WriteFile(wd+"/rust.scale", schema.Encode(), 0644)
	require.NoError(t, err)
	t.Cleanup(func() {
		err = os.Remove(wd + "/rust.scale")
		require.NoError(t, err)
	})

	cmd := exec.Command("npm", "run", "test", "--", "-t", "test-typescript-host-rust-guest")
	cmd.Dir = wd
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))
}

func TestGolangToGolang(t *testing.T) {
	t.Log("Starting TestGolangToGolang")
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
	t.Log("Starting TestRustToRust")
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

func TestTypescriptToTypescript(t *testing.T) {
	t.Log("Starting TestTypescriptToTypescript")
	wd, err := os.Getwd()
	require.NoError(t, err)

	s := new(signature.Schema)
	err = s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	typescriptDir := wd + "/typescript_tests/signature"
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
	t.Log("Starting TestGolangToRust")
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

func TestGolangToTypescript(t *testing.T) {
	t.Log("Starting TestGolangToTypescript")
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

	typescriptDir := wd + "/typescript_tests/signature"
	cmd = exec.Command("npm", "install", "--save-dev")
	cmd.Dir = typescriptDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("npm", "run", "test", "--", "-t", "test-input")
	cmd.Dir = typescriptDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))
}

func TestRustToGolang(t *testing.T) {
	t.Log("Starting TestRustToGolang")
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

func TestRustToTypescript(t *testing.T) {
	t.Log("Starting TestRustToTypescript")
	wd, err := os.Getwd()
	require.NoError(t, err)

	s := new(signature.Schema)
	err = s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

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

	typescriptDir := wd + "/typescript_tests/signature"
	cmd = exec.Command("npm", "install", "--save-dev")
	cmd.Dir = typescriptDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("npm", "run", "test", "--", "-t", "test-input")
	cmd.Dir = typescriptDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))
}

func TestTypescriptToGolang(t *testing.T) {
	t.Log("Starting TestTypescriptToGolang")
	wd, err := os.Getwd()
	require.NoError(t, err)

	s := new(signature.Schema)
	err = s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	typescriptSignatureDir := wd + "/typescript_tests/signature"
	cmd := exec.Command("npm", "install", "--save-dev")
	cmd.Dir = typescriptSignatureDir
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("npm", "run", "test", "--", "-t", "test-output")
	cmd.Dir = typescriptSignatureDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	golangSignatureDir := wd + "/golang_tests/signature"
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

func TestTypescriptToRust(t *testing.T) {
	t.Log("Starting TestTypescriptToRust")
	wd, err := os.Getwd()
	require.NoError(t, err)

	s := new(signature.Schema)
	err = s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	typescriptSignatureDir := wd + "/typescript_tests/signature"
	cmd := exec.Command("npm", "install", "--save-dev")
	cmd.Dir = typescriptSignatureDir
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	cmd = exec.Command("npm", "run", "test", "--", "-t", "test-output")
	cmd.Dir = typescriptSignatureDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err)
	t.Log(string(out))

	rustSignatureDir := wd + "/rust_tests/signature"
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
