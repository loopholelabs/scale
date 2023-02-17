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
	"errors"
	"fmt"
	"github.com/loopholelabs/scalefile/scalefunc"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	ErrInvalidName         = errors.New("invalid name")
	ErrInvalidTag          = errors.New("invalid tag")
	ErrInvalidOrganization = errors.New("invalid organization")
)

const (
	DefaultCacheDirectory = ".config/scale/functions"
)

var (
	Default *Storage
)

type Entry struct {
	ScaleFunc    *scalefunc.ScaleFunc
	Hash         string
	Organization string
}

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	Default, err = New(path.Join(homeDir, DefaultCacheDirectory))
	if err != nil {
		panic(err)
	}
}

// Storage is used to store and retrieve built Scale Functions
type Storage struct {
	BaseDirectory string
}

func New(baseDirectory string) (*Storage, error) {
	err := os.MkdirAll(baseDirectory, 0755)
	if err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}

	return &Storage{
		BaseDirectory: baseDirectory,
	}, nil
}

// Get returns the Scale Function with the given name, tag, and organization.
// The hash parameter is optional and can be used to check for a specific hash.
func (s *Storage) Get(name string, tag string, org string, hash string) (*Entry, error) {
	if name == "" || !scalefunc.ValidName(name) {
		return nil, ErrInvalidName
	}

	if tag == "" || !scalefunc.ValidName(tag) {
		return nil, ErrInvalidTag
	}

	if org == "" || !scalefunc.ValidName(org) {
		return nil, ErrInvalidOrganization
	}

	if hash != "" {
		f := s.functionName(name, tag, org, hash)
		p := s.fullPath(f)
		sf, err := scalefunc.Read(p)
		if err != nil {
			return nil, err
		}
		return &Entry{
			ScaleFunc:    sf,
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

	sf, err := scalefunc.Read(matches[0])
	if err != nil {
		return nil, err
	}

	return &Entry{
		ScaleFunc:    sf,
		Hash:         s.getHashFromFileName(matches[0]),
		Organization: s.getOrgFromFileName(matches[0]),
	}, nil
}

// Put stores the Scale Function with the given name, tag, organization, and hash
func (s *Storage) Put(name string, tag string, org string, hash string, sf *scalefunc.ScaleFunc) error {
	f := s.functionName(name, tag, org, hash)
	p := s.fullPath(f)
	return os.WriteFile(p, sf.Encode(), 0644)
}

// Delete removes the Scale Function with the given name, tag, org, and hash
func (s *Storage) Delete(name string, tag string, org string, hash string) error {
	return os.Remove(s.fullPath(s.functionName(name, tag, org, hash)))
}

// List returns all the Scale Functions stored in the storage
func (s *Storage) List() ([]Entry, error) {
	entries, err := os.ReadDir(s.BaseDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage directory %s: %w", s.BaseDirectory, err)
	}
	var scaleFuncEntries []Entry
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		scaleFunc, err := scalefunc.Read(s.fullPath(entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to decode scale function %s: %w", s.fullPath(entry.Name()), err)
		}
		scaleFuncEntries = append(scaleFuncEntries, Entry{
			ScaleFunc:    scaleFunc,
			Hash:         s.getHashFromFileName(entry.Name()),
			Organization: s.getOrgFromFileName(entry.Name()),
		})
	}
	return scaleFuncEntries, nil
}

func (s *Storage) fullPath(p string) string {
	return path.Join(s.BaseDirectory, p)
}

func (s *Storage) functionName(name string, tag string, org string, hash string) string {
	return fmt.Sprintf("%s.%s.%s.%s.scale", org, name, tag, hash)
}

func (s *Storage) functionSearch(name string, tag string, org string) string {
	return fmt.Sprintf("%s.%s.%s.*.scale", org, name, tag)
}

func (s *Storage) getHashFromFileName(fileName string) string {
	split := strings.Split(fileName, ".")
	if len(split) != 5 {
		return ""
	}

	return split[3]
}

func (s *Storage) getOrgFromFileName(fileName string) string {
	split := strings.Split(fileName, ".")
	if len(split) != 5 {
		return ""
	}

	return split[0]
}
