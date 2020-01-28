package main

import (
	"fmt"

	"github.com/Jon3123/go39/pkg/go39"
	"github.com/sirupsen/logrus"
)

func main() {

	connection := go39.Connection{}

	connection.TCPListen("127.0.0.1", 3223)
	id := connection.TCPAccept()
	readLoop(connection, id)
	for {

		logrus.Info(id)

	}
}

func readLoop(connection go39.Connection, id string) {
	for {
		bytesRead := connection.ReceiveMessage(id)
		fmt.Println(bytesRead)
		if bytesRead > 0 {
			fmt.Println(connection.PopByte(id))
			fmt.Println(connection.PopString(id))
			fmt.Println(connection.PopInt(id))
			connection.ClearWriteBuffer(id)
			connection.PushByte(id, 22)
			connection.PushInt(id, 2000)
			connection.PushString(id, "HI how are you")
			connection.SendMessage(id)
		}
	}
}
