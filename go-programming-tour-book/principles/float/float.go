package main

import (
	"fmt"
	"math/big"
)

func main() {
	test4()
}

// 测试打印的精度
func test1() {
	var f float32 = 0.333333333333
	fmt.Println(f) // 输出：0.33333334
}

func test2() {
	var n float64 = 0
	for i := 0; i < 1000; i++ {
		n += 0.01
	}
	// 输出: 9.999999999999831
	fmt.Println(n)
}

func test3() {
	var n float64 = 0.123456789999
	var m float64 = 0.034567833333

	var i float64 = 3.0

	// 加法 + 乘法
	// 期望的结果: 		0.474073869996
	// A x ( B + C )	0.47407386999600004
	// A x B + A x C	0.474073869996
	fmt.Println(i * (n + m))
	fmt.Println(i*n + i*m)
}

func test4() {
	a := 0.1
	b := 0.2
	c := a + b

	bigA := big.NewFloat(0.1)
	bigB := big.NewFloat(0.2)
	bigC := big.NewFloat(0)

	// 设置bigC的精度,值为53时，精度等效于float64
	bigC.SetPrec(106)
	// 相加
	bigC.Add(bigA, bigB)

	fmt.Println(c)
	fmt.Println(bigC)
}
