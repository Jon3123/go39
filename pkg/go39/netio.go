package go39

import (
	"bytes"
	"fmt"
)

//NetIO network input/ output
type NetIO struct {
	writeBuffer bytes.Buffer
	readBuffer  bytes.Buffer
}

//PushByte - Write a byte
func (n *NetIO) PushByte(value int) {
	log.Infof("writing byte %d", value)
	n.writeBuffer.WriteByte(byte(value & 0xFF))
}

//PushInt - push an int to the buffer
func (n *NetIO) PushInt(value int) {
	log.Infof("Writing int %d", value)
	n.writeBuffer.WriteRune(rune(value))
}

//PushString - push a string to the buffer
func (n *NetIO) PushString(str string) {
	log.Infof("Writing string %s", str)
	n.writeBuffer.WriteString(str)
	//write null terminating character
	n.writeBuffer.WriteByte(0)
}

//ClearWriteBuffer clear the write buffer
func (n *NetIO) ClearWriteBuffer() {
	log.Infof("Clearing write buffer")
	n.writeBuffer.Reset()
}

//PrintBuffer Print the buffer
func (n *NetIO) PrintBuffer() {
	fmt.Println(n.writeBuffer)
}

//Copy temp func
func (n *NetIO) Copy() {
	n.readBuffer = n.writeBuffer
}

//PopByte - Read byte
func (n *NetIO) PopByte() int {
	val, err := n.readBuffer.ReadByte()
	if err != nil {
		log.Warn("Read byte failed")
	}
	return int(val)
}

//PopInt - read int
func (n *NetIO) PopInt() int {
	val, _, err := n.readBuffer.ReadRune()
	if err != nil {
		log.Warn("Read int failed")
	}
	return int(val)
}

//PopString - read string
func (n *NetIO) PopString() string {
	str, err := n.readBuffer.ReadString(0)
	if err != nil {
		log.Warn("Read string failed")
		log.Warn(err)
	}
	return str
}

//ClearReadBuffer clear the read buffer
func (n *NetIO) ClearReadBuffer() {
	log.Infof("Clearing read buffer")
	n.readBuffer.Reset()
}
