package types

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
)

const MAX_SIZE = 50500

var sentFiles = make([]*File, 0)

type File struct {
	Name            string `json:"name"`
	Size            int64  `json:"size"`
	Content         []byte `json:"content"`
	AddressClient   string `json:"address"`
	ChannelPipeline int    `json:"pipeline"`
}

func NewFile(name string, size int64, content []byte, address string, pipeline int) File {
	return File{
		Name:            name,
		Size:            size,
		Content:         content,
		AddressClient:   address,
		ChannelPipeline: pipeline,
	}
}

func (file *File) appendToSentFiles() {
	sentFiles = append(sentFiles, file)
}

func SentFilesToJson() ([]byte, error) {
	return json.Marshal(sentFiles)
}

func ProcessingFile(connection net.Conn, path string, client *Client) (File, bool) {

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

	if fileInfo.Size() >= MAX_SIZE {
		fmt.Println("The file size cannot be greater than " + strconv.Itoa(MAX_SIZE) + " bytes")
		return File{}, true
	}

	// file's data can be read into a slice of bytes
	data := make([]byte, fileInfo.Size())
	count, err := file.Read(data)

	if err != nil {
		fmt.Println("Error counting the bytes of the file:", err)
		return File{}, true
	}

	newFile := NewFile(fileInfo.Name(), fileInfo.Size(), data[:count], client.Address, client.SubscribedToChannel)
	fmt.Printf("The file has been processed: %s -- %d -- %q\n", newFile.Name, newFile.Size, newFile.Content)

	return newFile, false
}
