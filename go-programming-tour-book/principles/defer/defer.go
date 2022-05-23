package main

import (
	"fmt"
	"net/http"
)

type myHandler struct{}

func (m myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 故意报错
	a := 0
	b := 3
	fmt.Println(b / a)

	fmt.Fprintln(w, "你好")
}

// defer实现函数中间件,此处案例是捕获服务器异常，并返回500
func middleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// 中间件
		defer func() {
			if errMsg := recover(); errMsg != nil {
				fmt.Println(errMsg)
				http.Error(w, "服务器异常", http.StatusInternalServerError)
			}
		}()

		// 真正的逻辑
		h.ServeHTTP(w, r)
	}
	// 将函数包装成HandlerFunc
	return http.HandlerFunc(fn)
}

func testMiddleware() {
	handler := myHandler{}
	server := &http.Server{
		Addr:    ":8080",
		Handler: middleware(handler),
	}

	server.ListenAndServe()
}

// LIFO 与栈的执行顺序是相同的
func LIFO() {
	// 压栈
	defer func() {
		fmt.Println("第一个defer")
	}()
	// 压栈
	defer func() {
		fmt.Println("第二个defer")
	}()
	// 压栈、
	defer func() {
		fmt.Println("第三个defer")
	}()
}

func testParam() {
	a := 1
	// 参数的预计算指当函数到达defer语句时，延迟调用的参数将立即求值，传递到defer函数中的参数将预先被固定，而不会等到函数执行完成后再传递参数到defer中。
	defer func(num int) {
		fmt.Println(num)
	}(a)
	a = 99
}

// return 并非原子操作，步骤：将返回值保存在栈上→执行defer函数→函数返回。

var g = 100

// 返回的是100
func returnVal1() int {
	defer func() {
		// 将g变为200
		g = 200
	}()
	fmt.Printf("returnVal1函数: g = %d\n", g)
	return g
}

// 返回的是300
func returnVal2() (r int) {
	// r = 100 -> r = 0 -> r = 200
	r = g
	defer func() {
		// return 100后将g变为200
		r = 300
	}()
	r = 0
	return r
}

func main() {
	i := returnVal1()
	// 【returnVal1】 main函数: g = 200, 函数返回值 = 100
	// 【returnVal2】 main函数: g = 100, 函数返回值 = 300
	fmt.Printf("main函数: g = %d, 函数返回值 = %d\n", g, i)
}
