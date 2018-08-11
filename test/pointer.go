package main

import (
	//"net"
	"fmt"
	"reflect"
	"unsafe"
)

func SizeStruct(data interface{}) int {
	return sizeof(reflect.ValueOf(data))
}

func sizeof(v reflect.Value) int {
	switch v.Kind(){
	case reflect.Map:
		sum :=0
		keys:= v.MapKeys()
		for i:=0;i< len(keys); i++ {
			mapKey := keys[i]
			s := sizeof(mapKey)
			if s < 0 {
				return -1
			}
			sum += s
			s = sizeof(v.MapIndex(mapKey))
			if s < 0{
				return -1
			}
			sum += s
		}
		return sum
	case reflect.Slice, reflect.Array:
		sum := 0
		for i, n := 0, v.Len(); i< n;i++{
			s := sizeof(v.Index(i))
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum
	case reflect.String:
		sum := 0
		for i,n:=0,v.Len();i<n;i++{
			s := sizeof(v.Index(i))
			if s < 0{
				return -1
			}
			sum += s
		}
		return sum
	case reflect.Ptr, reflect.Interface:
		p:= (*[]byte)(unsafe.Pointer(v.Pointer()))
		if p == nil {
			return 0
		}
		return sizeof(v.Elem())
	case reflect.Struct:
		sum :=0
		for i,n :=0,v.NumField(); i<n;i++{
			s:=sizeof(v.Field(i))
			if s< 0{
				return -1
			}
			sum += s
		}
		return sum
	case reflect.Uint8,reflect.Uint16,reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.Int:
		return int(v.Type().Size())
	default:
		fmt.Println("t.Kind() not found:", v.Kind())
	}
	return -1

}

type X uint8
type Y struct {
	T uint8
	TT *uint16
	K []string
}
func main(){
	var x  X
	y := &Y{
		//T:1,
		//TT: new(uint16),
		K: []string{"a"},
	}
	fmt.Println(SizeStruct(x))
	fmt.Println(SizeStruct(y))
	fmt.Println(SizeStruct(y))
}