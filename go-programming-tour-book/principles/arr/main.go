package main

import "fmt"

func main() {
	testArrayAddress()
	a := [3]int{1, 2, 3}
	fmt.Printf("a array address is : %p\n", &a)
	// 作为参数传递
	testArrayCopy(a)
	// 赋值
	c := a
	fmt.Printf("c array address is : %p\n", &c)
}

func testDeclare() {
	var arr1 [3]int
	arr2 := [3]int{1, 2, 3}
	arr3 := [...]int{1, 2, 3, 4}

	fmt.Printf("arr1: %T, arr2: %T, arr3: %T\n", arr1, arr2, arr3)

	fmt.Println(arr1)
	fmt.Println(arr2)
	fmt.Println(arr3)
}

func testArrayCopy(b [3]int) {
	fmt.Printf("b array address is : %p\n", &b)
}

func testArrayAddress() {
	arr1 := [3]int{1, 2, 3}
	arr2 := [3]int{1, 2, 3}
	arr3 := [5]int{1, 2, 3}
	arr4 := [3]int{1, 2, 3}
	fmt.Printf("arr1: %p, arr2: %p, arr3: %p, arr4: %p\n", &arr1, &arr2, &arr3, &arr4)
}
