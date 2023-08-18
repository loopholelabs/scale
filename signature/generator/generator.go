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
	"github.com/loopholelabs/scale/cli/version"
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/signature/generator/golang"
	"golang.org/x/mod/module"
	"golang.org/x/mod/zip"
)

type GuestRegistryPackage struct {
	GolangModule  *bytes.Buffer
	GolangModfile []byte
}

type GuestLocalPackage struct {
	GolangFiles []golang.File
}

type HostRegistryPackage struct {
	GolangModule  *bytes.Buffer
	GolangModfile []byte
}

type HostLocalPackage struct {
	GolangFiles []golang.File
}

type GeneratorOptions struct {
	Signature             *signature.Schema
	GolangImportPath      string
	GolangPackageName     string
	GolangPackageVersion  string
	GolangPolyglotVersion string
}

func GenerateGuestRegistry(options *GeneratorOptions) (*GuestRegistryPackage, error) {
	hash, err := options.Signature.Hash()
	if err != nil {
		return nil, err
	}
	hashString := hex.EncodeToString(hash)

	golangTypes, err := golang.Generate(options.Signature, options.GolangPackageName, version.Version)
	if err != nil {
		return nil, err
	}

	guest, err := golang.GenerateGuest(options.Signature, hashString, options.GolangPackageName, version.Version)
	if err != nil {
		return nil, err
	}

	modfile, err := golang.GenerateModfile(options.GolangImportPath, options.GolangPolyglotVersion)
	if err != nil {
		return nil, err
	}

	files := []zip.File{
		golang.NewFile("types.go", "types.go", golangTypes),
		golang.NewFile("guest.go", "guest.go", guest),
		golang.NewFile("go.mod", "go.mod", modfile),
	}

	buffer := new(bytes.Buffer)
	err = zip.Create(buffer, module.Version{
		Path:    options.GolangImportPath,
		Version: options.GolangPackageVersion,
	}, files)
	if err != nil {
		return nil, err
	}

	return &GuestRegistryPackage{
		GolangModule:  buffer,
		GolangModfile: modfile,
	}, nil
}

func GenerateGuestLocal(options *GeneratorOptions) (*GuestLocalPackage, error) {
	hash, err := options.Signature.Hash()
	if err != nil {
		return nil, err
	}
	hashString := hex.EncodeToString(hash)

	golangTypes, err := golang.Generate(options.Signature, options.GolangPackageName, version.Version)
	if err != nil {
		return nil, err
	}

	guest, err := golang.GenerateGuest(options.Signature, hashString, options.GolangPackageName, version.Version)
	if err != nil {
		return nil, err
	}

	modfile, err := golang.GenerateModfile(options.GolangImportPath, options.GolangPolyglotVersion)
	if err != nil {
		return nil, err
	}

	files := []golang.File{
		golang.NewFile("types.go", "types.go", golangTypes),
		golang.NewFile("guest.go", "guest.go", guest),
		golang.NewFile("go.mod", "go.mod", modfile),
	}

	return &GuestLocalPackage{
		GolangFiles: files,
	}, nil
}

func GenerateHostRegistry(options *GeneratorOptions) (*HostRegistryPackage, error) {
	hash, err := options.Signature.Hash()
	if err != nil {
		return nil, err
	}
	hashString := hex.EncodeToString(hash)

	sig, err := options.Signature.CloneWithDisabledAccessorsValidatorsAndModifiers()
	if err != nil {
		return nil, err
	}

	golangTypes, err := golang.Generate(sig, options.GolangPackageName, version.Version)
	if err != nil {
		return nil, err
	}

	host, err := golang.GenerateHost(sig, hashString, options.GolangPackageName, version.Version)
	if err != nil {
		return nil, err
	}

	modfile, err := golang.GenerateModfile(options.GolangImportPath, options.GolangPolyglotVersion)
	if err != nil {
		return nil, err
	}

	files := []zip.File{
		golang.NewFile("types.go", "types.go", golangTypes),
		golang.NewFile("host.go", "host.go", host),
		golang.NewFile("go.mod", "go.mod", modfile),
	}

	buffer := new(bytes.Buffer)
	err = zip.Create(buffer, module.Version{
		Path:    options.GolangImportPath,
		Version: options.GolangPackageVersion,
	}, files)
	if err != nil {
		return nil, err
	}

	return &HostRegistryPackage{
		GolangModule:  buffer,
		GolangModfile: modfile,
	}, nil
}

func GenerateHostLocal(options *GeneratorOptions) (*HostLocalPackage, error) {
	hash, err := options.Signature.Hash()
	if err != nil {
		return nil, err
	}
	hashString := hex.EncodeToString(hash)

	sig, err := options.Signature.CloneWithDisabledAccessorsValidatorsAndModifiers()
	if err != nil {
		return nil, err
	}

	golangTypes, err := golang.Generate(sig, options.GolangPackageName, version.Version)
	if err != nil {
		return nil, err
	}

	host, err := golang.GenerateHost(sig, hashString, options.GolangPackageName, version.Version)
	if err != nil {
		return nil, err
	}

	modfile, err := golang.GenerateModfile(options.GolangImportPath, options.GolangPolyglotVersion)
	if err != nil {
		return nil, err
	}

	files := []golang.File{
		golang.NewFile("types.go", "types.go", golangTypes),
		golang.NewFile("host.go", "host.go", host),
		golang.NewFile("go.mod", "go.mod", modfile),
	}

	return &HostLocalPackage{
		GolangFiles: files,
	}, nil
}
