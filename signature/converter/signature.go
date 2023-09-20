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

package converter

import (
	"encoding/json"

	interfaces "github.com/loopholelabs/scale-signature-interfaces"

	"github.com/loopholelabs/polyglot"

	"github.com/loopholelabs/scale/signature"
)

var _ interfaces.Signature = (*Signature)(nil)

type Signature struct {
	buffer    *polyglot.Buffer
	data      map[string]interface{}
	converter *Converter
}

func NewSignature(schema *signature.Schema) (*Signature, error) {
	converter, err := New(schema)
	if err != nil {
		return nil, err
	}
	return &Signature{
		buffer:    polyglot.NewBuffer(),
		data:      make(map[string]interface{}),
		converter: converter,
	}, nil
}

func (s *Signature) Signature() *Signature {
	return s
}

func (s *Signature) Read(b []byte) (err error) {
	d := polyglot.GetDecoder(b)
	s.data, err = s.converter.FromPolyglot(d)
	polyglot.ReturnDecoder(d)
	return
}

func (s *Signature) Write() []byte {
	s.buffer.Reset()
	err := s.converter.ToPolyglot(s.data, polyglot.Encoder(s.buffer))
	if err != nil {
		return s.Error(err)
	}
	return s.buffer.Bytes()
}

func (s *Signature) Error(err error) []byte {
	s.buffer.Reset()
	polyglot.Encoder(s.buffer).Error(err)
	return s.buffer.Bytes()
}

func (s *Signature) Hash() string {
	return ""
}

func (s *Signature) FromJSON(b []byte) error {
	return json.Unmarshal(b, &s.data)
}

func (s *Signature) ToJSON() ([]byte, error) {
	return json.Marshal(s.data)
}
