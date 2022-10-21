package types

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Odzen/TCPCustomFileServer/utils"
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
	log.Println("Client: " + client.Name + "--" + "is saving the file:" + file.Name)

	fmt.Fprintln(client.Connection, fmt.Sprintln("Received the file: ", file))

	err := os.MkdirAll(fmt.Sprintf("outFiles/%d", client.SuscribedToChannel), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Println("Error creating the folder", err)
		return err
	}
	fileToSave, err := os.Create(fmt.Sprintf("./outFiles/%d/%s", client.SuscribedToChannel, file.Name))

	if err != nil {
		log.Println("Error creating the file in the folder", err)
		return err
	}
	defer utils.CloseFile(fileToSave)

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
