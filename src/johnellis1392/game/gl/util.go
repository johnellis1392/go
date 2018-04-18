package gl

import (
	"reflect"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
)

// #include <stdlib.h>
import "C"

// //////////////////// //
// Conversion Functions //
// //////////////////// //

// Str converts a native-Golang string into a GL-compatible string
// (an unsigned byte array)
func Str(str string) *uint8 {
	if !strings.HasSuffix(str, "\x00") {
		panic("str argument missing null terminator: " + str)
	}
	header := (*reflect.StringHeader)(unsafe.Pointer(&str))
	return (*uint8)(unsafe.Pointer(header.Data))
}

// GoString converts the given byte pointer into a native-Golang string object.
func GoString(cstr *uint8) string {
	return C.GoString((*C.char)(unsafe.Pointer(cstr)))
}

// GetString gets a GL string using the given GL Object Pointer.
func GetString(n Enum) string {
	return GoString(gl.GetString(n.Unbox()))
}

// CString is a GL-Compatible string implementation that converts a Golang-Native
// string object into a *uint8 / byte array.
type CString struct {
	val *uint8
	// free func()
}

// AsCString converts a Native-Golang string into a GL-Compatible string.
func AsCString(s string) CString {
	header := (*reflect.StringHeader)(unsafe.Pointer(&s))
	val := (*uint8)(unsafe.Pointer(header.Data))
	return CString{
		val: val,
		// free: free,
	}
}

// Free frees the memory associated with the underlying string object pointed
// to by c.val
// func (c CString) Free() {
// 	if c.free != nil {
// 		c.free()
// 	}
// }

// String converts the given CString object into a Native Golang string.
func (c CString) String() string {
	cstr := unsafe.Pointer(c.val)
	return C.GoString((*C.char)(cstr))
}
