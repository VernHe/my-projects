package main

import "fmt"

func main() {
	a()
}

func a() {
	defer fmt.Println("defer a")
	b()
	fmt.Println("a exec")
}

func b() {
	//defer fmt.Println("defer b")
	defer func() {
		if x := recover(); x != nil {
			fmt.Printf("run time panic: %v\n", x)
		}
	}()
	c()
	fmt.Println("b exec")
}

func c() {
	defer fmt.Println("defer c")
	panic("this is panic")
	fmt.Println("c exec")
}
