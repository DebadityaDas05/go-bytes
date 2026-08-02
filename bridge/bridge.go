package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"encoding/json"
	"unsafe"

	bytesutil "github.com/DebadityaDas05/go-bytes/src"
)

//export Bytes
func Bytes(valJson *C.char, optsJson *C.char) *C.char {
	vStr := ""
	if valJson != nil {
		vStr = C.GoString(valJson)
	}
	oStr := ""
	if optsJson != nil {
		oStr = C.GoString(optsJson)
	}

	val := parseJSONVal(vStr)
	opts := parseJSONOpts(oStr)

	res := bytesutil.Bytes(val, opts)
	return C.CString(returnJSON(res))
}

//export Parse
func Parse(valJson *C.char) *C.char {
	vStr := ""
	if valJson != nil {
		vStr = C.GoString(valJson)
	}

	val := parseJSONVal(vStr)
	res := bytesutil.Parse(val)
	return C.CString(returnJSON(res))
}

//export Format
func Format(valJson *C.char, optsJson *C.char) *C.char {
	vStr := ""
	if valJson != nil {
		vStr = C.GoString(valJson)
	}
	oStr := ""
	if optsJson != nil {
		oStr = C.GoString(optsJson)
	}

	val := parseJSONVal(vStr)
	num, ok := val.(float64)
	if !ok {
		return C.CString(returnJSON(nil))
	}
	opts := parseJSONOpts(oStr)

	res := bytesutil.Format(num, opts)
	return C.CString(returnJSON(res))
}

//export FreeString
func FreeString(str *C.char) {
	if str != nil {
		C.free(unsafe.Pointer(str))
	}
}

func parseJSONVal(valJson string) interface{} {
	if valJson == "" || valJson == "undefined" {
		return nil
	}
	var val interface{}
	if err := json.Unmarshal([]byte(valJson), &val); err != nil {
		return nil
	}
	return val
}

func parseJSONOpts(optsJson string) *bytesutil.Options {
	if optsJson == "" || optsJson == "undefined" || optsJson == "null" {
		return nil
	}
	var opts bytesutil.Options
	if err := json.Unmarshal([]byte(optsJson), &opts); err != nil {
		return nil
	}
	return &opts
}

func returnJSON(val interface{}) string {
	if val == nil {
		return "null"
	}
	b, err := json.Marshal(val)
	if err != nil {
		return "null"
	}
	return string(b)
}

func main() {}
