package server

import (
	"log"

	"KVCache/broker"
	"KVCache/cache"
	"KVCache/verb_worker"
)

const (
	// default message queue size for holding incoming request from network
	_defaultProtocolMessageCapacity = 500
)

// KVCacheServer is the server for key value cache service.
// The server receivers request from network, then push the request
// in the back of he messages queue. The server reads requests from
// the front of the message queue and process them asynchronously.
type KVCacheServer struct {
	// Messages: request queue for holding all incoming request
	//TODO: Currently only one queue global, could be performance bottleleck
	//TODO: instead of using a queue globally,
	//TODO: We could use an array of queue to boost performance: [10]chan *broker.ProtocolMessage
	//TODO: when there's an incoming request, ramdomly select one queue to push.
	Messages chan *broker.ProtocolMessage

	// VerbWorkers: verb workers. Dispatch incoming request to different works
	VerbWorkers map[string]verb_worker.VerbWorker

	// Network server listens to new coming request and push it to the back of the message queue
	NetworkServer broker.NetworkServer

	//localCache: local cache engine
	localCache cache.KVCache

	// close: flag for server shuts down
	close chan struct{}
}

// Processing a request message.
// It calls the verb worker base on the verb.
// Gets the result from the verb woker and send to client
func (kv *KVCacheServer) processingMessage(message *broker.ProtocolMessage) {
	worker, found := kv.VerbWorkers[message.Verb]

	if found {
		request := verb_worker.Request{
			Key:   message.Key,
			Value: message.Data,
		}
		result := worker.Process(request)

		// TODO: add back off retry logic
		if err := message.ClientBroker.Send(result); err != nil {
			log.Println(err)
		}
	} else {
		log.Println("Unknow verb: ", message.Verb)
	}
}

// Read message request from request queue and call processMessage function asynchronously
func (kv *KVCacheServer) processingRequest() {
	for {
		select {
		case <-kv.close:
			// Server shut down
			return
		case protolMessage := <-kv.Messages:
			// process incoming request asynchronously.
			go kv.processingMessage(protolMessage)
		}
	}
}

// Start the server to receive incoming request
func (kv *KVCacheServer) Start() {
	go kv.processingRequest()
	kv.NetworkServer.Start()
}

// Stop the server
func (kv *KVCacheServer) Stop() {
	log.Println("Server is shutting down")

	// Stopping the processing goroutine
	kv.close <- struct{}{}

	// Stopping network server
	kv.NetworkServer.Stop()

	// Closing local cache
	kv.localCache.Close()
}

// Add a verbWorker for a verb
func (kv *KVCacheServer) AddVerWorker(verb string, worker verb_worker.VerbWorker) {
	kv.VerbWorkers[verb] = worker
}

// New returns a new KVCacheServer
func New(address string, kvCache cache.KVCache) *KVCacheServer {
	server := &KVCacheServer{
		Messages:    make(chan *broker.ProtocolMessage, _defaultProtocolMessageCapacity),
		VerbWorkers: make(map[string]verb_worker.VerbWorker),
		localCache:  kvCache,
		close:       make(chan struct{}),
	}

	server.VerbWorkers[verb_worker.Get] = verb_worker.New(verb_worker.Get, server.localCache)
	server.VerbWorkers[verb_worker.Set] = verb_worker.New(verb_worker.Set, server.localCache)
	server.NetworkServer = broker.New(address, server.Messages)

	return server
}
