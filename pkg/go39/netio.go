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
func (n *NetIO) PushByte(value int32) {
	log.Debugf("writing byte %d", value)
	n.writeBuffer.WriteByte(byte(value & 0xFF))
}

//PushInt - push an int to the buffer
func (n *NetIO) PushInt(value int32) {
	log.Debugf("Writing int %d", value)
	b1 := (value >> 24) & 0xFF
	b2 := (value >> 16) & 0xFF
	b3 := (value >> 8) & 0xFF
	b4 := (value & 0xFF)

	n.PushByte(b1)
	n.PushByte(b2)
	n.PushByte(b3)
	n.PushByte(b4)
}

//PushString - push a string to the buffer
func (n *NetIO) PushString(str string) {
	log.Debugf("Writing string %s", str)
	n.writeBuffer.WriteString(str)
	//write null terminating character
	n.writeBuffer.WriteByte(0)
}

//ClearWriteBuffer clear the write buffer
func (n *NetIO) ClearWriteBuffer() {
	log.Debugf("Clearing write buffer")
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
func (n *NetIO) PopByte() int32 {
	val, err := n.readBuffer.ReadByte()
	if err != nil {
		log.Warn("Read byte failed")
	}
	return int32(val)
}

//PopInt - read int
func (n *NetIO) PopInt() int32 {
	b1 := (n.PopByte() & 0xFF) << 24
	b2 := (n.PopByte() & 0xFF) << 16
	b3 := (n.PopByte() & 0xFF) << 8
	b4 := (n.PopByte() & 0xFF)

	res := b1 | b2 | b3 | b4

	return res
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
	log.Trace("Clearing read buffer")
	n.readBuffer.Reset()
}

//PrepWriteBuffer prep the writebuffer for sending
func (n *NetIO) PrepWriteBuffer() {
	s := n.writeBuffer.String()
	n.writeBuffer.Reset()
	n.PushByte(int32(len(s)))
	n.PushByte(0)
	n.writeBuffer.WriteString(s) //To prevent null terminating char that is in push string func
}
