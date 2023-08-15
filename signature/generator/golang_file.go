/*
	Copyright 2023 Loophole Labs

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

package generator

import (
	"bytes"
	"golang.org/x/mod/zip"
	"io"
	"io/fs"
	"os"
	"time"
)

var _ zip.File = (*GolangFile)(nil)
var _ os.FileInfo = (*GolangFile)(nil)

type GolangFile struct {
	name    string
	path    string
	content []byte
	reader  *bytes.Reader
	size    int64
}

func NewGolangFile(name string, path string, content []byte) GolangFile {
	return GolangFile{
		name:    name,
		path:    path,
		content: content,
		reader:  bytes.NewReader(content),
		size:    int64(len(content)),
	}
}

func (g GolangFile) Name() string {
	return g.name
}

func (g GolangFile) Size() int64 {
	return g.size
}

func (g GolangFile) Mode() fs.FileMode {
	return 0700
}

func (g GolangFile) ModTime() time.Time {
	return time.Now()
}

func (g GolangFile) IsDir() bool {
	return false
}

func (g GolangFile) Sys() any {
	return g.content
}

func (g GolangFile) Path() string {
	return g.path
}

func (g GolangFile) Lstat() (os.FileInfo, error) {
	return g, nil
}

func (g GolangFile) Open() (io.ReadCloser, error) {
	return io.NopCloser(g.reader), nil
}
