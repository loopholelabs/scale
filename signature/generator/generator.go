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
	"github.com/loopholelabs/scale/signature/generator/rust"
	"golang.org/x/mod/module"
	"golang.org/x/mod/zip"
)

type GuestRegistryPackage struct {
	GolangModule  *bytes.Buffer
	GolangModfile []byte
}

type GuestLocalPackage struct {
	GolangFiles []File
	RustFiles   []File
}

type HostRegistryPackage struct {
	GolangModule  *bytes.Buffer
	GolangModfile []byte
}

type HostLocalPackage struct {
	GolangFiles []File
}

type Options struct {
	Signature *signature.Schema

	GolangImportPath      string
	GolangPackageName     string
	GolangPackageVersion  string
	GolangScaleVersion    string
	GolangPolyglotVersion string

	RustPackageName     string
	RustPackageVersion  string
	RustScaleVersion    string
	RustPolyglotVersion string
}

func GenerateGuestRegistry(options *Options) (*GuestRegistryPackage, error) {
	hash, err := options.Signature.Hash()
	if err != nil {
		return nil, err
	}
	hashString := hex.EncodeToString(hash)

	golangTypes, err := golang.Generate(options.Signature, options.GolangPackageName, options.GolangScaleVersion)
	if err != nil {
		return nil, err
	}

	guest, err := golang.GenerateGuest(options.Signature, hashString, options.GolangPackageName, options.GolangScaleVersion)
	if err != nil {
		return nil, err
	}

	modfile, err := golang.GenerateModfile(options.GolangImportPath, options.GolangScaleVersion, options.GolangPolyglotVersion)
	if err != nil {
		return nil, err
	}

	files := []zip.File{
		NewFile("types.go", "types.go", golangTypes),
		NewFile("guest.go", "guest.go", guest),
		NewFile("go.mod", "go.mod", modfile),
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

func GenerateGuestLocal(options *Options) (*GuestLocalPackage, error) {
	hash, err := options.Signature.Hash()
	if err != nil {
		return nil, err
	}
	hashString := hex.EncodeToString(hash)

	golangTypes, err := golang.Generate(options.Signature, options.GolangPackageName, options.GolangScaleVersion)
	if err != nil {
		return nil, err
	}

	golangGuest, err := golang.GenerateGuest(options.Signature, hashString, options.GolangPackageName, options.GolangScaleVersion)
	if err != nil {
		return nil, err
	}

	modfile, err := golang.GenerateModfile(options.GolangImportPath, options.GolangScaleVersion, options.GolangPolyglotVersion)
	if err != nil {
		return nil, err
	}

	golangFiles := []File{
		NewFile("types.go", "types.go", golangTypes),
		NewFile("guest.go", "guest.go", golangGuest),
		NewFile("go.mod", "go.mod", modfile),
	}

	rustTypes, err := rust.Generate(options.Signature, options.RustPackageName, options.RustScaleVersion)
	if err != nil {
		return nil, err
	}

	rustGuest, err := rust.GenerateGuest(options.Signature, hashString, options.RustPackageName, options.RustScaleVersion)
	if err != nil {
		return nil, err
	}

	cargofile, err := rust.GenerateCargofile(options.RustPackageName, options.RustPackageVersion, options.RustScaleVersion, options.RustPolyglotVersion)
	if err != nil {
		return nil, err
	}

	rustFiles := []File{
		NewFile("types.rs", "types.rs", rustTypes),
		NewFile("guest.rs", "guest.rs", rustGuest),
		NewFile("Cargo.toml", "Cargo.toml", cargofile),
	}

	return &GuestLocalPackage{
		GolangFiles: golangFiles,
		RustFiles:   rustFiles,
	}, nil
}

func GenerateHostRegistry(options *Options) (*HostRegistryPackage, error) {
	hash, err := options.Signature.Hash()
	if err != nil {
		return nil, err
	}
	hashString := hex.EncodeToString(hash)

	sig, err := options.Signature.CloneWithDisabledAccessorsValidatorsAndModifiers()
	if err != nil {
		return nil, err
	}

	golangTypes, err := golang.Generate(sig, options.GolangPackageName, options.GolangScaleVersion)
	if err != nil {
		return nil, err
	}

	host, err := golang.GenerateHost(sig, hashString, options.GolangPackageName, options.GolangScaleVersion)

	if err != nil {
		return nil, err
	}

	modfile, err := golang.GenerateModfile(options.GolangImportPath, options.GolangScaleVersion, options.GolangPolyglotVersion)
	if err != nil {
		return nil, err
	}

	files := []zip.File{
		NewFile("types.go", "types.go", golangTypes),
		NewFile("host.go", "host.go", host),
		NewFile("go.mod", "go.mod", modfile),
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

func GenerateHostLocal(options *Options) (*HostLocalPackage, error) {
	hash, err := options.Signature.Hash()
	if err != nil {
		return nil, err
	}
	hashString := hex.EncodeToString(hash)

	sig, err := options.Signature.CloneWithDisabledAccessorsValidatorsAndModifiers()
	if err != nil {
		return nil, err
	}

	golangTypes, err := golang.Generate(sig, options.GolangPackageName, options.GolangScaleVersion)
	if err != nil {
		return nil, err
	}

	host, err := golang.GenerateHost(sig, hashString, options.GolangPackageName, options.GolangScaleVersion)
	if err != nil {
		return nil, err
	}

	modfile, err := golang.GenerateModfile(options.GolangImportPath, options.GolangScaleVersion, options.GolangPolyglotVersion)
	if err != nil {
		return nil, err
	}

	files := []File{
		NewFile("types.go", "types.go", golangTypes),
		NewFile("host.go", "host.go", host),
		NewFile("go.mod", "go.mod", modfile),
	}

	return &HostLocalPackage{
		GolangFiles: files,
	}, nil
}
