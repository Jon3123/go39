package go39

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
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

//PushShort push a short onto the buffer
func (n *NetIO) PushShort(value int32) {
	log.Debugf("Writing short %d", value)
	b1 := (value) & 0xFF
	b2 := (value >> 8 & 0xFF)

	n.PushByte(b1)
	n.PushByte(b2)
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

//PushFloat32 - push a float to buffer
func (n *NetIO) PushFloat32(value float32) {
	log.Debugf("Writing float 32 %f", value)
	b := math.Float32bits(value)
	n.PushByte(int32(b >> 24))
	n.PushByte(int32(b >> 16))
	n.PushByte(int32(b >> 8))
	n.PushByte(int32(b))
}

//PushFloat64 - push a float64 to buffer
func (n *NetIO) PushFloat64(value float64) {
	log.Debugf("Writing float 64 %f", value)
	b := math.Float64bits(value)
	n.PushByte(int32(b >> 56))
	n.PushByte(int32(b >> 48))
	n.PushByte(int32(b >> 40))
	n.PushByte(int32(b >> 32))
	n.PushByte(int32(b >> 24))
	n.PushByte(int32(b >> 16))
	n.PushByte(int32(b >> 8))
	n.PushByte(int32(b))
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

//PopShort - Read short
func (n *NetIO) PopShort() int32 {
	b1 := (n.PopByte() & 0xFF)
	b2 := (n.PopByte() & 0xFF) << 8 //weird idk to match what gm does

	res := (b1 | b2) & 0xFFFF

	return res
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
	return str[0 : len(str)-1] //the read string func is gonna include the null terminated case in the returned string
	//it is something that is not visibile but is in the actual string which will mess up regex expressions
}

//PopFloat32 - read float32
func (n *NetIO) PopFloat32() float32 {
	var buf [4]byte

	buf[0] = byte(n.PopByte())
	buf[1] = byte(n.PopByte())
	buf[2] = byte(n.PopByte())
	buf[3] = byte(n.PopByte())
	bits := binary.BigEndian.Uint32(buf[:])
	return math.Float32frombits(bits)
}

//PopFloat64 - read float64
func (n *NetIO) PopFloat64() float64 {
	var buf [8]byte
	buf[0] = byte(n.PopByte())
	buf[1] = byte(n.PopByte())
	buf[2] = byte(n.PopByte())
	buf[3] = byte(n.PopByte())
	buf[4] = byte(n.PopByte())
	buf[5] = byte(n.PopByte())
	buf[6] = byte(n.PopByte())
	buf[7] = byte(n.PopByte())

	bits := binary.BigEndian.Uint64(buf[:])
	return math.Float64frombits(bits)
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

//SkipBytes skip the amount of bytes in read buffer
func (n *NetIO) SkipBytes(count int32) {
	buf := make([]byte, count)
	n.readBuffer.Read(buf)
}
