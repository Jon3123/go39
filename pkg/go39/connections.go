package go39

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Jon3123/go39/pkg/utils"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

//ConfigLogger config log
func ConfigLogger(level logrus.Level, output io.Writer) {
	log.SetLevel(level)
	log.SetOutput(output)
}

//ConnectionType type of connection the connection represents
type ConnectionType int

const (
	//TCPClient this is a client
	TCPClient ConnectionType = 0
	//TCPServer this is a server
	TCPServer ConnectionType = 1
	//UDP udp type
	UDP ConnectionType = 2
	//MaxTransmitSize max count of bytes you are allowed to send
	MaxTransmitSize int = 1024
)

//Connection to fill out later
type Connection struct {
	netIO          NetIO
	writeMux       sync.Mutex
	readMux        sync.Mutex
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
	log.Info("New Connection")
	connectionID = c.addConnectionTCP(ln)
	return
}

//StartWrite start a write to lock buffer
func (c *Connection) StartWrite(connectionID string) {
	if val, ok := c.connections[connectionID]; ok {
		val.writeMux.Lock()
	}
}

//EndWrite to end write to allow other go rountines to have access
func (c *Connection) EndWrite(connectionID string) {
	if val, ok := c.connections[connectionID]; ok {
		val.writeMux.Unlock()
	}
}

//StartRead start a read to lock buffer
func (c *Connection) StartRead(connectionID string) {
	if val, ok := c.connections[connectionID]; ok {
		val.readMux.Lock()
	}
}

//EndRead start a read to lock buffer
func (c *Connection) EndRead(connectionID string) {
	if val, ok := c.connections[connectionID]; ok {
		val.readMux.Unlock()
	}
}

//Add tcp connection to map
func (c *Connection) addConnectionTCP(connection net.Conn) (connectionID string) {

	if c.connections == nil {
		log.Fatal("Cannot add connection! the map is nil!")
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

//Add self to connections map
func (c *Connection) addSelf() (connectionID string) {
	log.Tracef("Adding self to map ")
	if c.connections == nil {
		log.Fatal("Cannot add connection! the map is nil!")
	}
	connectionID = utils.GenerateConnectionID()
	c.connectionID = connectionID
	log.Tracef("Adding self with id %s", connectionID)
	c.connections[connectionID] = c
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
	log.Tracef("Reading from connection with ID %s\n", connectionID)
	conn, err := c.getConnection(connectionID)
	if err != nil {
		log.Fatalf("error getting connection: %s", err.Error())
	}
	b := make([]byte, MaxTransmitSize)
	conn.netConnection.SetReadDeadline(time.Now().Add(timeout))
	_, err = conn.netConnection.Read(b)
	if err != nil {
		if err.Error() == "EOF" {
			//TODO Add some disconnect stuff possibly ??
			return -1
		}

		if strings.Contains(err.Error(), "i/o timeout") {
			return 0
		}
		log.Warnf("error reading from connection %s: %s", connectionID, err.Error())
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

//PopFloat32 read float32
func (c *Connection) PopFloat32(connectionID string) (val float32) {
	log.Tracef("Reading float32 in conn with ID %s", connectionID)
	conn, err := c.getConnection(connectionID)

	if err != nil {
		log.Warnf("Failed to read float from conn with ID %s", connectionID)
		return
	}

	val = conn.netIO.PopFloat32()
	return
}

//PopFloat64 read float 64
func (c *Connection) PopFloat64(connectionID string) (val float64) {
	log.Tracef("Reading float64 in conn with ID %s", connectionID)
	conn, err := c.getConnection(connectionID)

	if err != nil {
		log.Warnf("Failed to read float from conn with ID %s", connectionID)
		return
	}

	val = conn.netIO.PopFloat64()
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

//PushString Write a string to buffer
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

//PushFloat32 write float 32 to buffer
func (c *Connection) PushFloat32(connectionID string, val float32) {
	log.Tracef("Pushing float32 %f to %s buffers", val, connectionID)
	conn, err := c.getConnection(connectionID)
	if err != nil {
		log.Warnf("Failed to write float32")
		return
	}

	conn.netIO.PushFloat32(val)

}

//PushFloat64 write float 64 to buffer
func (c *Connection) PushFloat64(connectionID string, val float64) {
	log.Tracef("Pushing float64 %f to %s buffers", val, connectionID)
	conn, err := c.getConnection(connectionID)
	if err != nil {
		log.Warnf("Failed to write float64")
		return
	}

	conn.netIO.PushFloat64(val)

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

//TCPConnect connect to tcp server
func (c *Connection) TCPConnect(ip string, port int) (connectionID string, err error) {
	log.Tracef("TCP Connect %s:%d", ip, port)
	netConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))

	if err != nil {
		log.Errorf("Failed to connect to %s:%d", ip, port)
		return
	}

	c.netConnection = netConn
	c.connectionType = TCPClient
	//add self to own map
	c.addSelf()
	connectionID = c.connectionID
	return
}

//UDPConnect connect to udp
func (c *Connection) UDPConnect(ip string, port int) (connectionID string, err error) {
	log.Tracef("UDP Connect %s:%d", ip, port)
	netConn, err := net.Dial("udp", fmt.Sprintf("%s:%d", ip, port))

	if err != nil {
		log.Errorf("Failed to connect to %s:%d", ip, port)
		return
	}

	c.netConnection = netConn
	c.connectionType = UDP
	c.addSelf()
	connectionID = c.connectionID
	return
}

//CloseConnection close connection
func (c *Connection) CloseConnection(connectionID string) {
	c.connections[connectionID].netConnection.Close()
	delete(c.connections, connectionID)
}
