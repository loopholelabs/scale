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

// Package scalefile implements the ScaleFile type, as well as any helper functions
// for interacting with ScaleFile types
package scalefile

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

var (
	VersionErr  = errors.New("unknown or invalid version")
	LanguageErr = errors.New("unknown or invalid language")
)

// Version is the Version of the ScaleFile definition
type Version string

const (
	// V1Alpha is the V1 Alpha definition of a ScaleFile
	V1Alpha Version = "v1alpha"
)

// Language is the Language the Scale Function's Source Language
type Language string

const (
	// Go is the Golang Source Language for Scale Functions
	Go Language = "go"
)

var (
	// AcceptedVersions is an array of acceptable Versions
	AcceptedVersions = []Version{V1Alpha}

	// AcceptedLanguages is an array of acceptable Languages
	AcceptedLanguages = []Language{Go}
)

// Dependency outlines the Dependency of a Scale Function
type Dependency struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
}

// ScaleFile describes the Scale Function and its dependencies
type ScaleFile struct {
	Version      Version      `json:"version" yaml:"version"`
	Name         string       `json:"name" yaml:"name"`
	Language     Language     `json:"language" yaml:"language"`
	Dependencies []Dependency `json:"dependencies" yaml:"dependencies"`
	Source       string       `json:"source" yaml:"source"`
}

// Read opens a file at the given path and returns a *ScaleFile
func Read(path string) (*ScaleFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	return Decode(file)
}

// Write opens a file at the given path and writes the given scalefile to it
func Write(path string, scalefile *ScaleFile) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	return Encode(file, scalefile)
}

// Decode reads the data stored in the given io.Reader and returns a *ScaleFile
func Decode(reader io.Reader) (*ScaleFile, error) {
	decoder := yaml.NewDecoder(reader)
	scalefile := new(ScaleFile)
	err := decoder.Decode(scalefile)
	if err != nil {
		return nil, err
	}

	invalid := true
	for _, v := range AcceptedVersions {
		if scalefile.Version == v {
			invalid = false
			break
		}
	}
	if invalid {
		return nil, VersionErr
	}

	invalid = true
	for _, l := range AcceptedLanguages {
		if scalefile.Language == l {
			invalid = false
			break
		}
	}
	if invalid {
		return nil, LanguageErr
	}

	return scalefile, nil
}

// Encode writes the given scalefile to the given io.Writer
func Encode(writer io.Writer, scalefile *ScaleFile) error {
	invalid := true
	for _, v := range AcceptedVersions {
		if scalefile.Version == v {
			invalid = false
			break
		}
	}
	if invalid {
		return VersionErr
	}

	invalid = true
	for _, l := range AcceptedLanguages {
		if scalefile.Language == l {
			invalid = false
			break
		}
	}
	if invalid {
		return LanguageErr
	}

	encoder := yaml.NewEncoder(writer)
	encoder.SetIndent(2)
	return encoder.Encode(scalefile)
}
