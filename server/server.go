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
	clients []Client
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
		go processClient(connection)
	}
}

func processClient(connection net.Conn) {

	client := newClient("", connection)
	suscribeToChannelGroup(*client, "first")

	fmt.Printf("Clients %v", channelGroup["first"])

	messages <- newMessage(" joined.", connection, "first")

	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		messages <- newMessage(": "+scanner.Text(), connection, "first")
	}
	// clients := channelGroup["first"]
	// delete(clients, connection.RemoteAddr().String())

	leaving <- newMessage(" has left.", connection, "first")

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
	//fmt.Printf("Clients %v", clients)
	for {
		fmt.Println("BROADCAST")
		select {
		case msg := <-messages:
			fmt.Println("CASE INCOMING")
			for _, client := range channelGroup[msg.channel] {
				fmt.Println("FOR")
				if msg.address == client.Address { // Checking if the user it's the same who sent the message
					fmt.Println("SAME")
					continue
				}
				fmt.Println("CLIENT: ", client.returnJSON())
				fmt.Fprintln(client.Connection, msg.text)
			}

		case msg := <-leaving:
			fmt.Println("CASE LEAVING")
			for _, client := range clients {
				fmt.Fprintln(client.Connection, msg.text)
			}

		}
	}
}
