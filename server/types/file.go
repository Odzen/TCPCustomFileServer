package types

import (
	"fmt"
	"net"
	"os"
)

const SIZE = 1024

func SendFileToClient(connection net.Conn, path string) {
	fmt.Println("Sending File to client", connection.LocalAddr().String())
	//defer utils.CloseConnectionClient(connection)

	file, err := os.Open(path)

	if err != nil {
		fmt.Println("Error reading the file:", err)
		return
	}
	defer file.Close()

	// file's data can be read into a slice of bytes
	data := make([]byte, 100)
	count, err := file.Read(data)
	if err != nil {
		fmt.Println("Error counting the bytes of the file:", err)
	}

	fmt.Printf("read %d bytes: %q\n", count, data[:count])

	//size := fillString()
}
