package main

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	// 网络类型
	network = "tcp4"
	// 超时强踢
	timeOutMillion = 5 * time.Second
)

// Server 服务器结构体
type Server struct {
	Ip   string
	Port int
	// 在线用户列表
	OnlineMap map[string]*User
	// 用户广播的channel
	BroadcastChannel chan string
	// 锁
	lock sync.RWMutex
}

// NewServer 创建服务器
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:               ip,
		Port:             port,
		OnlineMap:        make(map[string]*User),
		BroadcastChannel: make(chan string),
	}
	return server
}

// GetAddress 获取地址,ip:port
func (server *Server) GetAddress() string {
	return fmt.Sprintf("%s:%d", server.Ip, server.Port)
}

// Start 启动
func (server *Server) Start() {
	// 监听指定的端口
	listener, err := net.Listen(network, server.GetAddress())
	if HasError(err) {
		ErrorMessage(err)
		return
	}
	LogMessage("服务器启动")

	// 此方法执行结束后服务器将会关闭
	defer func() {
		err := listener.Close()
		if HasError(err) {
			ErrorMessage(err)
		}
	}()

	// 监听服务器的广播消息
	go server.ListenBroadcastMsg()

	// 不断地接收连接
	for {
		conn, err := listener.Accept()
		if HasError(err) {
			ErrorMessage(err)
			continue
		}
		// 创建一个goroutine去处理此次的连接
		go server.HandleConn(conn)
	}

}

// ListenBroadcastMsg 监听并发送广播的消息
func (server *Server) ListenBroadcastMsg() {
	LogMessage("正在监听广播信息")
	for {
		broadcastMsg := <-server.BroadcastChannel
		// 遍历每一个在线用户
		for _, user := range server.OnlineMap {
			user.receiveMsg(broadcastMsg)
		}
	}
}

// HandleConn 处理某个单独的连接
func (server *Server) HandleConn(conn net.Conn) {
	LogMessage(fmt.Sprintf("与%s建立了连接", conn.RemoteAddr()))

	// 增加新用户
	user := server.creatNewUser(conn)

	// 用户上线
	user.online()

	isAlive := make(chan bool)

	// 开启一个goroutine去监听用户发送的信息
	go server.listenUserInput(user, isAlive)

	// 阻塞
	// 超时强行踢下线
	for {
		select {
		// 当isAlive管道中读取到数据后，会执行下面的case，因此会重新倒计时
		case <-isAlive:
			// 不需要做任何处理，只为重新激活下面的计时
		case <-time.After(timeOutMillion):
			// 用户下线
			user.offline()
			// 发送通知
			user.doSendMsg("你已被踢下线了\n")
			// 关闭连接
			err := conn.Close()
			if err != nil {
				ErrorMessage(err)
				return
			}
			// 关闭资源
			close(isAlive)
			// 关闭专门处理当前用户的goroutine
			runtime.Goexit()
			//return
		}
	}
}

// creatNewUser 创建一个用户对象
func (server *Server) creatNewUser(conn net.Conn) *User {
	return NewUser(conn, server)
}

// sendBroadcastMsg 发送广播
func (server *Server) sendBroadcastMsg(msg string) {
	// 向广播Channel中添加信息
	server.BroadcastChannel <- msg
}

// addOnlineRecord 新增一条在线用户的记录
func (server *Server) addOnlineRecord(user *User) {
	server.lock.Lock()
	server.OnlineMap[user.Name] = user
	server.lock.Unlock()
}

// deleteOnlineRecord 删除已下线的用户的记录
func (server *Server) deleteOnlineRecord(user *User) {
	server.lock.Lock()
	delete(server.OnlineMap, user.Name)
	server.lock.Unlock()
}

// listenUserInput 监听客户端的输入
func (server *Server) listenUserInput(user *User, isAlive chan bool) {
	// 读取数据的缓存
	buf := make([]byte, 4096)
	// 用户的消息
	var userMsg string

	for {
		// num 表示读取到的字符数量
		num, err := user.getInputMsg(buf)

		// 用户下线了(如 Ctrl + C)
		if num == 0 {
			user.offline()
			// 停止监听
			return
		}

		// 发生异常
		if HasError(err) && err != io.EOF {
			ErrorMessage(err)
			// 停止监听
			return
		}

		// 去除回车
		if buf[num-1] == '\n' {
			userMsg = string(buf[:num-1])
		} else {
			userMsg = string(buf[:num])
		}

		// 处理用户的输入
		user.resolveMsg(userMsg)

		isAlive <- true
	}

}

// getOnlineUserList 获得当前在线用户
func (server *Server) getOnlineUserList() string {
	var buf strings.Builder
	server.lock.Lock()
	for name, user := range server.OnlineMap {
		buf.WriteString(fmt.Sprintf("[%s] %s 在线\n", name, user.Addr))
	}
	server.lock.Unlock()
	return buf.String()
}

// getOnlineUserByName 获取某个用户
func (server *Server) getOnlineUserByName(userName string) *User {
	user, ok := server.OnlineMap[userName]
	if ok {
		return user
	} else {
		user.doSendMsg(fmt.Sprintf("用户 %s 不存在\n", userName))
		return nil
	}
}
