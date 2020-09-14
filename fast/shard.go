package fast

import (
	"container/list"
	"sync"
)

type cacheShard struct {
	locker sync.RWMutex

	maxEntries int
	onEvicted  func(key string, value interface{})

	ll    *list.List
	cache map[string]*list.Element
}

type entry struct {
	key   string
	value interface{}
}

func newCacheShard(maxEntries int, onEvicted func(key string, value interface{})) *cacheShard {
	return &cacheShard{
		maxEntries: maxEntries,
		onEvicted:  onEvicted,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
	}
}

func (f *cacheShard) set(key string, value interface{}) {
	f.locker.Lock()
	defer f.locker.Unlock()

	if e, ok := f.cache[key]; ok {
		f.ll.MoveToBack(e)
		en := e.Value.(*entry)
		en.value = value
		return
	}

	en := &entry{key, value}
	e := f.ll.PushBack(en)
	f.cache[key] = e

	if f.maxEntries > 0 && f.ll.Len() > f.maxEntries {
		f.removeElement(f.ll.Front())
	}
}

func (f *cacheShard) get(key string) interface{} {
	f.locker.RLock()
	defer f.locker.RUnlock()

	if e, ok := f.cache[key]; ok {
		f.ll.MoveToBack(e)
		return e.Value.(*entry).value
	}
	return nil
}

func (f *cacheShard) del(key string) {
	f.locker.Lock()
	defer f.locker.Unlock()

	if e, ok := f.cache[key]; ok {
		f.removeElement(e)
	}
}

func (f *cacheShard) delOldest() {
	f.locker.Lock()
	defer f.locker.Unlock()

	f.removeElement(f.ll.Front())
}

func (f *cacheShard) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	f.ll.Remove(e)
	en := e.Value.(*entry)
	delete(f.cache, en.key)
	if f.onEvicted != nil {
		f.onEvicted(en.key, en.value)
	}
}

func (f *cacheShard) len() int {
	f.locker.RLock()
	defer f.locker.RUnlock()
	return f.ll.Len()
}
