package types

import (
	"fmt"
	"log"
	"net"
	"os"
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

func (client *Client) SaveFile(file File) error {
	log.Println("Saving file")
	fmt.Fprintln(client.Connection, file)

	err := os.MkdirAll(fmt.Sprintf("outFiles/%d", client.SuscribedToChannel), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Println(err)
		return err
	}
	fileToSave, err := os.Create(fmt.Sprintf("./files/%d/%s", client.SuscribedToChannel, file.Name))

	if err != nil {
		log.Println(err)
		return err
	}
	defer fileToSave.Close()

	if _, err := fileToSave.Write(file.Content); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
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
