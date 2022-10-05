package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	types "github.com/Odzen/TCPCustomFileServer/server/types"
	"github.com/joho/godotenv"
)

var (
	leaving      = make(chan types.Message)
	messages     = make(chan types.Message)
	channelGroup = map[string][]types.Client{
		"first":  {},
		"second": {},
	}
	clients    []types.Client
	numClients = 0
)

type ChannelGroup map[string][]types.Client

func suscribeToChannelGroup(client types.Client, channel string) {
	channelGroup[channel] = append(channelGroup[channel], client)
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
	client := types.NewClient("", connection)
	suscribeToChannelGroup(*client, chooseChannel)

	fmt.Printf("Clients %v", channelGroup[chooseChannel])

	messages <- types.NewMessage(" joined.", connection, chooseChannel)

	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		messages <- types.NewMessage(": "+scanner.Text(), connection, chooseChannel)
	}
	// clients := channelGroup["first"]
	// delete(clients, connection.RemoteAddr().String())

	leaving <- types.NewMessage(" has left.", connection, chooseChannel)

	connection.Close()

}

func broadcaster() {
	for {
		select {
		case msg := <-messages:
			for _, client := range channelGroup[msg.Channel] {

				if msg.Address == client.Address { // Checking if the user it's the same who sent the message

					continue
				}

				fmt.Fprintln(client.Connection, msg.Text)
			}

		case msg := <-leaving:
			for _, client := range clients {
				fmt.Fprintln(client.Connection, msg.Text)
			}

		}
	}
}
