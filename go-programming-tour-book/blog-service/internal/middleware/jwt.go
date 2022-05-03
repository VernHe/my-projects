package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
	"github.com/go-programming-tour-book/blog-service/pkg/errcode"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		if fullPath := c.FullPath(); fullPath != "/auth" {
			var token string
			var ecode = errcode.Success
			// 获取token
			if t, exist := c.GetQuery("token"); exist {
				token = t
			} else {
				token = c.GetHeader("token")
			}
			// 设置code
			if token == "" {
				ecode = errcode.Unauthorized
			} else {
				// 有token，校验合法性
				_, err := app.ParseToken(token)
				if err != nil {
					// err.(*jwt.ValidationError) 断言err是*jwt.ValidationError类型
					switch err.(*jwt.ValidationError).Errors {
					case jwt.ValidationErrorExpired:
						ecode = errcode.UnauthorizedTokenTimeout
					default:
						ecode = errcode.UnauthorizedTokenError
					}
				}
			}

			if ecode != errcode.Success {
				// 出现错误
				response := app.NewResponse(c)
				response.ToErrorResponse(ecode)
				c.Abort() // 确保不调用此请求的其余处理程序
				return
			}
		}

		c.Next()

	}
}
