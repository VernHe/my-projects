package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	//var s1 []int                // 长度和容量默认都是0
	//var s2 = []int{1, 2, 3}     // 长度和容量都为元素个数
	//var s3 = make([]int, 5)     // 长度和容量都为5
	//var s4 = make([]int, 5, 10) // 长度为5，容量为10
	//// s1: []int, s2: []int
	//fmt.Printf("【切片类型】 s1: %T, s2: %T, s3: %T, s4: %T\n", s1, s2, s3, s4)
	//fmt.Printf("【元素数量】 s1: %d, s2: %d, s3: %d, s4: %d\n", len(s1), len(s2), len(s3), len(s4))
	//fmt.Printf("【切片容量】 s1: %d, s2: %d, s3: %d, s4: %d\n", cap(s1), cap(s2), cap(s3), cap(s4))
	//testCutOut()
	//testSliceCopy()
	//testExpansion()

	testStruct()
}

func testCutOut() {
	nums1 := []int{1, 2, 3, 4, 5, 6, 7}
	nums2 := nums1[2:4] // 截取区间 [2,4)
	fmt.Printf("nums2: %v\n", nums2)
	// 修改nums2
	nums2[0] = 0
	nums2[1] = 0
	// 输出 nums1: [1 2 0 0 5 6 7]
	fmt.Printf("nums1: %v\n", nums1)
}

func testSliceCopy() {
	nums1 := []int{1, 2, 3, 4, 5, 6, 7}
	nums2 := nums1
	nums2[0] = 100
	fmt.Printf("nums1: %v\n", nums1)
}

func testExpansion() {
	//// 扩容后容量小于 < 1024
	//nums1 := []int{1, 2, 3, 4}
	//// 容量前后变化 4 -> 5
	//nums1 = append(nums1, 5)
	//// 新容量: 8 = 4 * 2 （增加1倍）
	//fmt.Printf("new cap: %d\n", cap(nums1)) // 输出 new cap: 8
	//
	//// 扩容后容量 > 1024
	//nums2 := make([]int, 1024, 1024)
	//// 容量前后变化 1024 -> 1025
	//nums2 = append(nums2, 1)
	//// 新容量: 1280 = 1024 + 256 （增加1/4）
	//fmt.Printf("new cap: %d\n", cap(nums2)) // 输出 new cap: 1280

	// 扩容后容量 > 扩容前容量的2倍
	nums3 := []int{1, 2, 3, 4}
	// 容量前后变化 5 -> 1025
	nums3 = append(nums3, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
	fmt.Printf("new cap: %d\n", cap(nums3))
	//nums3 = append(nums3, 1)
	//fmt.Printf("new cap: %d\n", cap(nums3))
	//nums3 = append(nums3, 1)
	//fmt.Printf("new cap: %d\n", cap(nums3))
	//nums3 = append(nums3, 1)
	//fmt.Printf("new cap: %d\n", cap(nums3))
	//nums3 = append(nums3, 1)
	//fmt.Printf("new cap: %d\n", cap(nums3))
	//nums3 = append(nums3, 1)
	//fmt.Printf("new cap: %d\n", cap(nums3))

}

func testStruct() {
	// nil
	var a = []int{1, 2, 3, 4, 5}
	// not nil
	b := make([]int, 3)
	fmt.Println(b)

	if a == nil {
		fmt.Println("a is nil")
	} else {
		fmt.Println("a is not nil")
	}

	if b == nil {
		fmt.Println("b is nil")
	} else {
		fmt.Println("b is not nil")
	}

	as := (*reflect.SliceHeader)(unsafe.Pointer(&a))
	bs := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	fmt.Printf("len = %d, cap = %d, type = %d\n", len(a), cap(a), as.Data)
	fmt.Printf("len = %d, cap = %d, type = %d\n", len(b), cap(b), bs.Data)

	// 改变指向的底层数组
	bs.Data = as.Data
	fmt.Println(b)
	// 改变b切片
	b[0] = 5
	fmt.Println(a)
}
