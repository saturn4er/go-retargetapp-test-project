package tcpWordsReceiver

import (
	"net"
	"fmt"
	"errors"
	"log"
)

type server struct {
	port              int
	started           bool
	onMessageCallback func(string)
}

func (this *server) Start() error {
	if this.started {
		return errors.New("Receiver already started")
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", this.port))
	if err != nil {
		return err
	}
	this.started = true
	log.Printf("[ TCP WORDS RECEIVER ] Started on %d port", this.port)
	go this.clientsWaitLoop(l)
	return nil
}
func (this *server) clientsWaitLoop(listener net.Listener) {
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[ TCP WORDS RECEIVER ] Error accepting connection: %v", err)
			continue
		}
		log.Printf("[ TCP WORDS RECEIVER ] New connection from: %v", conn.RemoteAddr())
		// Handle connections in a new goroutine.
		client := createClient(conn, this)
		go client.waitMessageLoop()
	}
}
func (this *server) OnMessage(cb func(string)) {
	this.onMessageCallback = cb
}
func (this *server) Stop() error {
	return nil
}
func NewTCPWordsReceiver(port int) *server {
	result := new(server)
	result.port = port
	return result
}