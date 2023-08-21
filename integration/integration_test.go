//go:build !integration

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
	"encoding/hex"
	"fmt"
	"github.com/loopholelabs/scale/build"
	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/signature/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestExampleSignature(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	s, err := signature.ReadSchema(wd + "/signatures/example.signature")
	require.NoError(t, err)
	hash, err := s.Hash()
	require.NoError(t, err)

	assert.Equal(t, "Example", s.Context)
	assert.Equal(t, 1, len(s.Models))
	assert.Equal(t, "Example", s.Models[0].Name)

	guest, err := generator.GenerateGuestLocal(&generator.GeneratorOptions{
		Signature:             s,
		GolangImportPath:      "signature",
		GolangPackageName:     "signature",
		GolangPackageVersion:  "v0.1.0",
		GolangPolyglotVersion: "v1.1.1",
	})
	require.NoError(t, err)

	golangSignatureDir := wd + "/golang_tests/signature"
	golangCompileDir := wd + "/golang_tests/compile"
	golangFunctionDir := wd + "/golang_tests/function"

	for _, file := range guest.GolangFiles {
		err = os.WriteFile(golangSignatureDir+"/"+file.Name(), file.Data(), 0644)
		require.NoError(t, err)
	}

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
		Signature:        s,
		SourceDirectory:  golangFunctionDir,
		StorageDirectory: golangCompileDir,
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
}
