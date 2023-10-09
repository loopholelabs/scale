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

package typescript

import (
	"encoding/json"
	"errors"
)

var (
	ErrInvalidManifest = errors.New("invalid manifest")
)

type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description,omitempty"`
	License         string            `json:"license,omitempty"`
	Files           []string          `json:"files,omitempty"`
	Main            string            `json:"main,omitempty"`
	Browser         string            `json:"browser,omitempty"`
	Dependencies    map[string]string `json:"dependencies,omitempty"`
	DevDependencies map[string]string `json:"devDependencies,omitempty"`

	internal map[string]interface{}
}

type Manifest struct {
	packageJSON *PackageJSON
}

func ParseManifest(data []byte) (*Manifest, error) {
	pkgJSON := new(PackageJSON)
	pkgJSON.internal = make(map[string]interface{})

	err := json.Unmarshal(data, &pkgJSON.internal)
	if err != nil {
		return nil, err
	}

	var ok bool
	pkgJSON.Name, ok = pkgJSON.internal["name"].(string)
	if !ok {
		return nil, ErrInvalidManifest
	}
	pkgJSON.Version, ok = pkgJSON.internal["version"].(string)
	if !ok {
		return nil, ErrInvalidManifest
	}
	pkgJSON.Description, _ = pkgJSON.internal["description"].(string)
	pkgJSON.License, _ = pkgJSON.internal["license"].(string)

	files := pkgJSON.internal["files"]
	if files != nil {
		pkgJSON.Files = make([]string, len(files.([]interface{})))
		for i, v := range files.([]interface{}) {
			pkgJSON.Files[i] = v.(string)
		}
	}

	pkgJSON.Main, _ = pkgJSON.internal["main"].(string)
	pkgJSON.Browser, _ = pkgJSON.internal["browser"].(string)

	dependencies := pkgJSON.internal["dependencies"]
	if dependencies != nil {
		pkgJSON.Dependencies = make(map[string]string, len(dependencies.(map[string]interface{})))
		for k, v := range dependencies.(map[string]interface{}) {
			pkgJSON.Dependencies[k] = v.(string)
		}
	}

	devDependencies := pkgJSON.internal["devDependencies"]
	if devDependencies != nil {
		pkgJSON.DevDependencies = make(map[string]string, len(devDependencies.(map[string]interface{})))
		for k, v := range devDependencies.(map[string]interface{}) {
			pkgJSON.DevDependencies[k] = v.(string)
		}
	}

	return &Manifest{
		packageJSON: pkgJSON,
	}, nil
}

func (m *Manifest) AddDependency(dependency string, version string) error {
	m.packageJSON.Dependencies[dependency] = version
	return nil
}

func (m *Manifest) HasDependency(dependency string) bool {
	_, ok := m.packageJSON.Dependencies[dependency]
	return ok
}

func (m *Manifest) GetDependency(dependency string) string {
	dep, ok := m.packageJSON.Dependencies[dependency]
	if !ok {
		return ""
	}
	return dep
}

func (m *Manifest) RemoveDependency(dependency string) error {
	delete(m.packageJSON.Dependencies, dependency)
	return nil
}

func (m *Manifest) Write() ([]byte, error) {
	m.packageJSON.internal["name"] = m.packageJSON.Name
	m.packageJSON.internal["version"] = m.packageJSON.Version

	if m.packageJSON.Description == "" {
		delete(m.packageJSON.internal, "description")
	} else {
		m.packageJSON.internal["description"] = m.packageJSON.Description
	}

	if m.packageJSON.License == "" {
		delete(m.packageJSON.internal, "license")
	} else {
		m.packageJSON.internal["license"] = m.packageJSON.License
	}

	if len(m.packageJSON.Files) == 0 {
		delete(m.packageJSON.internal, "files")
	} else {
		m.packageJSON.internal["files"] = m.packageJSON.Files
	}

	if m.packageJSON.Main == "" {
		delete(m.packageJSON.internal, "main")
	} else {
		m.packageJSON.internal["main"] = m.packageJSON.Main
	}

	if m.packageJSON.Browser == "" {
		delete(m.packageJSON.internal, "browser")
	} else {
		m.packageJSON.internal["browser"] = m.packageJSON.Browser
	}

	if len(m.packageJSON.Dependencies) == 0 {
		delete(m.packageJSON.internal, "dependencies")
	} else {
		m.packageJSON.internal["dependencies"] = m.packageJSON.Dependencies
	}

	if len(m.packageJSON.DevDependencies) == 0 {
		delete(m.packageJSON.internal, "devDependencies")
	} else {
		m.packageJSON.internal["devDependencies"] = m.packageJSON.DevDependencies
	}

	return json.MarshalIndent(m.packageJSON.internal, "", "\t")
}
