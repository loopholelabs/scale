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

//go:generate go test ./... -tags=generate -v

const testSchema = `
version = "v1alpha"
context = "Context"

enum GenericEnum {
	values = ["FirstValue", "SecondValue", "DefaultValue"]
}

model EmptyModel {}

model EmbeddedModel {
	string StringField {
		default = "DefaultValue"
	}
}

model Context {
    enum EnumField {
		default = "DefaultValue"
		reference = "GenericEnum"
	}

    enum_array EnumArrayField {
		reference = "GenericEnum"
		initial_size = 0
	}

	enum_map EnumMapField {
		reference = "GenericEnum"
		value = "string"
	}

    enum_map EnumMapModelField {
		reference = "GenericEnum"
		value = "EmbeddedModel"
	}

    model ModelField {
		reference = "EmbeddedModel"
	}
 
    model_array ModelArrayField {
		reference = "EmbeddedModel"
		initial_size = 0
	}

    model EmptyModelField {
		reference = "EmptyModel"
	}

    model_array EmptyModelArrayField {
		reference = "EmptyModel"
		initial_size = 0
	}

    string StringField {
		default = "DefaultValue"
	}

    string_array StringArrayField {
		initial_size = 0
	}

    string_map StringMapField {
		value = "string"
	}

    string_map StringModelMapField {
		value = "EmbeddedModel"
	}

    int32 Int32Field {
		default = 32
	}

    int32_array Int32ArrayField {
		initial_size = 0
	}

    int32_map Int32MapField {
		value = "int32"
	}

    int32_map Int32MapModelField {
		value = "EmbeddedModel"
	}

	int64 Int64Field {
		default = 64
	}

	int64_array Int64ArrayField {
		initial_size = 0
	}

	int64_map Int64MapField {
		value = "int64"
	}

	int64_map Int64MapModelField {
		value = "EmbeddedModel"
	}

    uint32 Uint32Field {
		default = 32
	}

	uint32_array Uint32ArrayField {
		initial_size = 0
	}

	uint32_map Uint32MapField {
		value = "uint32"
	}

	uint32_map Uint32MapModelField {
		value = "EmbeddedModel"
	}

	uint64 Uint64Field {
		default = 64
	}

	uint64_array Uint64ArrayField {
		initial_size = 0
	}

	uint64_map Uint64MapField {
		value = "uint64"
	}

	uint64_map Uint64MapModelField {
		value = "EmbeddedModel"
	}

    float32 Float32Field {
		default = 32.32
	}

	float32_array Float32ArrayField {
		initial_size = 0
	}

	float64 Float64Field {
		default = 64.64
	}

	float64_array Float64ArrayField {
		initial_size = 0
	}

    bool BoolField {
		default = true
	}

	bool_array BoolArrayField {
		initial_size = 0
	}

	bytes BytesField {
		initial_size = 512
	}

	bytes_array BytesArrayField {
		initial_size = 0
	}
}
`
