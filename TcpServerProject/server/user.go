package main

import (
	"fmt"
	"net"
	"runtime"
	"strings"
)

// User 用户对象
type User struct {
	Name       string
	Addr       string
	MsgChannel chan string
	Connection net.Conn
	server     *Server
}

// NewUser 新建一个用户
func NewUser(conn net.Conn, server *Server) *User {
	// 初始化用户对象
	newUser := &User{
		conn.RemoteAddr().String(),
		conn.RemoteAddr().String(),
		make(chan string),
		conn,
		server,
	}

	// 启动监听消息的goroutine
	go newUser.ListenMsgChanel()

	return newUser
}

// ListenMsgChanel 监听用户自己的Channel
func (user *User) ListenMsgChanel() {
	defer LogMessage("用户" + user.Name + "下线了")
	defer close(user.MsgChannel)
	defer user.closeConn()
	// 循环
	for {
		msg, ok := <-user.MsgChannel
		if ok {
			// 将消息发送给用户
			user.doSendMsg(msg)
		} else {
			runtime.Goexit()
		}
	}
}

// Online 用户上线
func (user *User) online() {
	// 向在线用户表中添加一条记录
	user.server.addOnlineRecord(user)
	// 发送上线通知
	user.sendOnlineMsg()
}

// Offline 用户下线
func (user *User) offline() {
	// 向在线用户表中添加一条记录
	user.server.deleteOnlineRecord(user)
	// 发送下线通知
	user.sendOfflineMsg()
}

// closeConn 关闭连接
func (user *User) closeConn() {
	err := user.Connection.Close()
	if err != nil {
		ErrorMessage(err)
	}
}

// SendOnlineMsg 发送上线通知
func (user *User) sendOnlineMsg() {
	msg := "用户 [" + user.Addr + "] 上线了"
	user.server.sendBroadcastMsg(msg)
}

// SendOfflineMsg 发送上线通知
func (user *User) sendOfflineMsg() {
	msg := "用户 [" + user.Addr + "] 下线了"
	user.server.sendBroadcastMsg(msg)
}

// receiveMsg 用户接收消息
func (user *User) receiveMsg(msg string) {
	user.MsgChannel <- msg
}

// getInputMsg 获取用户的输入
func (user *User) getInputMsg(buf []byte) (int, error) {
	return user.Connection.Read(buf)
}

// resolveMsg 解析用户的输入
func (user *User) resolveMsg(msg string) {
	if msg == "/users" {
		// 查询在线用户列表
		user.doSendMsg(user.server.getOnlineUserList())
	} else if len(msg) > 8 && msg[:8] == "/rename " {
		// 重命名
		user.rename(msg[8:])
	} else if len(msg) > 3 && msg[0:1] == "@" && msg[1:2] != "|" && len(strings.Split(msg, "|")) == 2 {
		// @张三|xxxx
		// 私聊某个用户
		user.sendPrivateMsg(strings.Split(msg, "|"))
	} else {
		user.server.sendBroadcastMsg(user.Name + " say : " + msg)
	}
}

// doSendMsg 将消息发送至用户
func (user *User) doSendMsg(msg string) {
	_, err := user.Connection.Write([]byte(msg + "\n"))
	if err != nil {
		ErrorMessage(err)
		return
	}
}

// sendBroadcastMsg 发送广播类型的消息
func (user *User) sendBroadcastMsg(msg string) {
	user.server.sendBroadcastMsg(user.Name + " say : " + msg)
}

// rename 用户修改昵称
func (user *User) rename(newName string) {
	user.server.lock.Lock()
	// 此名字有人使用
	if value, ok := user.server.OnlineMap[newName]; ok {
		user.doSendMsg(fmt.Sprintf("修改失败: %s 正在使用此昵称\n", value.Addr))
	} else {
		// 删除旧的记录，添加新的记录
		delete(user.server.OnlineMap, user.Name)
		user.server.OnlineMap[newName] = user
	}
	user.server.lock.Unlock()
	user.Name = newName
	user.doSendMsg("昵称修改成功，当前新昵称为 " + newName)

}

// sendPrivateMsg 私聊某个用户
func (user *User) sendPrivateMsg(msg []string) {
	// 检查是否存在此用户
	targetUser := user.server.getOnlineUserByName(msg[0][1:])

	// 目标用户不存在
	if targetUser != nil {
		targetUser.doSendMsg(fmt.Sprintf("%s 私聊你: %s\n", user.Name, msg[1]))
		return
	}
}
