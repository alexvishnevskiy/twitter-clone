package local

import (
	"container/list"
	"errors"
	"fmt"
	cc "github.com/alexvishnevskiy/twitter-clone/internal/cache"
	"strings"
	"sync"
)

type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	trie     *Trie
	items    *list.List
	mutex    sync.Mutex
}

type pair struct {
	key   string
	value []byte
}

var ErrorSplitKey = errors.New("failed to split key")

func splitKey(key string) (error, string, string) {
	parts := strings.Split(key, "_tweet_id_")
	if len(parts) != 2 {
		return ErrorSplitKey, "", ""
	}

	parts[1] = "tweet_id_" + parts[1]
	return nil, parts[0], parts[1]
}

func mergeKeys(key1 string, key2 string) string {
	return fmt.Sprintf("%s_%s", key1, key2)
}

func New(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		items:    list.New(),
		trie:     NewTrie(),
		mutex:    sync.Mutex{},
	}
}

func (lru *LRUCache) Get(key string) (bool, []byte) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if element, ok := lru.cache[key]; ok {
		lru.items.MoveToFront(element)
		return true, element.Value.(*pair).value
	}
	return false, []byte{}
}

func (lru *LRUCache) Put(key string, value []byte) error {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	err, key1, key2 := splitKey(key)
	if err == nil {
		lru.trie.Insert(key1, key2)
	}

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
	return nil
}

func (lru *LRUCache) Remove(key string) error {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if element, ok := lru.cache[key]; ok {
		delete(lru.cache, element.Value.(*pair).key)
		lru.items.Remove(element)
		err, key1, key2 := splitKey(key)
		if err == nil {
			lru.trie.Delete(key1, key2)
		}
	}
	return nil
}

// get by first word
func (lru *LRUCache) StartsWith(key string) (error, [][]byte) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if lru.trie == nil {
		return nil, [][]byte{}
	}
	//var result [][]byte
	second_words := lru.trie.StartsWith(key)
	result := make([][]byte, len(second_words))
	for i, second_word := range second_words {
		ok, res := lru.Get(mergeKeys(key, second_word))
		if ok {
			result[i] = res
		} else {
			return cc.CacheError, [][]byte{}
		}
	}
	return nil, result
}
