package logic

import (
	"github.com/goprogramming-tour-book/chatroom/global"
	"log"
)

// broadcaster 不导出，因为是单例模式，不允许其他地方调用
type broadcaster struct {
	// 用户列表
	Users map[string]*User

	// 接收通知的chan
	enteringChannel chan *User
	leavingChannel  chan *User
	messageChannel  chan *Message

	// 判断该昵称用户是否可以进入聊天室的
	checkUserChannel chan string
	// 结果
	checkUserCanInChannel chan bool
}

// Broadcaster 广播器，单例（饿汉模式）
var Broadcaster = broadcaster{
	Users:                 make(map[string]*User),
	enteringChannel:       make(chan *User),
	leavingChannel:        make(chan *User),
	messageChannel:        make(chan *Message, global.MsgChanSize),
	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),
}

func (b *broadcaster) Start() {
	for {

		select {
		case user := <-b.enteringChannel:
			// 新用户进入
			b.Users[user.NickName] = user
			// 发送最新的用户列表
			b.sendUserList()
			// 发送离线期间收到的消息
			OfflineProcessor.Send(user)
		case user := <-b.leavingChannel:
			// 用户退出
			delete(b.Users, user.NickName)
			// 避免 goroutine泄漏
			user.CloseMessageChannel()
			// 发送最新的用户列表
			b.sendUserList()
		case msg := <-b.messageChannel:
			if len(msg.Ats) == 0 {
				// 发送广播消息
				for _, user := range b.Users {
					// 避自己收到自己的消息
					if user.UID == msg.User.UID {
						continue
					}
					user.MessageChannel <- msg
				}
			} else {
				log.Printf("准备发送私信：%v\n", msg.Ats)
				for _, nickName := range msg.Ats {
					if user, ok := b.Users[nickName]; ok {
						user.MessageChannel <- msg
					}
					// @了不存在的用户
				}
			}
			// 保存消息
			OfflineProcessor.Save(msg)
		case nickName := <-b.checkUserChannel:
			// 检查昵称是否已经存在
			if _, ok := b.Users[nickName]; ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		}
	}
}

// CanEnterRoom 检查用户是否可以进入聊天室
func (b *broadcaster) CanEnterRoom(nickName string) bool {
	b.checkUserChannel <- nickName
	return <-b.checkUserCanInChannel
}

func (b *broadcaster) UserEntering(u *User) {
	b.enteringChannel <- u
}

func (b *broadcaster) UserLeaving(u *User) {
	b.leavingChannel <- u
}

func (b *broadcaster) Broadcast(msg *Message) {
	b.messageChannel <- msg
}

func (b *broadcaster) sendUserList() {
	userList := make([]*User, 0, len(b.Users))
	for _, user := range b.Users {
		userList = append(userList, user)
	}

	go func() {
		if len(b.messageChannel) < global.MsgChanSize {
			// 正常
			b.messageChannel <- NewUserListMessage(userList)
		} else {
			log.Println("消息并发量过大，导致MessageChannel拥堵。。。")
		}
	}()
}
