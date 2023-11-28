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

	"github.com/loopholelabs/scale/extension"
	"github.com/loopholelabs/scale/extension/generator"
)

const extensionSchema = `	version = "v1alpha"
	
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

func TestGenerateExtensionTestingSchema(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	s := new(extension.Schema)
	err = s.Decode([]byte(extensionSchema))
	require.NoError(t, err)

	guest, err := generator.GenerateGuestLocal(&generator.Options{
		Extension: s,

		GolangPackageImportPath: "extension",
		GolangPackageName:       "local_inttest_latest_guest",

		RustPackageName:    "local_inttest_latest_guest",
		RustPackageVersion: "v0.1.0",

		TypescriptPackageName:    "local_inttest_latest_guest",
		TypescriptPackageVersion: "v0.1.0",
	})
	require.NoError(t, err)

	golangExtensionDir := wd + "/golang_ext_tests/extension"
	for _, file := range guest.GolangFiles {
		err = os.WriteFile(golangExtensionDir+"/"+file.Name(), file.Data(), 0644)
		require.NoError(t, err)
	}

	rustExtensionDir := wd + "/rust_ext_tests/extension"
	for _, file := range guest.RustFiles {
		err = os.WriteFile(rustExtensionDir+"/"+file.Name(), file.Data(), 0644)
		require.NoError(t, err)
	}

	host, err := generator.GenerateHostLocal(&generator.Options{
		Extension: s,

		GolangPackageImportPath: "extension",
		GolangPackageName:       "local_inttest_latest_host",

		TypescriptPackageName:    "local-example-latest-host",
		TypescriptPackageVersion: "v0.1.0",
	})
	require.NoError(t, err)

	golangExtensionDir = wd + "/golang_ext_tests/host_extension"
	for _, file := range host.GolangFiles {
		if file.Name() != "go.mod" {
			err = os.WriteFile(golangExtensionDir+"/"+file.Name(), file.Data(), 0644)
			require.NoError(t, err)
		}
	}

	typescriptExtensionDir := wd + "/typescript_ext_tests/host_extension"
	for _, file := range host.TypescriptFiles {
		if file.Name() != "go.mod" {
			err = os.WriteFile(typescriptExtensionDir+"/"+file.Name(), file.Data(), 0644)
			require.NoError(t, err)
		}
	}
}
