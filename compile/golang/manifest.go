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
	"golang.org/x/mod/modfile"
)

type Manifest struct {
	modfile *modfile.File
}

func ParseManifest(data []byte) (*Manifest, error) {
	file, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, err
	}
	return &Manifest{
		modfile: file,
	}, nil
}

func (m *Manifest) AddReplacement(oldDependency string, oldVersion string, newDependency string, newVersion string) error {
	return m.modfile.AddReplace(oldDependency, oldVersion, newDependency, newVersion)
}

func (m *Manifest) AddRequire(dependency string, version string) error {
	return m.modfile.AddRequire(dependency, version)
}

func (m *Manifest) HasRequire(dependency string, version string, lax bool) bool {
	for _, v := range m.modfile.Require {
		if v.Mod.Path == dependency && (lax || v.Mod.Version == version) {
			return true
		}
	}
	return false
}

func (m *Manifest) HasReplacement(oldDependency string, oldVersion string, newDependency string, newVersion string, lax bool) bool {
	for _, v := range m.modfile.Replace {
		if v.Old.Path == oldDependency && (lax || (v.New.Path == newDependency && v.Old.Version == oldVersion && v.New.Version == newVersion)) {
			return true
		}
	}

	return false
}

func (m *Manifest) GetReplacement(dependency string) (string, string) {
	for _, v := range m.modfile.Replace {
		if v.Old.Path == dependency {
			return v.New.Version, v.New.Path
		}
	}
	return "", ""
}

func (m *Manifest) RemoveRequire(dependency string) error {
	return m.modfile.DropRequire(dependency)
}

func (m *Manifest) RemoveReplacement(dependency string, version string) error {
	return m.modfile.DropReplace(dependency, version)
}

func (m *Manifest) Write() ([]byte, error) {
	return m.modfile.Format()
}
