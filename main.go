package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Jon3123/go39/pkg/go39"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("S for server C for client")
		return
	}

	connection := go39.NewConnection()
	if arguments[1] == "S" || arguments[1] == "s" {
		serverLoop(connection)
	} else {
		clientLoop(connection)
	}
}

func clientLoop(connection *go39.Connection) {
	id, _ := connection.TCPConnect("127.0.0.1", 3223)

	connection.ClearWriteBuffer(id)
	connection.PushByte(id, 22)
	connection.PushString(id, "Hello from go client")
	connection.PushInt(id, 31231)
	connection.SendMessage(id)
}
func serverLoop(connection *go39.Connection) {
	connection.TCPListen("127.0.0.1", 3223)
	id := connection.TCPAccept()
	readLoop(connection, id)
}
func readLoop(connection *go39.Connection, id string) {
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
