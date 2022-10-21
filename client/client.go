package client

import (
	"io"
	"log"
	"net"
	"os"

	utils "github.com/Odzen/TCPCustomFileServer/utils"
	"github.com/joho/godotenv"
)

const BUFFERSIZE = 1024

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

	defer utils.CloseConnectionClient(connection)

	done := make(chan struct{})

	go func() {

		_, err := io.Copy(os.Stdout, connection)
		if err != nil {
			log.Fatal(err)
		}
		done <- struct{}{}
	}()

	copyContent(connection, os.Stdin)

	<-done

}

func copyContent(receiver io.Writer, source io.Reader) {
	_, err := io.Copy(receiver, source)

	if err != nil {
		log.Fatal(err)
	}
}
