package types

import (
	"fmt"
	"net"
	"os"
)

const SIZE = 1024

type File struct {
	Name            string
	Size            int64
	Content         []byte
	Address         string
	ChannelPipeline int
}

func NewFile(name string, size int64, content []byte, address string, pipeline int) File {
	return File{
		Name:            name,
		Size:            size,
		Content:         content,
		Address:         address,
		ChannelPipeline: pipeline,
	}
}

func ProccessingFile(connection net.Conn, path string, client *Client) (File, bool) {

	file, err := os.Open(path)

	if err != nil {
		fmt.Println("Error reading the file:", err)
		return File{}, true
	}

	fileInfo, err := file.Stat()

	if err != nil {
		fmt.Println("Error getting information of the file:", err)
		return File{}, true
	}
	// file's data can be read into a slice of bytes
	data := make([]byte, fileInfo.Size())
	count, err := file.Read(data)

	if err != nil {
		fmt.Println("Error counting the bytes of the file:", err)
		return File{}, true
	}

	newFile := NewFile(fileInfo.Name(), fileInfo.Size(), data[:count], client.Address, client.SuscribedToChannel)
	fmt.Printf("The file has been processed: %s -- %d -- %q\n", newFile.Name, newFile.Size, newFile.Content)

	return newFile, false
}
