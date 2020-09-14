package lfu

import (
	"testing"

	"github.com/matryer/is"
)

func TestOnEvicted(t *testing.T) {
	is := is.New(t)

	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}

	cache := New(32, onEvicted)

	cache.Set("k1", 1)
	cache.Set("k2", 2)
	// cache.Get("k1")
	// cache.Get("k1")
	// cache.Get("k2")
	cache.Set("k3", 3)
	cache.Set("k4", 4)

	expeced := []string{"k1", "k3"}
	is.Equal(expeced, keys) //not expeced
}
