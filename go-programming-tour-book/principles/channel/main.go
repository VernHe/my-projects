package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	//testChannel()
	testNil()
}

func testChannel() {
	links := []string{
		"http://www.baidu.com",
		"http://www.jd.com",
		"http://www.taobao.com",
	}
	c := make(chan string)
	fmt.Println(cap(c))
	for _, link := range links {
		go checkLink(link, c)
	}

	// 不好理解
	//for {
	//	go checkLink(<-c, c)
	//}

	// 会出现必闭包问题（引用传递）
	//for l := range c {
	//	go func() {
	//		time.Sleep(2 * time.Second)
	//		checkLink(l, c)
	//	}()
	//}

	// 通过参数复制一份（值传递），解决闭包问题
	for l := range c {
		go func(link string) {
			time.Sleep(2 * time.Second)
			checkLink(link, c)
		}(l)
	}
}

func checkLink(link string, c chan string) {
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, "might be down")
		c <- link
		return
	}
	fmt.Println(link, "is up")
	c <- link
}

func testSelect() {
	c := make(chan int, 1)
	c <- 1
	fmt.Println("开始执行select")
	select {
	case <-c:
		println("random 1")
	case <-c:
		println("random 2")
	}
}

func testTimeoutSelect() {
	c := make(chan int, 1)
	for {
		select {
		case <-c:
			println("get msg")
		case <-time.After(2 * time.Second):
			println("time out")
		}
	}
}

func testDefault() {
	c := make(chan int, 1)
	for {
		select {
		case <-c:
			println("get msg")
		default:
			println("no msg")
		}
	}
}

func testNil() {
	a := make(chan int)
	b := make(chan int)
	go func() {
		for i := 0; i < 4; i++ {
			// 当select语句的case对nil通道进行操作时，case分支将永远得不到执行。
			// 对一个nil的通道进行读写操作会阻塞
			select {
			case a <- 1:
				a = nil
			case b <- 2:
				b = nil
			}
		}
	}()
	// 实现交替写入ab的目的
	fmt.Println(<-a)
	fmt.Println(<-b)
}
