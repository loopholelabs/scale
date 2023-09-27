//go:build !integration && !generate

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
	"testing"

	"github.com/loopholelabs/polyglot"
	"github.com/stretchr/testify/require"

	"github.com/loopholelabs/scale/signature"
	generated "github.com/loopholelabs/scale/signature/converter/converter_tests"
)

const jsonData = `
{
  "Context": {
    "EnumField": "SecondValue",
	"EnumArrayField": ["FirstValue", "SecondValue", "DefaultValue"],
	"EnumMapField": {
		"FirstValue": "string1",
		"SecondValue": "string2",
		"DefaultValue": "string3"
	},
	"EnumMapModelField": {
		"FirstValue": {
			"StringField": "string1"
		},
		"SecondValue": {
			"StringField": "string2"
		},	
		"DefaultValue": {
			"StringField": "string3"
		}
	},
	"ModelField": {
		"StringField": "string1"
	},
	"ModelArrayField": [
		{
			"StringField": "string1"
		},
		{
			"StringField": "string2"
		},
		{	
			"StringField": "string3"
		}
	],
	"EmptyModelField": {},
	"EmptyModelArrayField": [],
	"StringField": "MyString",
	"StringArrayField": ["MyString1", "MyString2", "MyString3"],
	"StringMapField": {
		"key1": "MyString1",
		"key2": "MyString2",
		"key3": "MyString3"
	},
	"StringModelMapField": {
		"key1": {
			"StringField": "MyString1"
		},
		"key2": {	
			"StringField": "MyString2"
		},
		"key3": {
			"StringField": "MyString3"
		}
	},
	"Int32Field": -32,
	"Int32ArrayField": [-32, -64, -128],
	"Int32MapField": {
		"-32": -32,
		"-64": -64,
		"-128": -128
	},
	"Int32MapModelField": {
		"-32": {
			"StringField": "MyString1"
		},	
		"-64": {	
			"StringField": "MyString2"
		},
		"-128": {
			"StringField": "MyString3"
		}
	},
	"Int64Field": -64,
	"Int64ArrayField": [-64, -128, -256],
	"Int64MapField": {
		"-64": -64,
		"-128": -128,
		"-256": -256
	},
	"Int64MapModelField": {
		"-64": {
			"StringField": "MyString1"
		},
		"-128": {
			"StringField": "MyString2"
		},
		"-256": {	
			"StringField": "MyString3"
		}
	},
	"Uint32Field": 32,
	"Uint32ArrayField": [32, 64, 128],
	"Uint32MapField": {
		"32": 32,
		"64": 64,
		"128": 128
	},
	"Uint32MapModelField": {
		"32": {
			"StringField": "MyString1"
		},	
		"64": {	
			"StringField": "MyString2"
		},	
		"128": {
			"StringField": "MyString3"
		}
	},
	"Uint64Field": 64,
	"Uint64ArrayField": [64, 128, 256],
	"Uint64MapField": {
		"64": 64,
		"128": 128,
		"256": 256
	},
	"Uint64MapModelField": {
		"64": {	
			"StringField": "MyString1"
		},
		"128": {	
			"StringField": "MyString2"
		},
		"256": {
			"StringField": "MyString3"
		}
	},
	"Float32Field": 32.32,
	"Float32ArrayField": [32.32, 64.64, 128.128],
	"Float64Field": 64.64,
	"Float64ArrayField": [64.64, 128.128, 256.256],
	"BoolField": true,
	"BoolArrayField": [true, false, true],
	"BytesField": "dGVzdGluZzEyMw==",
	"BytesArrayField": ["dGVzdGluZzEyMw==", "dGVzdGluZzEyNA==", "dGVzdGluZzEyNQ=="]
  }
}
`

