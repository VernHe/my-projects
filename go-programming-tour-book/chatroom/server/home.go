package server

import (
	"fmt"
	"github.com/goprogramming-tour-book/chatroom/global"
	"html/template"
	"log"
	"net/http"
)

func homeHandleFunc(writer http.ResponseWriter, request *http.Request) {
	tpl, err := template.ParseFiles(global.RootDir + "/template/home.html")
	if err != nil {
		// 返回错误信息
		fmt.Fprint(writer, "模板解析错误!")
		log.Printf("err: %v\n", err)
		return
	}

	// 将渲染后的文件传输给前端
	err = tpl.Execute(writer, nil)
	if err != nil {
		// 返回错误信息
		fmt.Fprint(writer, "模板执行错误!")
		log.Printf("err: %v\n", err)
		return
	}
}
