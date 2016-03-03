package tcpWordsReceiver

import (
	"net"
	"bufio"
	"io"
	"log"
	"fmt"
)

type client struct {
	connection net.Conn
	server     *server
}

func (this *client) log(message string, params... interface{}) {
	msg_tmplate := fmt.Sprintf("[ TCP WORDS RECEIVER ] [ %v ] %s", this.connection.RemoteAddr(), message)
	log.Printf(msg_tmplate, params...)
}
func (this *client) waitMessageLoop() {
	defer this.connection.Close()
	conn_reader := bufio.NewReader(this.connection)
	for {
		line, _, err := conn_reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				this.log("Client disconnected")
			}else {
				this.log("Error reading from socket: %v", err)
			}
			break
		}
		if (this.server.onMessageCallback != nil) {
			this.server.onMessageCallback(string(line))
		}

	}
}
func createClient(conn net.Conn, server *server) *client {
	result := new(client)
	result.connection = conn
	result.server = server
	return result
}
