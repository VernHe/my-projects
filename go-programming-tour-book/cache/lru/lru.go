package lru

import (
	"container/list"
	"github.com/go-programming-tour-book/cache"
)

type lru struct {
	// 最大容量(字节Byte)，0表示无限制
	maxBytes int
	// 使用容量(字节Byte)
	usedBytes int

	// 元素被删除时的回调函数
	onEvicted func(key string, value interface{})

	// 链表(存储元素，保证其顺序)
	ll *list.List
	// Map用于查找
	cache map[string]*list.Element
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

func (l *lru) Set(key string, value interface{}) {
	// k-v已存在
	if element, ok := l.cache[key]; ok {
		// 移动到末尾
		l.ll.MoveToBack(element)
		entry := element.Value.(*entry)
		// 重新计算使用的容量(只计算Value的，Key的不计算)
		l.usedBytes = l.usedBytes - entry.Len() + cache.CalcLen(value)
		// 修改value
		entry.value = value
		return
	}

	newEntry := &entry{
		key:   key,
		value: value,
	}
	// 添加到链表末尾
	element := l.ll.PushBack(newEntry)
	// 修改容量
	l.usedBytes += newEntry.Len()
	// 往map中添加
	l.cache[key] = element
	// 检查是否需要删除旧数据
	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.DelOldest()
	}
}

func (l *lru) Get(key string) interface{} {
	// 检查map，看是否存在
	if element, ok := l.cache[key]; ok {
		l.ll.MoveToBack(element)
		return element.Value.(*entry).value
	}
	// 不存在
	return nil
}

func (l *lru) Del(key string) {
	// 存在
	if element, ok := l.cache[key]; ok {
		l.removeElement(element)
	}
}

func (l *lru) DelOldest() {
	// 移除队首元素
	l.removeElement(l.ll.Front())
}

func (l *lru) Len() int {
	return l.ll.Len()
}

func (l *lru) removeElement(e *list.Element) {
	if e == nil {
		return
	}
	// 从List中删除
	en := l.ll.Remove(e).(*entry)
	// 更新使用容量
	l.usedBytes -= en.Len()
	// 从Map中删除
	delete(l.cache, en.key)
	// 回调函数
	if l.onEvicted != nil {
		l.onEvicted(en.key, en.value)
	}
}

// New 创建一个fifo
func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	return &lru{
		maxBytes:  maxBytes,
		usedBytes: 0,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}
