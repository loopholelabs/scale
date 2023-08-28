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

// Package storage is used to store and retrieve built Scale Functions
package storage

import (
	"encoding/hex"
	"fmt"
	"github.com/loopholelabs/polyglot/version"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/signature/generator"
	"os"
	"path"
	"path/filepath"
)

const (
	SignatureDirectory = "signatures"
)

var (
	DefaultSignature *SignatureStorage
)

type Signature struct {
	Name         string
	Tag          string
	Schema       *signature.Schema
	Hash         string
	Organization string
}

type SignatureStorage struct {
	Directory string
}

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DefaultSignature, err = NewSignature(path.Join(homeDir, DefaultDirectory, SignatureDirectory))
	if err != nil {
		panic(err)
	}
}

func NewSignature(baseDirectory string) (*SignatureStorage, error) {
	dir := path.Join(baseDirectory, SignatureDirectory)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}

	return &SignatureStorage{
		Directory: dir,
	}, nil
}

// Get returns the Scale Signature with the given name, tag, and organization.
// The hash parameter is optional and can be used to check for a specific hash.
func (s *SignatureStorage) Get(name string, tag string, org string, hash string) (*Signature, error) {
	if name == "" || !scalefunc.ValidString(name) {
		return nil, ErrInvalidName
	}

	if tag == "" || !scalefunc.ValidString(tag) {
		return nil, ErrInvalidTag
	}

	if org == "" || !scalefunc.ValidString(org) {
		return nil, ErrInvalidOrganization
	}

	if hash != "" {
		f := s.signatureName(name, tag, org, hash)
		p := s.fullPath(f)

		stat, err := os.Stat(p)
		if err != nil {
			return nil, err
		}

		if !stat.IsDir() {
			return nil, fmt.Errorf("found signature is a file not a directory %s/%s:%s", org, name, tag)
		}

		sig, err := signature.ReadSchema(path.Join(p, "signature"))
		if err != nil {
			return nil, err
		}

		return &Signature{
			Name:         name,
			Tag:          tag,
			Schema:       sig,
			Hash:         hash,
			Organization: org,
		}, nil
	}

	f := s.signatureSearch(name, tag, org)
	p := s.fullPath(f)

	matches, err := filepath.Glob(p)
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, nil
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple matches found for %s/%s:%s", org, name, tag)
	}

	stat, err := os.Stat(matches[0])
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("found signature is a file not a directory %s/%s:%s", org, name, tag)
	}

	sig, err := signature.ReadSchema(path.Join(matches[0], "signature"))
	if err != nil {
		return nil, err
	}

	return &Signature{
		Name:         name,
		Tag:          tag,
		Schema:       sig,
		Hash:         getHashFromName(matches[0]),
		Organization: getOrgFromName(matches[0]),
	}, nil
}

func (s *SignatureStorage) Path(name string, tag string, org string, hash string) (string, error) {
	if name == "" || !scalefunc.ValidString(name) {
		return "", ErrInvalidName
	}

	if tag == "" || !scalefunc.ValidString(tag) {
		return "", ErrInvalidTag
	}

	if org == "" || !scalefunc.ValidString(org) {
		return "", ErrInvalidOrganization
	}

	if hash != "" {
		f := s.signatureName(name, tag, org, hash)
		p := s.fullPath(f)

		stat, err := os.Stat(p)
		if err != nil {
			return "", err
		}

		if !stat.IsDir() {
			return "", fmt.Errorf("found signature is a file not a directory %s/%s:%s", org, name, tag)
		}

		return p, nil
	}

	f := s.signatureSearch(name, tag, org)
	p := s.fullPath(f)

	matches, err := filepath.Glob(p)
	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "", nil
	}

	if len(matches) > 1 {
		return "", fmt.Errorf("multiple matches found for %s/%s:%s", org, name, tag)
	}

	stat, err := os.Stat(matches[0])
	if err != nil {
		return "", err
	}

	if !stat.IsDir() {
		return "", fmt.Errorf("found signature is a file not a directory %s/%s:%s", org, name, tag)
	}

	return matches[0], nil
}

