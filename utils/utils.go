package utils

import (
	"fmt"
	"net"
	"os"
)

func CloseConnectionServer(serverConnection net.Listener) {
	if err := serverConnection.Close(); err != nil {
		fmt.Println("Error closing the server connection ")
		panic(err)
	}
}

func CloseConnectionClient(clientConnection net.Conn) {
	if err := clientConnection.Close(); err != nil {
		fmt.Println("Error closing the client connection ")
		panic(err)
	}
}

func CloseFile(file *os.File) {
	if err := file.Close(); err != nil {
		fmt.Println("Error closing the file ", err)
	}
}
