package rapid

// #include <stdlib.h>
// #include <stdbool.h>
// #include "rapid.h"
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"unsafe"

	"github.com/anantn/jog"
)

type rapidValue struct {
	clean unsafe.Pointer
	value unsafe.Pointer
}

// Constructor by string.
func New(val string) (jog.Value, error) {
	var cerr *C.char
	cval := C.CString(val)
	doc := C.NewDocument(cval, &cerr)
	if doc == nil {
		msg := C.GoString(cerr)
		C.free(unsafe.Pointer(cval))
		C.free(unsafe.Pointer(cerr))
		return nil, errors.New(msg)
	}

	obj := &rapidValue{unsafe.Pointer(cval), doc}
	runtime.SetFinalizer(obj, cleanupDocument)
	return obj, nil
}

// Data Getters.
func (j *rapidValue) Get(path ...string) (jog.Value, error) {
	if len(path) == 0 {
		return j, nil
	}

	pathPtr := convertPath(path)
	if pathPtr != nil {
		defer C.free(unsafe.Pointer(pathPtr.keys))
	}
	childval := C.Get(j.value, pathPtr)
	if childval == nil {
		return nil, fmt.Errorf("Could not find a child at %s", strings.Join(path, "/"))
	}
	return &rapidValue{nil, childval}, nil
}

func (j *rapidValue) GetInt(path ...string) (int, error) {
	pathPtr := convertPath(path)
	if pathPtr != nil {
		defer C.free(unsafe.Pointer(pathPtr.keys))
	}

	intval, err := C.GetInt(j.value, pathPtr)
	if err != nil {
		return 0, fmt.Errorf("Could not find int value at %s", strings.Join(path, "/"))
	}
	return int(intval), nil
}

func (j *rapidValue) GetUInt(path ...string) (uint, error) {
	pathPtr := convertPath(path)
	if pathPtr != nil {
		defer C.free(unsafe.Pointer(pathPtr.keys))
	}

	uintval, err := C.GetUInt(j.value, pathPtr)
	if err != nil {
		return 0, fmt.Errorf("Could not find uint value at %s", strings.Join(path, "/"))
	}
	return uint(uintval), nil
}

func (j *rapidValue) GetFloat(path ...string) (float64, error) {
	pathPtr := convertPath(path)
	if pathPtr != nil {
		defer C.free(unsafe.Pointer(pathPtr.keys))
	}

	dval, err := C.GetDouble(j.value, pathPtr)
	if err != nil {
		return 0, fmt.Errorf("Could not find float value at %s", strings.Join(path, "/"))
	}
	return float64(dval), nil
}

func (j *rapidValue) GetBool(path ...string) (bool, error) {
	pathPtr := convertPath(path)
	if pathPtr != nil {
		defer C.free(unsafe.Pointer(pathPtr.keys))
	}

	bval, err := C.GetBool(j.value, pathPtr)
	if err != nil {
		return false, fmt.Errorf("Could not find bool value at %s", strings.Join(path, "/"))
	}
	return bool(bval), nil
}

func (j *rapidValue) GetString(path ...string) (string, error) {
	pathPtr := convertPath(path)
	if pathPtr != nil {
		defer C.free(unsafe.Pointer(pathPtr.keys))
	}

	strval := C.GetString(j.value, pathPtr)
	if strval == nil {
		return "", fmt.Errorf("Could not find string value at %s", strings.Join(path, "/"))
	}
	return C.GoString(strval), nil
}

func (j *rapidValue) GetArray(path ...string) ([]jog.Value, error) {
	pathPtr := convertPath(path)
	if pathPtr != nil {
		defer C.free(unsafe.Pointer(pathPtr.keys))
	}

	var arrlen C.size_t
	arrval := C.GetArray(j.value, pathPtr, &arrlen)

	if arrval == nil {
		return []jog.Value{}, fmt.Errorf("Could not find array value at %s", strings.Join(path, "/"))
	}

	length := int(arrlen)
	if length == 0 {
		return []jog.Value{}, nil
	}

	var void unsafe.Pointer
	ptrSize := unsafe.Sizeof(void)

	array := make([]jog.Value, length)
	for i := 0; i < length; i++ {
		ptr := (*unsafe.Pointer)(unsafe.Pointer(uintptr(unsafe.Pointer(arrval)) + uintptr(i)*ptrSize))
		array[i] = &rapidValue{unsafe.Pointer(arrval), *ptr}
	}
	runtime.SetFinalizer(&array, cleanupArray)
	return array, nil
}

