package client

import (
	"io"
	"log"
	"net"
	"os"

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
	defer connection.Close()

	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, connection)
		done <- struct{}{}
	}()

	copyContent(connection, os.Stdin)
	<-done
}

func copyContent(destino io.Writer, fuente io.Reader) {
	_, err := io.Copy(destino, fuente)

	if err != nil {
		log.Fatal(err)
	}
}
