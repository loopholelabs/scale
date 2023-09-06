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
var __reExport = (target, mod, secondTarget) => (__copyProps(target, mod, "default"), secondTarget && __copyProps(secondTarget, mod, "default"));
var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);
var stdin_exports = {};
__export(stdin_exports, {
  Error: () => Error2,
  Hash: () => Hash,
  Next: () => Next,
  Read: () => Read,
  Resize: () => Resize,
  Write: () => Write
});
module.exports = __toCommonJS(stdin_exports);
var import_scale_signature_interfaces = require("@loopholelabs/scale-signature-interfaces");
var import_polyglot = require("@loopholelabs/polyglot");
__reExport(stdin_exports, require("./types"), module.exports);
let WRITE_BUFFER = new Uint8Array().buffer;
let READ_BUFFER = new Uint8Array().buffer;
const hash = "3a592aa345d412faa2e6285ee048ca2ab5aa64b0caa2f9ca67b2c1e0792101e5";
function Write(ctx) {
  WRITE_BUFFER = ctx.encode(new Uint8Array()).buffer;
  const addrof = global[import_scale_signature_interfaces.Signature.TYPESCRIPT_ADDRESS_OF];
  const ptr = addrof(WRITE_BUFFER);
  const len = WRITE_BUFFER.byteLength;
  return [ptr, len];
}
function Read() {
  const dec = new import_polyglot.Decoder(new Uint8Array(READ_BUFFER));
  return ModelWithAllFieldTypes.decode(dec).value;
}
function Error2(err) {
  const enc = new import_polyglot.Encoder();
  enc.error(err);
  WRITE_BUFFER = enc.buffer;
  const addrof = global[interfaces.TYPESCRIPT_ADDRESS_OF];
  const ptr = addrof(WRITE_BUFFER);
  const len = WRITE_BUFFER.byteLength;
  return [ptr, len];
}
function Resize(size) {
  READ_BUFFER = new Uint8Array(size).buffer;
  const addrof = global[interfaces.TYPESCRIPT_ADDRESS_OF];
  return addrof(READ_BUFFER);
}
function Hash() {
  const enc = new import_polyglot.Encoder();
  enc.string(hash);
  WRITE_BUFFER = enc.buffer;
  const addrof = global[interfaces.TYPESCRIPT_ADDRESS_OF];
  const ptr = addrof(WRITE_BUFFER);
  const len = WRITE_BUFFER.byteLength;
  return [ptr, len];
}
function Next(ctx) {
  const [ptr, len] = Write(ctx);
  const next = global[interfaces.TYPESCRIPT_NEXT];
  next([ptr, len]);
  return Read();
}
//# sourceMappingURL=index.js.map