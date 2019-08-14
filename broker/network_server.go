// Package broker provides functionality for listening incoming message from network
package broker

import (
	"log"
	"net"
	"sync"
)

type ProtocolMessage struct {
	Key          string
	Data         string
	Verb         string
	ClientBroker ClientBroker
}

// NetworkServer is the interface wraps a network server.
// The network server could be a TCP server listen to a port,
// or a UDP server
type NetworkServer interface {
	// Start starts the server and begin to listen for new connection request
	Start()

	//Stop stops the all open connections
	Stop()
}

// tcpServer is a tcp server wrapper and listen to a port
type tcpServer struct {
	address  string // local address to listen to
	listener net.Listener
	messages chan *ProtocolMessage // messages holds all incoming message from network

	clients     map[string]ClientBroker //clients maintains all open collections
	clientMutex sync.Mutex

	//TODO use sync/atomi
	closed     bool // flag for shutting down the server
	closeMutex sync.Mutex
}

// Start starts the server and begin to listen for new connection request
func (s *tcpServer) Start() {
	var err error
	s.listener, err = net.Listen("tcp", s.address)

	if err != nil {
		log.Fatal("Error starting TCP server.")
	}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			{
				s.clientMutex.Lock()
				defer s.clientMutex.Unlock()

				if s.closed {
					return
				}
			}

			log.Println(err)
			continue
		}

		remoteAddr := conn.RemoteAddr().String()
		client := NewClientBroker(conn, remoteAddr, s)
		s.clientMutex.Lock()
		s.clients[remoteAddr] = client
		s.clientMutex.Unlock()

		// TODO: currently unlimitted connection number.
		// TODO: we should limit the connection number
		go client.Start(s.messages)
	}
}

// Close a remote connection and remove it
func (s *tcpServer) RemoveClient(address string) {
	s.clientMutex.Lock()
	_, ok := s.clients[address]
	if ok {
		s.clients[address].Close()
		delete(s.clients, address)
	}
	s.clientMutex.Unlock()
}

// stop the all open connections
func (s *tcpServer) Stop() {
	s.clientMutex.Lock()
	s.closed = true
	s.listener.Close()
	s.clientMutex.Unlock()

	s.clientMutex.Lock()
	s.clientMutex.Unlock()

	for k, v := range s.clients {
		v.Close()
		delete(s.clients, k)
	}
}

// Creates new tcp server instance
// address is the port the tcp server listens to
// protocolMessage is the message queue to holds all incoming message from network.
func New(address string, protocolMessage chan *ProtocolMessage) NetworkServer {
	log.Println("Creating TCP server in address: ", address)
	server := &tcpServer{
		address:  address,
		messages: protocolMessage,
		clients:  make(map[string]ClientBroker),
		closed:   false,
	}
	return server
}
