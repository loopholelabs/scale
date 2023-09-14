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
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
)

const (
	V1AlphaVersion = "v1alpha"
)

var (
	ErrInvalidName = errors.New("invalid name")
	ErrInvalidTag  = errors.New("invalid tag")
)

type SignatureSchema struct {
	Organization string `hcl:"organization,attr"`
	Name         string `hcl:"name,attr"`
	Tag          string `hcl:"tag,attr"`
}

type ExtensionSchema struct {
	Organization string `hcl:"organization,attr"`
	Name         string `hcl:"name,attr"`
	Tag          string `hcl:"tag,attr"`
}

type ExtensionSchema struct {
	Organization string `hcl:"organization,optional"`
	Name         string `hcl:"name,attr"`
	Tag          string `hcl:"tag,attr"`
}

type Schema struct {
	Version     string            `hcl:"version,attr"`
	Name        string            `hcl:"name,attr"`
	Tag         string            `hcl:"tag,attr"`
	Language    string            `hcl:"language,attr"`
	Signature   SignatureSchema   `hcl:"signature,block"`
	Stateless   bool              `hcl:"stateless,optional"`
	Function    string            `hcl:"function,attr"`
	Initialize  string            `hcl:"initialize,attr"`
	Description string            `hcl:"description,optional"`
	Extensions  []ExtensionSchema `hcl:"extension,optional"`
}

func ReadSchema(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	s := new(Schema)
	return s, s.Decode(data)
}

func (s *Schema) Decode(data []byte) error {
	file, diag := hclsyntax.ParseConfig(data, "", hcl.Pos{Line: 1, Column: 1})
	if diag.HasErrors() {
		return diag.Errs()[0]
	}

	diag = gohcl.DecodeBody(file.Body, nil, s)
	if diag.HasErrors() {
		return diag.Errs()[0]
	}

	return nil
}

func (s *Schema) Encode() ([]byte, error) {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(s, f.Body())
	return f.Bytes(), nil
}

func (s *Schema) Validate() error {
	switch s.Version {
	case V1AlphaVersion:
		if !scalefunc.ValidString(s.Name) {
			return ErrInvalidName
		}

		if scalefunc.ValidString(s.Tag) {
			return ErrInvalidTag
		}

		switch scalefunc.Language(s.Language) {
		case scalefunc.Go, scalefunc.Rust, scalefunc.TypeScript:
		default:
			return fmt.Errorf("unknown or invalid language: %s", s.Language)
		}

		if !scalefunc.ValidString(s.Signature.Organization) {
			return fmt.Errorf("invalid organization: %s", s.Signature.Organization)
		}

		if !signature.ValidLabel.MatchString(s.Signature.Name) {
			return fmt.Errorf("invalid name: %s", s.Signature.Name)
		}

		if signature.InvalidString.MatchString(s.Signature.Tag) {
			return fmt.Errorf("invalid tag: %s", s.Signature.Tag)
		}

		if len(s.Function) == 0 {
			return fmt.Errorf("function must be defined")
		}

		return nil
	default:
		return fmt.Errorf("unknown schema version: %s", s.Version)
	}
}
