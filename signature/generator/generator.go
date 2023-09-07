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
	"path"
	"strings"

	"github.com/loopholelabs/scale/signature/generator/typescript"

	"golang.org/x/mod/module"
	"golang.org/x/mod/zip"

	"github.com/loopholelabs/scale/signature"
	"github.com/loopholelabs/scale/signature/generator/golang"
	"github.com/loopholelabs/scale/signature/generator/rust"
)

type GuestRegistryPackage struct {
	GolangModule          *bytes.Buffer
	GolangModfile         []byte
	RustCrate             *bytes.Buffer
	RustCargofile         []byte
	TypescriptPackage     *bytes.Buffer
	TypescriptPackageJSON []byte
}

type GuestLocalPackage struct {
	GolangFiles     []File
	RustFiles       []File
	TypescriptFiles []File
}

type HostRegistryPackage struct {
	GolangModule          *bytes.Buffer
	GolangModfile         []byte
	TypescriptPackage     *bytes.Buffer
	TypescriptPackageJSON []byte
}

type HostLocalPackage struct {
	GolangFiles     []File
	TypescriptFiles []File
}

type Options struct {
	Signature *signature.Schema

	GolangPackageImportPath string
	GolangPackageName       string
	GolangPackageVersion    string

	RustPackageName    string
	RustPackageVersion string

	TypescriptPackageName    string
	TypescriptPackageVersion string
}

func GenerateGuestRegistry(options *Options) (*GuestRegistryPackage, error) {
	hash, err := options.Signature.Hash()
	if err != nil {
		return nil, err
	}
	hashString := hex.EncodeToString(hash)

	golangTypes, err := golang.GenerateTypes(options.Signature, options.GolangPackageName)
	if err != nil {
		return nil, err
	}

	golangGuest, err := golang.GenerateGuest(options.Signature, hashString, options.GolangPackageName)
	if err != nil {
		return nil, err
	}

	modfile, err := golang.GenerateModfile(options.GolangPackageImportPath)
	if err != nil {
		return nil, err
	}

	files := []zip.File{
		NewFile("types.go", "types.go", golangTypes),
		NewFile("guest.go", "guest.go", golangGuest),
		NewFile("go.mod", "go.mod", modfile),
	}

	golangBuffer := new(bytes.Buffer)
	err = zip.Create(golangBuffer, module.Version{
		Path:    options.GolangPackageImportPath,
		Version: options.GolangPackageVersion,
	}, files)
	if err != nil {
		return nil, err
	}

	rustTypes, err := rust.GenerateTypes(options.Signature, options.RustPackageName)
	if err != nil {
		return nil, err
	}

	rustGuest, err := rust.GenerateGuest(options.Signature, hashString, options.RustPackageName)
	if err != nil {
		return nil, err
	}

	cargofile, err := rust.GenerateCargofile(options.RustPackageName, options.RustPackageVersion)
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

		header.Name = path.Join(fmt.Sprintf("%s-%s", options.RustPackageName, strings.TrimPrefix(options.RustPackageVersion, "v")), header.Name)

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

	typescriptTypes, err := typescript.GenerateTypesTranspiled(options.Signature, options.TypescriptPackageName, "types.js")
	if err != nil {
		return nil, err
	}

	typescriptGuest, err := typescript.GenerateGuestTranspiled(options.Signature, hashString, options.TypescriptPackageName, "index.js")
	if err != nil {
		return nil, err
	}

	packageJSON, err := typescript.GeneratePackageJSON(options.TypescriptPackageName, options.TypescriptPackageVersion)
	if err != nil {
		return nil, err
	}

	typescriptFiles := []File{
		NewFile("types.ts", "types.ts", typescriptTypes.Typescript),
		NewFile("types.js", "types.js", typescriptTypes.Javascript),
		NewFile("types.js.map", "types.js.map", typescriptTypes.SourceMap),
		NewFile("types.d.ts", "types.d.ts", typescriptTypes.Declaration),
		NewFile("index.ts", "index.ts", typescriptGuest.Typescript),
		NewFile("index.js", "index.js", typescriptGuest.Javascript),
		NewFile("index.js.map", "index.js.map", typescriptGuest.SourceMap),
		NewFile("index.d.ts", "index.d.ts", typescriptGuest.Declaration),
		NewFile("package.json", "package.json", packageJSON),
	}

	typescriptBuffer := new(bytes.Buffer)
	gzipTypescriptWriter := gzip.NewWriter(typescriptBuffer)
	tarTypescriptWriter := tar.NewWriter(gzipTypescriptWriter)

	for _, file := range typescriptFiles {
		header, err = tar.FileInfoHeader(file, file.Name())
		if err != nil {
			_ = tarTypescriptWriter.Close()
			_ = gzipTypescriptWriter.Close()
			return nil, fmt.Errorf("failed to create tar header for %s: %w", file.Name(), err)
		}

		header.Name = path.Join("package", header.Name)

		err = tarTypescriptWriter.WriteHeader(header)
		if err != nil {
			_ = tarTypescriptWriter.Close()
			_ = gzipTypescriptWriter.Close()
			return nil, fmt.Errorf("failed to write tar header for %s: %w", file.Name(), err)
		}
		_, err = tarTypescriptWriter.Write(file.Data())
		if err != nil {
			_ = tarTypescriptWriter.Close()
			_ = gzipTypescriptWriter.Close()
			return nil, fmt.Errorf("failed to write tar data for %s: %w", file.Name(), err)
		}
	}

	err = tarTypescriptWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close tar writer: %w", err)
	}

	err = gzipTypescriptWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	return &GuestRegistryPackage{
		GolangModule:          golangBuffer,
		GolangModfile:         modfile,
		RustCrate:             rustBuffer,
		RustCargofile:         cargofile,
		TypescriptPackage:     typescriptBuffer,
		TypescriptPackageJSON: packageJSON,
	}, nil
}

