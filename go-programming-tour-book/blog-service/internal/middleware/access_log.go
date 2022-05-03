package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/pkg/logger"
	"time"
)

type AccessLogWriter struct {
	gin.ResponseWriter               // 在写入流时，调用的是 http.ResponseWriter
	body               *bytes.Buffer // 保存相应的数据
}

// Write 记录相应的数据，在这里完成双写，
func (w AccessLogWriter) Write(p []byte) (int, error) {
	// 备份一份响应的内容
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p) // 写入HTTP的响应内容
}

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessLogWriter := &AccessLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		// 替换gin默认的
		c.Writer = accessLogWriter

		// 记录时间
		beginTime := time.Now().Unix()
		c.Next()
		endTime := time.Now().Unix()

		fields := logger.Fields{
			"request":  c.Request.Form.Encode(),       // 请求体
			"response": accessLogWriter.body.String(), // 响应内容
		}

		// 输出日志
		global.Logger.WithFields(fields).Infof("access log: method: %s, status_code: %d, begin_time: %d, end_time: %d",
			c.Request.Method,         // 请求方法
			accessLogWriter.Status(), // 响应吗
			beginTime,                // 开始时间
			endTime,                  // 结束时间
		)
	}
}
