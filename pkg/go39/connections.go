package go39

import (
	"bytes"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

//ConnectionType type of connection the connection represents
type ConnectionType int

const (
	//TCPClient this is a client
	TCPClient ConnectionType = 0
	//TCPServer this is a server
	TCPServer ConnectionType = 1
)

//Connection to fill out later
type Connection struct {
	writeBuffer    bytes.Buffer
	readBuffer     bytes.Buffer
	socket         net.Listener
	connectionType ConnectionType
}

//TCPListen - listen
func (c *Connection) TCPListen(host string, port int) {
	l, err := net.Listen("tcp", host+":"+string(port))

	if err != nil {
		fmt.Print("ERROR")
	}
	c.socket = l
	c.connectionType = TCPServer
}
