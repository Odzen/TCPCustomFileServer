package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

var clients = make(map[string]net.Conn)
var leaving = make(chan message)
var messages = make(chan message)

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

	clients[connection.RemoteAddr().String()] = connection

	messages <- newMessage(" joined.", connection)

	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		messages <- newMessage(": "+scanner.Text(), connection)
	}

	delete(clients, connection.RemoteAddr().String())

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
	for {
		select {
		case msg := <-messages:
			for _, conn := range clients {
				if msg.address == conn.RemoteAddr().String() {
					continue
				}
				fmt.Fprintln(conn, msg.text) // NOTE: ignoring network errors
			}

		case msg := <-leaving:
			for _, conn := range clients {
				fmt.Fprintln(conn, msg.text) // NOTE: ignoring network errors
			}

		}
	}
}
