package lfu

import (
	"github.com/matryer/is"
	"testing"
)

func TestSet(t *testing.T) {
	is := is.New(t)
	lfu := New(24, nil)
	lfu.DelOldest()
	lfu.Set("k1", 1)
	lfu.Set("k2", 2)
	v1 := lfu.Get("k1")
	is.Equal(v1, 1)

	lfu.Del("k1")
	is.Equal(lfu.Len(), 1)
}

func TestOnEvicted(t *testing.T) {
	is := is.New(t)

	k := make([]string, 0, 8)

	onEvicted := func(key string, value interface{}) {
		k = append(k, key)
	}

	lfu := New(32, onEvicted)
	lfu.Set("k1", 1)
	lfu.Set("k2", 2)
	//lfu.Get("k1")
	//lfu.Get("k1")
	//lfu.Get("k1")
	lfu.Set("k3", 3)
	lfu.Set("k4", 4)

	is.Equal(lfu.Len(), 2)
	is.Equal(k, []string{"k1", "k3"})

}
