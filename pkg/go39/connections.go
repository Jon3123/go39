package go39

import (
	"bufio"
	"fmt"
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
	netIO          NetIO
	socket         net.Listener
	netConnection  net.Conn
	connectionType ConnectionType
	connectionID   string
	connections    map[string]*Connection
}

//TCPListen - listen
func (c *Connection) TCPListen(host string, port int) {
	log.SetLevel(logrus.TraceLevel)

	l, err := net.Listen("tcp", host+":"+strconv.Itoa(port))

	if err != nil {
		log.Fatalf("Failed to listen error: %s", err)
	}
	c.socket = l
	c.connectionType = TCPServer
	c.connections = make(map[string]*Connection)
	log.Infof("Listening on host, %s , port, %d", host, port)
}

//TCPAccept - accept a connection
func (c *Connection) TCPAccept() (connectionID string) {
	ln, err := c.socket.Accept()

	if err != nil {
		log.Warn("FAILED TO ACCEPT SOCKET")
		return
	}
	log.Info(ln)
	log.Info("New Connection")
	connectionID = c.addConnectionTCP(ln)
	return
}

func (c *Connection) addConnectionTCP(connection net.Conn) (connectionID string) {

	if c.connections == nil {
		log.Fatal("Cannot add socket! the map is nil!")
	}

	connectionID = utils.GenerateConnectionID()
	log.Tracef("Adding socket with id %s", connectionID)
	c.connections[connectionID] = &Connection{
		connectionType: TCPClient,
		connectionID:   connectionID,
		netConnection:  connection,
	}
	return
}

func (c *Connection) getConnection(connectionID string) (connection *Connection, err error) {
	if val, ok := c.connections[connectionID]; ok {
		connection = val
		return
	}
	err = fmt.Errorf("%s connection does not exist", connectionID)
	return
}

//ReceiveMessage TODO
func (c *Connection) ReceiveMessage(connectionID string) (bytesRead int) {
	log.Tracef("Reading from connection with ID %s", connectionID)
	conn, err := c.getConnection(connectionID)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(conn.netConnection)
	conn.netIO.readBuffer.Reset()
	conn.netIO.readBuffer.ReadFrom(reader)
	bytesRead = c.PopByte(connectionID)
	c.PopByte(connectionID)
	return
}

//PopString Readstring
func (c *Connection) PopString(connectionID string) (str string) {
	log.Tracef("Reading string in conn with ID %s ", connectionID)
	conn, err := c.getConnection(connectionID)
	if err != nil {
		log.Warnf("Failed to read string from conn with ID %s", connectionID)
		return
	}

	str = conn.netIO.PopString()
	return
}

//PopInt Readint
func (c *Connection) PopInt(connectionID string) (val int) {
	log.Tracef("Reading int in conn with ID %s ", connectionID)
	conn, err := c.getConnection(connectionID)
	if err != nil {
		log.Warnf("Failed to read int from conn with ID %s", connectionID)
		return
	}

	val = conn.netIO.PopInt()
	return
}

//PopByte readbyte
func (c *Connection) PopByte(connectionID string) (val int) {
	log.Tracef("Reading byte in conn with ID %s ", connectionID)
	conn, err := c.getConnection(connectionID)
	if err != nil {
		log.Warnf("Failed to read byte from conn with ID %s", connectionID)
		return
	}

	val = conn.netIO.PopByte()
	return
}

//PushString TODO
func (c *Connection) PushString(connectionID string, str string) {
	log.Tracef("Pushing string %s to %s buffers", str, connectionID)
	conn, err := c.getConnection(connectionID)
	if err != nil {
		log.Warnf("Failed to write string conn with ID %s", connectionID)
		return
	}

	conn.netIO.PushString(str)
}

//PushInt TODO
func (c *Connection) PushInt(connectionID string, val int) {
	log.Tracef("Pushing int %d to %s buffers", val, connectionID)
	conn, err := c.getConnection(connectionID)
	if err != nil {
		log.Warnf("Failed to write int conn with ID %s", connectionID)
		return
	}

	conn.netIO.PushInt(val)
}

//PushByte TODO
func (c *Connection) PushByte(connectionID string, val int) {
	log.Tracef("Pushing byte %d to %s buffers", val, connectionID)
	conn, err := c.getConnection(connectionID)
	if err != nil {
		log.Warnf("Failed to write byte conn with ID %s", connectionID)
		return
	}

	conn.netIO.PushByte(val)
}

//ClearWriteBuffer TODO
func (c *Connection) ClearWriteBuffer(connectionID string) {
	log.Tracef("Clearing %s write buffers", connectionID)
	conn, err := c.getConnection(connectionID)
	if err != nil {
		log.Warnf("Failed to clear writebuffer conn with ID %s", connectionID)
		return
	}

	conn.netIO.ClearWriteBuffer()
}
