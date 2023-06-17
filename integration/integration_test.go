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
	golangCompile "github.com/loopholelabs/scale/compile/golang"
	"github.com/loopholelabs/scale/scalefile"
	"github.com/loopholelabs/scale/scalefunc"
	signatureSchema "github.com/loopholelabs/scale/signature"
	golangSignature "github.com/loopholelabs/scale/signature/generator/golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestExampleSignature(t *testing.T) {
	s, err := signatureSchema.ReadSchema("signatures/example.signature")
	require.NoError(t, err)

	require.NoError(t, s.Validate())

	assert.Equal(t, "example", s.Name)
	assert.Equal(t, "latest", s.Tag)
	assert.Equal(t, "Example", s.Context)
	assert.Equal(t, 1, len(s.Models))
	assert.Equal(t, "Example", s.Models[0].Name)

	g, err := golangSignature.New()
	require.NoError(t, err)

	formatted, err := g.Generate(s, "signature", "v0.1.0")
	require.NoError(t, err)

	guest, err := g.GenerateGuest(s, "signature", "v0.1.0")
	require.NoError(t, err)

	signatureModfile, err := g.GenerateModfile("signature", "v1.1.1")
	require.NoError(t, err)

	const golangSignatureDir = "golang_tests/signature"

	err = os.WriteFile(golangSignatureDir+"/generated.go", formatted, 0644)
	require.NoError(t, err)

	err = os.WriteFile(golangSignatureDir+"/guest.go", guest, 0644)
	require.NoError(t, err)

	err = os.WriteFile(golangSignatureDir+"/go.mod", signatureModfile, 0644)
	require.NoError(t, err)

	gc := golangCompile.New()
	scf := &scalefile.Schema{
		Version:  scalefile.V1AlphaVersion,
		Name:     "example",
		Tag:      "latest",
		Language: "go",
		Signature: scalefile.SignatureSchema{
			Organization: "",
			Name:         "example",
			Tag:          "latest",
		},
		Function: "Example",
	}

	mainFile, err := gc.GenerateGoMain(s, scf, "v0.1.0")
	require.NoError(t, err)

	const golangCompileDir = "golang_tests/compile"

	err = os.WriteFile(golangCompileDir+"/main.go", mainFile, 0644)
	require.NoError(t, err)

	//signatureImport string, signatureVersion string, dependencies []*scalefunc.Dependency, version string
	dependencies := []*scalefunc.Dependency{
		{
			Name:     "signature",
			Version:  "v0.1.0",
			Metadata: nil,
		},
		{
			Name:     scf.Name,
			Version:  "v0.1.0",
			Metadata: nil,
		},
	}
	modfile, err := gc.GenerateGoModfile(scf, "../signature", "", "../function", dependencies, "compile", "v0.1.0")
	require.NoError(t, err)

	err = os.WriteFile(golangCompileDir+"/go.mod", modfile, 0644)
	require.NoError(t, err)
}
