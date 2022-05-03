package main

import (
	"flag"
	"fmt"
)

var serverIp string
var serverPort int

func init() {
	// 初始化命令行的参数，可通过-h查看
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "目标服务器的IP地址(Ipv4)")
	flag.IntVar(&serverPort, "port", 8888, "目标服务器的端口号")
}

func main() {

	// 解析命令行
	flag.Parse()

	// 创建client
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("客户端启动失败")
		return
	}

	// 启动一个goroutine去监听服务器的响应，并显示在终端
	go client.printMsgOfServer()
	
	fmt.Println("客户端启动成功")
	client.run()
}
