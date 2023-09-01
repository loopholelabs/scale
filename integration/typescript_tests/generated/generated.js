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
  Context: () => Context
});
module.exports = __toCommonJS(stdin_exports);
class Context {
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
      if (err !== void 0) {
        throw err;
      }
      this.a = decoder.int32();
      this.b = decoder.int32();
      this.c = decoder.int32();
    } else {
      this.a = 0;
      this.b = 0;
      this.c = 0;
    }
  }
  /**
  * @throws {Error}
  */
  encode(encoder) {
    encoder.int32(this.a);
    encoder.int32(this.b);
    encoder.int32(this.c);
  }
  /**
  * @throws {Error}
  */
  static decode(decoder) {
    if (decoder.null()) {
      return void 0;
    }
    return new Context(decoder);
  }
  /**
  * @throws {Error}
  */
  static encode_undefined(encoder) {
    encoder.null();
  }
}
//# sourceMappingURL=generated.js.map