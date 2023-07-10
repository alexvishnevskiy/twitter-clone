package local

import (
	"container/list"
	"sync"
)

type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	items    *list.List
	mutex    sync.RWMutex
}

type pair struct {
	key   string
	value []byte
}

func New(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		items:    list.New(),
		mutex:    sync.RWMutex{},
	}
}

func (lru *LRUCache) Get(key string) (bool, []byte) {
	lru.mutex.RLock()
	defer lru.mutex.RUnlock()

	if element, ok := lru.cache[key]; ok {
		lru.items.MoveToFront(element)
		return true, element.Value.(*pair).value
	}
	return false, []byte{}
}

func (lru *LRUCache) Put(key string, value []byte) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if element, ok := lru.cache[key]; ok {
		lru.items.MoveToFront(element)
		element.Value.(*pair).value = value
	} else {
		if lru.items.Len() >= lru.capacity {
			back := lru.items.Back()
			delete(lru.cache, back.Value.(*pair).key)
			lru.items.Remove(back)
		}
		pair := &pair{key, value}
		element := lru.items.PushFront(pair)
		lru.cache[key] = element
	}
}
