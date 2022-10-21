package utils

import (
	"log"
	"net"
	"os"
)

func CloseConnectionServer(serverConnection net.Listener) {
	if err := serverConnection.Close(); err != nil {
		log.Println("Error closing the server connection ")
		panic(err)
	}
}

func CloseConnectionClient(clientConnection net.Conn) {
	if err := clientConnection.Close(); err != nil {
		log.Println("Error closing the client connection ")
		panic(err)
	}
}

func CloseFile(file *os.File) {
	if err := file.Close(); err != nil {
		log.Println("Error closing the file ", err)
	}
}
