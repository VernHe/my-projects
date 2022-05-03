package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"time"
)

type LimiterIface interface {
	Key(c *gin.Context) string                      // 根据URI获取对应的Key
	GetBucket(key string) (*ratelimit.Bucket, bool) // 获取对应的Bucket
	AddBuckets(rules ...BucketRule) LimiterIface    // 添加对应的Bucket
}

// Limiter 封装了ratelimit.Bucket
type Limiter struct {
	limiterBuckets map[string]*ratelimit.Bucket
}

type BucketRule struct {
	Key          string        // Bucket对应的Key
	FillInterval time.Duration // 间隔多久时间放 N 个令牌。
	Capacity     int64         // 容量
	Quantum      int64         // 每次到达间隔时间后所放的具体令牌数量
}
