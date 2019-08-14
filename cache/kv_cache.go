// Package cache provides functionality for a in-memory cache.
package cache

import (
	"time"
)

// Cache type based on eviction policy
// LRU: When the cache reached its capacity,
//      it should invalidate the least recently used item before inserting a new item.
// SLRU: a segmented LRU.
//       See http://highscalability.com/blog/2016/1/25/design-of-a-modern-cache.html
const (
	LRU = "lru"

	// TODO: implemet slru
	SLRU = "slru"
)

// Key is any value.
type Key interface{}

// Value is any value.
type Value interface{}

// KVCache is the interface for a in-memory key-value cache.
// Methods Get and Set may be called concurrently from multiple goroutines,
type KVCache interface {
	// Get the value for a given the key and flag for found or not
	// if cache miss, return (value, true), otherwith return (nil, false)
	Get(Key) (Value, bool)

	// Set or insert the value if the key is not already present.
	// When the cache reached its capacity, it should invalidate the item based on the eviction policy.
	Set(Key, Value)

	// Close key-value cache.
	Close()

	// Initialize the in-memory key-value cache with capacity.
	Init(int)
}

// key-value entry in the key-value cache
type entry struct {
	key   Key
	value Value

	// Accessed is the last time this entry was accessed.
	// Current Not Used
	accessTS time.Time
	// Updated is the last time this entry was updated.
	// Current Not Used
	updatedTS time.Time
}

// Creates new key-value cache instance according to cache type.
func NewCache(name string) KVCache {
	switch name {
	//TODO: implememnt segment lru
	case LRU:
		return &lru_cache{}
	default:
		panic("cache: unsupported " + name)
	}
}
