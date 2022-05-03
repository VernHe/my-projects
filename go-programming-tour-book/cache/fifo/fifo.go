package fifo

import (
	"container/list"
	"github.com/go-programming-tour-book/cache"
)

type fifo struct {
	// 最大容量(字节Byte)
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

// New 创建一个fifo
func New(maxBytes int, onEvicted func(key string, value interface{})) *fifo {
	return &fifo{
		maxBytes:  maxBytes,
		usedBytes: 0,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}

// Set 向fifo中添加一个entry
func (f *fifo) Set(key string, value interface{}) {
	// k-v已存在
	if element, ok := f.cache[key]; ok {
		// 移动到末尾
		f.ll.MoveToBack(element)
		entry := element.Value.(*entry)
		// 重新计算使用的容量(只计算Value的，Key的不计算)
		f.usedBytes = f.usedBytes - entry.Len() + cache.CalcLen(value)
		// 修改value
		entry.value = value
		return
	}

	newEntry := &entry{
		key:   key,
		value: value,
	}
	// 添加到链表末尾
	element := f.ll.PushBack(newEntry)
	// 修改容量
	f.usedBytes += newEntry.Len()
	// 往map中添加
	f.cache[key] = element
	// 检查是否需要删除旧数据
	if f.maxBytes > 0 && f.usedBytes > f.maxBytes {
		f.DelOldest()
	}
}

// Get 根据key获取一个元素
func (f *fifo) Get(key string) interface{} {
	// 检查map，看是否存在
	if element, ok := f.cache[key]; ok {
		return element.Value.(*entry).value
	}
	// 不存在
	return nil
}

// Del 根据key删除元素
func (f fifo) Del(key string) {
	// 存在
	if element, ok := f.cache[key]; ok {
		f.removeElement(element)
	}
}

// DelOldest 删除旧数据
func (f *fifo) DelOldest() {
	// 移除队首元素
	f.removeElement(f.ll.Front())
}

func (f *fifo) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	// 从List中删除
	en := f.ll.Remove(e).(*entry)
	// 更新使用容量
	f.usedBytes -= en.Len()
	// 从Map中删除
	delete(f.cache, en.key)
	// 回调函数
	if f.onEvicted != nil {
		f.onEvicted(en.key, en.value)
	}
}

// Len 返回Element个数（用于测试）
func (f fifo) Len() int {
	return f.ll.Len()
}
