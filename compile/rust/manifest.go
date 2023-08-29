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

package rust

import (
	"bytes"
	"github.com/BurntSushi/toml"
)

type Cargo struct {
	CargoFeatures     interface{}            `toml:"cargo-features,omitempty"`
	Package           interface{}            `toml:"package,omitempty"`
	Lib               interface{}            `toml:"lib,omitempty"`
	Bin               interface{}            `toml:"bin,omitempty"`
	Example           interface{}            `toml:"example,omitempty"`
	Test              interface{}            `toml:"test,omitempty"`
	Bench             interface{}            `toml:"bench,omitempty"`
	Dependencies      map[string]interface{} `toml:"dependencies,omitempty"`
	DevDependencies   interface{}            `toml:"dev-dependencies,omitempty"`
	BuildDependencies interface{}            `toml:"build-dependencies,omitempty"`
	Target            interface{}            `toml:"target,omitempty"`
	Badges            interface{}            `toml:"badges,omitempty"`
	Features          interface{}            `toml:"features,omitempty"`
	Patch             interface{}            `toml:"patch,omitempty"`
	Replace           interface{}            `toml:"replace,omitempty"`
	Profile           interface{}            `toml:"profile,omitempty"`
	Workspace         interface{}            `toml:"workspace,omitempty"`
}

type DependencyPath struct {
	Path    string `toml:"path"`
	Package string `toml:"package"`
}

type DependencyVersion struct {
	Version  string `toml:"version"`
	Package  string `toml:"package"`
	Registry string `toml:"registry"`
}

type Manifest struct {
	cargo *Cargo
}

func ParseManifest(data []byte) (*Manifest, error) {
	cargo := new(Cargo)
	_, err := toml.Decode(string(data), cargo)
	if err != nil {
		return nil, err
	}
	return &Manifest{
		cargo: cargo,
	}, nil
}

func (m *Manifest) AddDependencyWithVersion(dependency string, version DependencyVersion) error {
	m.cargo.Dependencies[dependency] = version
	return nil
}

func (m *Manifest) AddDependencyWithPath(dependency string, path DependencyPath) error {
	m.cargo.Dependencies[dependency] = path
	return nil
}

func (m *Manifest) HasDependency(dependency string) bool {
	_, ok := m.cargo.Dependencies[dependency]
	return ok
}

func (m *Manifest) RemoveDependency(dependency string) error {
	delete(m.cargo.Dependencies, dependency)
	return nil
}

func (m *Manifest) Write() ([]byte, error) {
	b := bytes.NewBuffer(nil)
	err := toml.NewEncoder(b).Encode(m.cargo)
	return b.Bytes(), err
}
