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

package signature

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

var (
	VersionErr = errors.New("unknown or invalid version")
)

// Version is the Version of the ScaleFile definition
type Version string

const (
	// V1Alpha is the V1 Alpha definition of a Scale Signature
	V1Alpha Version = "v1alpha"
)

var (
	// AcceptedVersions is an array of acceptable Versions
	AcceptedVersions = []Version{V1Alpha}
)

// PublishedVersion is the published version of the signature specific to a language
type PublishedVersion struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
}

// PublishedVersions is a list of published versions of the signature specific to a language
type PublishedVersions struct {
	Go PublishedVersion `json:"go" yaml:"go"`
}

// Definition is the definition of a signature as well as the published versions
type Definition struct {
	DefinitionVersion Version           `json:"definition" yaml:"definition"`
	Name              string            `json:"name" yaml:"name"`
	Version           string            `json:"version" yaml:"version"`
	PublishedVersions PublishedVersions `json:"published" yaml:"published"`
}

// Read opens a file at the given path and returns a *Definition
func Read(path string) (*Definition, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	return Decode(file)
}

// Write opens a file at the given path and writes the given definition to it
func Write(path string, definition *Definition) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	return Encode(file, definition)
}

// Decode reads the data stored in the given io.Reader and returns a *Definition
func Decode(reader io.Reader) (*Definition, error) {
	decoder := yaml.NewDecoder(reader)
	definition := new(Definition)
	err := decoder.Decode(definition)
	if err != nil {
		return nil, err
	}

	invalid := true
	for _, v := range AcceptedVersions {
		if definition.DefinitionVersion == v {
			invalid = false
			break
		}
	}
	if invalid {
		return nil, VersionErr
	}

	return definition, nil
}

// Encode writes the given definition to the given io.Writer
func Encode(writer io.Writer, definition *Definition) error {
	invalid := true
	for _, v := range AcceptedVersions {
		if definition.DefinitionVersion == v {
			invalid = false
			break
		}
	}
	if invalid {
		return VersionErr
	}

	encoder := yaml.NewEncoder(writer)
	encoder.SetIndent(2)
	return encoder.Encode(definition)
}
