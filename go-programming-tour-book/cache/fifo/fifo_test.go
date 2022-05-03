package fifo

import (
	"github.com/matryer/is"
	"testing"
)

func TestGetSet(t *testing.T) {
	is := is.New(t)
	cache := New(24, nil)
	cache.DelOldest()
	cache.Set("k1", 1)
	v1 := cache.Get("k1")
	is.Equal(v1, 1)

	cache.Del("k1")
	is.Equal(cache.Len(), 0)
}

func TestOnEvicted(t *testing.T) {
	is := is.New(t)

	oldKeys := make(map[string]interface{})

	onEvicted := func(key string, value interface{}) {
		oldKeys[key] = value
	}

	fifo := New(16, onEvicted)

	fifo.Set("k1", 1)
	fifo.Set("k2", 2)
	fifo.Set("k3", 3)
	fifo.Set("k4", 4)

	is.Equal(fifo.Len(), 2)
	is.Equal(len(oldKeys), 2)
	is.Equal(oldKeys["k1"], 1)
	is.Equal(oldKeys["k2"], 2)
	is.Equal(oldKeys["k3"], nil)
}
