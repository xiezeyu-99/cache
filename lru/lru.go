package lru

import (
	"container/list"

	"github.com/xiezeyu-99/cache"
)

type lru struct {
	//缓存最大容量
	maxBytes  int
	onEvicted func(key string, value interface{})
	//已使用字节数，key不算
	usedBytes int
	ll        *list.List
	cache     map[string]*list.Element
}

type entry struct {
	key   string
	value interface{}
}

func (e *entry) Len() int {
	return cache.CaclLen(e.value)
}

func New(maxByte int, onEvicted func(key string, value interface{})) cache.Cache {
	return &lru{
		maxBytes:  maxByte,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}

func (f *lru) Set(key string, value interface{}) {
	if e, ok := f.cache[key]; ok {
		f.ll.MoveToBack(e)
		en := e.Value.(*entry)
		f.usedBytes = f.usedBytes - cache.CaclLen(en.value) + cache.CaclLen(value)
		en.value = value
		return
	}

	en := &entry{key, value}
	e := f.ll.PushBack(en)
	f.cache[key] = e

	f.usedBytes += en.Len()
	if f.maxBytes > 0 && f.usedBytes > f.maxBytes {
		f.DelOldest()
	}
}

func (f *lru) Get(key string) interface{} {
	if e, ok := f.cache[key]; ok {
		f.ll.MoveToBack(e)
		return e.Value.(*entry).value
	}
	return nil
}

func (f *lru) Del(key string) {
	if e, ok := f.cache[key]; ok {
		f.removeElement(e)
	}
}

func (f *lru) DelOldest() {
	f.removeElement(f.ll.Front())
}

func (f *lru) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	f.ll.Remove(e)
	en := e.Value.(*entry)
	f.usedBytes -= en.Len()
	delete(f.cache, en.key)
	if f.onEvicted != nil {
		f.onEvicted(en.key, en.value)
	}
}

func (f *lru) Len() int {
	return f.ll.Len()
}
