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

package typescript

import (
	"path/filepath"
)

// SignatureInfo is the import information for a Signature
type SignatureInfo struct {
	// Local specifies whether the signature is a locally importable signature
	Local bool

	// ImportPath is the import path for the signature, if the signature is local this will be a path on disk
	ImportPath string
}

func (s *SignatureInfo) normalize() {
	if s.Local {
		if !filepath.IsAbs(s.ImportPath) && !filepath.IsLocal(s.ImportPath) {
			s.ImportPath = "./" + s.ImportPath
		}
	}
}

// FunctionInfo is the import information for a Function
type FunctionInfo struct {
	// PackageName is the name of the function package
	PackageName string

	// ImportPath is the import path for the function, always a path on disk
	ImportPath string
}

func (f *FunctionInfo) normalize() {
	if !filepath.IsAbs(f.ImportPath) && !filepath.IsLocal(f.ImportPath) {
		f.ImportPath = "./" + f.ImportPath
	}
}
