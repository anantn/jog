package yajl

// #include <string.h>
// #include "api/yajl_gen.h"
// #include "api/yajl_tree.h"
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"

	"github.com/anantn/jog"
)

type yajlValue struct {
	ptr *C.struct_yajl_val_s
}

const (
	yajl_t_string = iota + 1
	yajl_t_number
	yajl_t_object
	yajl_t_array
	yajl_t_true
	yajl_t_false
	yajl_t_null
	yajl_t_any
)

var ptr *C.char
var ptrSize = unsafe.Sizeof(ptr)
var ll C.longlong
var llSize = unsafe.Sizeof(ll)
var d C.double
var dSize = unsafe.Sizeof(d)

type yajlNumber struct {
	i     C.longlong
	d     C.double
	r     *C.char
	flags C.uint
}

func unionToNumber(u [32]byte) *yajlNumber {
	return &yajlNumber{
		*(*C.longlong)(unsafe.Pointer(&u[0:llSize][0])),
		*(*C.double)(unsafe.Pointer(&u[llSize : llSize+dSize][0])),
		*(**C.char)(unsafe.Pointer(&u[llSize+dSize : llSize+dSize+ptrSize][0])),
		*(*C.uint)(unsafe.Pointer(&u[llSize+dSize+ptrSize : llSize+dSize+ptrSize*2][0])),
	}
}

type yajlObject struct {
	keys   **C.char
	values **C.struct_yajl_val_s
	len    C.size_t
}

func unionToObject(u [32]byte) *yajlObject {
	return &yajlObject{
		*(***C.char)(unsafe.Pointer(&u[0:ptrSize][0])),
		*(***C.struct_yajl_val_s)(unsafe.Pointer(&u[ptrSize : ptrSize*2][0])),
		*(*C.size_t)(unsafe.Pointer(&u[ptrSize*2 : ptrSize*2+ptrSize][0])),
	}
}

type yajlArray struct {
	values **C.struct_yajl_val_s
	len    C.size_t
}

func unionToArray(u [32]byte) *yajlArray {
	return &yajlArray{
		*(***C.struct_yajl_val_s)(unsafe.Pointer(&u[0:ptrSize][0])),
		*(*C.size_t)(unsafe.Pointer(&u[ptrSize : ptrSize*2][0])),
	}
}

// Finalizer to call the yajl_tree_free.
func cleanupTree(j jog.Value) {
	yajlValue, ok := j.(*yajlValue)
	if !ok {
		panic("cleanupTree called on non yajlValue object!")
	}
	C.yajl_tree_free(yajlValue.ptr)
}

// Constructor by string.
func New(val string) (jog.Value, error) {
	cval := C.CString(val)
	yval := C.yajl_tree_parse(cval, nil, 0)
	if yval == nil {
		return nil, errors.New("Could not parse JSON!")
	}
	obj := &yajlValue{yval}
	runtime.SetFinalizer(obj, cleanupTree)
	return obj, nil
}

func (j *yajlValue) get(path ...string) (*C.struct_yajl_val_s, error) {
	if len(path) == 0 {
		return j.ptr, nil
	}
	n := j.ptr
	for _, part := range path {
		if int(n._type) != yajl_t_object {
			return nil, errors.New("Get called on a non-object value!")
		}
		i := 0
		obj := unionToObject(n.u)
		for ; i < int(obj.len); i++ {
			keyPtr := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(obj.keys)) + uintptr(i)*ptrSize))
			child := C.GoString(*keyPtr)
			if child == part {
				valPtr := (**C.struct_yajl_val_s)(unsafe.Pointer(uintptr(unsafe.Pointer(obj.values)) + uintptr(i)*ptrSize))
				n = *valPtr
				break
			}
		}
		if i == int(obj.len) {
			return nil, errors.New("Could not find child at path!")
		}
	}
	if n == nil {
		return nil, errors.New("Could not find valid value at path!")
	}
	return n, nil
}

func (j *yajlValue) Get(path ...string) (jog.Value, error) {
	n, err := j.get(path...)
	if err != nil {
		return nil, err
	}
	return &yajlValue{n}, nil
}

func (j *yajlValue) getNumber(path ...string) (*yajlNumber, error) {
	n := j.ptr
	var err error
	if len(path) != 0 {
		n, err = j.get(path...)
		if err != nil {
			return nil, err
		}
	}
	if int(n._type) != yajl_t_number {
		return nil, errors.New("GetInt called on a non-number value!")
	}
	obj := unionToNumber(n.u)
	return obj, nil
}

func (j *yajlValue) GetInt(path ...string) (int, error) {
	obj, err := j.getNumber(path...)
	if err != nil {
		return 0, err
	}
	return int(obj.i), nil
}

func (j *yajlValue) GetUInt(path ...string) (uint, error) {
	obj, err := j.getNumber(path...)
	if err != nil {
		return 0, err
	}
	return uint(obj.i), nil
}

func (j *yajlValue) GetFloat(path ...string) (float64, error) {
	obj, err := j.getNumber(path...)
	if err != nil {
		return 0, err
	}
	return float64(obj.d), nil
}

func (j *yajlValue) GetBool(path ...string) (bool, error) {
	n := j.ptr
	var err error
	if len(path) != 0 {
		n, err = j.get(path...)
		if err != nil {
			return false, err
		}
	}
	if int(n._type) == yajl_t_true {
		return true, nil
	}
	if int(n._type) == yajl_t_false {
		return false, nil
	}
	return false, errors.New("GetBool called on non-boolean type!")
}

func (j *yajlValue) GetString(path ...string) (string, error) {
	n := j.ptr
	var err error
	if len(path) != 0 {
		n, err = j.get(path...)
		if err != nil {
			return "", err
		}
	}
	if int(n._type) != yajl_t_string {
		return "", errors.New("GetString called on a non-string value!")
	}
	str := (**C.char)(unsafe.Pointer(&n.u))
	return C.GoString(*str), nil
}

