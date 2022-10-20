package client

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"os"
	"strings"

	"github.com/Odzen/TCPCustomFileServer/types"
	utils "github.com/Odzen/TCPCustomFileServer/utils"
	"github.com/joho/godotenv"
)

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func EstablishConnection() {
	connection, err := net.Dial(os.Getenv("PROTOCOL_TYPE"), os.Getenv("HOST")+":"+os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Client connected using the port, ", connection.LocalAddr().String())
	defer utils.CloseConnectionClient(connection)

	reader := bufio.NewReader(os.Stdin)

	for {

		newLine, err := reader.ReadString('\n')

		if err != nil {
			log.Println(err)
		}

		// Process commands
		newLine = strings.Trim(newLine, "\r\n")

		args := strings.Split(newLine, " ")
		command := strings.TrimSpace(args[0])

		log.Println("Command, ", command)

		commandStruct, isValid := types.ProcessCommand(command, args)
		if !isValid {
			log.Printf("-> The command `%s` was not accepted. Use the command `=instructions` to see the available commands \n", command)
		}
		log.Println("Command: ", commandStruct)

		body := types.BodyMessage{
			TypedCommand: commandStruct,
		}

		// Deserialized the information sent by the client
		err = gob.NewEncoder(connection).Encode(&body)

		if err != nil {
			return
		}

		//ProcessCommand(command, args)

	}

}
