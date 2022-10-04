package common

import "os"
import "strings"
import "log"
import "unsafe"
import "reflect"
import "fmt"

func GetProjectDir() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get project working directory")
	}
	paths := strings.Split(wd, "/")
	i := 0
	j := 0
	for i = 0; i < len(paths); i++{
		cur := paths[i]
		if cur == "dieselfluid"{
			if i < len(paths)-1{
				j = i+1
				break
			}else{
				j = i
				break
			}
		}
	}
	return strings.Join(paths[0:j], "/")
}

func ProjectRelativePath(relative_path string) string {
	return GetProjectDir() + "/" + relative_path
}

//Linux Cd command
func Cd(path string)string{
	paths := strings.Split(path, "/")
	return strings.Join(paths[0:len(paths)-1], "/")
}

// Ptr takes a slice or pointer (to a singular scalar value or the first
// element of an array or slice) and returns its GL-compatible address.
//
// For example:
//
// 	var data []uint8
// 	...
// 	gl.TexImage2D(gl.TEXTURE_2D, ..., gl.UNSIGNED_BYTE, gl.Ptr(&data[0]))
func Ptr(data interface{}) unsafe.Pointer {
	if data == nil {
		return unsafe.Pointer(nil)
	}
	var addr unsafe.Pointer
	v := reflect.ValueOf(data)
	switch v.Type().Kind() {
	case reflect.Ptr:
		e := v.Elem()
		switch e.Kind() {
		case
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			addr = unsafe.Pointer(e.UnsafeAddr())
		default:
			panic(fmt.Errorf("unsupported pointer to type %s; must be a slice or pointer to a singular scalar value or the first element of an array or slice", e.Kind()))
		}
	case reflect.Uintptr:
		addr = unsafe.Pointer(data.(uintptr))
	case reflect.Slice:
		addr = unsafe.Pointer(v.Index(0).UnsafeAddr())
	default:
		panic(fmt.Errorf("unsupported type %s; must be a slice or pointer to a singular scalar value or the first element of an array or slice", v.Type()))
	}
	return addr
}
