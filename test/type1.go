package main

import (
	"fmt"
	"reflect"
)

type T1 []string
type T2 []string

func (t1 *T1) Append(t []string) {
	*t1 = t
}

func main() {
	foo0 := []string{}
	foo1 := T1{}
	foo2 := T2{}

	fmt.Println(reflect.TypeOf(foo0))
	fmt.Println(reflect.TypeOf(foo1))
	fmt.Println(reflect.TypeOf(foo2))

	foo1 = foo0
	foo0 = foo1

	//foo1 = (T1)foo2
	foo1.Append(foo2)
}
