package fast

type fastCache struct {
	shards    []*cacheShard
	shardMask uint64
	hash      fnv64a
}

func NewFastCache(maxEntries, shardNum int, onEvicted func(key string, val interface{})) *fastCache {
	fastCache := &fastCache{
		shards: make([]*cacheShard, shardNum),
		// 用于取模得到下标,比如 n%8 >> n&7，注：n必须是2的x次方
		shardMask: uint64(shardNum - 1),
		hash:      newDefaultHasher(),
	}
	// 初始化shards
	for i := 0; i < shardNum; i++ {
		fastCache.shards[i] = newCacheShard(maxEntries, onEvicted)
	}

	return fastCache
}

// getShard 根据key获取对应的shard
func (c *fastCache) getShard(key string) *cacheShard {
	hashedKey := c.hash.Sum64(key)
	return c.shards[hashedKey&c.shardMask]
}

func (c *fastCache) Set(key string, val interface{}) {
	// 算得所在的shards，再set
	c.getShard(key).set(key, val)
}

func (c *fastCache) Del(key string) {
	c.getShard(key).del(key)
}

func (c *fastCache) Len() int {
	length := 0
	// 统计所有shard的总长
	for _, shard := range c.shards {
		length += shard.len()
	}
	return length
}
