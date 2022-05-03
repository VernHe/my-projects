package main

import (
	"context"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func main() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	conn, _, err := websocket.Dial(ctx, "ws://localhost:2022/ws", nil)
	if err != nil {
		log.Printf("dial err: %v\n", err)
	}
	defer conn.Close(websocket.StatusInternalError, "the sky is falling")

	err = wsjson.Write(ctx, conn, "Hello Websocket Server")
	if err != nil {
		log.Printf("send err: %v\n", err)
	}

	var v interface{}
	err = wsjson.Read(ctx, conn, &v)
	if err != nil {
		log.Printf("receive err: %v\n", err)
	}
	log.Printf("接收到消息: %v\n", v)

	conn.Close(websocket.StatusNormalClosure, "")
}
