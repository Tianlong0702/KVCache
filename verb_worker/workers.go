package verb_worker

import (
	"KVCache/cache"
)

const (
	_success  = "Success\t\n"
	_notFound = "Not Found"
	_end      = "End\t\n"
	_tab      = "\t\n"
)

// set_worker is the worker for verb: set
type set_worker struct {
	cache cache.KVCache
}

// Process calls local cache to set the key-value pair
func (sw *set_worker) Process(request Request) string {
	sw.cache.Set(request.Key, request.Value)
	return _success
}

// get_worker is the worker for verb: get
type get_worker struct {
	cache cache.KVCache
}

// Process calls local cache to get the value and return the result
func (gw *get_worker) Process(request Request) string {
	val, found := gw.cache.Get(request.Key)
	message := _tab + _end
	if found {
		return val.(string) + message
	} else {
		return _notFound + message
	}
}
