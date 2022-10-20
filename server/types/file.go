package types

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const SIZE = 1024

type File struct {
	Name    string
	Size    int
	Content []byte
}

func NewFile(name string, size int, content []byte) File {
	return File{
		Name:    name,
		Size:    size,
		Content: content,
	}
}

func SendFileToClient(connection net.Conn, path string) {
	fmt.Println("Sending File to client", connection.LocalAddr().String())
	//defer utils.CloseConnectionClient(connection)

	file, err := os.Open(path)

	if err != nil {
		fmt.Println("Error reading the file:", err)
		return
	}
	//defer file.Close()

	// file's data can be read into a slice of bytes
	data := make([]byte, 100)
	count, err := file.Read(data)
	if err != nil {
		fmt.Println("Error counting the bytes of the file:", err)
		return
	}

	fmt.Printf("read %d bytes: %q\n", count, data[:count])

	newFile := NewFile(strings.Split(file.Name(), "/")[1], count, data[:count])
	fmt.Printf("File to send over: %s -- %d -- %q", newFile.Name, newFile.Size, newFile.Content)

	// fileInfo, err := file.Stat()
	// fmt.Println("File information: ", fileInfo)

	//size := fillString()
}
