package go39

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"time"

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
	//MaxTransmitSize max count of bytes you are allowed to send
	MaxTransmitSize int = 1024
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

//NewConnection Create a new Connection
func NewConnection() *Connection {
	c := &Connection{
		connectionID: utils.GenerateConnectionID(),
	}
	c.connections = make(map[string]*Connection)
	c.connections[c.connectionID] = c
	return c
}

//TCPListen - listen
func (c *Connection) TCPListen(host string, port int) {
	log.SetLevel(logrus.DebugLevel)

	l, err := net.Listen("tcp", host+":"+strconv.Itoa(port))

	if err != nil {
		log.Fatalf("Failed to listen error: %s", err)
	}
	c.socket = l
	c.connectionType = TCPServer
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

//ReceiveMessage receives message from the connection with the given ID and set a timeout duration return -1 when disconnect
func (c *Connection) ReceiveMessage(connectionID string, timeout time.Duration) (bytesRead int32) {
	log.Tracef("Reading from connection with ID %s", connectionID)
	conn, err := c.getConnection(connectionID)
	if err != nil {
		log.Fatalf("error getting connection: %s", err.Error())
	}
	b := make([]byte, MaxTransmitSize)
	conn.netConnection.SetReadDeadline(time.Now().Add(timeout))
	_, err = conn.netConnection.Read(b)
	if err != nil {
		if err.Error() == "EOF" {
			return -1
		}
		log.Warnf("error reading from connection %s: %s", connectionID, err.Error())
	} else {
		fmt.Println("not nil")
	}
	reader := bytes.NewReader(b)
	conn.netIO.ClearReadBuffer()

	_, err = conn.netIO.readBuffer.ReadFrom(reader)
	if err != nil {
		log.Warnf("error while reading %s: %s", connectionID, err.Error())

	}

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
func (c *Connection) PopInt(connectionID string) (val int32) {
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
func (c *Connection) PopByte(connectionID string) (val int32) {
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

//PushInt write int to buffer
func (c *Connection) PushInt(connectionID string, val int32) {
	log.Tracef("Pushing int %d to %s buffers", val, connectionID)
	conn, err := c.getConnection(connectionID)
	if err != nil {
		log.Warnf("Failed to write int conn with ID %s", connectionID)
		return
	}

	conn.netIO.PushInt(val)
}

//PushByte write byte to buffer
func (c *Connection) PushByte(connectionID string, val int32) {
	log.Tracef("Pushing byte %d to %s buffers", val, connectionID)
	conn, err := c.getConnection(connectionID)
	if err != nil {
		log.Warnf("Failed to write byte conn with ID %s", connectionID)
		return
	}

	conn.netIO.PushByte(val)
}

//ClearWriteBuffer clear
func (c *Connection) ClearWriteBuffer(connectionID string) {
	log.Tracef("Clearing %s write buffers", connectionID)
	conn, err := c.getConnection(connectionID)
	if err != nil {
		log.Warnf("Failed to clear writebuffer conn with ID %s", connectionID)
		return
	}

	conn.netIO.ClearWriteBuffer()
}

//SendMessage send message
func (c *Connection) SendMessage(connectionID string) {
	log.Tracef("Sending message to %s", connectionID)
	conn, err := c.getConnection(connectionID)
	if err != nil {
		log.Warnf("Failed to send message to %s", connectionID)
	}
	conn.netIO.PrepWriteBuffer()
	_, err = conn.netConnection.Write(conn.netIO.writeBuffer.Bytes())
	if err != nil {
		log.Warnf("error sending message %s", err.Error())
	}

}
