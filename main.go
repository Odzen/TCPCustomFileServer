package main

import (
	"fmt"
	"os"

	client "github.com/Odzen/TCPCustomFileServer/client"
	server "github.com/Odzen/TCPCustomFileServer/server"
)

func main() {

	program := os.Args[1:]

	if len(program) == 0 {
		fmt.Println("Select please which program do you want to execute, first run the server, then the client")
		return
	}

	switch program[0] {
	case "client":
		client.EstablishConnection()
	case "server":
		server.RunServer()
	default:
		fmt.Println("Program not found.")
	}

}
