package go39

import (
	"bufio"
	"bytes"
	"net"
	"strconv"

	"github.com/Jon3123/go39/pkg/utils"

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
	sockets        map[string]*net.Conn
}

//TCPListen - listen
func (c *Connection) TCPListen(host string, port int) {
	l, err := net.Listen("tcp", host+":"+strconv.Itoa(port))

	if err != nil {
		log.Info(err)
		log.Fatal("Failed to listen")
	}
	c.socket = l
	c.connectionType = TCPServer
	c.sockets = make(map[string]*net.Conn)
	log.Infof("Listening on host, %s , port, %d", host, port)
}

//TCPAccept - accept a connection
func (c *Connection) TCPAccept() (socketID string) {
	ln, err := c.socket.Accept()

	if err != nil {
		log.Warn("FAILED TO ACCEPT SOCKET")
		return
	}
	log.Info(ln)
	log.Info("New Connection")
	socketID = c.addSocket(&ln)
	return
}

func (c *Connection) addSocket(socket *net.Conn) (socketID string) {

	if c.sockets == nil {
		log.Fatal("Cannot add socket! the map is nil!")
	}

	socketID = utils.GenerateSocketID()
	log.Tracef("Adding socket with id %s", socketID)
	c.sockets[socketID] = socket
	return
}

//ReceiveMessage TODO
func (c *Connection) ReceiveMessage(socketID string) {
	log.Infof("Reading from socket with ID %s", socketID)
	socket := c.sockets[socketID]
	reader := bufio.NewReader(*socket)
	c.readBuffer.ReadFrom(reader)
}
