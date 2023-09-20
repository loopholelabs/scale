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

	"github.com/loopholelabs/scale/extension"
	"github.com/loopholelabs/scale/extension/generator/golang"
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
	Extension *extension.Schema

	GolangPackageImportPath string
	GolangPackageName       string
	GolangPackageVersion    string
}

func GenerateGuestLocal(options *Options) (*GuestLocalPackage, error) {
	golangTypes, err := golang.GenerateTypes(options.Extension, options.GolangPackageName)
	if err != nil {
		return nil, err
	}

	golangGuest, err := golang.GenerateGuest(options.Extension, options.GolangPackageName, options.GolangPackageVersion)
	if err != nil {
		return nil, err
	}

	golangInterfaces, err := golang.GenerateInterfaces(options.Extension, options.GolangPackageName, options.GolangPackageVersion)
	if err != nil {
		return nil, err
	}

	modfile, err := golang.GenerateModfile(options.GolangPackageName)
	if err != nil {
		return nil, err
	}

	golangFiles := []File{
		NewFile("types.go", "types.go", golangTypes),
		NewFile("guest.go", "guest.go", golangGuest),
		NewFile("interfaces.go", "interfaces.go", golangInterfaces),
		NewFile("go.mod", "go.mod", modfile),
	}

	return &GuestLocalPackage{
		GolangFiles: golangFiles,
	}, nil
}

func GenerateHostLocal(options *Options) (*HostLocalPackage, error) {
	golangTypes, err := golang.GenerateTypes(options.Extension, options.GolangPackageName)
	if err != nil {
		return nil, err
	}

	golangHost, err := golang.GenerateHost(options.Extension, options.GolangPackageName, options.GolangPackageVersion)
	if err != nil {
		return nil, err
	}

	golangInterfaces, err := golang.GenerateInterfaces(options.Extension, options.GolangPackageName, options.GolangPackageVersion)
	if err != nil {
		return nil, err
	}

	modfile, err := golang.GenerateModfile(options.GolangPackageName)
	if err != nil {
		return nil, err
	}

	golangFiles := []File{
		NewFile("types.go", "types.go", golangTypes),
		NewFile("host.go", "host.go", golangHost),
		NewFile("interfaces.go", "interfaces.go", golangInterfaces),
		NewFile("go.mod", "go.mod", modfile),
	}

	return &HostLocalPackage{
		GolangFiles: golangFiles,
	}, nil
}
