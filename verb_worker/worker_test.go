package verb_worker

import (
	"KVCache/cache"
	"KVCache/test_util"
	"testing"
)

func TestWorker(t *testing.T) {
	localCache := cache.NewCache(cache.LRU)
	localCache.Init(4)

	getWorker := New(Get, localCache)
	putWorker := New(Set, localCache)

	reques := Request{
		Key:   1,
		Value: "aa",
	}
	result := getWorker.Process(reques)
	test_util.VerifyStringValue(t, _notFound+_tab+_end, result)

	reques = Request{
		Key:   1,
		Value: "aa",
	}
	result = putWorker.Process(reques)
	test_util.VerifyStringValue(t, _success, result)

	result = getWorker.Process(reques)
	test_util.VerifyStringValue(t, "aa"+_tab+_end, result)
}
