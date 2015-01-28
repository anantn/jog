package jog

// #include <stdlib.h>
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
func (j *Jog) GetInt(path ...string) (int, error) {
	pathPtr := convertPath(path)
	if pathPtr != nil {
		defer C.free(unsafe.Pointer(pathPtr.keys))
	}

	intval, err := C.GetInt(j.value, pathPtr)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Could not find int value at %s", strings.Join(path, "/")))
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
		return 0, errors.New(fmt.Sprintf("Could not find uint value at %s", strings.Join(path, "/")))
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
		return 0, errors.New(fmt.Sprintf("Could not find float value at %s", strings.Join(path, "/")))
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
		return false, errors.New(fmt.Sprintf("Could not find bool value at %s", strings.Join(path, "/")))
	}
	if int(bval) == 1 {
		return true, nil
	}
	return false, nil
}

func (j *Jog) GetString(path ...string) (string, error) {
	pathPtr := convertPath(path)
	if pathPtr != nil {
		defer C.free(unsafe.Pointer(pathPtr.keys))
	}

	strval := C.GetString(j.value, pathPtr)
	if strval == nil {
		return "", errors.New(fmt.Sprintf("Could not find string value at %s", strings.Join(path, "/")))
	}
	return C.GoString(strval), nil
}

func (j *Jog) GetObject(path ...string) (*Jog, error) {
	pathPtr := convertPath(path)
	if pathPtr != nil {
		defer C.free(unsafe.Pointer(pathPtr.keys))
	}

	objval := C.GetObject(j.value, pathPtr)
	if objval == nil {
		return nil, errors.New(fmt.Sprintf("Could not find object value at %s", strings.Join(path, "/")))
	}
	return &Jog{nil, objval}, nil
}

func (j *Jog) GetArray(path ...string) ([]*Jog, error) {
	pathPtr := convertPath(path)
	if pathPtr != nil {
		defer C.free(unsafe.Pointer(pathPtr.keys))
	}

	var arrlen C.size_t
	arrval := C.GetArray(j.value, pathPtr, &arrlen)
	length := int(arrlen)

	if arrval == nil {
		return []*Jog{}, errors.New(fmt.Sprintf("Could not find array value at %s", strings.Join(path, "/")))
	}
	if length == 0 {
		return []*Jog{}, nil
	}

	var void unsafe.Pointer
	ptrSize := unsafe.Sizeof(void)

	array := make([]*Jog, int(arrlen))
	for i := 0; i < length; i++ {
		ptr := (*unsafe.Pointer)(unsafe.Pointer(uintptr(unsafe.Pointer(arrval)) + uintptr(i)*ptrSize))
		array[i] = &Jog{unsafe.Pointer(arrval), *ptr}
	}
	runtime.SetFinalizer(&array, cleanupArray)
	return array, nil
}
