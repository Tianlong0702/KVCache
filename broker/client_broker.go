// Package broker provides functionality for listening incoming message from network
package broker

import (
	"KVCache/verb_worker"
	"bufio"
	"net"
	"strings"
)

// ClientBroker is the interface wraps a server-client connection.
type ClientBroker interface {
	// Start: start to listen incoming message from the connection
	// after get a message from connection,
	// it push the message to the end protocol message queue.
	Start(protocolMessageChan chan *ProtocolMessage)

	//Send: send message to the client
	Send(message string) error

	// Close: close the collection
	Close() error
}

//Send: send message to the client
func (c *clientBroker) Send(message string) error {
	_, err := c.conn.Write([]byte(message))
	return err
}

// Start: start to listen incoming message from the connection
func (cb *clientBroker) Start(protocolMessageChan chan *ProtocolMessage) {
	reader := bufio.NewReader(cb.conn)
	//setKey hold the key for set command from last message
	setKey := ""
	for {
		// TODO: the implemetation here is a litte ugly.
		// TODO: extra the txt protocal to a component.
		message, err := reader.ReadString('\n')
		if err != nil {
			// TODO: need more elegant way to handle error
			// TODO: probably retry before closing the client
			cb.tcpServer.RemoveClient(cb.remoteAddr)
			return
		}
		message = strings.TrimSpace(message)
		command := strings.Split(message, " ")

		if len(command) < 1 {
			// invalid message
			setKey = ""
			continue
		}

		if setKey != "" {
			// the message is data
			protocolMessage := &ProtocolMessage{
				Key:          setKey,
				Data:         message,
				Verb:         verb_worker.Set,
				ClientBroker: cb,
			}
			setKey = ""
			// Send request to the queue
			protocolMessageChan <- protocolMessage
			continue
		}

		if setKey == "" {
			// the message is a command
			if len(command) < 2 {
				// invalid command
				continue
			}
			if command[0] == verb_worker.Set {
				setKey = strings.TrimSpace(command[1])
				continue
			} else if command[0] == verb_worker.Get {
				protocolMessage := &ProtocolMessage{
					Key:          strings.TrimSpace(command[1]),
					Data:         "",
					Verb:         verb_worker.Get,
					ClientBroker: cb,
				}
				// Send request to the queue
				protocolMessageChan <- protocolMessage
				continue
			}
		}
	}
}

// Close: close the collection
func (cb *clientBroker) Close() error {
	return cb.conn.Close()
}

// a tcp connection wrapper
type clientBroker struct {
	conn       net.Conn   // the connection
	remoteAddr string     // remote address
	tcpServer  *tcpServer // tcp server
}

// NewClientBroker: returns a tcp connection wrapper
func NewClientBroker(conn net.Conn, remoteAddr string, tcpServer *tcpServer) ClientBroker {
	client := &clientBroker{
		conn:       conn,
		tcpServer:  tcpServer,
		remoteAddr: remoteAddr,
	}
	return client
}
