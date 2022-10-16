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

package scalefunc

import (
	"github.com/loopholelabs/polyglot-go"
	"github.com/loopholelabs/scale-go/scalefile"
)

type ScaleFunc struct {
	ScaleFile scalefile.ScaleFile
	Tag       string
	Function  []byte
}

func (s *ScaleFunc) Encode() []byte {
	b := polyglot.GetBuffer()
	defer polyglot.PutBuffer(b)
	e := polyglot.Encoder(b)
	e.String(s.ScaleFile.Name)
	e.String(s.Tag)
	e.String(s.ScaleFile.Build.Language)
	e.Uint32(uint32(len(s.ScaleFile.Build.Dependencies)))
	for _, dep := range s.ScaleFile.Build.Dependencies {
		e.String(dep.Name).String(dep.Version)
	}
	e.String(s.ScaleFile.Source)
	e.Bytes(s.Function)

	return b.Bytes()
}

func (s *ScaleFunc) Decode(data []byte) error {
	d := polyglot.GetDecoder(data)
	defer d.Return()

	var err error
	s.ScaleFile.Name, err = d.String()
	if err != nil {
		return err
	}

	s.Tag, err = d.String()
	if err != nil {
		return err
	}

	s.ScaleFile.Build.Language, err = d.String()
	if err != nil {
		return err
	}

	var size uint32
	size, err = d.Uint32()
	if err != nil {
		return err
	}

	s.ScaleFile.Build.Dependencies = make([]scalefile.Dependency, size)
	for i := uint32(0); i < size; i++ {
		s.ScaleFile.Build.Dependencies[i].Name, err = d.String()
		if err != nil {
			return err
		}
		s.ScaleFile.Build.Dependencies[i].Version, err = d.String()
		if err != nil {
			return err
		}
	}

	s.ScaleFile.Source, err = d.String()
	if err != nil {
		return err
	}

	s.Function, err = d.Bytes(nil)
	if err != nil {
		return err
	}

	return nil
}
