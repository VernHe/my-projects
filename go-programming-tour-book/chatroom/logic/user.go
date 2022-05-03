package logic

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"math/rand"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"regexp"
	"strings"
	"sync/atomic"
	"time"
)

var (
	// System 系统用户
	System    = &User{}
	globalUid = rand.Uint32()
)

// User 用户对象
type User struct {
	UID            int           `json:"uid"`
	NickName       string        `json:"nickname"`
	EnterAt        time.Time     `json:"enter_at"`
	Addr           string        `json:"addr"`
	MessageChannel chan *Message `json:"-"`

	conn *websocket.Conn

	IsNew bool   `json:"is_new"`
	Token string `json:"token"`
}

func NewUser(conn *websocket.Conn, token, nickName, Addr string) *User {
	user := &User{
		NickName:       nickName,
		EnterAt:        time.Now(),
		Addr:           Addr,
		MessageChannel: make(chan *Message),
		conn:           conn,
		Token:          token,
	}

	if user.Token != "" {
		// 校验token，如果是老用户则更新其uid
		uid, err := parseTokenAndValidate(user.Token, user.NickName)
		if err != nil {
			user.UID = uid
		}
	}

	// 如果是新用户，生成uid与token
	if user.UID == 0 {
		// 自增
		user.UID = int(atomic.AddUint32(&globalUid, 1))
		user.Token = getToken(user.UID, user.NickName)
		user.IsNew = true
	}

	return user
}

// SendMessage 监听channel并向用户发送消息
func (u *User) SendMessage(ctx context.Context) {
	// 使用的是range操作，因此如果用户退出，此channel必须关闭，否则goroutine将会一直存在
	for msg := range u.MessageChannel {
		wsjson.Write(ctx, u.conn, msg)
	}
}

// CloseMessageChannel 关闭msg channel避免goroutine泄漏
func (u *User) CloseMessageChannel() {
	close(u.MessageChannel)
}

// ReceiveMessage 不断接收用户发送的消息
func (u *User) ReceiveMessage(ctx context.Context) error {
	var (
		msg map[string]string
		err error
	)
	for {
		// 从conn中读取数据
		err = wsjson.Read(ctx, u.conn, &msg)

		// 检查是否出现错误
		if err != nil {
			// 判断是否是用户关闭连接
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				// 关闭连接
				return nil
			}
			return err
		}

		// 从请求中获取消息内容并过滤敏感词汇
		message := NewUserMessage(u, FilterSensitive(msg["content"]))
		// 检查是否有@符号且长度大于2，如果不匹配则会panic
		compile := regexp.MustCompile(`@[^\s@]{2,20}`)
		// 设置@的人
		message.Ats = compile.FindAllString(message.Content, -1)
		// 去除'@'
		for i, at := range message.Ats {
			message.Ats[i] = at[1:]
		}
		// 发送广播消息
		Broadcaster.messageChannel <- message
	}
}

// getToken 计算token
func getToken(uid int, nickname string) string {
	// 用于加密的密钥
	secret := viper.GetString("token-secret")
	// 要加密的内容
	message := fmt.Sprintf("%s%s%d", nickname, secret, uid)
	messageMAC := macSha256([]byte(message), []byte(secret))
	return fmt.Sprintf("%suid%d", base64.StdEncoding.EncodeToString(messageMAC), uid)
}

// macSha256 sha256加密
func macSha256(message, secret []byte) []byte {
	hash := hmac.New(sha256.New, secret)
	hash.Write(message)
	return hash.Sum(nil)
}

// 解析token
func parseTokenAndValidate(token, nickname string) (int, error) {
	// 根据 'uid' 截取获得 uid
	index := strings.LastIndex(token, "uid")
	uid := cast.ToInt(token[index+3:])

	// 截取messageMac
	messageMac, err := base64.StdEncoding.DecodeString(token[:index])
	if err != nil {
		return 0, nil
	}

	secret := viper.GetString("token-secret")
	message := fmt.Sprintf("%s%s%d", nickname, secret, uid)

	// 对比
	isValid := validateMac([]byte(message), messageMac, []byte(secret))
	if isValid {
		return uid, nil
	}

	return 0, errors.New("token is illegal")
}

// 判断是否相等
func validateMac(message, messageMac, secret []byte) bool {
	hash := hmac.New(sha256.New, secret)
	hash.Write(message)
	return hmac.Equal(messageMac, hash.Sum(nil))
}
