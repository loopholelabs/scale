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
	"os"
	"path"
	"path/filepath"

	"github.com/loopholelabs/scale/extension"
	"github.com/loopholelabs/scale/extension/generator"
	"github.com/loopholelabs/scale/scalefunc"
)

const (
	ExtensionDirectory = "extensions"
)

var (
	DefaultExtension *ExtensionStorage
)

type Extension struct {
	Name         string
	Tag          string
	Schema       *extension.Schema
	Hash         string
	Organization string
}

type ExtensionStorage struct {
	Directory string
}

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DefaultExtension, err = NewExtension(path.Join(homeDir, DefaultDirectory, ExtensionDirectory))
	if err != nil {
		panic(err)
	}
}

func NewExtension(baseDirectory string) (*ExtensionStorage, error) {
	err := os.MkdirAll(baseDirectory, 0755)
	if err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}

	return &ExtensionStorage{
		Directory: baseDirectory,
	}, nil
}

// Get returns the Scale Extension with the given name, tag, and organization.
// The hash parameter is optional and can be used to check for a specific hash.
func (s *ExtensionStorage) Get(name string, tag string, org string, hash string) (*Extension, error) {
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
		f := s.extensionName(name, tag, org, hash)
		p := s.fullPath(f)

		stat, err := os.Stat(p)
		if err != nil {
			return nil, err
		}

		if !stat.IsDir() {
			return nil, fmt.Errorf("found extension is a file not a directory %s/%s:%s", org, name, tag)
		}

		sig, err := extension.ReadSchema(path.Join(p, "extension"))
		if err != nil {
			return nil, err
		}

		return &Extension{
			Name:         name,
			Tag:          tag,
			Schema:       sig,
			Hash:         hash,
			Organization: org,
		}, nil
	}

	f := s.extensionSearch(name, tag, org)
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
		return nil, fmt.Errorf("found extension is a file not a directory %s/%s:%s", org, name, tag)
	}

	sig, err := extension.ReadSchema(path.Join(matches[0], "extension"))
	if err != nil {
		return nil, err
	}

	return &Extension{
		Name:         name,
		Tag:          tag,
		Schema:       sig,
		Hash:         getHashFromName(filepath.Base(matches[0])),
		Organization: getOrgFromName(filepath.Base(matches[0])),
	}, nil
}

func (s *ExtensionStorage) Path(name string, tag string, org string, hash string) (string, error) {
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
		f := s.extensionName(name, tag, org, hash)
		p := s.fullPath(f)

		stat, err := os.Stat(p)
		if err != nil {
			return "", err
		}

		if !stat.IsDir() {
			return "", fmt.Errorf("found extension is a file not a directory %s/%s:%s", org, name, tag)
		}

		return p, nil
	}

	f := s.extensionSearch(name, tag, org)
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
		return "", fmt.Errorf("found extension is a file not a directory %s/%s:%s", org, name, tag)
	}

	return matches[0], nil
}

// Put stores the Scale Extension with the given name, tag, organization
func (s *ExtensionStorage) Put(name string, tag string, org string, sig *extension.Schema) error {
	hash, err := sig.Hash()
	if err != nil {
		return err
	}

	hashString := hex.EncodeToString(hash)

	f := s.extensionName(name, tag, org, hashString)
	directory := s.fullPath(f)
	err = os.MkdirAll(directory, 0755)
	if err != nil {
		return err
	}

	err = GenerateExtension(sig, name, tag, org, directory)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes the Scale Extension with the given name, tag, org, and hash
func (s *ExtensionStorage) Delete(name string, tag string, org string, hash string) error {
	return os.RemoveAll(s.fullPath(s.extensionName(name, tag, org, hash)))
}

// List returns all the Scale Extensions stored in the storage
func (s *ExtensionStorage) List() ([]Extension, error) {
	entries, err := os.ReadDir(s.Directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage directory %s: %w", s.Directory, err)
	}
	var scaleExtensionEntries []Extension
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		sig, err := extension.ReadSchema(path.Join(s.fullPath(entry.Name()), "extension"))
		if err != nil {
			return nil, fmt.Errorf("failed to decode scale extension %s: %w", s.fullPath(entry.Name()), err)
		}
		scaleExtensionEntries = append(scaleExtensionEntries, Extension{
			Name:         getNameFromName(entry.Name()),
			Tag:          getTagFromName(entry.Name()),
			Schema:       sig,
			Hash:         getHashFromName(entry.Name()),
			Organization: getOrgFromName(entry.Name()),
		})
	}
	return scaleExtensionEntries, nil
}

func (s *ExtensionStorage) fullPath(p string) string {
	return path.Join(s.Directory, p)
}

func (s *ExtensionStorage) extensionName(name string, tag string, org string, hash string) string {
	return fmt.Sprintf("%s_%s_%s_%s_extension", org, name, tag, hash)
}

func (s *ExtensionStorage) extensionSearch(name string, tag string, org string) string {
	return fmt.Sprintf("%s_%s_%s_*_extension", org, name, tag)
}

// GenerateExtension generates the extension files and writes them to
// the given path.
func GenerateExtension(ext *extension.Schema, name string, tag string, org string, directory string) error {
	encoded, err := ext.Encode()
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(directory, "extension"), encoded, 0644)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(directory, "golang", "guest"), 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(directory, "golang", "host"), 0755)
	if err != nil {
		return err
	}

	guestPackage, err := generator.GenerateGuestLocal(&generator.Options{
		Extension:               ext,
		GolangPackageImportPath: "extension",
		GolangPackageName:       ext.Name,
		GolangPackageVersion:    "v0.1.0",
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

	hostPackage, err := generator.GenerateHostLocal(&generator.Options{
		Extension:               ext,
		GolangPackageImportPath: "extension",
		GolangPackageName:       ext.Name,
		GolangPackageVersion:    "v0.1.0",
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
