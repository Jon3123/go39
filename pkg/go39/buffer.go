package go39

import "fmt"

//PushByte - Write a byte
func (c *Connection) PushByte(value int) {
	log.Infof("writing byte %d", value)
	c.writeBuffer.WriteByte(byte(value & 0xFF))
}

//PushInt - push an int to the buffer
func (c *Connection) PushInt(value int) {
	log.Infof("Writing int %d", value)
	c.writeBuffer.WriteRune(rune(value))
}

//PushString - push a string to the buffer
func (c *Connection) PushString(str string) {
	log.Infof("Writing string %s", str)
	c.writeBuffer.WriteString(str)
	//write null terminating character
	c.writeBuffer.WriteByte(0)
}

//ClearWriteBuffer clear the write buffer
func (c *Connection) ClearWriteBuffer() {
	log.Infof("Clearing write buffer")
	c.writeBuffer.Reset()
}

//PrintBuffer Print the buffer
func (c *Connection) PrintBuffer() {
	fmt.Println(c.writeBuffer)
}

//Copy temp func
func (c *Connection) Copy() {
	c.readBuffer = c.writeBuffer
}

//PopByte - Read byte
func (c *Connection) PopByte() int {
	val, err := c.readBuffer.ReadByte()
	if err != nil {
		log.Warn("Read byte failed")
	}
	return int(val)
}

//PopInt - read int
func (c *Connection) PopInt() int {
	val, _, err := c.readBuffer.ReadRune()
	if err != nil {
		log.Warn("Read int failed")
	}
	return int(val)
}

//PopString - read string
func (c *Connection) PopString() string {
	str, err := c.readBuffer.ReadString(0)
	if err != nil {
		log.Warn("Read string failed")
		log.Warn(err)
	}
	return str
}

//ClearReadBuffer clear the read buffer
func (c *Connection) ClearReadBuffer() {
	log.Infof("Clearing read buffer")
	c.readBuffer.Reset()
}
