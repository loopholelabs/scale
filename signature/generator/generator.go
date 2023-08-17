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
	"encoding/hex"
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/signature/generator/golang"
	"golang.org/x/mod/module"
	"golang.org/x/mod/zip"
)

type GuestPackage struct {
	GolangModule  *bytes.Buffer
	GolangModfile []byte
}

type HostPackage struct {
	GolangModule  *bytes.Buffer
	GolangModfile []byte
}

type Options struct {
	Signature *signature.Schema

	GolangImportPath      string
	GolangPackageName     string
	GolangPackageVersion  string
	GolangPolyglotVersion string
}

func GenerateGuest(options *Options) (*GuestPackage, error) {
	hash, err := options.Signature.Hash()
	if err != nil {
		return nil, err
	}
	hashString := hex.EncodeToString(hash)

	golangTypes, err := golang.Generate(options.Signature, options.GolangPackageName, options.GolangPackageVersion)
	if err != nil {
		return nil, err
	}

	guest, err := golang.GenerateGuest(options.Signature, hashString, options.GolangPackageName, options.GolangPackageVersion)
	if err != nil {
		return nil, err
	}

	modfile, err := golang.GenerateModfile(options.GolangImportPath, options.GolangPolyglotVersion)
	if err != nil {
		return nil, err
	}

	files := []zip.File{
		NewGolangFile("types.go", "types.go", golangTypes),
		NewGolangFile("guest.go", "guest.go", guest),
		NewGolangFile("go.mod", "go.mod", modfile),
	}

	buffer := new(bytes.Buffer)
	err = zip.Create(buffer, module.Version{
		Path:    options.GolangImportPath,
		Version: options.GolangPackageVersion,
	}, files)
	if err != nil {
		return nil, err
	}

	return &GuestPackage{
		GolangModule:  buffer,
		GolangModfile: modfile,
	}, nil
}

func GenerateHost(options *Options) (*HostPackage, error) {
	hash, err := options.Signature.Hash()
	if err != nil {
		return nil, err
	}
	hashString := hex.EncodeToString(hash)

	sig, err := options.Signature.CloneWithDisabledAccessorsValidatorsAndModifiers()
	if err != nil {
		return nil, err
	}

	golangTypes, err := golang.Generate(sig, options.GolangPackageName, options.GolangPackageVersion)
	if err != nil {
		return nil, err
	}

	host, err := golang.GenerateHost(sig, hashString, options.GolangPackageName, options.GolangPackageVersion)
	if err != nil {
		return nil, err
	}

	modfile, err := golang.GenerateModfile(options.GolangImportPath, options.GolangPolyglotVersion)
	if err != nil {
		return nil, err
	}

	files := []zip.File{
		NewGolangFile("types.go", "types.go", golangTypes),
		NewGolangFile("host.go", "host.go", host),
		NewGolangFile("go.mod", "go.mod", modfile),
	}

	buffer := new(bytes.Buffer)
	err = zip.Create(buffer, module.Version{
		Path:    options.GolangImportPath,
		Version: options.GolangPackageVersion,
	}, files)
	if err != nil {
		return nil, err
	}

	return &HostPackage{
		GolangModule:  buffer,
		GolangModfile: modfile,
	}, nil
}
