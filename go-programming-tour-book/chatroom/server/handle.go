package server

import (
	"github.com/goprogramming-tour-book/chatroom/logic"
	"net/http"
)

func RegisterHandle() {
	// 推断根目录
	//inferRootDir()

	// 启动广播器
	go logic.Broadcaster.Start()

	// 注册Handler
	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)

}
