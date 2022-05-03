package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

func ContextTimeout(t time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		timeout, cancelFunc := context.WithTimeout(c.Request.Context(), t)
		defer cancelFunc()

		c.Request.WithContext(timeout) // 替换上下文
		c.Next()
	}
}
