package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

const (
	network = "tcp4"
)

// Client 客户端
type Client struct {
	ServerIp   string
	ServerPort int
	Conn       net.Conn
	Name       string
	flag       int // 用户所当前选择的功能
}

func NewClient(ip string, port int) *Client {
	addr := fmt.Sprintf("%s:%d", ip, port)
	// 建立连接
	conn, err := net.Dial(network, addr)
	if err != nil {
		ErrorMessage(err)
		return nil
	}
	// 构建Client对象
	return &Client{
		ServerIp:   ip,
		ServerPort: port,
		Conn:       conn,
		Name:       addr,
		flag:       999,
	}
}

// menu 打印菜单，return true: 选择完成, false: 选择失败
func (client *Client) menu() bool {
	var inputVal int

	fmt.Println("1、公共聊天")
	fmt.Println("2、私密聊天")
	fmt.Println("3、更改昵称")
	fmt.Println("0、退   出")

	// 监听用户的输入
	_, err := fmt.Scanln(&inputVal)
	if err != nil {
		ErrorMessage(err)
		return false
	}
	// 合法性校验
	if inputVal >= 0 && inputVal <= 3 {
		client.flag = inputVal
		return true
	} else {
		fmt.Println("请输入正确的选项....")
		return false
	}
}

func (client *Client) run() {
	// 如果用户未退出
	for client.flag != 0 {
		// 如果用户没有做合法的选择
		for client.menu() != true {
		}

		// 当用户做了合法的选择后，做相应的处理
		switch client.flag {
		case 0:
			fmt.Println("退   出")
			break
		case 1:
			client.publicChat()
			break
		case 2:
			client.privateChat()
			break
		case 3:
			client.rename()
			break
		}
	}
}

// printMsgOfServer 接收并打印服务器响应的内容
func (client *Client) printMsgOfServer() {
	// 永久阻塞在这里，并不是执行一次
	_, err := io.Copy(os.Stdout, client.Conn)
	if err != nil {
		ErrorMessage(err)
		return
	}

	// 效果等同于下面

	//buf := make([]byte, 4096)
	//for {
	//	num, err := client.Conn.Read(buf)
	//	if err != nil {
	//		ErrorMessage(err)
	//		continue
	//	}
	//	fmt.Println(string(buf[:num]))
	//}
}

// sendMstToServer 向服务器发送数据
func (client *Client) sendMstToServer(msg string) {
	_, err := client.Conn.Write([]byte(msg))
	if err != nil {
		ErrorMessage(err)
		return
	}
}

// 查询当前在线用户
func (client *Client) getOnlineUserList() {
	client.sendMstToServer("/users")
}

// 发送私密消息
func (client Client) sendPrivateMsg(targetName string, msg string) {
	client.sendMstToServer(fmt.Sprintf("@%s|%s", targetName, msg))
}

// publicChat 公聊模式
func (client *Client) publicChat() {

	var inputStr string

	// 提示
	fmt.Println("进入公聊模式，请输入要发送的消息(/exit 退出此模式)")

	// 读取用户输入
	_, err := fmt.Scanln(&inputStr)
	if err != nil {
		ErrorMessage(err)
		return
	}

	// 一直监听用户的输入（除非退出公聊模式）
	for inputStr != "/exit" {
		// 发送数据
		if len(inputStr) > 0 {
			client.sendMstToServer(inputStr)
		}
		// 重置inputStr
		inputStr = ""
		// 继续接收输入
		_, err := fmt.Scanln(&inputStr)
		if err != nil {
			ErrorMessage(err)
			return
		}
		fmt.Println("请输入要发送的消息(/exit 退出此模式)")
	}

}

// privateChat 私聊模式
func (client *Client) privateChat() {
	// 查询在线用户
	client.getOnlineUserList()

	// 选择对应用户
	fmt.Println("请输入要私聊的人的名字(输入/exit 退出):")
	var targetUserName string
	_, err := fmt.Scanln(&targetUserName)
	if err != nil {
		ErrorMessage(err)
		return
	}

	// 发送消息
	for targetUserName != "/exit" {
		// 用户的输入
		var inputStr string
		for inputStr != "/exit" {
			// 发送用户信息给指定用户
			if len(inputStr) > 0 {
				client.sendPrivateMsg(targetUserName, inputStr)
			}

			inputStr = ""
			// 继续监听用户输入
			fmt.Println("请输入聊天内容(输入/exit 退出):")
			_, err := fmt.Scanln(&inputStr)
			if err != nil {
				ErrorMessage(err)
				return
			}
		}
		fmt.Println("请输入要私聊的人的名字(输入/exit 退出):")
		_, err = fmt.Scanln(&targetUserName)
		if err != nil {
			ErrorMessage(err)
			return
		}
	}
}

// rename 重命名
func (client *Client) rename() {
	var newName string
	// 提示
	fmt.Println("请输入新的用户名:")
	// 接收用户输入
	_, err := fmt.Scanln(&newName)
	if err != nil {
		ErrorMessage(err)
		return
	}
	// 发送给服务器
	client.sendMstToServer(fmt.Sprintf("/rename %s", newName))
}
