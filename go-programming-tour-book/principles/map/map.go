package main

import "fmt"

// Map只支持并发读
func main() {
	testHash()
}

// 测试并发读写
// 并发读写会报错 fatal error: concurrent map read and map write
func test1() {
	m := make(map[int]int)
	go func() {
		for {
			m[0] = 5
		}
	}()
	go func() {
		for {
			_ = m[1]
		}
	}()
	select {}
}

// 测试并发读
func test2() {
	m := make(map[int]int)
	go func() {
		for {
			_ = m[0]
		}
	}()
	go func() {
		for {
			_ = m[1]
		}
	}()
	select {}
}

// 测试并发写
// 报错： fatal error: concurrent map writes
func test3() {
	m := make(map[int]int)
	go func() {
		for {
			m[0] = 5
		}
	}()
	go func() {
		for {
			m[1] = 5
		}
	}()
	select {}
}

// 计算hash值是截取后8位
func testHash() {
	// 对应二进制: 11011001000000111
	hash := 111111
	// uint8 8bit，范围是 0 ~ 255
	// 强制转换时只截取后8位，即00000111 = 7
	fmt.Println(uint8(hash))
}
