package types

import (
	"fmt"
	"log"
	"net"
)

type IClient interface {
	ReturnJSON() string
	ChangeName(newName string)
}

type Client struct {
	Name               string   `json:"name"`
	Address            string   `json:"address"`
	Connection         net.Conn `json:"connection"`
	SuscribedToChannel int
	Commands           chan<- Command
	ChannelForFile     chan File
}

// TODO : Return clients in JSON format
// func (client *Client) ReturnJSON() string {
// 	clientJSON, _ := json.Marshal(client)
// 	return string(clientJSON)
// }

func (client *Client) ChangeName(newName string) {
	client.Name = newName
}

func (client *Client) equals(otherClient Client) bool {
	return (client.Address == otherClient.Address) && (client.Name == otherClient.Name) && (client.Connection == otherClient.Connection)
}

func (client *Client) GetCurrentChannel() int {
	return client.SuscribedToChannel
}

func (client *Client) VerifiyingFiles(file File) {
	log.Println("Verifiying files")
	fmt.Fprintln(client.Connection, file)
	// log.Println("Verifiying files")
	// for file := range client.ChannelForFile {
	// 	fmt.Fprintln(client.Connection, file)
	// }
}

func NewClient(name string, connection net.Conn, commands chan Command, channelFile chan File) *Client {
	return &Client{
		Name:           name,
		Address:        connection.RemoteAddr().String(),
		Connection:     connection,
		Commands:       commands,
		ChannelForFile: channelFile,
	}
}
