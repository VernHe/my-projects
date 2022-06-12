package main

// 测试闭包
func fa(a int) func(i int) int {
	// 函数引用到了函数之外的变量
	return func(i int) int {
		println("a:", &a, a)
		a += i
		return a
	}
}

func main() {
	//a: 0xc000043f60 1
	//2
	//a: 0xc000043f60 2
	//4
	f := fa(1)
	println(f(1))
	println(f(2))
}
