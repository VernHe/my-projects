package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
	"github.com/go-programming-tour-book/blog-service/pkg/email"
	"github.com/go-programming-tour-book/blog-service/pkg/errcode"
	"time"
)

func Recovery() gin.HandlerFunc {
	// 创建邮件对象
	e := email.NewEmail(&email.SMTPInfo{
		Host:     global.EmailSetting.Host,
		Port:     global.EmailSetting.Port,
		IsSSL:    global.EmailSetting.IsSSL,
		UserName: global.EmailSetting.UserName,
		Password: global.EmailSetting.Password,
		From:     global.EmailSetting.From,
	})
	// 返回处理panic的方法
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 写日志
				global.Logger.WithCallersFrames().Infof("panic recover err: %v", err)
				// 发邮件
				err = e.SendMail(
					global.EmailSetting.To,
					fmt.Sprintf("异常抛出，发生时间：%d", time.Now().Unix()),
					fmt.Sprintf("错误信息: %v", err),
				)
				if err != nil {
					global.Logger.Panicf("mail.SendMail err: %v", err)
				}
				// 响应
				app.NewResponse(c).ToErrorResponse(errcode.ServerError)
				// 不再处理此次请求
				c.Abort()
			}
		}()
		c.Next()
	}
}
