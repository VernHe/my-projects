package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "HTTP, Hello")
	})

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		// 建立连接
		conn, err := websocket.Accept(writer, request, nil)
		if err != nil {
			log.Printf("conn err: %v\n", err)
		}
		defer conn.Close(websocket.StatusInternalError, "the sky is falling")

		// 设置超时时间
		ctx, cancelFunc := context.WithTimeout(request.Context(), time.Second*10)
		defer cancelFunc()

		// 读取数据
		var v interface{}
		err = wsjson.Read(ctx, conn, &v)
		if err != nil {
			log.Printf("read err: %v\n", err)
		}
		log.Printf("接收到客户端数据: %v\n", v)

		// 写数据
		err = wsjson.Write(ctx, conn, "Hello Websocket Client")
		if err != nil {
			log.Printf("send err: %v\n", err)
		}
		conn.Close(websocket.StatusNormalClosure, "")
	})

	log.Fatalln(http.ListenAndServe(":2022", nil))
}
