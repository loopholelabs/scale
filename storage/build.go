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
	"github.com/google/uuid"
	"os"
	"path"
)

const (
	BuildDirectory = "builds"
)

var (
	DefaultBuild *BuildStorage
)

type Build struct {
	Path string
}

type BuildStorage struct {
	Directory string
}

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DefaultBuild, err = NewBuild(path.Join(homeDir, DefaultDirectory, BuildDirectory))
	if err != nil {
		panic(err)
	}
}

func NewBuild(baseDirectory string) (*BuildStorage, error) {
	dir := path.Join(baseDirectory, BuildDirectory)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}

	return &BuildStorage{
		Directory: dir,
	}, nil
}

// Mkdir creates and returns a new build directory
func (s *BuildStorage) Mkdir() (*Build, error) {
	p := path.Join(s.Directory, uuid.New().String())
	return &Build{
		Path: p,
	}, os.MkdirAll(p, 0755)
}

func (s *BuildStorage) Delete(b *Build) error {
	return os.RemoveAll(b.Path)
}
