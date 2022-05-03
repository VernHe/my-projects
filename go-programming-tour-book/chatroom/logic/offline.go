package logic

import (
	"container/ring"
	"github.com/spf13/viper"
)

type offlineProcessor struct {
	// 缓冲的消息数量
	n int
	// 存放最近消息的循环双向链表
	recentRing *ring.Ring
	// 用户最近被@而收到的消息
	userRing map[string]*ring.Ring
}

var OfflineProcessor = newOfflineProcessor()

// newOfflineProcessor 创建离线消息处理器
func newOfflineProcessor() *offlineProcessor {
	n := viper.GetInt("offline-num")

	return &offlineProcessor{
		n:          n,
		recentRing: ring.New(n),
		userRing:   make(map[string]*ring.Ring),
	}
}

// Save 存消息
func (o *offlineProcessor) Save(msg *Message) {
	// 只保存用户发送的普通消息
	if msg.Type != MsgTypeNormal {
		return
	}

	// 保存至最近消息中
	o.recentRing.Value = msg
	o.recentRing = o.recentRing.Next()

	// 检查是否有@消息，如果有，为特定用户保存这类消息
	for _, at := range msg.Ats {
		// 检查该用户是否有创建对应的ring，如果没有则创建
		r, ok := o.userRing[at]
		if !ok {
			// 创建ring
			o.userRing[at] = ring.New(o.n)
			r = o.userRing[at]
		}
		// 存入私聊消息
		r.Value = msg.Content
		o.userRing[at] = r.Next()
	}
}

func (o *offlineProcessor) Send(user *User) {
	// 给用户发送最近的普通消息
	o.recentRing.Do(func(v interface{}) {
		if v != nil {
			user.MessageChannel <- v.(*Message)
		}
	})

	// 如果是新用户则跳过下面的步骤
	if user.IsNew {
		return
	}

	// 给用户发送最近的私聊消息，完成后清空
	r, ok := o.userRing[user.NickName]
	if ok {
		// 发送消息
		r.Do(func(v interface{}) {
			user.MessageChannel <- v.(*Message)
		})

		// 清空
		delete(o.userRing, user.NickName)
	}

}
