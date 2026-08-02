const koffi = require("koffi");
const path = require("path");

const libPath = process.platform === "win32"
    ? path.join(__dirname, "./libbytes.dll")
    : path.join(__dirname, "./libbytes.so");

const lib = koffi.load(libPath);

const nativeBytes = lib.func("const char *Bytes(const char *valJson, const char *optsJson)");
const nativeParse = lib.func("const char *Parse(const char *valJson)");
const nativeFormat = lib.func("const char *Format(const char *valJson, const char *optsJson)");

function toRawJson(val) {
    return val === undefined ? "undefined" : JSON.stringify(val);
}

function callNative(func, ...args) {
    const str = func(...args);
    return str ? JSON.parse(str) : null;
}

function bytes(val, options) {
    return callNative(nativeBytes, toRawJson(val), toRawJson(options));
}

function parse(val) {
    return callNative(nativeParse, toRawJson(val));
}

function format(val, options) {
    return callNative(nativeFormat, toRawJson(val), toRawJson(options));
}

bytes.parse = parse;
bytes.format = format;

module.exports = bytes;
module.exports.format = format;
module.exports.parse = parse;
