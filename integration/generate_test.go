//go:build integration && generate

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
	"github.com/loopholelabs/polyglot/version"
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/signature/generator"
	"github.com/loopholelabs/scale/signature/generator/golang"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
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
		Signature:             s,
		GolangImportPath:      "signature",
		GolangPackageName:     "signature",
		GolangPackageVersion:  "v0.1.0",
		GolangPolyglotVersion: version.Version(),
	})
	require.NoError(t, err)

	golangSignatureDir := wd + "/golang_tests/signature"
	for _, file := range guest.GolangFiles {
		err = os.WriteFile(golangSignatureDir+"/"+file.Name(), file.Data(), 0644)
		require.NoError(t, err)
	}
}

func TestGenerateSimpleSchema(t *testing.T) {
	s := new(signature.Schema)
	err := s.Decode([]byte(simpleSchema))
	require.NoError(t, err)

	formatted, err := golang.Generate(s, "generated", "v0.1.0")
	require.NoError(t, err)

	err = os.WriteFile("./golang_tests/generated/generated.go", formatted, 0644)
	require.NoError(t, err)
}
