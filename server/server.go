package server

import (
	"bufio"
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
)

type ChannelGroup map[string][]Client

type Client struct {
	name       string
	address    string
	connection net.Conn
}

func newClient(name string, connection net.Conn) *Client {
	return &Client{
		name:       name,
		address:    connection.RemoteAddr().String(),
		connection: connection,
	}
}

func suscribeToChannelGroup(client Client, channel string) {
	channelGroup[channel] = append(channelGroup[channel], client)
}

type message struct {
	text    string
	address string
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

	messages <- newMessage(" joined.", connection)

	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		messages <- newMessage(": "+scanner.Text(), connection)
	}
	// clients := channelGroup["first"]
	// delete(clients, connection.RemoteAddr().String())

	leaving <- newMessage(" has left.", connection)

	connection.Close()

}

func newMessage(msg string, conn net.Conn) message {
	addr := conn.RemoteAddr().String()
	return message{
		text:    addr + msg,
		address: addr,
	}
}

func broadcaster() {
	clients := channelGroup["first"]
	for {
		select {
		case msg := <-messages:
			for _, client := range clients {
				if msg.address == client.address { // Checking if the user it's the same who sent the message
					continue
				}
				fmt.Fprintln(client.connection, msg.text)
			}

		case msg := <-leaving:
			for _, client := range clients {
				fmt.Fprintln(client.connection, msg.text)
			}

		}
	}
}
