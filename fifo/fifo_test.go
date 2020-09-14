package fifo

import (
	"testing"

	"github.com/matryer/is"
)

func TestSetGet(t *testing.T) {
	is := is.New(t)

	cache := New(24, nil)

	cache.DelOldest()

	cache.Set("k1", 1)
	v := cache.Get("k1")
	is.Equal(v, 1) // expect to be the same

	cache.Del("k1")
	is.Equal(0, cache.Len()) // expect to be the same
}
