package lfu

import "github.com/go-programming-tour-book/cache"

type entry struct {
	key    string
	value  interface{}
	weight int
	index  int
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value) + 4 + 4
}

type queue []*entry

func (q queue) Len() int {
	return len(q)
}

// Less 最小堆
func (q queue) Less(i, j int) bool {
	return q[i].weight < q[j].weight
}

func (q queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = j
	q[j].index = i
}

func (q *queue) Push(v interface{}) {
	en := v.(*entry)
	en.index = len(*q)
	// 更新原slice
	*q = append(*q, en)
}

func (q *queue) Pop() interface{} {
	// 旧slice
	old := *q
	n := len(old)
	// 删除使用最少的数据
	en := old[n-1]
	old[n-1] = nil
	// 以防万一
	en.index = -1
	// 更新原slice
	*q = old[:n-1]
	return en
}
