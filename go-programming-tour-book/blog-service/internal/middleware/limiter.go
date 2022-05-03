package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
	"github.com/go-programming-tour-book/blog-service/pkg/errcode"
	"github.com/go-programming-tour-book/blog-service/pkg/limiter"
)

func RateLimiter(l limiter.LimiterIface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取key
		key := l.Key(c)
		// 检查是否被限流
		if bucket, exist := l.GetBucket(key); exist {
			// 如果存在对应的bucket，就尝试去从中获取令牌
			count := bucket.TakeAvailable(1)
			if count == 0 {
				// 没有获取到令牌，说明请求的速度超过了限制
				response := app.NewResponse(c)
				response.ToErrorResponse(errcode.TooManyRequests)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
