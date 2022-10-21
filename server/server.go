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

		// ln := scanner.Text()
		// m := strings.Fields(ln)[0] // method

		// if m == "GET" {
		// 	fmt.Println("HTTP REQUEST")
		// } else {
		// 	fmt.Println("Command")
		// }
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
			types.SendFile(command.Client, command.Args, channelGroup)

		case types.EXIT:
			clientLeft = true
			types.Exit(command.Client, channelGroup)

		}
	}
}

func handleHttpRequest(conn net.Conn) {
	i := 0
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
		if i == 0 {
			mux(conn, ln)
		}
		if ln == "" {
			// headers are done
			break
		}
		i++
	}
}

func mux(conn net.Conn, ln string) {
	// request line
	m := strings.Fields(ln)[0] // method
	u := strings.Fields(ln)[1] // uri
	fmt.Println("***METHOD", m)
	fmt.Println("***URI", u)

	// multiplexer
	if m == "GET" && u == "/clients" {
		getClients(conn)
	}
}

func getClients(conn net.Conn) {
	body := channelGroup.ToJson()
	fmt.Println("Body", body)

	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: application/json\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, body)
}
