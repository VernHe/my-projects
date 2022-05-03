package main

import (
	"github.com/allegro/bigcache"
	"log"
	"time"
)

func main1() {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		log.Println("new", err)
		return
	}

	bytes, err := cache.Get("my-key")
	if err != nil {
		log.Println("get", err)
	}

	if bytes == nil {
		log.Println("缓存中未找到，从数据库查找")
		bytes = []byte("my-val")
		cache.Set("my-key", bytes)
	}

	bytes, err = cache.Get("my-key")
	log.Println("再次查询:", string(bytes))

}
