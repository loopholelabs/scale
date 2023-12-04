// Code generated by scale-signature 0.4.5, DO NOT EDIT.
// output: local-example-latest-guest

"use strict";
var __defProp = Object.defineProperty;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __export = (target, all) => {
  for (var name in all)
    __defProp(target, name, { get: all[name], enumerable: true });
};
var __copyProps = (to, from, except, desc) => {
  if (from && typeof from === "object" || typeof from === "function") {
    for (let key of __getOwnPropNames(from))
      if (!__hasOwnProp.call(to, key) && key !== except)
        __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
  }
  return to;
};
var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);
var stdin_exports = {};
__export(stdin_exports, {
  EmptyModel: () => EmptyModel,
  EmptyModelWithDescription: () => EmptyModelWithDescription,
  GenericEnum: () => GenericEnum,
  ModelWithAllFieldTypes: () => ModelWithAllFieldTypes,
  ModelWithEmbeddedModels: () => ModelWithEmbeddedModels,
  ModelWithEmbeddedModelsAccessor: () => ModelWithEmbeddedModelsAccessor,
  ModelWithEmbeddedModelsAccessorAndDescription: () => ModelWithEmbeddedModelsAccessorAndDescription,
  ModelWithEmbeddedModelsAndDescription: () => ModelWithEmbeddedModelsAndDescription,
  ModelWithEnum: () => ModelWithEnum,
  ModelWithEnumAccessor: () => ModelWithEnumAccessor,
  ModelWithEnumAccessorAndDescription: () => ModelWithEnumAccessorAndDescription,
  ModelWithEnumAndDescription: () => ModelWithEnumAndDescription,
  ModelWithMultipleFields: () => ModelWithMultipleFields,
  ModelWithMultipleFieldsAccessor: () => ModelWithMultipleFieldsAccessor,
  ModelWithMultipleFieldsAccessorAndDescription: () => ModelWithMultipleFieldsAccessorAndDescription,
  ModelWithMultipleFieldsAndDescription: () => ModelWithMultipleFieldsAndDescription,
  ModelWithSingleInt32Field: () => ModelWithSingleInt32Field,
  ModelWithSingleInt32FieldAndDescription: () => ModelWithSingleInt32FieldAndDescription,
  ModelWithSingleStringField: () => ModelWithSingleStringField,
  ModelWithSingleStringFieldAndDescription: () => ModelWithSingleStringFieldAndDescription
});
module.exports = __toCommonJS(stdin_exports);
var import_polyglot = require("@loopholelabs/polyglot");
var GenericEnum = /* @__PURE__ */ ((GenericEnum2) => {
  GenericEnum2[GenericEnum2["FirstValue"] = 0] = "FirstValue";
  GenericEnum2[GenericEnum2["SecondValue"] = 1] = "SecondValue";
  GenericEnum2[GenericEnum2["DefaultValue"] = 2] = "DefaultValue";
  return GenericEnum2;
})(GenericEnum || {});
class EmptyModel {
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
    } else {
    }
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new EmptyModel(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
class EmptyModelWithDescription {
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
    } else {
    }
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new EmptyModelWithDescription(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
class ModelWithSingleStringField {
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
      this.stringField = decoder.string();
    } else {
      this.stringField = "DefaultValue";
    }
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
    encoder.string(this.stringField);
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new ModelWithSingleStringField(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
class ModelWithSingleStringFieldAndDescription {
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
      this.stringField = decoder.string();
    } else {
      this.stringField = "DefaultValue";
    }
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
    encoder.string(this.stringField);
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new ModelWithSingleStringFieldAndDescription(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
class ModelWithSingleInt32Field {
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
      this.int32Field = decoder.int32();
    } else {
      this.int32Field = 32;
    }
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
    encoder.int32(this.int32Field);
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new ModelWithSingleInt32Field(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
class ModelWithSingleInt32FieldAndDescription {
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
      this.int32Field = decoder.int32();
    } else {
      this.int32Field = 32;
    }
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
    encoder.int32(this.int32Field);
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new ModelWithSingleInt32FieldAndDescription(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
class ModelWithMultipleFields {
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
      this.stringField = decoder.string();
      this.int32Field = decoder.int32();
    } else {
      this.stringField = "DefaultValue";
      this.int32Field = 32;
    }
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
    encoder.string(this.stringField);
    encoder.int32(this.int32Field);
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new ModelWithMultipleFields(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
class ModelWithMultipleFieldsAndDescription {
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
      this.stringField = decoder.string();
      this.int32Field = decoder.int32();
    } else {
      this.stringField = "DefaultValue";
      this.int32Field = 32;
    }
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
    encoder.string(this.stringField);
    encoder.int32(this.int32Field);
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new ModelWithMultipleFieldsAndDescription(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
class ModelWithEnum {
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
      this.enumField = decoder.uint32();
    } else {
      this.enumField = 2 /* DefaultValue */;
    }
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
    encoder.uint32(this.enumField);
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new ModelWithEnum(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
class ModelWithEnumAndDescription {
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
      this.enumField = decoder.uint32();
    } else {
      this.enumField = 2 /* DefaultValue */;
    }
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
    encoder.uint32(this.enumField);
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new ModelWithEnumAndDescription(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
class ModelWithEnumAccessor {
  #enumField;
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
      this.#enumField = decoder.uint32();
    } else {
      this.#enumField = 2 /* DefaultValue */;
    }
  }
  get enumField() {
    return this.#enumField;
  }
  set enumField(val) {
    this.#enumField = val;
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
    encoder.uint32(this.#enumField);
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new ModelWithEnumAccessor(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
class ModelWithEnumAccessorAndDescription {
  #enumField;
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
      this.#enumField = decoder.uint32();
    } else {
      this.#enumField = 2 /* DefaultValue */;
    }
  }
  get enumField() {
    return this.#enumField;
  }
  set enumField(val) {
    this.#enumField = val;
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
    encoder.uint32(this.#enumField);
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new ModelWithEnumAccessorAndDescription(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
class ModelWithMultipleFieldsAccessor {
  #stringField;
  #int32Field;
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
      this.#stringField = decoder.string();
      this.#int32Field = decoder.int32();
    } else {
      this.#stringField = "DefaultValue";
      this.#int32Field = 32;
    }
  }
  get stringField() {
    return this.#stringField;
  }
  set stringField(val) {
    if (!/^[a-zA-Z0-9]*$/.test(val)) {
      throw new Error("value must match ^[a-zA-Z0-9]*$");
    }
    if (val.length > 20 || val.length < 1) {
      throw new Error("length must be between 1 and 20");
    }
    val = val.toUpperCase();
    this.#stringField = val;
  }
  get int32Field() {
    return this.#int32Field;
  }
  set int32Field(val) {
    if (val > 100 || val < 0) {
      throw new Error("value must be between 0 and 100");
    }
    this.#int32Field = val;
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
    encoder.string(this.#stringField);
    encoder.int32(this.#int32Field);
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new ModelWithMultipleFieldsAccessor(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
class ModelWithMultipleFieldsAccessorAndDescription {
  #stringField;
  #int32Field;
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
      this.#stringField = decoder.string();
      this.#int32Field = decoder.int32();
    } else {
      this.#stringField = "DefaultValue";
      this.#int32Field = 32;
    }
  }
  get stringField() {
    return this.#stringField;
  }
  set stringField(val) {
    this.#stringField = val;
  }
  get int32Field() {
    return this.#int32Field;
  }
  set int32Field(val) {
    this.#int32Field = val;
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
    encoder.string(this.#stringField);
    encoder.int32(this.#int32Field);
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new ModelWithMultipleFieldsAccessorAndDescription(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
class ModelWithEmbeddedModels {
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
      this.embeddedEmptyModel = EmptyModel.decode(decoder);
      const embeddedModelArrayWithMultipleFieldsAccessorSize = decoder.array(import_polyglot.Kind.Any);
      this.embeddedModelArrayWithMultipleFieldsAccessor = new Array(embeddedModelArrayWithMultipleFieldsAccessorSize);
      for (let i = 0; i < embeddedModelArrayWithMultipleFieldsAccessorSize; i += 1) {
        const x = ModelWithMultipleFieldsAccessor.decode(decoder);
        if (typeof x !== "undefined") {
          this.embeddedModelArrayWithMultipleFieldsAccessor[i] = x;
        }
      }
    } else {
      this.embeddedEmptyModel = new EmptyModel();
      this.embeddedModelArrayWithMultipleFieldsAccessor = [];
    }
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
    if (typeof this.embeddedEmptyModel === "undefined") {
      encoder.null();
    } else {
      this.embeddedEmptyModel.encode(encoder);
    }
    const embeddedModelArrayWithMultipleFieldsAccessorLength = this.embeddedModelArrayWithMultipleFieldsAccessor.length;
    encoder.array(embeddedModelArrayWithMultipleFieldsAccessorLength, import_polyglot.Kind.Any);
    for (let i = 0; i < embeddedModelArrayWithMultipleFieldsAccessorLength; i += 1) {
      const el = this.embeddedModelArrayWithMultipleFieldsAccessor[i];
      el.encode(encoder);
    }
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new ModelWithEmbeddedModels(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
class ModelWithEmbeddedModelsAndDescription {
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
      this.embeddedEmptyModel = EmptyModel.decode(decoder);
      const embeddedModelArrayWithMultipleFieldsAccessorSize = decoder.array(import_polyglot.Kind.Any);
      this.embeddedModelArrayWithMultipleFieldsAccessor = new Array(embeddedModelArrayWithMultipleFieldsAccessorSize);
      for (let i = 0; i < embeddedModelArrayWithMultipleFieldsAccessorSize; i += 1) {
        const x = ModelWithMultipleFieldsAccessor.decode(decoder);
        if (typeof x !== "undefined") {
          this.embeddedModelArrayWithMultipleFieldsAccessor[i] = x;
        }
      }
    } else {
      this.embeddedEmptyModel = new EmptyModel();
      this.embeddedModelArrayWithMultipleFieldsAccessor = [];
    }
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
    if (typeof this.embeddedEmptyModel === "undefined") {
      encoder.null();
    } else {
      this.embeddedEmptyModel.encode(encoder);
    }
    const embeddedModelArrayWithMultipleFieldsAccessorLength = this.embeddedModelArrayWithMultipleFieldsAccessor.length;
    encoder.array(embeddedModelArrayWithMultipleFieldsAccessorLength, import_polyglot.Kind.Any);
    for (let i = 0; i < embeddedModelArrayWithMultipleFieldsAccessorLength; i += 1) {
      const el = this.embeddedModelArrayWithMultipleFieldsAccessor[i];
      el.encode(encoder);
    }
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new ModelWithEmbeddedModelsAndDescription(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
class ModelWithEmbeddedModelsAccessor {
  #embeddedEmptyModel;
  #embeddedModelArrayWithMultipleFieldsAccessor;
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
      this.#embeddedEmptyModel = EmptyModel.decode(decoder);
      const embeddedModelArrayWithMultipleFieldsAccessorSize = decoder.array(import_polyglot.Kind.Any);
      this.#embeddedModelArrayWithMultipleFieldsAccessor = new Array(embeddedModelArrayWithMultipleFieldsAccessorSize);
      for (let i = 0; i < embeddedModelArrayWithMultipleFieldsAccessorSize; i += 1) {
        const x = ModelWithMultipleFieldsAccessor.decode(decoder);
        if (typeof x !== "undefined") {
          this.#embeddedModelArrayWithMultipleFieldsAccessor[i] = x;
        }
      }
    } else {
      this.#embeddedEmptyModel = new EmptyModel();
      this.#embeddedModelArrayWithMultipleFieldsAccessor = [];
    }
  }
  get embeddedEmptyModel() {
    return this.#embeddedEmptyModel;
  }
  set embeddedEmptyModel(val) {
    this.#embeddedEmptyModel = val;
  }
  get embeddedModelArrayWithMultipleFieldsAccessor() {
    return this.#embeddedModelArrayWithMultipleFieldsAccessor;
  }
  set EmbeddedModelArrayWithMultipleFieldsAccessor(val) {
    this.#embeddedModelArrayWithMultipleFieldsAccessor = val;
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
    if (typeof this.#embeddedEmptyModel === "undefined") {
      encoder.null();
    } else {
      this.#embeddedEmptyModel.encode(encoder);
    }
    const embeddedModelArrayWithMultipleFieldsAccessorLength = this.#embeddedModelArrayWithMultipleFieldsAccessor.length;
    encoder.array(embeddedModelArrayWithMultipleFieldsAccessorLength, import_polyglot.Kind.Any);
    for (let i = 0; i < embeddedModelArrayWithMultipleFieldsAccessorLength; i += 1) {
      const el = this.#embeddedModelArrayWithMultipleFieldsAccessor[i];
      el.encode(encoder);
    }
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new ModelWithEmbeddedModelsAccessor(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
class ModelWithEmbeddedModelsAccessorAndDescription {
  #embeddedEmptyModel;
  #embeddedModelArrayWithMultipleFieldsAccessor;
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
      this.#embeddedEmptyModel = EmptyModel.decode(decoder);
      const embeddedModelArrayWithMultipleFieldsAccessorSize = decoder.array(import_polyglot.Kind.Any);
      this.#embeddedModelArrayWithMultipleFieldsAccessor = new Array(embeddedModelArrayWithMultipleFieldsAccessorSize);
      for (let i = 0; i < embeddedModelArrayWithMultipleFieldsAccessorSize; i += 1) {
        const x = ModelWithMultipleFieldsAccessor.decode(decoder);
        if (typeof x !== "undefined") {
          this.#embeddedModelArrayWithMultipleFieldsAccessor[i] = x;
        }
      }
    } else {
      this.#embeddedEmptyModel = new EmptyModel();
      this.#embeddedModelArrayWithMultipleFieldsAccessor = [];
    }
  }
  get embeddedEmptyModel() {
    return this.#embeddedEmptyModel;
  }
  set embeddedEmptyModel(val) {
    this.#embeddedEmptyModel = val;
  }
  get embeddedModelArrayWithMultipleFieldsAccessor() {
    return this.#embeddedModelArrayWithMultipleFieldsAccessor;
  }
  set EmbeddedModelArrayWithMultipleFieldsAccessor(val) {
    this.#embeddedModelArrayWithMultipleFieldsAccessor = val;
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
    if (typeof this.#embeddedEmptyModel === "undefined") {
      encoder.null();
    } else {
      this.#embeddedEmptyModel.encode(encoder);
    }
    const embeddedModelArrayWithMultipleFieldsAccessorLength = this.#embeddedModelArrayWithMultipleFieldsAccessor.length;
    encoder.array(embeddedModelArrayWithMultipleFieldsAccessorLength, import_polyglot.Kind.Any);
    for (let i = 0; i < embeddedModelArrayWithMultipleFieldsAccessorLength; i += 1) {
      const el = this.#embeddedModelArrayWithMultipleFieldsAccessor[i];
      el.encode(encoder);
    }
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new ModelWithEmbeddedModelsAccessorAndDescription(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
class ModelWithAllFieldTypes {
  /**
  * @throws {Error}
  */
  constructor(decoder) {
    if (decoder) {
      let err;
      try {
        err = decoder.error();
      } catch (_) {
      }
      if (typeof err !== "undefined") {
        throw err;
      }
      this.modelField = EmptyModel.decode(decoder);
      const modelArrayFieldSize = decoder.array(import_polyglot.Kind.Any);
      this.modelArrayField = new Array(modelArrayFieldSize);
      for (let i = 0; i < modelArrayFieldSize; i += 1) {
        const x = EmptyModel.decode(decoder);
        if (typeof x !== "undefined") {
          this.modelArrayField[i] = x;
        }
      }
      this.stringField = decoder.string();
      const stringArrayFieldSize = decoder.array(import_polyglot.Kind.String);
      this.stringArrayField = new Array(stringArrayFieldSize);
      for (let i = 0; i < stringArrayFieldSize; i += 1) {
        this.stringArrayField[i] = decoder.string();
      }
      this.stringMapField = /* @__PURE__ */ new Map();
      let stringMapFieldSize = decoder.map(import_polyglot.Kind.String, import_polyglot.Kind.String);
      for (let i = 0; i < stringMapFieldSize; i++) {
        let key = decoder.string();
        let val = decoder.string();
        this.stringMapField.set(key, val);
      }
      this.stringMapFieldEmbedded = /* @__PURE__ */ new Map();
      let stringMapFieldEmbeddedSize = decoder.map(import_polyglot.Kind.String, import_polyglot.Kind.Any);
      for (let i = 0; i < stringMapFieldEmbeddedSize; i++) {
        let key = decoder.string();
        let val = EmptyModel.decode(decoder);
        if (typeof val !== "undefined") {
          this.stringMapFieldEmbedded.set(key, val);
        }
      }
      this.int32Field = decoder.int32();
      const int32ArrayFieldSize = decoder.array(import_polyglot.Kind.Int32);
      this.int32ArrayField = new Array(int32ArrayFieldSize);
      for (let i = 0; i < int32ArrayFieldSize; i += 1) {
        this.int32ArrayField[i] = decoder.int32();
      }
      this.int32MapField = /* @__PURE__ */ new Map();
      let int32MapFieldSize = decoder.map(import_polyglot.Kind.Int32, import_polyglot.Kind.Int32);
      for (let i = 0; i < int32MapFieldSize; i++) {
        let key = decoder.int32();
        let val = decoder.int32();
        this.int32MapField.set(key, val);
      }
      this.int32MapFieldEmbedded = /* @__PURE__ */ new Map();
      let int32MapFieldEmbeddedSize = decoder.map(import_polyglot.Kind.Int32, import_polyglot.Kind.Any);
      for (let i = 0; i < int32MapFieldEmbeddedSize; i++) {
        let key = decoder.int32();
        let val = EmptyModel.decode(decoder);
        if (typeof val !== "undefined") {
          this.int32MapFieldEmbedded.set(key, val);
        }
      }
      this.int64Field = decoder.int64();
      const int64ArrayFieldSize = decoder.array(import_polyglot.Kind.Int64);
      this.int64ArrayField = new Array(int64ArrayFieldSize);
      for (let i = 0; i < int64ArrayFieldSize; i += 1) {
        this.int64ArrayField[i] = decoder.int64();
      }
      this.int64MapField = /* @__PURE__ */ new Map();
      let int64MapFieldSize = decoder.map(import_polyglot.Kind.Int64, import_polyglot.Kind.Int64);
      for (let i = 0; i < int64MapFieldSize; i++) {
        let key = decoder.int64();
        let val = decoder.int64();
        this.int64MapField.set(key, val);
      }
      this.int64MapFieldEmbedded = /* @__PURE__ */ new Map();
      let int64MapFieldEmbeddedSize = decoder.map(import_polyglot.Kind.Int64, import_polyglot.Kind.Any);
      for (let i = 0; i < int64MapFieldEmbeddedSize; i++) {
        let key = decoder.int64();
        let val = EmptyModel.decode(decoder);
        if (typeof val !== "undefined") {
          this.int64MapFieldEmbedded.set(key, val);
        }
      }
      this.uint32Field = decoder.uint32();
      const uint32ArrayFieldSize = decoder.array(import_polyglot.Kind.Uint32);
      this.uint32ArrayField = new Array(uint32ArrayFieldSize);
      for (let i = 0; i < uint32ArrayFieldSize; i += 1) {
        this.uint32ArrayField[i] = decoder.uint32();
      }
      this.uint32MapField = /* @__PURE__ */ new Map();
      let uint32MapFieldSize = decoder.map(import_polyglot.Kind.Uint32, import_polyglot.Kind.Uint32);
      for (let i = 0; i < uint32MapFieldSize; i++) {
        let key = decoder.uint32();
        let val = decoder.uint32();
        this.uint32MapField.set(key, val);
      }
      this.uint32MapFieldEmbedded = /* @__PURE__ */ new Map();
      let uint32MapFieldEmbeddedSize = decoder.map(import_polyglot.Kind.Uint32, import_polyglot.Kind.Any);
      for (let i = 0; i < uint32MapFieldEmbeddedSize; i++) {
        let key = decoder.uint32();
        let val = EmptyModel.decode(decoder);
        if (typeof val !== "undefined") {
          this.uint32MapFieldEmbedded.set(key, val);
        }
      }
      this.uint64Field = decoder.uint64();
      const uint64ArrayFieldSize = decoder.array(import_polyglot.Kind.Uint64);
      this.uint64ArrayField = new Array(uint64ArrayFieldSize);
      for (let i = 0; i < uint64ArrayFieldSize; i += 1) {
        this.uint64ArrayField[i] = decoder.uint64();
      }
      this.uint64MapField = /* @__PURE__ */ new Map();
      let uint64MapFieldSize = decoder.map(import_polyglot.Kind.Uint64, import_polyglot.Kind.Uint64);
      for (let i = 0; i < uint64MapFieldSize; i++) {
        let key = decoder.uint64();
        let val = decoder.uint64();
        this.uint64MapField.set(key, val);
      }
      this.uint64MapFieldEmbedded = /* @__PURE__ */ new Map();
      let uint64MapFieldEmbeddedSize = decoder.map(import_polyglot.Kind.Uint64, import_polyglot.Kind.Any);
      for (let i = 0; i < uint64MapFieldEmbeddedSize; i++) {
        let key = decoder.uint64();
        let val = EmptyModel.decode(decoder);
        if (typeof val !== "undefined") {
          this.uint64MapFieldEmbedded.set(key, val);
        }
      }
      this.float32Field = decoder.float32();
      const float32ArrayFieldSize = decoder.array(import_polyglot.Kind.Float32);
      this.float32ArrayField = new Array(float32ArrayFieldSize);
      for (let i = 0; i < float32ArrayFieldSize; i += 1) {
        this.float32ArrayField[i] = decoder.float32();
      }
      this.float64Field = decoder.float64();
      const float64ArrayFieldSize = decoder.array(import_polyglot.Kind.Float64);
      this.float64ArrayField = new Array(float64ArrayFieldSize);
      for (let i = 0; i < float64ArrayFieldSize; i += 1) {
        this.float64ArrayField[i] = decoder.float64();
      }
      this.enumField = decoder.uint32();
      const enumArrayFieldSize = decoder.array(import_polyglot.Kind.Uint32);
      this.enumArrayField = new Array(enumArrayFieldSize);
      for (let i = 0; i < enumArrayFieldSize; i += 1) {
        this.enumArrayField[i] = decoder.uint32();
      }
      this.enumMapField = /* @__PURE__ */ new Map();
      let enumMapFieldSize = decoder.map(import_polyglot.Kind.Uint32, import_polyglot.Kind.String);
      for (let i = 0; i < enumMapFieldSize; i++) {
        let key = decoder.uint32();
        let val = decoder.string();
        this.enumMapField.set(key, val);
      }
      this.enumMapFieldEmbedded = /* @__PURE__ */ new Map();
      let enumMapFieldEmbeddedSize = decoder.map(import_polyglot.Kind.Uint32, import_polyglot.Kind.Any);
      for (let i = 0; i < enumMapFieldEmbeddedSize; i++) {
        let key = decoder.uint32();
        let val = EmptyModel.decode(decoder);
        if (typeof val !== "undefined") {
          this.enumMapFieldEmbedded.set(key, val);
        }
      }
      this.bytesField = decoder.uint8Array();
      const bytesArrayFieldSize = decoder.array(import_polyglot.Kind.Uint8Array);
      this.bytesArrayField = new Array(bytesArrayFieldSize);
      for (let i = 0; i < bytesArrayFieldSize; i += 1) {
        this.bytesArrayField[i] = decoder.uint8Array();
      }
      this.boolField = decoder.boolean();
      const boolArrayFieldSize = decoder.array(import_polyglot.Kind.Boolean);
      this.boolArrayField = new Array(boolArrayFieldSize);
      for (let i = 0; i < boolArrayFieldSize; i += 1) {
        this.boolArrayField[i] = decoder.boolean();
      }
    } else {
      this.modelField = new EmptyModel();
      this.modelArrayField = [];
      this.stringField = "DefaultValue";
      this.stringArrayField = [];
      this.stringMapField = /* @__PURE__ */ new Map();
      this.stringMapFieldEmbedded = /* @__PURE__ */ new Map();
      this.int32Field = 32;
      this.int32ArrayField = [];
      this.int32MapField = /* @__PURE__ */ new Map();
      this.int32MapFieldEmbedded = /* @__PURE__ */ new Map();
      this.int64Field = 64n;
      this.int64ArrayField = [];
      this.int64MapField = /* @__PURE__ */ new Map();
      this.int64MapFieldEmbedded = /* @__PURE__ */ new Map();
      this.uint32Field = 32;
      this.uint32ArrayField = [];
      this.uint32MapField = /* @__PURE__ */ new Map();
      this.uint32MapFieldEmbedded = /* @__PURE__ */ new Map();
      this.uint64Field = 64n;
      this.uint64ArrayField = [];
      this.uint64MapField = /* @__PURE__ */ new Map();
      this.uint64MapFieldEmbedded = /* @__PURE__ */ new Map();
      this.float32Field = 32.32;
      this.float32ArrayField = [];
      this.float64Field = 64.64;
      this.float64ArrayField = [];
      this.enumField = 2 /* DefaultValue */;
      this.enumArrayField = [];
      this.enumMapField = /* @__PURE__ */ new Map();
      this.enumMapFieldEmbedded = /* @__PURE__ */ new Map();
      this.bytesField = new Uint8Array(512);
      this.bytesArrayField = [];
      this.boolField = true;
      this.boolArrayField = [];
    }
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
    if (typeof this.modelField === "undefined") {
      encoder.null();
    } else {
      this.modelField.encode(encoder);
    }
    const modelArrayFieldLength = this.modelArrayField.length;
    encoder.array(modelArrayFieldLength, import_polyglot.Kind.Any);
    for (let i = 0; i < modelArrayFieldLength; i += 1) {
      const el = this.modelArrayField[i];
      el.encode(encoder);
    }
    encoder.string(this.stringField);
    const stringArrayFieldLength = this.stringArrayField.length;
    encoder.array(stringArrayFieldLength, import_polyglot.Kind.String);
    for (let i = 0; i < stringArrayFieldLength; i += 1) {
      encoder.string(this.stringArrayField[i]);
    }
    encoder.map(this.stringMapField.size, import_polyglot.Kind.String, import_polyglot.Kind.String);
    this.stringMapField.forEach((val, key) => {
      encoder.string(key);
      encoder.string(val);
    });
    encoder.map(this.stringMapFieldEmbedded.size, import_polyglot.Kind.String, import_polyglot.Kind.Any);
    this.stringMapFieldEmbedded.forEach((val, key) => {
      encoder.string(key);
      val.encode(encoder);
    });
    encoder.int32(this.int32Field);
    const int32ArrayFieldLength = this.int32ArrayField.length;
    encoder.array(int32ArrayFieldLength, import_polyglot.Kind.Int32);
    for (let i = 0; i < int32ArrayFieldLength; i += 1) {
      encoder.int32(this.int32ArrayField[i]);
    }
    encoder.map(this.int32MapField.size, import_polyglot.Kind.Int32, import_polyglot.Kind.Int32);
    this.int32MapField.forEach((val, key) => {
      encoder.int32(key);
      encoder.int32(val);
    });
    encoder.map(this.int32MapFieldEmbedded.size, import_polyglot.Kind.Int32, import_polyglot.Kind.Any);
    this.int32MapFieldEmbedded.forEach((val, key) => {
      encoder.int32(key);
      val.encode(encoder);
    });
    encoder.int64(this.int64Field);
    const int64ArrayFieldLength = this.int64ArrayField.length;
    encoder.array(int64ArrayFieldLength, import_polyglot.Kind.Int64);
    for (let i = 0; i < int64ArrayFieldLength; i += 1) {
      encoder.int64(this.int64ArrayField[i]);
    }
    encoder.map(this.int64MapField.size, import_polyglot.Kind.Int64, import_polyglot.Kind.Int64);
    this.int64MapField.forEach((val, key) => {
      encoder.int64(key);
      encoder.int64(val);
    });
    encoder.map(this.int64MapFieldEmbedded.size, import_polyglot.Kind.Int64, import_polyglot.Kind.Any);
    this.int64MapFieldEmbedded.forEach((val, key) => {
      encoder.int64(key);
      val.encode(encoder);
    });
    encoder.uint32(this.uint32Field);
    const uint32ArrayFieldLength = this.uint32ArrayField.length;
    encoder.array(uint32ArrayFieldLength, import_polyglot.Kind.Uint32);
    for (let i = 0; i < uint32ArrayFieldLength; i += 1) {
      encoder.uint32(this.uint32ArrayField[i]);
    }
    encoder.map(this.uint32MapField.size, import_polyglot.Kind.Uint32, import_polyglot.Kind.Uint32);
    this.uint32MapField.forEach((val, key) => {
      encoder.uint32(key);
      encoder.uint32(val);
    });
    encoder.map(this.uint32MapFieldEmbedded.size, import_polyglot.Kind.Uint32, import_polyglot.Kind.Any);
    this.uint32MapFieldEmbedded.forEach((val, key) => {
      encoder.uint32(key);
      val.encode(encoder);
    });
    encoder.uint64(this.uint64Field);
    const uint64ArrayFieldLength = this.uint64ArrayField.length;
    encoder.array(uint64ArrayFieldLength, import_polyglot.Kind.Uint64);
    for (let i = 0; i < uint64ArrayFieldLength; i += 1) {
      encoder.uint64(this.uint64ArrayField[i]);
    }
    encoder.map(this.uint64MapField.size, import_polyglot.Kind.Uint64, import_polyglot.Kind.Uint64);
    this.uint64MapField.forEach((val, key) => {
      encoder.uint64(key);
      encoder.uint64(val);
    });
    encoder.map(this.uint64MapFieldEmbedded.size, import_polyglot.Kind.Uint64, import_polyglot.Kind.Any);
    this.uint64MapFieldEmbedded.forEach((val, key) => {
      encoder.uint64(key);
      val.encode(encoder);
    });
    encoder.float32(this.float32Field);
    const float32ArrayFieldLength = this.float32ArrayField.length;
    encoder.array(float32ArrayFieldLength, import_polyglot.Kind.Float32);
    for (let i = 0; i < float32ArrayFieldLength; i += 1) {
      encoder.float32(this.float32ArrayField[i]);
    }
    encoder.float64(this.float64Field);
    const float64ArrayFieldLength = this.float64ArrayField.length;
    encoder.array(float64ArrayFieldLength, import_polyglot.Kind.Float64);
    for (let i = 0; i < float64ArrayFieldLength; i += 1) {
      encoder.float64(this.float64ArrayField[i]);
    }
    encoder.uint32(this.enumField);
    const enumArrayFieldLength = this.enumArrayField.length;
    encoder.array(enumArrayFieldLength, import_polyglot.Kind.Uint32);
    for (let i = 0; i < enumArrayFieldLength; i += 1) {
      encoder.uint32(this.enumArrayField[i]);
    }
    encoder.map(this.enumMapField.size, import_polyglot.Kind.Uint32, import_polyglot.Kind.String);
    this.enumMapField.forEach((val, key) => {
      encoder.uint32(key);
      encoder.string(val);
    });
    encoder.map(this.enumMapFieldEmbedded.size, import_polyglot.Kind.Uint32, import_polyglot.Kind.Any);
    this.enumMapFieldEmbedded.forEach((val, key) => {
      encoder.uint32(key);
      val.encode(encoder);
    });
    encoder.uint8Array(this.bytesField);
    const bytesArrayFieldLength = this.bytesArrayField.length;
    encoder.array(bytesArrayFieldLength, import_polyglot.Kind.Uint8Array);
    for (let i = 0; i < bytesArrayFieldLength; i += 1) {
      encoder.uint8Array(this.bytesArrayField[i]);
    }
    encoder.boolean(this.boolField);
    const boolArrayFieldLength = this.boolArrayField.length;
    encoder.array(boolArrayFieldLength, import_polyglot.Kind.Boolean);
    for (let i = 0; i < boolArrayFieldLength; i += 1) {
      encoder.boolean(this.boolArrayField[i]);
    }
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new ModelWithAllFieldTypes(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
//# sourceMappingURL=types.js.map