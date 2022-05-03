package logic

import (
	"time"
)

// 消息类型
const (
	MsgTypeNormal   = iota // 用户消息
	MsgTypeSystem          // 系统消息
	MsgTypeError           // 错误消息
	MsgTypeUserList        // 发送当前用户列表
)

// Message 消息
type Message struct {
	User           *User     `json:"user"`
	Type           int       `json:"type"`
	Content        string    `json:"content"`
	MsgTime        time.Time `json:"msg_time"`
	ClientSendTime time.Time `json:"client_send_time"`

	// @的人
	Ats []string `json:"ats"`

	Users map[string]*User `json:"users"`
}

func NewUserMessage(u *User, msg string) *Message {
	return &Message{
		User:           u,
		Type:           MsgTypeNormal,
		Content:        msg,
		MsgTime:        time.Now(),
		ClientSendTime: time.Now(),
	}
}

func NewErrorMessage(msg string) *Message {
	return &Message{
		User:    System,
		Type:    MsgTypeError,
		Content: msg,
		MsgTime: time.Now(),
	}
}

func NewWelcomeMessage(nickName string) *Message {
	return &Message{
		User:    System,
		Type:    MsgTypeSystem,
		Content: "欢迎加入聊天室！" + nickName,
		MsgTime: time.Now(),
		Ats:     nil,
	}
}

func NewNoticeMessage(msg string) *Message {
	return &Message{
		User:    System,
		Type:    MsgTypeSystem,
		Content: msg,
		MsgTime: time.Now(),
		Ats:     nil,
	}
}

func NewUserListMessage(userList []*User) *Message {
	var strList string
	for _, user := range userList {
		strList += user.NickName + " "
	}
	return &Message{
		User:    System,
		Type:    MsgTypeSystem,
		Content: strList,
		MsgTime: time.Now(),
		Ats:     nil,
	}
}
