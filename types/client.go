package types

import (
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
}

type BodyMessage struct {
	TypedCommand Command
	TextMessage  Message
	FileMessage  File
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

func NewClient(name string, connection net.Conn, commands chan Command) *Client {
	return &Client{
		Name:       name,
		Address:    connection.RemoteAddr().String(),
		Connection: connection,
		Commands:   commands,
	}
}