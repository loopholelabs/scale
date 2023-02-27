module.exports = function (buffer: Buffer) {
    let out = "const scale = require('@loopholelabs/scale');";
    out += "const scaleFile = require('@loopholelabs/scalefile');";
    out += "let buffer = new ArrayBuffer(" + buffer.length + ");";
    out += "let uint8Buffer = new Uint8Array(buffer);";
    out += "uint8Buffer.set([";
    for(let i = 0; i < buffer.length; i++) {
        out += buffer[i] + ","
    }
    out += "]);"
    out += "const sf = scaleFile.ScaleFunc.Decode(uint8Buffer);"
    out += "const mod = new WebAssembly.Module(sf.Function);"
    out += "const fn = new scale.Func(sf, mod);"
    out += "module.exports = fn;"
    // @ts-ignore
    this.callback(null, out);
}
module.exports.raw = true