func GenerateGuestLocal(options *Options) (*GuestLocalPackage, error) {
	hash, err := options.Signature.Hash()
	if err != nil {
		return nil, err
	}
	hashString := hex.EncodeToString(hash)

	golangTypes, err := golang.GenerateTypes(options.Signature, options.GolangPackageName)
	if err != nil {
		return nil, err
	}

	golangGuest, err := golang.GenerateGuest(options.Signature, hashString, options.GolangPackageName)
	if err != nil {
		return nil, err
	}

	modfile, err := golang.GenerateModfile(options.GolangPackageImportPath)
	if err != nil {
		return nil, err
	}

	golangFiles := []File{
		NewFile("types.go", "types.go", golangTypes),
		NewFile("guest.go", "guest.go", golangGuest),
		NewFile("go.mod", "go.mod", modfile),
	}

	rustTypes, err := rust.GenerateTypes(options.Signature, options.RustPackageName)
	if err != nil {
		return nil, err
	}

	rustGuest, err := rust.GenerateGuest(options.Signature, hashString, options.RustPackageName)
	if err != nil {
		return nil, err
	}

	cargofile, err := rust.GenerateCargofile(options.RustPackageName, options.RustPackageVersion)
	if err != nil {
		return nil, err
	}

	rustFiles := []File{
		NewFile("types.rs", "types.rs", rustTypes),
		NewFile("guest.rs", "guest.rs", rustGuest),
		NewFile("Cargo.toml", "Cargo.toml", cargofile),
	}

	typescriptTypes, err := typescript.GenerateTypesTranspiled(options.Signature, options.TypescriptPackageName, "types.js")
	if err != nil {
		return nil, err
	}

	typescriptGuest, err := typescript.GenerateGuestTranspiled(options.Signature, hashString, options.TypescriptPackageName, "index.js")
	if err != nil {
		return nil, err
	}

	packageJSON, err := typescript.GeneratePackageJSON(options.TypescriptPackageName, options.TypescriptPackageVersion)
	if err != nil {
		return nil, err
	}

	typescriptFiles := []File{
		NewFile("types.ts", "types.ts", typescriptTypes.Typescript),
		NewFile("types.js", "types.js", typescriptTypes.Javascript),
		NewFile("types.js.map", "types.js.map", typescriptTypes.SourceMap),
		NewFile("types.d.ts", "types.d.ts", typescriptTypes.Declaration),
		NewFile("index.ts", "index.ts", typescriptGuest.Typescript),
		NewFile("index.js", "index.js", typescriptGuest.Javascript),
		NewFile("index.js.map", "index.js.map", typescriptGuest.SourceMap),
		NewFile("index.d.ts", "index.d.ts", typescriptGuest.Declaration),
		NewFile("package.json", "package.json", packageJSON),
	}

	return &GuestLocalPackage{
		GolangFiles:     golangFiles,
		RustFiles:       rustFiles,
		TypescriptFiles: typescriptFiles,
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

	golangTypes, err := golang.GenerateTypes(sig, options.GolangPackageName)
	if err != nil {
		return nil, err
	}

	golangHost, err := golang.GenerateHost(sig, hashString, options.GolangPackageName)

	if err != nil {
		return nil, err
	}

	modfile, err := golang.GenerateModfile(options.GolangPackageImportPath)
	if err != nil {
		return nil, err
	}

	files := []zip.File{
		NewFile("types.go", "types.go", golangTypes),
		NewFile("host.go", "host.go", golangHost),
		NewFile("go.mod", "go.mod", modfile),
	}

	golangBuffer := new(bytes.Buffer)
	err = zip.Create(golangBuffer, module.Version{
		Path:    options.GolangPackageImportPath,
		Version: options.GolangPackageVersion,
	}, files)
	if err != nil {
		return nil, err
	}

	typescriptTypes, err := typescript.GenerateTypesTranspiled(sig, options.TypescriptPackageName, "types.js")
	if err != nil {
		return nil, err
	}

	typescriptHost, err := typescript.GenerateHostTranspiled(sig, hashString, options.TypescriptPackageName, "index.js")
	if err != nil {
		return nil, err
	}

	packageJSON, err := typescript.GeneratePackageJSON(options.TypescriptPackageName, options.TypescriptPackageVersion)
	if err != nil {
		return nil, err
	}

	typescriptFiles := []File{
		NewFile("types.ts", "types.ts", typescriptTypes.Typescript),
		NewFile("types.js", "types.js", typescriptTypes.Javascript),
		NewFile("types.js.map", "types.js.map", typescriptTypes.SourceMap),
		NewFile("types.d.ts", "types.d.ts", typescriptTypes.Declaration),
		NewFile("index.ts", "index.ts", typescriptHost.Typescript),
		NewFile("index.js", "index.js", typescriptHost.Javascript),
		NewFile("index.js.map", "index.js.map", typescriptHost.SourceMap),
		NewFile("index.d.ts", "index.d.ts", typescriptHost.Declaration),
		NewFile("package.json", "package.json", packageJSON),
	}

	typescriptBuffer := new(bytes.Buffer)
	gzipTypescriptWriter := gzip.NewWriter(typescriptBuffer)
	tarTypescriptWriter := tar.NewWriter(gzipTypescriptWriter)

	var header *tar.Header
	for _, file := range typescriptFiles {
		header, err = tar.FileInfoHeader(file, file.Name())
		if err != nil {
			_ = tarTypescriptWriter.Close()
			_ = gzipTypescriptWriter.Close()
			return nil, fmt.Errorf("failed to create tar header for %s: %w", file.Name(), err)
		}

		header.Name = path.Join("package", header.Name)

		err = tarTypescriptWriter.WriteHeader(header)
		if err != nil {
			_ = tarTypescriptWriter.Close()
			_ = gzipTypescriptWriter.Close()
			return nil, fmt.Errorf("failed to write tar header for %s: %w", file.Name(), err)
		}
		_, err = tarTypescriptWriter.Write(file.Data())
		if err != nil {
			_ = tarTypescriptWriter.Close()
			_ = gzipTypescriptWriter.Close()
			return nil, fmt.Errorf("failed to write tar data for %s: %w", file.Name(), err)
		}
	}

	err = tarTypescriptWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close tar writer: %w", err)
	}

	err = gzipTypescriptWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	return &HostRegistryPackage{
		GolangModule:          golangBuffer,
		GolangModfile:         modfile,
		TypescriptPackage:     typescriptBuffer,
		TypescriptPackageJSON: packageJSON,
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

	golangTypes, err := golang.GenerateTypes(sig, options.GolangPackageName)
	if err != nil {
		return nil, err
	}

	golangHost, err := golang.GenerateHost(sig, hashString, options.GolangPackageName)
	if err != nil {
		return nil, err
	}

	modfile, err := golang.GenerateModfile(options.GolangPackageImportPath)
	if err != nil {
		return nil, err
	}

	golangFiles := []File{
		NewFile("types.go", "types.go", golangTypes),
		NewFile("host.go", "host.go", golangHost),
		NewFile("go.mod", "go.mod", modfile),
	}

	typescriptTypes, err := typescript.GenerateTypesTranspiled(sig, options.TypescriptPackageName, "types.js")
	if err != nil {
		return nil, err
	}

	typescriptHost, err := typescript.GenerateHostTranspiled(sig, hashString, options.TypescriptPackageName, "index.js")
	if err != nil {
		return nil, err
	}

	packageJSON, err := typescript.GeneratePackageJSON(options.TypescriptPackageName, options.TypescriptPackageVersion)
	if err != nil {
		return nil, err
	}

	typescriptFiles := []File{
		NewFile("types.ts", "types.ts", typescriptTypes.Typescript),
		NewFile("types.js", "types.js", typescriptTypes.Javascript),
		NewFile("types.js.map", "types.js.map", typescriptTypes.SourceMap),
		NewFile("types.d.ts", "types.d.ts", typescriptTypes.Declaration),
		NewFile("index.ts", "index.ts", typescriptHost.Typescript),
		NewFile("index.js", "index.js", typescriptHost.Javascript),
		NewFile("index.js.map", "index.js.map", typescriptHost.SourceMap),
		NewFile("index.d.ts", "index.d.ts", typescriptHost.Declaration),
		NewFile("package.json", "package.json", packageJSON),
	}

	return &HostLocalPackage{
		GolangFiles:     golangFiles,
		TypescriptFiles: typescriptFiles,
	}, nil
}
