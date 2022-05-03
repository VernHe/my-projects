package server

import (
	"github.com/goprogramming-tour-book/chatroom/logic"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func WebSocketHandleFunc(writer http.ResponseWriter, request *http.Request) {
	// Accept 从客户端接收 WebSocket 握手，并将连接升级到 WebSocket。
	// 如果 Origin 域与主机不同，Accept 将拒绝握手，除非设置了 InsecureSkipVerify 选项（通过第三个参数 AcceptOptions 设置）。
	// 换句话说，默认情况下，它不允许跨源请求。如果发生错误，Accept 将始终写入适当的响应
	conn, err := websocket.Accept(writer, request, nil)
	if err != nil {
		log.Println("websocket accept error:", err)
		return
	}

	// 1、获取昵称, ws://host:port/ws?nickname=xxx&token=xxx
	nickname := request.FormValue("nickname")
	token := request.FormValue("token")
	// 昵称合法性校验
	if l := len(nickname); l < 2 || l > 20 {
		log.Println("nickname illegal: ", nickname)
		wsjson.Write(request.Context(), conn, logic.NewErrorMessage("非法昵称，昵称长度：4-20"))
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal!")
		return
	}
	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("昵称已经存在：", nickname)
		wsjson.Write(request.Context(), conn, logic.NewErrorMessage("该昵称已经已存在！"))
		conn.Close(websocket.StatusUnsupportedData, "nickname exists!")
		return
	}

	user := logic.NewUser(conn, token, nickname, request.RemoteAddr)

	// 2、为用户开启一个goroutine，用于给用户发送消息
	go user.SendMessage(request.Context())

	// 3、加入用户列表
	logic.Broadcaster.UserEntering(user)
	log.Printf("user %s joins chat\n", nickname)
	user.MessageChannel <- logic.NewWelcomeMessage(nickname)
	logic.Broadcaster.Broadcast(logic.NewNoticeMessage(nickname + "加入了聊天室"))

	// 4、开启一个goroutine用于接收用户的消息
	err = user.ReceiveMessage(request.Context())

	// 5、用户断开连接
	logic.Broadcaster.UserLeaving(user)
	log.Printf("user %s leaves chat\n", nickname)
	logic.Broadcaster.Broadcast(logic.NewNoticeMessage(nickname + "离开了聊天室"))

	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Printf("close err: %v\n", err)
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}

}
