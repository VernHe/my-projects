package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/goprogramming-tour-book/chatroom/logic"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"strconv"
	"time"
)

var (
	userNum       int
	loginInterval time.Duration
	msgInterval   time.Duration
)

func init() {
	flag.IntVar(&userNum, "u", 500, "用户数量")
	flag.DurationVar(&loginInterval, "l", 5e9, "用户登录时间间隔")
	flag.DurationVar(&msgInterval, "m", time.Minute, "用户发送消息时间间隔")
}

func main() {
	flag.Parse()
	for i := 0; i < userNum; i++ {
		go UserConnect("user" + strconv.Itoa(i))
		time.Sleep(loginInterval)
	}
	// 防止主线程退出
	select {}
}

// UserConnect 开启一个用户连接
func UserConnect(nickname string) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFunc()
	// 建立连接
	conn, _, err := websocket.Dial(ctx, "ws://127.0.0.1:2022/ws?nickname="+nickname, nil)
	if err != nil {
		log.Printf("Dail error: %v\n", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "内部错误")

	// 发送消息
	go sendMsg(conn, nickname)
	// 接收消息
	ctx = context.Background()
	for {
		var message logic.Message
		err := wsjson.Read(ctx, conn, &message)
		if err != nil {
			log.Printf("receving msg error: %v\n", err)
			continue
		}
		if message.ClientSendTime.IsZero() {
			continue
		}
		if d := time.Now().Sub(message.ClientSendTime); d > time.Second {
			// 记录大于1s的消息
			fmt.Printf("接收到响应的消息(%d): %v\n", d.Milliseconds(), message.Content)
		}
	}
	conn.Close(websocket.StatusNormalClosure, "")
}

// sendMsg 不断发送消息
func sendMsg(conn *websocket.Conn, nickName string) {
	ctx := context.Background()
	msg := make(map[string]string)
	i := 1
	for {
		msg["content"] = "来自" + nickName + "的消息: " + strconv.Itoa(i)
		msg["send_time"] = strconv.FormatInt(time.Now().Unix(), 10)
		err := wsjson.Write(ctx, conn, msg)
		if err != nil {
			log.Println("send msg error:", err, "nickname:", nickName, "no:", i)
		}
		i++
		time.Sleep(msgInterval)
	}
}
