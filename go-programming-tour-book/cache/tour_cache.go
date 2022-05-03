package cache

// Getter 从数据源获取数据
type Getter interface {
	Get(key string) interface{}
}

// GetFunc 为Getter接口提供一个默认的实现，类似于http.Handler与http.HandlerFunc
type GetFunc func(key string) interface{}

// Get 这样任意一个函数，只要签名和 Get(key string) interface{} 一致，通过转为 GetFunc 类型，就实现了 Getter 接口。
func (f GetFunc) Get(key string) interface{} {
	// 默认实现
	return f(key)
}

type TourCache struct {
	mainCache *safeCache
	getter    Getter
}

func NewTourCache(getter Getter, cache Cache) *TourCache {
	return &TourCache{
		mainCache: newSafeCache(cache),
		getter:    getter,
	}
}

func (t TourCache) Get(key string) interface{} {
	val := t.mainCache.get(key)
	if val != nil {
		return val
	}

	if t.getter != nil {
		val = t.getter.Get(key)
		if val != nil {
			t.mainCache.set(key, val)
			return val
		}
		return nil
	}

	return nil
}

func (t TourCache) Set(key string, value interface{}) {
	if value == nil {
		return
	}
	t.mainCache.set(key, value)
}

func (t TourCache) Stat() *Stat {
	return t.mainCache.stat()
}
