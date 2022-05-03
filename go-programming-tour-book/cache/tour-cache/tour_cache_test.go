package tour_cache

import (
	"github.com/go-programming-tour-book/cache"
	"github.com/go-programming-tour-book/cache/lru"
	"github.com/matryer/is"
	"log"
	"sync"
	"testing"
)

func TestTourCacheGet(t *testing.T) {
	// 模拟数据库
	db := map[string]string{
		"key1": "v1",
		"key2": "v2",
		"key3": "v3",
		"key4": "v4",
	}

	// 将匿名func转换成GetFunc类型，也即实现Getter接口
	getter := cache.GetFunc(func(key string) interface{} {
		log.Println("[From DB] find key", key)
		if val, ok := db[key]; ok {
			return val
		}
		return nil
	})

	tourCache := cache.NewTourCache(getter, lru.New(0, nil))

	is := is.New(t)

	waitGroup := sync.WaitGroup{}

	for k, v := range db {
		waitGroup.Add(1)
		go func(k, v string) {
			defer waitGroup.Done()
			// 未命中
			is.Equal(tourCache.Get(k), v)
			// 命中
			is.Equal(tourCache.Get(k), v)
		}(k, v)
	}

	// 等待所有线程执行结束
	waitGroup.Wait()

	is.Equal(tourCache.Get("unknown"), nil)
	is.Equal(tourCache.Get("unknown"), nil)

	is.Equal(tourCache.Stat().Nget, 10)
	is.Equal(tourCache.Stat().Nhit, 4)
}
