package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
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

// TCP

func RunServer() {
	server, err := net.Listen(os.Getenv("PROTOCOL_TYPE"), os.Getenv("HOST")+":"+os.Getenv("PORT_TCP"))

	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer utils.CloseConnectionServer(server)

	fmt.Println("Server Running! Waiting for connections...")
	fmt.Println("Listening on " + os.Getenv("HOST") + ":" + os.Getenv("PORT_TCP"))
	fmt.Println("Waiting for client...")

	go handleHttpRequest() // Go routine to handle the http requests

	go handleCommands() // Go routine to handle the incoming commands

	for {
		connection, err := server.Accept()
		fmt.Println("Connection: ", connection.RemoteAddr().String())
		if err != nil {
			fmt.Println(err)
			continue
		}
		numClients++
		go processClient(connection) // Go routine to handle the incoming clients

	}
}

func processClient(connection net.Conn) {
	client := types.NewClient("anonymous", connection, commands)

	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		// Process commands
		newLine := strings.Trim(scanner.Text(), "\r\n")
		args := strings.Split(newLine, " ")
		command := strings.TrimSpace(args[0])
		types.ProcessCommand(command, args, client)

	}

	// Checking the flag to know if the client already left using the command `=exit`, or is trying to leave forcing the program to stop
	if !clientLeft {
		types.Exit(client, channelGroup)
	}

}

func handleCommands() {
	for command := range commands {
		switch command.Id {
		case types.USERNAME:
			if len(command.Args) == 1 {
				fmt.Fprintln(command.Client.Connection, "-> "+"You have to type a name")
				break
			}
			types.CreateUsername(command.Client, command.Args)

		case types.SUBSCRIBE:
			if len(command.Args) == 1 {
				fmt.Fprintln(command.Client.Connection, "-> "+"You have to type a channel, remember that the channel must be a number")
				break
			}
			types.SuscribeToChannel(command.Client, command.Args, channelGroup)

		case types.CHANNELS:
			types.ShowChannels(command.Client, command.Args, channelGroup)

		case types.CURRENT_CHANNEL:
			types.CurrentChannel(command.Client)

		case types.INSTRUCTIONS:
			types.Instructions(command.Client)

		case types.MESSAGE:
			if len(command.Args) == 1 {
				fmt.Fprintln(command.Client.Connection, "-> "+"You have to type a message")
				break
			}
			types.SendMessage(command.Client, command.Args, channelGroup)

		case types.FILE:
			if len(command.Args) == 1 {
				fmt.Fprintln(command.Client.Connection, "-> "+"You have to type the absolute path or the name of the file")
				break
			}
			types.SendFile(command.Client, command.Args, channelGroup)

		case types.EXIT:
			clientLeft = true
			types.Exit(command.Client, channelGroup)

		}
	}
}

// HTTP

func handleHttpRequest() {
	http.HandleFunc("/clients", serveHTTPClients)
	http.HandleFunc("/files", serveHTTPFiles)
	err := http.ListenAndServe(":"+os.Getenv("PORT_WEB"), nil)
	if err != nil {
		log.Fatal("Error Listening to port: "+os.Getenv("PORT_WEB")+" ", err)
	}
}

func serveHTTPClients(res http.ResponseWriter, req *http.Request) {
	channelGroupJson, err := channelGroup.ToJson()

	//Allow CORS here By * or specific origin
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if err != nil {
		fmt.Println("Error parsing channel group to JSON", err)
	}

	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(channelGroupJson)

	if err != nil {
		fmt.Println("Error writing the response", err)
	}
}

func serveHTTPFiles(res http.ResponseWriter, req *http.Request) {
	filesJson, err := types.SentFilesToJson()

	//Allow CORS here By * or specific origin
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if err != nil {
		fmt.Println("Error parsing the files to JSON", err)
	}

	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(filesJson)

	if err != nil {
		fmt.Println("Error writing the response", err)
	}
}
