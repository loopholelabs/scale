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

package format

import (
	"context"
	_ "embed"
	"github.com/loopholelabs/scale"
	"github.com/loopholelabs/scale/scalefunc"
	"github.com/loopholelabs/scale/signature/generator/rust/format/signature"
)

//go:embed local-rustfmt-latest.scale
var localRustfmtLatest []byte

type Formatter struct {
	runtime *scale.Scale[*signature.Signature]
}

func New() (*Formatter, error) {
	s := new(scalefunc.Schema)
	err := s.Decode(localRustfmtLatest)
	if err != nil {
		return nil, err
	}

	cfg := scale.NewConfig(signature.New).WithFunction(s)
	runtime, err := scale.New(cfg)
	if err != nil {
		return nil, err
	}

	return &Formatter{
		runtime: runtime,
	}, nil
}

func (f *Formatter) Format(ctx context.Context, code string) (string, error) {
	instance, err := f.runtime.PersistentInstance(ctx)
	if err != nil {
		return "", err
	}

	sig := signature.New()
	sig.Context.Data = code

	err = instance.Run(ctx, sig)
	if err != nil {
		return "", err
	}
	return sig.Context.Data, nil
}
