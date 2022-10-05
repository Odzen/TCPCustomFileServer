package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

var (
	leaving      = make(chan message)
	messages     = make(chan message)
	channelGroup = map[string][]Client{
		"first":  {},
		"second": {},
	}
	clients    []Client
	numClients = 0
)

type ChannelGroup map[string][]Client

type Client struct {
	Name       string   `json:"name"`
	Address    string   `json:"address"`
	Connection net.Conn `json:"connection"`
}

func (client *Client) returnJSON() string {
	clientJSON, _ := json.Marshal(client)
	return string(clientJSON)
}

func newClient(name string, connection net.Conn) *Client {
	return &Client{
		Name:       name,
		Address:    connection.RemoteAddr().String(),
		Connection: connection,
	}
}

func suscribeToChannelGroup(client Client, channel string) {
	channelGroup[channel] = append(channelGroup[channel], client)
}

type message struct {
	text    string
	address string
	channel string
}

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func RunServer() {
	server, err := net.Listen(os.Getenv("PROTOCOL_TYPE"), os.Getenv("HOST")+":"+os.Getenv("PORT"))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("Server Running...")
	fmt.Println("Listening on " + os.Getenv("HOST") + ":" + os.Getenv("PORT"))
	fmt.Println("Waiting for client...")

	go broadcaster()
	for {
		connection, err := server.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println("Client connected")
		numClients++
		go processClient(connection)
	}
}

func ChooseChannel() string {
	if numClients%2 == 0 {

		fmt.Printf("Clientes par\n")
		return "first"
	} else {

		fmt.Printf("Clientes impar\n")
		return "second"
	}
}

func processClient(connection net.Conn) {
	chooseChannel := ChooseChannel()
	client := newClient("", connection)
	suscribeToChannelGroup(*client, chooseChannel)

	fmt.Printf("Clients %v", channelGroup[chooseChannel])

	messages <- newMessage(" joined.", connection, chooseChannel)

	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		messages <- newMessage(": "+scanner.Text(), connection, chooseChannel)
	}
	// clients := channelGroup["first"]
	// delete(clients, connection.RemoteAddr().String())

	leaving <- newMessage(" has left.", connection, chooseChannel)

	connection.Close()

}

func newMessage(msg string, conn net.Conn, channel string) message {
	addr := conn.RemoteAddr().String()
	return message{
		text:    addr + msg,
		address: addr,
		channel: channel,
	}
}

func broadcaster() {
	for {
		select {
		case msg := <-messages:
			for _, client := range channelGroup[msg.channel] {

				if msg.address == client.Address { // Checking if the user it's the same who sent the message

					continue
				}
				fmt.Fprintln(client.Connection, msg.text)
			}

		case msg := <-leaving:
			for _, client := range clients {
				fmt.Fprintln(client.Connection, msg.text)
			}

		}
	}
}
