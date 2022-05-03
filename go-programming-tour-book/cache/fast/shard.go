package fast

import (
	"container/list"
	"github.com/go-programming-tour-book/cache"
	"sync"
)

// cacheShard 分片的实现
type cacheShard struct {
	locker sync.RWMutex
	// 最大个数
	maxEntries int
	onEvicted  func(key string, val interface{})

	ll    *list.List
	cache map[string]*list.Element
}

func newCacheShard(maxEntries int, onEvicted func(key string, val interface{})) *cacheShard {
	return &cacheShard{
		maxEntries: maxEntries,
		onEvicted:  onEvicted,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
	}
}

// 实现了Value接口
type entry struct {
	key   string
	value interface{}
}

// Len 计算大小(字节Byte)
func (e *entry) Len() int {
	return cache.CalcLen(e.value)
}

func (c *cacheShard) get(key string) interface{} {
	c.locker.RLock()
	defer c.locker.RUnlock()

	if e, ok := c.cache[key]; ok {
		c.ll.MoveToBack(e)
		return e.Value.(*entry).value
	}

	return nil
}

func (c *cacheShard) set(key string, value interface{}) {
	c.locker.Lock()
	defer c.locker.Unlock()

	// k-v已存在
	if element, ok := c.cache[key]; ok {
		// 移动到末尾
		c.ll.MoveToBack(element)
		entry := element.Value.(*entry)
		// 修改value
		entry.value = value
		return
	}

	newEntry := &entry{
		key:   key,
		value: value,
	}
	// 添加到链表末尾
	element := c.ll.PushBack(newEntry)
	// 往map中添加
	c.cache[key] = element
}

func (c *cacheShard) del(key string) {
	// 检查是否存在此key
	if e, ok := c.cache[key]; ok {
		c.locker.Lock()
		defer c.locker.Unlock()
		// 从ll中删除
		c.ll.Remove(e)
		// 从map中删除
		delete(c.cache, key)
		// 回调函数
		if c.onEvicted != nil {
			c.onEvicted(key, e.Value.(entry).value)
		}
	}
}

func (c *cacheShard) len() int {
	return c.ll.Len()
}
