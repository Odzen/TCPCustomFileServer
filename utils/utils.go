package utils

import (
	"log"
	"net"
)

func CloseConnectionServer(serverConnection net.Listener) {
	if err := serverConnection.Close(); err != nil {
		panic(err)
	}
}

func CloseConnectionClient(clientConnection net.Conn) {
	if err := clientConnection.Close(); err != nil {
		log.Println(err)
	}
}
