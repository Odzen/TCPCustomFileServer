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

var (
	channelGroup = types.NewChannelGroup(make(map[int][]*types.Client))
	commands     = make(chan types.Command)
	numClients   = 0
	clientLeft   = false
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
	fmt.Println("Server Running! Waiting for connections...")
	fmt.Println("Listening on " + os.Getenv("HOST") + ":" + os.Getenv("PORT"))
	fmt.Println("Waiting for client...")

	go handleCommands()
	for {
		connection, err := server.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println("Client connected: ", connection.LocalAddr().String())
		numClients++
		go processClient(connection)
	}
}

func processClient(connection net.Conn) {
	//defer utils.CloseConnectionClient(connection)
	client := types.NewClient("anonymous", connection, commands)

	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		// Process commands
		newLine := strings.Trim(scanner.Text(), "\r\n")
		args := strings.Split(newLine, " ")
		command := strings.TrimSpace(args[0])
		types.ProcessCommand(command, args, client)

	}

	// Check the flag to know if the client already left using the command `=exit`, or is trying to leave forcing the program to stop
	if !clientLeft {
		types.Exit(client, channelGroup)
	}

}

func handleCommands() {
	for command := range commands {
		switch command.Id {
		case types.USERNAME:
			types.CreateUsername(command.Client, command.Args)

		case types.SUBSCRIBE:
			types.SuscribeToChannel(command.Client, command.Args, channelGroup)

		case types.CHANNELS:
			types.ShowChannels(command.Client, command.Args, channelGroup)

		case types.CURRENT_CHANNEL:
			types.CurrentChannel(command.Client)

		case types.INSTRUCTIONS:
			types.Instructions(command.Client)

		case types.MESSAGE:
			types.SendMessage(command.Client, command.Args, channelGroup)

		case types.FILE:
			types.SendFile(command.Client, command.Args)

		case types.EXIT:
			clientLeft = true
			types.Exit(command.Client, channelGroup)

		}
	}
}
