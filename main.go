package main

import (
	"fmt"
	"time"

	"github.com/Jon3123/go39/pkg/go39"
)

func main() {

	connection := *go39.NewConnection()

	connection.TCPListen("127.0.0.1", 3223)
	id := connection.TCPAccept()
	readLoop(connection, id)
}

func readLoop(connection go39.Connection, id string) {
	for {
		bytesRead := connection.ReceiveMessage(id, time.Second)
		if bytesRead > 0 {
			fmt.Printf("read byte %d\n", connection.PopByte(id))
			fmt.Printf("read string %s\n", connection.PopString(id))
			fmt.Printf("read int %d\n", connection.PopInt(id))
			connection.ClearWriteBuffer(id)
			connection.PushByte(id, 22)
			connection.PushInt(id, 2000)
			connection.PushString(id, "HI how are you")
			connection.SendMessage(id)
		}

		fmt.Println(bytesRead)
		if bytesRead == -1 {
			break
		}
	}
	fmt.Println("BREAK!")
}
