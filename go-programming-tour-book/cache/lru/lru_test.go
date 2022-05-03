package lru

import (
	"github.com/matryer/is"
	"testing"
)

func TestOnEvicted(t *testing.T) {
	is := is.New(t)

	keys := make([]string, 0, 8)

	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}

	lru := New(24, onEvicted)

	// k1 k2 k3
	lru.Set("k1", 1)
	lru.Set("k2", 2)
	lru.Set("k3", 3)
	// k2 k3 k1
	lru.Get("k1")
	// k3 k1 k4
	lru.Set("k4", 4)

	is.Equal(keys, []string{"k2"})
}
