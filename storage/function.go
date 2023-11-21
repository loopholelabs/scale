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

	"github.com/loopholelabs/scale/scalefunc"
)

const (
	FunctionDirectory = "functions"
)

var (
	DefaultFunction *FunctionStorage
)

type Function struct {
	Schema       *scalefunc.V1AlphaSchema
	Hash         string
	Organization string
}

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DefaultFunction, err = NewFunction(path.Join(homeDir, DefaultDirectory, FunctionDirectory))
	if err != nil {
		panic(err)
	}
}

// FunctionStorage is used to store and retrieve built Scale Functions
type FunctionStorage struct {
	Directory string
}

func NewFunction(baseDirectory string) (*FunctionStorage, error) {
	err := os.MkdirAll(baseDirectory, 0755)
	if err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}

	return &FunctionStorage{
		Directory: baseDirectory,
	}, nil
}

// Get returns the Scale Function with the given name, tag, and organization.
// The hash parameter is optional and can be used to check for a specific hash.
func (s *FunctionStorage) Get(name string, tag string, org string, hash string) (*Function, error) {
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
		f := s.functionName(name, tag, org, hash)
		p := s.fullPath(f)

		stat, err := os.Stat(p)
		if err != nil {
			return nil, err
		}

		if stat.IsDir() {
			return nil, fmt.Errorf("found function is a directory not a file %s/%s:%s", org, name, tag)
		}

		sf, err := scalefunc.Read(p)
		if err != nil {
			return nil, err
		}
		return &Function{
			Schema:       sf,
			Hash:         hash,
			Organization: org,
		}, nil
	}

	f := s.functionSearch(name, tag, org)
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

	if stat.IsDir() {
		return nil, fmt.Errorf("found function is a directory not a file %s/%s:%s", org, name, tag)
	}

	sf, err := scalefunc.Read(matches[0])
	if err != nil {
		return nil, err
	}

	return &Function{
		Schema:       sf,
		Hash:         getHashFromName(path.Base(matches[0])),
		Organization: getOrgFromName(path.Base(matches[0])),
	}, nil
}

// Put stores the Scale Function with the given name, tag, organization, and hash
func (s *FunctionStorage) Put(name string, tag string, org string, sf *scalefunc.V1AlphaSchema) error {
	f := s.functionName(name, tag, org, hex.EncodeToString(sf.GetHash()))
	p := s.fullPath(f)
	return os.WriteFile(p, sf.Encode(), 0644)
}

// Delete removes the Scale Function with the given name, tag, org, and hash
func (s *FunctionStorage) Delete(name string, tag string, org string, hash string) error {
	return os.Remove(s.fullPath(s.functionName(name, tag, org, hash)))
}

// List returns all the Scale Functions stored in the storage
func (s *FunctionStorage) List() ([]Function, error) {
	entries, err := os.ReadDir(s.Directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage directory %s: %w", s.Directory, err)
	}
	var scaleFuncEntries []Function
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		scaleFunc, err := scalefunc.Read(s.fullPath(entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to decode scale function %s: %w", s.fullPath(entry.Name()), err)
		}
		scaleFuncEntries = append(scaleFuncEntries, Function{
			Schema:       scaleFunc,
			Hash:         getHashFromName(entry.Name()),
			Organization: getOrgFromName(entry.Name()),
		})
	}
	return scaleFuncEntries, nil
}

func (s *FunctionStorage) fullPath(p string) string {
	return path.Join(s.Directory, p)
}

func (s *FunctionStorage) functionName(name string, tag string, org string, hash string) string {
	return fmt.Sprintf("%s_%s_%s_%s_scale", org, name, tag, hash)
}

func (s *FunctionStorage) functionSearch(name string, tag string, org string) string {
	return fmt.Sprintf("%s_%s_%s_*_scale", org, name, tag)
}