// Put stores the Scale Signature with the given name, tag, organization
func (s *SignatureStorage) Put(name string, tag string, org string, sig *signature.Schema) error {
	hash, err := sig.Hash()
	if err != nil {
		return err
	}

	hashString := hex.EncodeToString(hash)

	f := s.signatureName(name, tag, org, hashString)
	directory := s.fullPath(f)
	err = os.MkdirAll(directory, 0755)
	if err != nil {
		return err
	}

	err = GenerateSignature(sig, name, tag, org, directory)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes the Scale Signature with the given name, tag, org, and hash
func (s *SignatureStorage) Delete(name string, tag string, org string, hash string) error {
	return os.RemoveAll(s.fullPath(s.signatureName(name, tag, org, hash)))
}

// List returns all the Scale Signatures stored in the storage
func (s *SignatureStorage) List() ([]Signature, error) {
	entries, err := os.ReadDir(s.Directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage directory %s: %w", s.Directory, err)
	}
	var scaleSignatureEntries []Signature
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		sig, err := signature.ReadSchema(path.Join(s.fullPath(entry.Name()), "signature"))
		if err != nil {
			return nil, fmt.Errorf("failed to decode scale signature %s: %w", s.fullPath(entry.Name()), err)
		}
		scaleSignatureEntries = append(scaleSignatureEntries, Signature{
			Name:         getNameFromName(entry.Name()),
			Tag:          getTagFromName(entry.Name()),
			Schema:       sig,
			Hash:         getHashFromName(entry.Name()),
			Organization: getOrgFromName(entry.Name()),
		})
	}
	return scaleSignatureEntries, nil
}

func (s *SignatureStorage) fullPath(p string) string {
	return path.Join(s.Directory, p)
}

func (s *SignatureStorage) signatureName(name string, tag string, org string, hash string) string {
	return fmt.Sprintf("%s_%s_%s_%s_signature", org, name, tag, hash)
}

func (s *SignatureStorage) signatureSearch(name string, tag string, org string) string {
	return fmt.Sprintf("%s_%s_%s_*_signature", org, name, tag)
}

// GenerateSignature generates the signature files and writes them to
// the given path.
func GenerateSignature(sig *signature.Schema, name string, tag string, org string, directory string) error {
	encoded, err := sig.Encode()
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(directory, "signature"), encoded, 0644)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(directory, "golang", "guest"), 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(directory, "rust", "guest"), 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(directory, "golang", "host"), 0755)
	if err != nil {
		return err
	}

	guestPackage, err := generator.GenerateGuestLocal(&generator.Options{
		Signature:             sig,
		GolangImportPath:      "signature",
		GolangPackageName:     "signature",
		GolangPackageVersion:  "v0.1.0",
		GolangPolyglotVersion: version.Version(),
		RustPackageName:       fmt.Sprintf("%s_%s_%s_guest", org, name, tag),
		RustPackageVersion:    "0.1.0",
		RustPolyglotVersion:   version.Version(),
	})
	if err != nil {
		return err
	}

	for _, file := range guestPackage.GolangFiles {
		err = os.WriteFile(path.Join(directory, "golang", "guest", file.Path()), file.Data(), 0644)
		if err != nil {
			return err
		}
	}

	for _, file := range guestPackage.RustFiles {
		err = os.WriteFile(path.Join(directory, "rust", "guest", file.Path()), file.Data(), 0644)
		if err != nil {
			return err
		}
	}

	hostPackage, err := generator.GenerateHostLocal(&generator.Options{
		Signature:             sig,
		GolangImportPath:      "signature",
		GolangPackageName:     "signature",
		GolangPackageVersion:  "v0.1.0",
		GolangPolyglotVersion: version.Version(),
	})
	if err != nil {
		return err
	}

	for _, file := range hostPackage.GolangFiles {
		err = os.WriteFile(path.Join(directory, "golang", "host", file.Path()), file.Data(), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
