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
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/signature/generator"
	"github.com/loopholelabs/scale/signature/generator/golang"
	"github.com/loopholelabs/scale/signature/generator/rust"
	"github.com/loopholelabs/scale/signature/generator/typescript"
)

const simpleSchema = `
version = "v1alpha"
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

func TestGenerateMasterTestingSchema(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	s := new(signature.Schema)
	err = s.Decode([]byte(signature.MasterTestingSchema))
	require.NoError(t, err)
	require.Equal(t, "ModelWithAllFieldTypes", s.Context)

	guest, err := generator.GenerateGuestLocal(&generator.Options{
		Signature: s,

		GolangPackageImportPath: "signature",
		GolangPackageVersion:    "v0.1.0",

		RustPackageName:    "local_example_latest_guest",
		RustPackageVersion: "v0.1.0",

		TypescriptPackageName:    "local-example-latest-guest",
		TypescriptPackageVersion: "v0.1.0",
	})
	require.NoError(t, err)

	golangSignatureDir := wd + "/golang_tests/signature"
	for _, file := range guest.GolangFiles {
		err = os.WriteFile(golangSignatureDir+"/"+file.Name(), file.Data(), 0644)
		require.NoError(t, err)
	}

	rustSignatureDir := wd + "/rust_tests/signature"
	for _, file := range guest.RustFiles {
		err = os.WriteFile(rustSignatureDir+"/"+file.Name(), file.Data(), 0644)
		require.NoError(t, err)
	}

	typescriptSignatureDir := wd + "/typescript_tests/signature"
	for _, file := range guest.TypescriptFiles {
		err = os.WriteFile(typescriptSignatureDir+"/"+file.Name(), file.Data(), 0644)
		require.NoError(t, err)
	}

	host, err := generator.GenerateHostLocal(&generator.Options{
		Signature: s,

		GolangPackageImportPath: "signature",
		GolangPackageVersion:    "v0.1.0",

		TypescriptPackageName:    "local-example-latest-host",
		TypescriptPackageVersion: "v0.1.0",
	})
	require.NoError(t, err)

	golangSignatureDir = wd + "/golang_tests/host_signature"
	for _, file := range host.GolangFiles {
		if file.Name() != "go.mod" {
			err = os.WriteFile(golangSignatureDir+"/"+file.Name(), file.Data(), 0644)
			require.NoError(t, err)
		}
	}

	typescriptSignatureDir = wd + "/typescript_tests/host_signature"
	for _, file := range host.TypescriptFiles {
		err = os.WriteFile(typescriptSignatureDir+"/"+file.Name(), file.Data(), 0644)
		require.NoError(t, err)
	}
}

func TestGenerateSimpleSchema(t *testing.T) {
	s := new(signature.Schema)
	err := s.Decode([]byte(simpleSchema))
	require.NoError(t, err)

	formatted, err := golang.GenerateTypes(s, "generated")
	require.NoError(t, err)

	err = os.WriteFile("./golang_tests/generated/generated.go", formatted, 0644)
	require.NoError(t, err)

	formatted, err = rust.GenerateTypes(s, "generated")
	require.NoError(t, err)

	err = os.WriteFile("./rust_tests/generated/generated.rs", formatted, 0644)
	require.NoError(t, err)

	transpiled, err := typescript.GenerateTypesTranspiled(s, "generated", "generated.js")
	require.NoError(t, err)

	err = os.WriteFile("./typescript_tests/generated/generated.js", transpiled.Javascript, 0644)
	require.NoError(t, err)

	err = os.WriteFile("./typescript_tests/generated/generated.js.map", transpiled.SourceMap, 0644)
	require.NoError(t, err)

	err = os.WriteFile("./typescript_tests/generated/generated.d.ts", transpiled.Declaration, 0644)
	require.NoError(t, err)
}
