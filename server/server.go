package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	types "github.com/Odzen/TCPCustomFileServer/server/types"
	utils "github.com/Odzen/TCPCustomFileServer/utils"
	"github.com/joho/godotenv"
)

var defaultChannels = map[int][]*types.Client{
	1: {},
	2: {},
}

var (
	// leaving      = make(chan types.Message)
	// messages     = make(chan types.Message)
	channelGroup = types.NewChannelGroup(defaultChannels)
	commands     = make(chan types.Command)
	numClients   = 0
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
	defer utils.CloseConnectionServer(server)
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
	defer utils.CloseConnectionClient(connection)

	// selectedChannel := ChooseChannel()
	client := types.NewClient("anonymous", connection, commands)
	// channelGroup.SuscribeToChannelGroup(*client, selectedChannel)

	// messages <- types.NewMessage(" joined.", connection, selectedChannel)

	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {

		// Process commands
		newLine := strings.Trim(scanner.Text(), "\r\n")
		args := strings.Split(newLine, " ")
		command := strings.TrimSpace(args[0])
		types.ProcessCommand(command, args, client)

		// messages <- types.NewMessage(": "+scanner.Text(), connection, selectedChannel)
	}

	//fmt.Println("CLIENT LEFT")
	//channelGroup.DeleteClientFromChannel(*client, selectedChannel)

	// leaving <- types.NewMessage(" has left.", connection, selectedChannel)

}

func broadcaster() {
	for command := range commands {
		// select {
		// case command := <-commands:
		switch command.Id {
		case types.USERNAME:
			types.CreateUsername(command.Client, command.Args)

		case types.SUSCRIBE:
			types.SuscribeToChannel(command.Client, command.Args, channelGroup)

		case types.CHANNELS:
			types.ShowChannels(command.Client, command.Args, channelGroup)

		case types.MESSAGE:
			types.SendMessage(command.Client, command.Args, channelGroup)

		case types.FILE:
			types.SendFile(command.Client, command.Args)

		case types.EXIT:
			types.Exit(command.Client, command.Args, channelGroup)
		}

		// case msg := <-messages:
		// 	for _, client := range channelGroup.Channels[msg.ChannelPipeline] {

		// 		if msg.Address == client.Address { // Checking if the user it's the same who sent the message
		// 			continue
		// 		}

		// 		fmt.Fprintln(client.Connection, "-> "+msg.Text)
		// 	}

		// case msg := <-leaving:
		// 	for _, client := range channelGroup.Channels[msg.ChannelPipeline] {
		// 		fmt.Fprintln(client.Connection, "-> "+msg.Text)
		// 	}

		// }
	}
}
