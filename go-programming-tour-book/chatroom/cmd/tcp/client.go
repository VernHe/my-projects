package main

import (
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:2022", time.Duration(time.Second*5))
	if err != nil {
		log.Fatalf("connection err: %v\n", err)
	}

	done := make(chan struct{})
	go func() {
		// 打印相应结果
		io.Copy(os.Stdout, conn)
		log.Println("done")
		// 结束通知
		done <- struct{}{}
	}()

	// 处理用户输入
	mustCopy(conn, os.Stdin)

	// 关闭连接
	conn.Close()
	// 等待连接关闭
	<-done
}

func mustCopy(conn net.Conn, stdin *os.File) {
	if _, err := io.Copy(conn, stdin); err != nil {
		log.Fatal(err)
	}
}
