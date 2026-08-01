package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"bytes-go/bytesutil"
	"unsafe"
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
	res := bytesutil.Bytes(vStr, oStr)
	return C.CString(res)
}

//export Parse
func Parse(valJson *C.char) *C.char {
	vStr := ""
	if valJson != nil {
		vStr = C.GoString(valJson)
	}
	res := bytesutil.Parse(vStr)
	return C.CString(res)
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
	res := bytesutil.Format(vStr, oStr)
	return C.CString(res)
}

//export FreeString
func FreeString(str *C.char) {
	if str != nil {
		C.free(unsafe.Pointer(str))
	}
}

func main() {}
