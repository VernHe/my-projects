package cache

import (
	"fmt"
	"log"
	"runtime"
	"sync"
)

var DefaultMaxBytes = 1 << 29

type safeCache struct {
	m     sync.RWMutex
	cache Cache
	// 记录命中率
	nhit, nget int
}

type Cache interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Del(key string)
	DelOldest()
	Len() int
}

// Value 对于存储的Value需要知道其所占容量
type Value interface {
	Len() int
}

type Stat struct {
	Nhit, Nget int
}

func CalcLen(value interface{}) int {
	var l int
	switch v := value.(type) {
	case Value:
		l = v.Len()
	case string:
		if runtime.GOARCH == "amd64" {
			l = 16 + len(v)
		} else {
			l = 8 + len(v)
		}
	case bool, uint8, int8:
		l = 1
	case int16, uint16:
		l = 2
	case int32, uint32, float32:
		l = 4
	case int64, uint64, float64:
		l = 8
	case int, uint:
		if runtime.GOARCH == "amd64" {
			l = 8
		} else {
			l = 4
		}
	case complex64:
		l = 8
	case complex128:
		l = 16
	default:
		panic(fmt.Sprintf("%T is not implement cache.Value", value))
	}
	return l
}

func newSafeCache(cache Cache) *safeCache {
	return &safeCache{
		cache: cache,
	}
}

func (sc *safeCache) set(key string, value interface{}) {
	sc.m.Lock()
	defer sc.m.Unlock()
	sc.cache.Set(key, value)
}

func (sc *safeCache) get(key string) interface{} {
	sc.m.RLock()
	defer sc.m.RUnlock()
	sc.nget++

	if sc.cache == nil {
		return nil
	}

	v := sc.cache.Get(key)
	// 命中
	if v != nil {
		log.Println("[Tour Cache] hit")
		sc.nhit++
	}
	return v
}

func (sc *safeCache) stat() *Stat {
	return &Stat{
		Nget: sc.nget,
		Nhit: sc.nhit,
	}
}
