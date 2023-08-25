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

package parser

const simpleSignatureJSON = `
{
  "ContextModel": {
    "StringField": "MyString"
  }
}
`

//func TestSimpleSignature(t *testing.T) {
//	s := new(signature.Schema)
//	err := s.Decode([]byte(integration.SimpleExampleSignature))
//	require.NoError(t, err)
//
//	p, err := New(s)
//	require.NoError(t, err)
//
//	d := make(map[string]interface{})
//
//	err = json.Unmarshal([]byte(simpleSignatureJSON), &d)
//	require.NoError(t, err)
//
//	buf := polyglot.NewBuffer()
//	enc := polyglot.Encoder(buf)
//	err = p.Parse(d, enc)
//	require.NoError(t, err)
//
//	ctx := sig.NewContextModel()
//
//	ctx, err = sig.DecodeContextModel(ctx, buf.Bytes())
//	require.NoError(t, err)
//
//	require.Equal(t, "MyString", ctx.StringField)
//
//}
