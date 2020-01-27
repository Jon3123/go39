package main

import (
	"fmt"

	"github.com/Jon3123/go39/pkg/go39"
	"github.com/sirupsen/logrus"
)

func main() {
	connection := go39.Connection{}

	connection.TCPListen("127.0.0.1", 3223)

	for {
		id := connection.TCPAccept()
		logrus.Info(id)
		connection.ReceiveMessage(id)
		fmt.Println(connection.PopString(id))
	}
}
