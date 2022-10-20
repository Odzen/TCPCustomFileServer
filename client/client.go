package client

import (
	"io"
	"log"
	"net"
	"os"

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

	done := make(chan struct{})
	log.Println("Antes de la rutina")
	go func() {
		log.Println("En rutina")
		bytes, err := io.Copy(os.Stdout, connection)
		log.Println("Bytes read from console and written to connection: ", bytes)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Antes del done")
		done <- struct{}{}
	}()
	log.Println("Sali de la rutina")
	copyContent(connection, os.Stdin)
	log.Println("Antes del Done cerrado")
	<-done
	log.Println("Done cerrado")
}

func copyContent(receiver io.Writer, source io.Reader) {
	log.Println("En copy content")
	bytes, err := io.Copy(receiver, source)
	log.Println("Bytes read from connection and written to console: ", bytes)
	if err != nil {
		log.Fatal(err)
	}
}
