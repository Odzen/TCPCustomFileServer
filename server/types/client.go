package types

import (
	"encoding/json"
	"net"
)

type IClient interface {
	ReturnJSON() string
	ChangeName(newName string)
}

type Client struct {
	Name       string   `json:"name"`
	Address    string   `json:"address"`
	Connection net.Conn `json:"connection"`
}

func (client *Client) ReturnJSON() string {
	clientJSON, _ := json.Marshal(client)
	return string(clientJSON)
}

func (client *Client) ChangeName(newName string) {
	client.Name = newName
}

func NewClient(name string, connection net.Conn) *Client {
	return &Client{
		Name:       name,
		Address:    connection.RemoteAddr().String(),
		Connection: connection,
}
}