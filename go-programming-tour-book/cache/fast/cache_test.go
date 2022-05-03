package fast

import (
	"github.com/spf13/cast"
	"math/rand"
	"testing"
	"time"
)

var (
	maxEntrySize = 1024
)

func BenchmarkTourFastCacheSetParallel(b *testing.B) {
	cache := NewFastCache(b.N, maxEntrySize, nil)
	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		counter := 0
		for pb.Next() {
			cache.Set(parallelKey(id, counter), value())
			counter++
		}
	})
}

func parallelKey(id, counter int) string {
	return cast.ToString(id + counter)
}

func value() interface{} {
	return ""
}
