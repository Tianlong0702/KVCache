// Package cache provides functionality for a in-memory cache.
package cache

import (
	"container/list"
	"sync"
)

// lru_cache is a thread-safe fixed size LRU cache.
type lru_cache struct {
	capacity int
	linkList list.List             // list keeps track of the order items based on least recently used
	dataMap  map[Key]*list.Element // data map keep all key-value pairs in memory
	mu       sync.RWMutex
}

// Initialize the lru cache with the capacity
func (l *lru_cache) Init(capacity int) {
	l.dataMap = make(map[Key]*list.Element)

	l.capacity = capacity
	l.linkList.Init()
}

// Get the value for a given key.
func (l *lru_cache) Get(k Key) (Value, bool) {
	l.mu.RLock()
	el, found := l.dataMap[k]
	l.mu.RUnlock()
	if !found {
		return nil, false
	}
	en := getEntry(el)
	v := en.value

	// move the entry to the front of the list
	l.mu.Lock()
	l.linkList.MoveToFront(el)
	l.mu.Unlock()
	return v, true
}

// Set the value for a key
func (l *lru_cache) Set(k Key, v Value) {
	l.mu.RLock()
	el, found := l.dataMap[k]
	l.mu.RUnlock()

	if found {
		// Update list element value
		getEntry(el).value = v

		// move the entry to the front of the list
		l.mu.Lock()
		l.linkList.MoveToFront(el)
		l.mu.Unlock()
		return
	}
	en := &entry{
		key:   k,
		value: v,
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.capacity <= 0 || l.linkList.Len() < l.capacity {
		// Add this entry
		el = l.linkList.PushFront(en)
		l.dataMap[en.key] = el
		return
	}

	// Replace with the last one
	el = l.linkList.Back()
	if el == nil {
		// Can happen if cap is zero
		return
	}
	remEn := getEntry(el)
	el.Value = en
	l.linkList.MoveToFront(el)

	delete(l.dataMap, remEn.key)
	l.dataMap[en.key] = el
	return
}

func getEntry(el *list.Element) *entry {
	return el.Value.(*entry)
}

func (l *lru_cache) Close() {
}
