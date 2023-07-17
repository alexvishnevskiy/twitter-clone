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

// split user_id_tweet_id keys
func splitKey(key string) (error, string, string) {
	parts := strings.Split(key, "_tweet_id_")
	if len(parts) != 2 {
		return ErrorSplitKey, "", ""
	}

	parts[1] = "tweet_id_" + parts[1]
	return nil, parts[0], parts[1]
}

// merge user_id_tweet_id keys
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

// get value for the key
func (lru *LRUCache) Get(key string) (bool, []byte) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if element, ok := lru.cache[key]; ok {
		lru.items.MoveToFront(element)
		return true, element.Value.(*pair).value
	}
	return false, []byte{}
}

// put value for specific key
func (lru *LRUCache) Put(key string, value []byte) error {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	// if the key is user_id_tweet_id
	err, key1, key2 := splitKey(key)
	if err == nil {
		lru.trie.Insert(key1, key2)
	}

	// lru logic
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

// remove key from cache
func (lru *LRUCache) Remove(key string) error {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if element, ok := lru.cache[key]; ok {
		// delete from dict, list
		delete(lru.cache, element.Value.(*pair).key)
		lru.items.Remove(element)
		// delete from trie if key is user_id_tweet_id
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
	// get all values for keys: first_word*
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
