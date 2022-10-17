package types

import (
	"fmt"
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
	Commands   chan<- Command
}

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

func NewClient(name string, connection net.Conn, commands chan Command) *Client {
	return &Client{
		Name:       name,
		Address:    connection.RemoteAddr().String(),
		Connection: connection,
		Commands:   commands,
	}
}

func (client *Client) ProcessCommand(command string, args []string) {
	switch command {
	case "=username":
		client.Commands <- Command{
			Id:     USERNAME,
			Client: client,
			Args:   args,
		}
	case "=suscribe":
		client.Commands <- Command{
			Id:     SUSCRIBE,
			Client: client,
			Args:   args,
		}
	case "=channels":
		client.Commands <- Command{
			Id:     CHANNELS,
			Client: client,
			Args:   args,
		}
	case "=message":
		client.Commands <- Command{
			Id:     MESSAGE,
			Client: client,
			Args:   args,
		}
	case "=file":
		client.Commands <- Command{
			Id:     FILE,
			Client: client,
			Args:   args,
		}
	case "=exit":
		client.Commands <- Command{
			Id:     EXIT,
			Client: client,
			Args:   args,
		}
	default:
		fmt.Fprintln(client.Connection, "Unknown Command: "+command)
	}
}
