package lfu

import (
	"container/heap"
	"github.com/go-programming-tour-book/cache"
)

type lfu struct {
	// 缓存最大的容量，单位字节；
	// groupcache 使用的是最大存放 entry 个数
	maxBytes int
	// 当一个 entry 从缓存中移除是调用该回调函数，默认为 nil
	// groupcache 中的 key 是任意的可比较类型；value 是 interface{}
	onEvicted func(key string, value interface{})

	// 已使用的字节数，只包括值，key 不算
	usedBytes int

	queue *queue
	cache map[string]*entry
}

func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	q := make(queue, 0, 1024)
	l := &lfu{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		usedBytes: 0,
		queue:     &q,
		cache:     make(map[string]*entry),
	}
	return l
}

func (l *lfu) Set(key string, value interface{}) {
	if en, ok := l.cache[key]; ok {
		// 更新usedBytes
		l.usedBytes = l.usedBytes - cache.CalcLen(en.value) + cache.CalcLen(value)
		// 更新queue
		l.queue.update(en, value, en.weight+1)
		return
	}

	// 添加新的entry
	e := &entry{
		key:   key,
		value: value,
	}
	// 添加到堆中
	heap.Push(l.queue, e)
	// 添加到map中
	l.cache[key] = e
	// 更新usedBytes
	l.usedBytes += e.Len()
	// 检查是否需要清理
	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.removeElement(heap.Pop(l.queue).(*entry))
	}
}

func (l *lfu) Get(key string) interface{} {
	var v interface{}
	if en, ok := l.cache[key]; ok {
		v = en.value
		l.queue.update(en, en.value, en.weight+1)
	}
	return v
}

func (l *lfu) Del(key string) {
	if en, ok := l.cache[key]; ok {
		// 从map中删除，从queue中删除，更新usedBytes
		l.removeElement(heap.Remove(l.queue, en.index))
	}
}

func (l *lfu) DelOldest() {
	if l.queue.Len() == 0 {
		return
	}
	l.removeElement(heap.Pop(l.queue))
}

func (l *lfu) Len() int {
	return l.queue.Len()
}

func (q *queue) update(en *entry, value interface{}, weight int) {
	en.value = value
	en.weight = weight
	// 更新堆
	heap.Fix(q, en.index)
}

func (l *lfu) removeElement(v interface{}) {
	if v == nil {
		return
	}

	en := v.(*entry)
	delete(l.cache, en.key)
	l.usedBytes -= cache.CalcLen(en.value)

	if l.onEvicted != nil {
		l.onEvicted(en.key, en.value)
	}
}
