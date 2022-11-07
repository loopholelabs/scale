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

// Package scalefunc implements the ScaleFunc type, as well as any helper functions
// for interacting with ScaleFunc types
package scalefunc

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/loopholelabs/polyglot-go"
)

var (
	VersionErr  = errors.New("unknown or invalid version")
	LanguageErr = errors.New("unknown or invalid language")
	ChecksumErr = errors.New("error while verifying checksum")
)

// Version is the Version of the ScaleFunc definition
type Version string

const (
	// V1Alpha is the V1 Alpha definition of a ScaleFunc
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

// ScaleFunc is the type used to define the requirements of a
// scale function for a Scale Runtime
type ScaleFunc struct {
	Version  Version  `json:"version" yaml:"version"`
	Name     string   `json:"name" yaml:"name"`
	Language Language `json:"language" yaml:"language"`
	Function []byte   `json:"function" yaml:"function"`
	Size     uint32   `json:"size" yaml:"size"`
	Checksum string   `json:"checksum" yaml:"checksum"`
}

func (s *ScaleFunc) Encode() []byte {
	b := polyglot.GetBuffer()
	defer polyglot.PutBuffer(b)
	e := polyglot.Encoder(b)
	e.String(string(s.Version))
	e.String(s.Name)
	e.String(string(s.Language))
	e.Bytes(s.Function)

	size := uint32(len(b.Bytes()))
	hash := sha256.New()
	hash.Write(b.Bytes())
	checksum := hex.EncodeToString(hash.Sum(nil))

	e.Uint32(size)
	e.String(checksum)

	return b.Bytes()
}

func (s *ScaleFunc) Decode(data []byte) error {
	d := polyglot.GetDecoder(data)
	defer d.Return()

	version, err := d.String()
	if err != nil {
		return err
	}
	s.Version = Version(version)

	invalid := true
	for _, v := range AcceptedVersions {
		if s.Version == v {
			invalid = false
			break
		}
	}
	if invalid {
		return VersionErr
	}

	s.Name, err = d.String()
	if err != nil {
		return err
	}

	language, err := d.String()
	if err != nil {
		return err
	}
	s.Language = Language(language)

	invalid = true
	for _, l := range AcceptedLanguages {
		if l == s.Language {
			invalid = false
			break
		}
	}
	if invalid {
		return LanguageErr
	}

	s.Function, err = d.Bytes(nil)
	if err != nil {
		return err
	}

	s.Size, err = d.Uint32()
	if err != nil {
		return err
	}

	s.Checksum, err = d.String()
	if err != nil {
		return err
	}

	hash := sha256.New()
	hash.Write(data[:s.Size])

	if hex.EncodeToString(hash.Sum(nil)) != s.Checksum {
		return ChecksumErr
	}

	return nil
}
