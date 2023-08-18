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

package golang

import (
	"bytes"
	"golang.org/x/mod/zip"
	"io"
	"io/fs"
	"os"
	"time"
)

var _ zip.File = (*File)(nil)
var _ os.FileInfo = (*File)(nil)

type File struct {
	name    string
	path    string
	content []byte
	reader  *bytes.Reader
	size    int64
}

func NewFile(name string, path string, content []byte) File {
	return File{
		name:    name,
		path:    path,
		content: content,
		reader:  bytes.NewReader(content),
		size:    int64(len(content)),
	}
}

func (g File) Name() string {
	return g.name
}

func (g File) Size() int64 {
	return g.size
}

func (g File) Mode() fs.FileMode {
	return 0700
}

func (g File) ModTime() time.Time {
	return time.Now()
}

func (g File) IsDir() bool {
	return false
}

func (g File) Sys() any {
	return g.content
}

func (g File) Path() string {
	return g.path
}

func (g File) Lstat() (os.FileInfo, error) {
	return g, nil
}

func (g File) Open() (io.ReadCloser, error) {
	return io.NopCloser(g.reader), nil
}

func (g File) Data() []byte {
	return g.content
}
