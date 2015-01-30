package jog

// #include <stdlib.h>
// #include <stdbool.h>
// #include "wrapper.h"
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"unsafe"
)

type Jog struct {
	clean unsafe.Pointer
	value unsafe.Pointer
}

type Type int

const (
	TypeBool Type = iota
	TypeNull
	TypeArray
	TypeNumber
	TypeString
	TypeObject
	TypeUnknown
)

// Finalizer to call the Document destructor.
func cleanupDocument(j *Jog) {
	C.free(j.clean)
	C.DeleteDocument(j.value)
}

// Finalizer to call the Array destructor.
func cleanupArray(a *[]*Jog) {
	arr := *a
	C.free(arr[0].clean)
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

// Constructor by string.
func New(val string) (*Jog, error) {
	var cerr *C.char
	cval := C.CString(val)
	doc := C.NewDocument(cval, &cerr)
	if doc == nil {
		msg := C.GoString(cerr)
		C.free(unsafe.Pointer(cval))
		C.free(unsafe.Pointer(cerr))
		return nil, errors.New(msg)
	}

	obj := &Jog{unsafe.Pointer(cval), doc}
	runtime.SetFinalizer(obj, cleanupDocument)
	return obj, nil
}

// Data Getters.
func (j *Jog) Get(path ...string) (*Jog, error) {
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
	return &Jog{nil, childval}, nil
}

func (j *Jog) GetInt(path ...string) (int, error) {
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

func (j *Jog) GetUInt(path ...string) (uint, error) {
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

func (j *Jog) GetFloat(path ...string) (float64, error) {
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

func (j *Jog) GetBool(path ...string) (bool, error) {
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

func (j *Jog) GetString(path ...string) (string, error) {
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

func (j *Jog) GetArray(path ...string) ([]*Jog, error) {
	pathPtr := convertPath(path)
	if pathPtr != nil {
		defer C.free(unsafe.Pointer(pathPtr.keys))
	}

	var arrlen C.size_t
	arrval := C.GetArray(j.value, pathPtr, &arrlen)

	if arrval == nil {
		return []*Jog{}, fmt.Errorf("Could not find array value at %s", strings.Join(path, "/"))
	}

	length := int(arrlen)
	if length == 0 {
		return []*Jog{}, nil
	}

	var void unsafe.Pointer
	ptrSize := unsafe.Sizeof(void)

	array := make([]*Jog, length)
	for i := 0; i < length; i++ {
		ptr := (*unsafe.Pointer)(unsafe.Pointer(uintptr(unsafe.Pointer(arrval)) + uintptr(i)*ptrSize))
		array[i] = &Jog{unsafe.Pointer(arrval), *ptr}
	}
	runtime.SetFinalizer(&array, cleanupArray)
	return array, nil
}

func (j *Jog) GetObject(path ...string) (map[string]*Jog, error) {
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
		return map[string]*Jog{}, nil
	}

	var void unsafe.Pointer
	ptrSize := unsafe.Sizeof(void)
	charSize := unsafe.Sizeof(keys)
	members := make(map[string]*Jog, length)
	for i := 0; i < length; i++ {
		ptr := (*unsafe.Pointer)(unsafe.Pointer(uintptr(unsafe.Pointer(objval)) + uintptr(i)*ptrSize))
		keyPtr := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(keys)) + uintptr(i)*charSize))
		keyVal := C.GoString(*keyPtr)
		members[keyVal] = &Jog{nil, *ptr}
	}
	C.free(unsafe.Pointer(objval))
	C.free(unsafe.Pointer(keys))
	return members, nil
}

func (j *Jog) Type(path ...string) Type {
	pathPtr := convertPath(path)
	if pathPtr != nil {
		defer C.free(unsafe.Pointer(pathPtr.keys))
	}

	ct := C.Type(j.value, pathPtr)
	switch C.GoString(ct) {
	case "BOOL":
		return TypeBool
	case "NULL":
		return TypeNull
	case "ARRAY":
		return TypeArray
	case "STRING":
		return TypeString
	case "OBJECT":
		return TypeObject
	case "NUMBER":
		return TypeNumber
	}

	return TypeUnknown
}

func (j *Jog) Stringify(path ...string) (string, error) {
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