func (j *yajlValue) GetArray(path ...string) ([]jog.Value, error) {
	n := j.ptr
	var err error
	if len(path) != 0 {
		n, err = j.get(path...)
		if err != nil {
			return nil, err
		}
	}
	if int(n._type) != yajl_t_array {
		return nil, errors.New("GetArray called on a non-array value!")
	}
	obj := unionToArray(n.u)
	l := int(obj.len)
	arr := make([]jog.Value, l)
	for i := 0; i < l; i++ {
		valPtr := (**C.struct_yajl_val_s)(unsafe.Pointer(uintptr(unsafe.Pointer(obj.values)) + uintptr(i)*ptrSize))
		arr[i] = &yajlValue{*valPtr}
	}
	return arr, nil
}

func (j *yajlValue) GetObject(path ...string) (map[string]jog.Value, error) {
	n := j.ptr
	var err error
	if len(path) != 0 {
		n, err = j.get(path...)
		if err != nil {
			return nil, err
		}
	}
	if int(n._type) != yajl_t_object {
		return nil, errors.New("GetObject called on a non-object value!")
	}
	obj := unionToObject(n.u)
	l := int(obj.len)
	bag := make(map[string]jog.Value, l)
	keySize := unsafe.Sizeof(obj.keys)
	valueSize := unsafe.Sizeof(obj.values)
	for i := 0; i < l; i++ {
		keyPtr := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(obj.keys)) + uintptr(i)*keySize))
		valPtr := (**C.struct_yajl_val_s)(unsafe.Pointer(uintptr(unsafe.Pointer(obj.values)) + uintptr(i)*valueSize))
		bag[C.GoString(*keyPtr)] = &yajlValue{*valPtr}
	}
	return bag, nil
}

func (j *yajlValue) Type(path ...string) jog.Type {
	n, err := j.get(path...)
	if err != nil {
		return jog.TypeUnknown
	}
	switch int(n._type) {
	case yajl_t_string:
		return jog.TypeString
	case yajl_t_number:
		return jog.TypeNumber
	case yajl_t_object:
		return jog.TypeObject
	case yajl_t_array:
		return jog.TypeArray
	case yajl_t_true:
		fallthrough
	case yajl_t_false:
		return jog.TypeBool
	case yajl_t_null:
		return jog.TypeNull
	}
	return jog.TypeUnknown
}

func (j *yajlValue) Stringify(path ...string) (string, error) {
	n, err := j.get(path...)
	if err != nil {
		return "", err
	}
	h := C.yajl_gen_alloc(nil)
	defer C.yajl_gen_free(h)

	err = toString(n, h)
	if err != nil {
		return "", err
	}

	var buf *C.uchar
	var length C.size_t
	if int(C.yajl_gen_get_buf(h, &buf, &length)) != 0 {
		return "", errors.New("Could not get final encoded buffer!")
	}

	return C.GoString((*C.char)(unsafe.Pointer(buf))), nil
}

func toString(n *C.struct_yajl_val_s, h *C.struct_yajl_gen_t) error {
	switch int(n._type) {
	case yajl_t_string:
		str := (**C.char)(unsafe.Pointer(&n.u))
		length := C.strlen(*str)
		if int(C.yajl_gen_string(h, *(**C.uchar)(unsafe.Pointer(str)), length)) != 0 {
			errors.New("Could not encode string!")
		}
	case yajl_t_number:
		num := unionToNumber(n.u)
		if int(C.yajl_gen_number(h, num.r, C.strlen(num.r))) != 0 {
			errors.New("Could not encode number!")
		}
	case yajl_t_object:
		obj := unionToObject(n.u)
		if int(C.yajl_gen_map_open(h)) != 0 {
			errors.New("Could not start encoding object!")
		}
		for i := 0; i < int(obj.len); i++ {
			keyPtr := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(obj.keys)) + uintptr(i)*ptrSize))
			if int(C.yajl_gen_string(h, *(**C.uchar)(unsafe.Pointer(keyPtr)), C.strlen(*keyPtr))) != 0 {
				return errors.New("Could not encode map key!")
			}
			valPtr := (**C.struct_yajl_val_s)(unsafe.Pointer(uintptr(unsafe.Pointer(obj.values)) + uintptr(i)*ptrSize))
			err := toString(*valPtr, h)
			if err != nil {
				return err
			}
		}
		if int(C.yajl_gen_map_close(h)) != 0 {
			errors.New("Could not end encoding object!")
		}
	case yajl_t_array:
		arr := unionToArray(n.u)
		if int(C.yajl_gen_array_open(h)) != 0 {
			errors.New("Could not start encoding array!")
		}
		for i := 0; i < int(arr.len); i++ {
			valPtr := (**C.struct_yajl_val_s)(unsafe.Pointer(uintptr(unsafe.Pointer(arr.values)) + uintptr(i)*ptrSize))
			err := toString(*valPtr, h)
			if err != nil {
				return err
			}
		}
		if int(C.yajl_gen_array_close(h)) != 0 {
			errors.New("Could not end encoding array!")
		}
	case yajl_t_true:
		if int(C.yajl_gen_bool(h, 1)) != 0 {
			errors.New("Could not encode true!")
		}
	case yajl_t_false:
		if int(C.yajl_gen_bool(h, 0)) != 0 {
			errors.New("Could not encode false!")
		}
	case yajl_t_null:
		if int(C.yajl_gen_null(h)) != 0 {
			errors.New("Could not encode null!")
		}
	default:
		return fmt.Errorf("Could not encode unknown value type! %d", int(n._type))
	}
	return nil
}