func TestConverter(t *testing.T) {
	s := new(signature.Schema)
	err := s.Decode([]byte(testSchema))
	require.NoError(t, err)

	d := make(map[string]interface{})

	err = json.Unmarshal([]byte(jsonData), &d)
	require.NoError(t, err)

	buf := polyglot.NewBuffer()
	enc := polyglot.Encoder(buf)
	err = ToPolyglot(s, d, enc)
	require.NoError(t, err)

	ctx := generated.NewContext()

	ctx, err = generated.DecodeContext(ctx, buf.Bytes())
	require.NoError(t, err)

	require.Equal(t, generated.GenericEnumSecondValue, ctx.EnumField)
	require.Equal(t, []generated.GenericEnum{
		generated.GenericEnumFirstValue,
		generated.GenericEnumSecondValue,
		generated.GenericEnumDefaultValue,
	}, ctx.EnumArrayField)
	require.Equal(t, map[generated.GenericEnum]string{
		generated.GenericEnumFirstValue:   "string1",
		generated.GenericEnumSecondValue:  "string2",
		generated.GenericEnumDefaultValue: "string3",
	}, ctx.EnumMapField)
	require.Equal(t, map[generated.GenericEnum]generated.EmbeddedModel{
		generated.GenericEnumFirstValue:   {StringField: "string1"},
		generated.GenericEnumSecondValue:  {StringField: "string2"},
		generated.GenericEnumDefaultValue: {StringField: "string3"},
	}, ctx.EnumMapModelField)
	require.Equal(t, &generated.EmbeddedModel{StringField: "string1"}, ctx.ModelField)
	require.Equal(t, []generated.EmbeddedModel{
		{StringField: "string1"},
		{StringField: "string2"},
		{StringField: "string3"},
	}, ctx.ModelArrayField)
	require.Equal(t, &generated.EmptyModel{}, ctx.EmptyModelField)
	require.Equal(t, []generated.EmptyModel{}, ctx.EmptyModelArrayField)
	require.Equal(t, "MyString", ctx.StringField)
	require.Equal(t, []string{"MyString1", "MyString2", "MyString3"}, ctx.StringArrayField)
	require.Equal(t, map[string]string{
		"key1": "MyString1",
		"key2": "MyString2",
		"key3": "MyString3",
	}, ctx.StringMapField)
	require.Equal(t, map[string]generated.EmbeddedModel{
		"key1": {StringField: "MyString1"},
		"key2": {StringField: "MyString2"},
		"key3": {StringField: "MyString3"},
	}, ctx.StringModelMapField)
	require.Equal(t, int32(-32), ctx.Int32Field)
	require.Equal(t, []int32{-32, -64, -128}, ctx.Int32ArrayField)
	require.Equal(t, map[int32]int32{
		-32:  -32,
		-64:  -64,
		-128: -128,
	}, ctx.Int32MapField)
	require.Equal(t, map[int32]generated.EmbeddedModel{
		-32:  {StringField: "MyString1"},
		-64:  {StringField: "MyString2"},
		-128: {StringField: "MyString3"},
	}, ctx.Int32MapModelField)
	require.Equal(t, int64(-64), ctx.Int64Field)
	require.Equal(t, []int64{-64, -128, -256}, ctx.Int64ArrayField)
	require.Equal(t, map[int64]int64{
		-64:  -64,
		-128: -128,
		-256: -256,
	}, ctx.Int64MapField)
	require.Equal(t, map[int64]generated.EmbeddedModel{
		-64:  {StringField: "MyString1"},
		-128: {StringField: "MyString2"},
		-256: {StringField: "MyString3"},
	}, ctx.Int64MapModelField)
	require.Equal(t, uint32(32), ctx.Uint32Field)
	require.Equal(t, []uint32{32, 64, 128}, ctx.Uint32ArrayField)
	require.Equal(t, map[uint32]uint32{
		32:  32,
		64:  64,
		128: 128,
	}, ctx.Uint32MapField)
	require.Equal(t, map[uint32]generated.EmbeddedModel{
		32:  {StringField: "MyString1"},
		64:  {StringField: "MyString2"},
		128: {StringField: "MyString3"},
	}, ctx.Uint32MapModelField)
	require.Equal(t, uint64(64), ctx.Uint64Field)
	require.Equal(t, []uint64{64, 128, 256}, ctx.Uint64ArrayField)
	require.Equal(t, map[uint64]uint64{
		64:  64,
		128: 128,
		256: 256,
	}, ctx.Uint64MapField)
	require.Equal(t, map[uint64]generated.EmbeddedModel{
		64:  {StringField: "MyString1"},
		128: {StringField: "MyString2"},
		256: {StringField: "MyString3"},
	}, ctx.Uint64MapModelField)
	require.Equal(t, float32(32.32), ctx.Float32Field)
	require.Equal(t, []float32{32.32, 64.64, 128.128}, ctx.Float32ArrayField)
	require.Equal(t, 64.64, ctx.Float64Field)
	require.Equal(t, []float64{64.64, 128.128, 256.256}, ctx.Float64ArrayField)
	require.Equal(t, true, ctx.BoolField)
	require.Equal(t, []bool{true, false, true}, ctx.BoolArrayField)
	require.Equal(t, []byte("testing123"), ctx.BytesField)
	require.Equal(t, [][]byte{[]byte("testing123"), []byte("testing124"), []byte("testing125")}, ctx.BytesArrayField)

	dec := polyglot.GetDecoder(buf.Bytes())
	data, err := FromPolyglot(s, dec)
	require.NoError(t, err)

	require.Equal(t, d["Context"].(map[string]interface{})["ModelField"], data["Context"].(map[string]interface{})["ModelField"])
	require.Equal(t, d["Context"].(map[string]interface{})["ModelArrayField"], data["Context"].(map[string]interface{})["ModelArrayField"])
	require.Equal(t, d["Context"].(map[string]interface{})["StringField"], data["Context"].(map[string]interface{})["StringField"])
	require.Equal(t, d["Context"].(map[string]interface{})["StringArrayField"], data["Context"].(map[string]interface{})["StringArrayField"])
	require.Equal(t, d["Context"].(map[string]interface{})["StringMapField"], data["Context"].(map[string]interface{})["StringMapField"])
	require.Equal(t, d["Context"].(map[string]interface{})["StringModelMapField"], data["Context"].(map[string]interface{})["StringModelMapField"])
	require.Equal(t, d["Context"].(map[string]interface{})["Int32Field"], data["Context"].(map[string]interface{})["Int32Field"])
	require.Equal(t, d["Context"].(map[string]interface{})["Int32ArrayField"], data["Context"].(map[string]interface{})["Int32ArrayField"])
	require.Equal(t, d["Context"].(map[string]interface{})["Int32MapField"], data["Context"].(map[string]interface{})["Int32MapField"])
	require.Equal(t, d["Context"].(map[string]interface{})["Int32MapModelField"], data["Context"].(map[string]interface{})["Int32MapModelField"])
	require.Equal(t, d["Context"].(map[string]interface{})["Int64Field"], data["Context"].(map[string]interface{})["Int64Field"])
	require.Equal(t, d["Context"].(map[string]interface{})["Int64ArrayField"], data["Context"].(map[string]interface{})["Int64ArrayField"])
	require.Equal(t, d["Context"].(map[string]interface{})["Int64MapField"], data["Context"].(map[string]interface{})["Int64MapField"])
	require.Equal(t, d["Context"].(map[string]interface{})["Int64MapModelField"], data["Context"].(map[string]interface{})["Int64MapModelField"])
	require.Equal(t, d["Context"].(map[string]interface{})["Uint32Field"], data["Context"].(map[string]interface{})["Uint32Field"])
	require.Equal(t, d["Context"].(map[string]interface{})["Uint32ArrayField"], data["Context"].(map[string]interface{})["Uint32ArrayField"])
	require.Equal(t, d["Context"].(map[string]interface{})["Uint32MapField"], data["Context"].(map[string]interface{})["Uint32MapField"])
	require.Equal(t, d["Context"].(map[string]interface{})["Uint32MapModelField"], data["Context"].(map[string]interface{})["Uint32MapModelField"])
	require.Equal(t, d["Context"].(map[string]interface{})["Uint64Field"], data["Context"].(map[string]interface{})["Uint64Field"])
	require.Equal(t, d["Context"].(map[string]interface{})["Uint64ArrayField"], data["Context"].(map[string]interface{})["Uint64ArrayField"])
	require.Equal(t, d["Context"].(map[string]interface{})["Uint64MapField"], data["Context"].(map[string]interface{})["Uint64MapField"])
	require.Equal(t, d["Context"].(map[string]interface{})["Uint64MapModelField"], data["Context"].(map[string]interface{})["Uint64MapModelField"])
	require.Equal(t, float32(d["Context"].(map[string]interface{})["Float32Field"].(float64)), float32(data["Context"].(map[string]interface{})["Float32Field"].(float64)))
	require.Equal(t, float32(d["Context"].(map[string]interface{})["Float32ArrayField"].([]interface{})[0].(float64)), float32(data["Context"].(map[string]interface{})["Float32ArrayField"].([]interface{})[0].(float64)))
	require.Equal(t, float32(d["Context"].(map[string]interface{})["Float32ArrayField"].([]interface{})[1].(float64)), float32(data["Context"].(map[string]interface{})["Float32ArrayField"].([]interface{})[1].(float64)))
	require.Equal(t, float32(d["Context"].(map[string]interface{})["Float32ArrayField"].([]interface{})[2].(float64)), float32(data["Context"].(map[string]interface{})["Float32ArrayField"].([]interface{})[2].(float64)))
	require.Equal(t, d["Context"].(map[string]interface{})["Float64Field"], data["Context"].(map[string]interface{})["Float64Field"])
	require.Equal(t, d["Context"].(map[string]interface{})["Float64ArrayField"], data["Context"].(map[string]interface{})["Float64ArrayField"])
	require.Equal(t, d["Context"].(map[string]interface{})["BoolField"], data["Context"].(map[string]interface{})["BoolField"])
	require.Equal(t, d["Context"].(map[string]interface{})["BoolArrayField"], data["Context"].(map[string]interface{})["BoolArrayField"])
	require.Equal(t, d["Context"].(map[string]interface{})["BytesField"], data["Context"].(map[string]interface{})["BytesField"])
	require.Equal(t, d["Context"].(map[string]interface{})["BytesArrayField"], data["Context"].(map[string]interface{})["BytesArrayField"])
}
