package broker

import (
	"KVCache/test_util"
	"fmt"
	"net"
	"testing"
)

func TestBroker(t *testing.T) {
	protocolMessage := make(chan *ProtocolMessage, 10)
	server := New("127.0.0.1:12345", protocolMessage)
	go server.Start()
	defer server.Stop()

	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	test_util.IsNil(t, err)

	fmt.Fprintf(conn, "get ab \t\n")
	fmt.Fprintf(conn, "get abc 1 2 3\t\n")
	fmt.Fprintf(conn, "set ab \t\n")
	fmt.Fprintf(conn, "somevalue \t\n")

	message := <-protocolMessage
	test_util.VerifyStringValue(t, "get", message.Verb)
	test_util.VerifyStringValue(t, "ab", message.Key)
	test_util.VerifyStringValue(t, "", message.Data)

	message = <-protocolMessage
	test_util.VerifyStringValue(t, "get", message.Verb)
	test_util.VerifyStringValue(t, "abc", message.Key)
	test_util.VerifyStringValue(t, "", message.Data)

	message = <-protocolMessage
	test_util.VerifyStringValue(t, "set", message.Verb)
	test_util.VerifyStringValue(t, "ab", message.Key)
	test_util.VerifyStringValue(t, "somevalue", message.Data)
}
