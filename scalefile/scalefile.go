/*
	Copyright 2022 Loophole Labs

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

package scalefile

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Build struct {
	Language     string       `json:"language"`
	Dependencies []Dependency `json:"dependencies,omitempty"`
}

type ScaleFile struct {
	Version    string `json:"version"`
	Name       string `json:"name"`
	Build      Build  `json:"build"`
	Source     string `json:"source"`
	Middleware bool   `json:"middleware"`
}

func Read(path string) (ScaleFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return ScaleFile{}, err
	}
	defer func() {
		_ = file.Close()
	}()

	return decode(file)
}

func decode(data io.Reader) (ScaleFile, error) {
	decoder := yaml.NewDecoder(data)
	manifest := ScaleFile{}
	err := decoder.Decode(&manifest)
	if err != nil {
		return ScaleFile{}, err
	}

	return manifest, nil
}

func Write(path string, scalefile ScaleFile) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	return encode(file, scalefile)
}

func encode(data io.Writer, scalefile ScaleFile) error {
	encoder := yaml.NewEncoder(data)
	encoder.SetIndent(2)
	return encoder.Encode(scalefile)
}
