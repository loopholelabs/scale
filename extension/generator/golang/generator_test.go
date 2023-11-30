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

package golang

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/loopholelabs/scale/extension"
)

func TestGenerator(t *testing.T) {
	s := new(extension.Schema)
	err := s.Decode([]byte(extension.MasterTestingSchema))
	require.NoError(t, err)

	packageName := "extfetch"

	interfaces, err := GenerateInterfaces(s, packageName)
	require.NoError(t, err)
	// os.WriteFile("./interfaces.txt", interfaces, 0644)

	expInterfaces, err := os.ReadFile("./interfaces.txt")
	require.NoError(t, err)
	require.Equal(t, string(expInterfaces), string(interfaces))

	formatted, err := GenerateTypes(s, "types")
	require.NoError(t, err)
	// os.WriteFile("./generated.txt", formatted, 0644)

	expTypes, err := os.ReadFile("./generated.txt")
	require.NoError(t, err)
	require.Equal(t, string(expTypes), string(formatted))

	sHash, err := s.Hash()
	require.NoError(t, err)

	hash := hex.EncodeToString(sHash)

	host, err := GenerateHost(s, hash, packageName)
	require.NoError(t, err)
	// os.WriteFile("./host.txt", host, 0644)
	expHost, err := os.ReadFile("./host.txt")
	require.NoError(t, err)
	require.Equal(t, string(expHost), string(host))

	guest, err := GenerateGuest(s, hash, packageName)
	require.NoError(t, err)
	// os.WriteFile("./guest.txt", guest, 0644)
	expGuest, err := os.ReadFile("./guest.txt")
	require.NoError(t, err)
	require.Equal(t, string(expGuest), string(guest))

	mod, err := GenerateModfile(packageName)
	require.NoError(t, err)
	// os.WriteFile("./modfile.txt", mod, 0644)
	expMod, err := os.ReadFile("./modfile.txt")
	require.NoError(t, err)
	require.Equal(t, string(expMod), string(mod))

}
