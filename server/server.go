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
	leaving                         = make(chan types.Message)
	messages                        = make(chan types.Message)
	channelGroup types.ChannelGroup = map[int][]types.Client{
		1: {},
		2: {},
	}
	clients    []types.Client
	numClients = 0
)

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

func ChooseChannel() int {
	if numClients%2 == 0 {
		return 1
	} else {
		return 2
	}
}

func processClient(connection net.Conn) {
	chooseChannel := ChooseChannel()
	client := types.NewClient("", connection)
	channelGroup.SuscribeToChannelGroup(*client, chooseChannel)

	fmt.Printf("Clients %v", channelGroup[chooseChannel])

	messages <- types.NewMessage(" joined.", connection, chooseChannel)

	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		messages <- types.NewMessage(": "+scanner.Text(), connection, chooseChannel)
	}
	// clients := channelGroup["1"]
	// delete(clients, connection.RemoteAddr().String())

	leaving <- types.NewMessage(" has left.", connection, chooseChannel)

	connection.Close()

}

func broadcaster() {
	for {
		select {
		case msg := <-messages:
			for _, client := range channelGroup[msg.ChannelPipeline] {

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
