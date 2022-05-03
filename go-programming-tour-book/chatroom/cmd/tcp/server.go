package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

var (
	// 新用户上线通知
	enteringChannel = make(chan *User)
	// 用户离线通知
	leavingChannel = make(chan *User)
	messageChannel = make(chan *Message, 8)
)

type Message struct {
	Uid     int
	Content string
}

// User 用户对象
type User struct {
	ID             int
	Addr           string
	EnterAt        time.Time
	MessageChannel chan string
}

func (u *User) String() string {
	return fmt.Sprintf("{ID:%d , Addr:%s, EnterTime: %v}", u.ID, u.Addr, u.EnterAt.Format("2006/01/02 15:04:05"))
}

func (u *User) message(content string) *Message {
	return &Message{
		Uid:     u.ID,
		Content: content,
	}
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:2022")
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}

	go broadcaster()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("connection err: %v\n", err)
			continue
		}

		go handleConn(conn)
	}
}

// handleConn 处理连接
func handleConn(conn net.Conn) {
	defer conn.Close()

	// 创建一个user对象
	user := &User{
		ID:             GenUserId(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}

	// 用户每次活跃都会用于通知重置定时器
	active := make(chan struct{})
	// 超时强踢
	go func() {
		duration := time.Second * 30
		// 定时器
		timer := time.NewTimer(duration)
		select {
		case <-timer.C:
			// 30s没反应就强制踢下线
			conn.Close()
		case <-active:
			// 如果活跃就重置定时器
			timer.Reset(duration)
		}
	}()

	// 发送欢迎通知
	user.MessageChannel <- "Welcome, " + user.String()

	// 广播上线通知
	messageChannel <- user.message("user `" + strconv.Itoa(user.ID) + "` has enter")
	enteringChannel <- user

	// 开启一个goroutine用于给用户发送消息
	go sendMessage(conn, user.MessageChannel)

	// 监听用户的输入
	input := bufio.NewScanner(conn)
	for input.Scan() {
		active <- struct{}{}
		messageChannel <- user.message(strconv.Itoa(user.ID) + ":" + input.Text())
	}

	// 用户断开连接
	if err := input.Err(); err != nil {
		log.Printf("%s 发生错误,错误信息: %v\n", user.String(), err)
	}
	leavingChannel <- user

	// 广播离线通知
	messageChannel <- user.message("user " + strconv.Itoa(user.ID) + "has left")
}

// sendMessage 监听chan并给用户发送消息
// <-chan string 表示此channel只能读，不能写
func sendMessage(conn net.Conn, ch <-chan string) {
	// 通过range读取channel中的数据，当channel被close后会自动跳出循环
	for msg := range ch {
		// 给用户发送消息
		fmt.Fprintln(conn, msg)
	}
}

// GenUserId 生成用户ID
func GenUserId() int {
	return rand.Int()
}

// broadcaster 用于记录聊天室用户，并进行消息广播：
// 1. 新用户进来；2. 用户普通消息；3. 用户离开
func broadcaster() {
	onlineUsers := make(map[*User]struct{})
	for {
		select {
		case user := <-enteringChannel:
			// 新用户进入
			onlineUsers[user] = struct{}{}
			log.Println(user.String() + "上线了")
		case user := <-leavingChannel:
			// 用户离开
			delete(onlineUsers, user)
			// 避免信息泄露
			close(user.MessageChannel)
			log.Println(user.String() + "下线了")
		case msg := <-messageChannel:
			// 发送广播消息
			for onlineUser := range onlineUsers {
				// 避免发给自己
				if onlineUser.ID == msg.Uid {
					continue
				}
				onlineUser.MessageChannel <- msg.Content
			}
		}
	}
}