func (j *rapidValue) GetObject(path ...string) (map[string]jog.Value, error) {
	pathPtr := convertPath(path)
	if pathPtr != nil {
		defer C.free(unsafe.Pointer(pathPtr.keys))
	}

	var keys **C.char
	var memlen C.size_t
	objval := C.GetObject(j.value, pathPtr, &memlen, &keys)
	if objval == nil {
		return nil, fmt.Errorf("Could not find object value at %s", strings.Join(path, "/"))
	}

	length := int(memlen)
	if length == 0 {
		return map[string]jog.Value{}, nil
	}

	var void unsafe.Pointer
	ptrSize := unsafe.Sizeof(void)
	charSize := unsafe.Sizeof(keys)
	members := make(map[string]jog.Value, length)
	for i := 0; i < length; i++ {
		ptr := (*unsafe.Pointer)(unsafe.Pointer(uintptr(unsafe.Pointer(objval)) + uintptr(i)*ptrSize))
		keyPtr := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(keys)) + uintptr(i)*charSize))
		keyVal := C.GoString(*keyPtr)
		members[keyVal] = &rapidValue{nil, *ptr}
	}
	C.free(unsafe.Pointer(objval))
	C.free(unsafe.Pointer(keys))
	return members, nil
}

func (j *rapidValue) Type(path ...string) jog.Type {
	pathPtr := convertPath(path)
	if pathPtr != nil {
		defer C.free(unsafe.Pointer(pathPtr.keys))
	}

	ct := C.Type(j.value, pathPtr)
	switch C.GoString(ct) {
	case "BOOL":
		return jog.TypeBool
	case "NULL":
		return jog.TypeNull
	case "ARRAY":
		return jog.TypeArray
	case "STRING":
		return jog.TypeString
	case "OBJECT":
		return jog.TypeObject
	case "NUMBER":
		return jog.TypeNumber
	}

	return jog.TypeUnknown
}

func (j *rapidValue) Stringify(path ...string) (string, error) {
	pathPtr := convertPath(path)
	if pathPtr != nil {
		defer C.free(unsafe.Pointer(pathPtr.keys))
	}

	strval := C.Stringify(j.value, pathPtr)
	if strval == nil {
		return "", fmt.Errorf("Could not stringify value because it's not an object or array (%s)", strings.Join(path, "/"))
	}
	ret := C.GoString(strval)
	C.free(unsafe.Pointer(strval))
	return ret, nil
}

// Private methods.

// Finalizer to call the Document destructor.
func cleanupDocument(j jog.Value) {
	rapidValue, ok := j.(*rapidValue)
	if !ok {
		panic("cleanupDocument called on non rapidValue object!")
	}
	C.free(rapidValue.clean)
	C.DeleteDocument(rapidValue.value)
}

// Finalizer to call the Array destructor.
func cleanupArray(a *[]jog.Value) {
	arr := *a
	rapidValue, ok := arr[0].(*rapidValue)
	if !ok {
		panic("cleanupArray called on non rapidValue object!")
	}
	C.free(rapidValue.clean)
}

// Utility function to convert a Go slice to C struct.
// Caller must free C.struct_Path.keys!
func convertPath(path []string) *C.struct_Path {
	if len(path) == 0 {
		return nil
	}

	var c *C.char
	ptrSize := unsafe.Sizeof(c)
	ptr := C.malloc(C.size_t(len(path)) * C.size_t(ptrSize))
	for i := 0; i < len(path); i++ {
		element := (**C.char)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*ptrSize))
		*element = C.CString(path[i])
	}
	return &C.struct_Path{(**C.char)(ptr), C.size_t(len(path))}
}
