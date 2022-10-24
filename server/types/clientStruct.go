package types

import (
	"fmt"
	"net"
	"os"

	"github.com/Odzen/TCPCustomFileServer/utils"
)

type IClient interface {
	ChangeName(newName string)
	equals(otherClient Client) bool
	getCurrentChannel() int
	saveFile(file File) error
}

type Client struct {
	Name               string         `json:"name"`
	Address            string         `json:"address"`
	Connection         net.Conn       `json:"-"`
	SuscribedToChannel int            `json:"channel"`
	Commands           chan<- Command `json:"-"`
}

func (client *Client) changeName(newName string) {
	client.Name = newName
}

func (client *Client) equals(otherClient Client) bool {
	return (client.Address == otherClient.Address) && (client.Name == otherClient.Name) && (client.Connection == otherClient.Connection)
}

func (client *Client) getCurrentChannel() int {
	return client.SuscribedToChannel
}

func (client *Client) saveFile(file File) error {

	fmt.Println("Client: " + client.Name + "--" + client.Address + " is saving the file:" + file.Name)

	fmt.Fprintln(client.Connection, fmt.Sprintln("-> Received the file: ", file))

	err := os.MkdirAll(fmt.Sprintf("outFiles/%d/%s", client.SuscribedToChannel, client.Name), os.ModePerm) // path,  Unix permission bits, 0o777

	if err != nil && !os.IsExist(err) {
		fmt.Println("Error creating the folder", err)
		return err
	}

	fileToSave, err := os.Create(fmt.Sprintf("./outFiles/%d/%s/%s", client.SuscribedToChannel, client.Name, file.Name))

	if err != nil {
		fmt.Println("Error creating the file in the folder", err)
		return err
	}
	defer utils.CloseFile(fileToSave)

	if _, err := fileToSave.Write(file.Content); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func NewClient(name string, connection net.Conn, commands chan Command) *Client {
	return &Client{
		Name:       name,
		Address:    connection.RemoteAddr().String(),
		Connection: connection,
		Commands:   commands,
	}
}
