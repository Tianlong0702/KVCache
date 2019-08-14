package verb_worker

import (
	"KVCache/cache"
)

// Verb types
// currently only supports get and set
const (
	Get = "get"
	Set = "set"
)

// Action request.
// For verb: get, only Key is used, and value is nil.
// For verb: set, both key and value are used.
type Request struct {
	Key   cache.Key
	Value cache.Value
}

// Action interface to dispatch incoming network request
// e.g. "get" => GetVerbWorker
//      "set" => SetVerbWorker
type VerbWorker interface {
	// Process the incoming request and return result.
	Process(Request) string
}

// New returns a verb worker base on the verb type.
func New(verb string, cache cache.KVCache) VerbWorker {
	switch verb {
	case Get:
		return &get_worker{
			cache: cache,
		}
	case Set:
		return &set_worker{
			cache: cache,
		}
	default:
		panic("cache: unsupported " + verb)
	}
}
