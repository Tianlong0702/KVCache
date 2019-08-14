// Package cache provides functionality for a in-memory cache.
package cache

import (
	"testing"

	"KVCache/test_util"
)

func TestCache(t *testing.T) {
	cache := NewCache(LRU)
	cache.Init(4)

	cache.Set(1, "1")
	cache.Set(2, "2")
	cache.Set(3, "3")
	cache.Set(4, "4")
	val, found := cache.Get(1)
	test_util.VerifyBoolValue(t, true, found)
	test_util.VerifyStringValue(t, "1", val.(string))

	cache.Set(5, "5")

	val, found = cache.Get(2)
	test_util.VerifyBoolValue(t, false, found)
	test_util.IsNil(t, val)

	val, found = cache.Get(1)
	test_util.VerifyBoolValue(t, true, found)
	test_util.VerifyStringValue(t, "1", val.(string))

	val, found = cache.Get(3)
	test_util.VerifyBoolValue(t, true, found)
	test_util.VerifyStringValue(t, "3", val.(string))

	val, found = cache.Get(4)
	test_util.VerifyBoolValue(t, true, found)
	test_util.VerifyStringValue(t, "4", val.(string))

	val, found = cache.Get(5)
	test_util.VerifyBoolValue(t, true, found)
	test_util.VerifyStringValue(t, "5", val.(string))

	cache.Set(3, "33")

	val, found = cache.Get(3)
	test_util.VerifyBoolValue(t, true, found)
	test_util.VerifyStringValue(t, "33", val.(string))

	cache.Set(8, "88")

	val, found = cache.Get(1)
	test_util.VerifyBoolValue(t, false, found)
	test_util.IsNil(t, val)
}
