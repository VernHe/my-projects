package main

import "fmt"

type foo interface {
	fooFunc()
}
type foo1 struct{}

func (f1 foo1) fooFunc() {}
func main() {
	var f foo
	f1 := foo1{}
	f = foo(f1) // foo1{} escapes to heap
	fmt.Println(f)
	//f.fooFunc() // 调用方法时，f发生逃逸，因为方法是动态分配的
}
