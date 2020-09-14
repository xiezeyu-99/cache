package cache_test

import (
	"log"
	"sync"
	"testing"

	"github.com/matryer/is"
	"github.com/xiezeyu-99/cache"
	"github.com/xiezeyu-99/cache/lru"
)

func TestTourCacheGet(t *testing.T) {
	db := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
	}
	getter := cache.GetFunc(func(key string) interface{} {
		if val, ok := db[key]; ok {
			log.Println("[From db] find key", key)
			return val
		}
		return nil
	})
	tourCache := cache.NewTourCache(getter, lru.New(0, nil))

	is := is.New(t)

	var wg sync.WaitGroup

	for k, v := range db {
		wg.Add(1)
		go func(k, v string) {
			defer wg.Done()
			is.Equal(tourCache.Get(k), v)
			is.Equal(tourCache.Get(k), v)
		}(k, v)
	}
	wg.Wait()

	is.Equal(tourCache.Get("unknown"), nil)
	is.Equal(tourCache.Get("unknown"), nil)

	is.Equal(tourCache.Stat().NGet, 10)
	is.Equal(tourCache.Stat().NHit, 4)

}
