package server

import (
	"KVCache/broker"
	"KVCache/cache"
	"KVCache/test_util"
	"KVCache/verb_worker"
	"testing"
)

func TestServer(t *testing.T) {
	messageChan := make(chan *broker.ProtocolMessage, 10)
	testServer := KVCacheServer{
		Messages:      messageChan,
		VerbWorkers:   make(map[string]verb_worker.VerbWorker),
		NetworkServer: &TestNetWorkServer{},
		close:         make(chan struct{}),
		localCache:    &TestCache{},
	}
	testServer.VerbWorkers[verb_worker.Get] = &TestGetworker{}
	testServer.VerbWorkers[verb_worker.Set] = &TestSetworker{}
	testServer.Start()
	defer testServer.Stop()

	testClientBroker1 := &TestClientBroker{
		Data: make(chan string, 1),
	}
	testClientBroker2 := &TestClientBroker{
		Data: make(chan string, 1),
	}
	message1 := &broker.ProtocolMessage{
		Key:          "key1",
		Data:         "",
		ClientBroker: testClientBroker1,
		Verb:         "get",
	}

	message2 := &broker.ProtocolMessage{
		Key:          "key2",
		Data:         "data2",
		ClientBroker: testClientBroker2,
		Verb:         "set",
	}
	messageChan <- message1
	messageChan <- message2

	data1 := <-testClientBroker1.Data
	test_util.VerifyStringValue(t, "get key1 ", data1)

	data2 := <-testClientBroker2.Data
	test_util.VerifyStringValue(t, "set key2 data2", data2)
}

type TestNetWorkServer struct {
}

func (s *TestNetWorkServer) Start() {}
func (s *TestNetWorkServer) Stop()  {}

type TestSetworker struct {
}

func (pw *TestSetworker) Process(request verb_worker.Request) string {

	return "set " + request.Key.(string) + " " + request.Value.(string)
}

type TestGetworker struct {
	cache cache.KVCache
}

func (gw *TestGetworker) Process(request verb_worker.Request) string {

	return "get " + request.Key.(string) + " " + request.Value.(string)
}

type TestClientBroker struct {
	Data chan string
}

func (cb *TestClientBroker) Start(protocolMessageChan chan *broker.ProtocolMessage) {

}
func (cb *TestClientBroker) Send(message string) error {
	cb.Data <- message
	return nil
}
func (cb *TestClientBroker) Close() error {
	return nil
}

type TestCache struct {
}

func (c *TestCache) Get(key cache.Key) (cache.Value, bool) {
	return nil, true
}

func (c *TestCache) Set(cache.Key, cache.Value) {
}

func (c *TestCache) Close() {
}

func (c *TestCache) Init(cap int) {
}
