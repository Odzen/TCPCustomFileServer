package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
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
	err := connection.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Fatalln("Timeout error")
	}
	defer connection.Close()

	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
		fmt.Fprintf(connection, "I heard you say: %s\n", ln)
	}

}
