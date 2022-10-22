package server

import (
	"bufio"
	"encoding/json"
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
	sentFiles    = make([]*types.File, 0)
)

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

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

	go handleHttpRequest()
	go handleCommands()

	for {
		connection, err := server.Accept()
		fmt.Println("Connection: ", connection.RemoteAddr().String())
		if err != nil {
			fmt.Println(err)
			continue
		}
		numClients++
		go processClient(connection)

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
			types.SendFile(command.Client, command.Args, channelGroup, sentFiles)

		case types.EXIT:
			clientLeft = true
			types.Exit(command.Client, channelGroup)

		}
	}
}

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
	filesJson, err := json.Marshal(sentFiles)

	if err != nil {
		fmt.Println("Error parsing the files to JSON", err)
	}

	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(filesJson)

	if err != nil {
		fmt.Println("Error writing the response", err)
	}
}

// TODO -- Endpoint for file stastistics
// 1. Create a global variable which will be an array of Files -- DONE
// 2. Send than array when calling SendFile function -- DONE
// 3. In that function add the files after processing it
// 4. Parse to JSON and open the enpoint
