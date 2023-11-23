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
	"path/filepath"
	"testing"

	"github.com/loopholelabs/scale"
	"github.com/loopholelabs/scale/build"
	"github.com/loopholelabs/scale/extension"
	hostExtension "github.com/loopholelabs/scale/integration/golang_ext_tests/host_extension"
	hostSignature "github.com/loopholelabs/scale/integration/golang_ext_tests/host_signature"
	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const extension_schema = `	version = "v1alpha"
	
function New {
	params = "stringval"
	return = "Example"	
}

interface Example {
	function Hello {
		params = "stringval"
		return = "stringval"
	}
}

function World {
	params = "stringval"
	return = "stringval"
}

model stringval {
	string value {
		default = ""
	}
}	
`

func compileExtGolangGuest(t *testing.T) *scalefunc.V1BetaSchema {
	wd, err := os.Getwd()
	require.NoError(t, err)

	s := new(signature.Schema)
	err = s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	hash, err := s.Hash()
	require.NoError(t, err)

	golangCompileDir := wd + "/golang_ext_tests/compile"
	err = os.MkdirAll(golangCompileDir, 0755)
	require.NoError(t, err)

	t.Cleanup(func() {
		err = os.RemoveAll(golangCompileDir)
		require.NoError(t, err)
	})

	ex := new(extension.Schema)
	err = ex.Decode([]byte(extension_schema))

	require.NoError(t, err)

	golangFunctionDir := wd + "/golang_ext_tests/function"
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
		Extensions: []scalefile.ExtensionSchema{
			{
				Organization: "local",
				Name:         "example",
				Tag:          "latest",
			},
		},
	}

	stb, err := storage.NewBuild(golangCompileDir)
	require.NoError(t, err)

	ext_dir, err := filepath.Abs("golang_ext_tests/extension")
	require.NoError(t, err)

	extensionInfo := []extension.Info{
		{
			Name:    "extension",
			Path:    ext_dir,
			Version: "v0.1.0",
		},
	}

	extensionSchemas := []*extension.Schema{
		ex,
	}

	schema, err := build.LocalGolang(&build.LocalGolangOptions{
		Stdout:           os.Stdout,
		Scalefile:        scf,
		SourceDirectory:  golangFunctionDir,
		SignatureSchema:  s,
		Storage:          stb,
		Release:          false,
		Target:           build.WASITarget,
		Extensions:       extensionInfo,
		ExtensionSchemas: extensionSchemas,
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

	assert.Equal(t, 1, len(schema.Extensions))

	return schema
}

func compileExtRustGuest(t *testing.T) *scalefunc.V1BetaSchema {
	wd, err := os.Getwd()
	require.NoError(t, err)

	s := new(signature.Schema)
	err = s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)

	hash, err := s.Hash()
	require.NoError(t, err)

	rustCompileDir := wd + "/rust_ext_tests/compile"
	err = os.MkdirAll(rustCompileDir, 0755)
	require.NoError(t, err)

	t.Cleanup(func() {
		err = os.RemoveAll(rustCompileDir)
		require.NoError(t, err)
	})

	ex := new(extension.Schema)
	err = ex.Decode([]byte(extension_schema))

	rustFunctionDir := wd + "/rust_ext_tests/function"
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
		Extensions: []scalefile.ExtensionSchema{
			{
				Organization: "local",
				Name:         "example",
				Tag:          "latest",
			},
		},
	}

	stb, err := storage.NewBuild(rustCompileDir)
	require.NoError(t, err)

	extensionSchemas := []*extension.Schema{
		ex,
	}

	schema, err := build.LocalRust(&build.LocalRustOptions{
		Stdout:           os.Stdout,
		Scalefile:        scf,
		SourceDirectory:  rustFunctionDir,
		SignatureSchema:  s,
		Storage:          stb,
		Release:          false,
		Target:           build.WASITarget,
		ExtensionSchemas: extensionSchemas,
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
	assert.Equal(t, 1, len(schema.Extensions))

	return schema
}

/**
 * Implementation of the simple extension
 *
 */
type ExtensionImpl struct{}

func (ei ExtensionImpl) New(p *hostExtension.Stringval) (hostExtension.Example, error) {
	return &ExtensionExample{}, nil
}

func (ei ExtensionImpl) World(p *hostExtension.Stringval) (hostExtension.Stringval, error) {
	return hostExtension.Stringval{Value: "Return World"}, nil
}

type ExtensionExample struct{}

func (ee ExtensionExample) Hello(p *hostExtension.Stringval) (hostExtension.Stringval, error) {
	return hostExtension.Stringval{Value: "Return Hello"}, nil
}

func TestExtGolangHostGolangGuest(t *testing.T) {

	ext_impl := &ExtensionImpl{}

	e := hostExtension.New(ext_impl)

	t.Log("Starting TestExtGolangHostGolangGuest")
	schema := compileExtGolangGuest(t)
	cfg := scale.NewConfig(hostSignature.New).WithFunction(schema).WithStdout(os.Stdout).WithStderr(os.Stderr).WithExtension(e)
	runtime, err := scale.New(cfg)
	require.NoError(t, err)

	instance, err := runtime.Instance()
	require.NoError(t, err)

	sig := hostSignature.New()

	ctx := context.Background()
	err = instance.Run(ctx, sig)
	require.NoError(t, err)

	require.Equal(t, "This is a Golang Function. Extension New().Hello()=Return Hello World()=Return World", sig.Context.StringField)
}

func TestExtGolangHostRustGuest(t *testing.T) {
	ext_impl := &ExtensionImpl{}

	e := hostExtension.New(ext_impl)

	t.Log("Starting TestGolangHostRustGuest")
	schema := compileExtRustGuest(t)
	cfg := scale.NewConfig(hostSignature.New).WithFunction(schema).WithStdout(os.Stdout).WithStderr(os.Stderr).WithExtension(e)
	runtime, err := scale.New(cfg)
	require.NoError(t, err)

	instance, err := runtime.Instance()
	require.NoError(t, err)

	sig := hostSignature.New()

	ctx := context.Background()
	err = instance.Run(ctx, sig)
	require.NoError(t, err)

	require.Equal(t, "This is a Rust Function. Extension New().Hello()=Return Hello World()=Return World", sig.Context.StringField)
}
