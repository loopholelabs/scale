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
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"fmt"
	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/signature/generator/golang"
	"github.com/loopholelabs/scale/signature/generator/rust"
	"golang.org/x/mod/module"
	"golang.org/x/mod/zip"
	"path"
)

type GuestRegistryPackage struct {
	GolangModule  *bytes.Buffer
	GolangModfile []byte
	RustCrate     *bytes.Buffer
	RustCargofile []byte
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

	golangBuffer := new(bytes.Buffer)
	err = zip.Create(golangBuffer, module.Version{
		Path:    options.GolangImportPath,
		Version: options.GolangPackageVersion,
	}, files)
	if err != nil {
		return nil, err
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

	rustBuffer := new(bytes.Buffer)
	gzipRustWriter := gzip.NewWriter(rustBuffer)
	tarRustWriter := tar.NewWriter(gzipRustWriter)

	var header *tar.Header
	for _, file := range rustFiles {
		header, err = tar.FileInfoHeader(file, file.Name())
		if err != nil {
			_ = tarRustWriter.Close()
			_ = gzipRustWriter.Close()
			return nil, fmt.Errorf("failed to create tar header for %s: %w", file.Name(), err)
		}

		header.Name = path.Join(fmt.Sprintf("%s-%s", options.RustPackageName, options.RustPackageVersion), header.Name)

		err = tarRustWriter.WriteHeader(header)
		if err != nil {
			_ = tarRustWriter.Close()
			_ = gzipRustWriter.Close()
			return nil, fmt.Errorf("failed to write tar header for %s: %w", file.Name(), err)
		}
		_, err = tarRustWriter.Write(file.Data())
		if err != nil {
			_ = tarRustWriter.Close()
			_ = gzipRustWriter.Close()
			return nil, fmt.Errorf("failed to write tar data for %s: %w", file.Name(), err)
		}
	}

	err = tarRustWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close tar writer: %w", err)
	}

	err = gzipRustWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	return &GuestRegistryPackage{
		GolangModule:  golangBuffer,
		GolangModfile: modfile,
		RustCrate:     rustBuffer,
		RustCargofile: cargofile,
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
