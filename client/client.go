package client

import (
	"io"
	"log"
	"net"
	"os"

	utils "github.com/Odzen/TCPCustomFileServer/utils"
	"github.com/joho/godotenv"
)

// var messagesGlobales = make(chan string)

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

	//go handleFiles(connection)

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

// func handleFiles(connection net.Conn) {
// 	log.Println("Handleling files")
// 	bufferFileName := make([]byte, 64)
// 	bufferFileSize := make([]byte, 10)

// 	// log.Println("File name", bufferFileName)
// 	// log.Println("File size", bufferFileSize)
// 	connection.Read(bufferFileSize)
// 	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

// 	connection.Read(bufferFileName)
// 	fileName := strings.Trim(string(bufferFileName), ":")

// 	log.Println("Before creating files")
// 	newFile, err := os.Create(fileName)

// 	if err != nil {
// 		panic(err)
// 	}
// 	defer newFile.Close()
// 	var receivedBytes int64
// 	log.Println("Before For files")

// 	for {
// 		if (fileSize - receivedBytes) < BUFFERSIZE {
// 			io.CopyN(newFile, connection, (fileSize - receivedBytes))
// 			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
// 			break
// 		}
// 		io.CopyN(newFile, connection, BUFFERSIZE)
// 		receivedBytes += BUFFERSIZE
// 	}
// 	fmt.Println("Received file completely!")
// }

func copyContent(receiver io.Writer, source io.Reader) {
	_, err := io.Copy(receiver, source)

	if err != nil {
		log.Fatal(err)
	}
}

// // Escribe todos los mensajes que se van recibiendo
// // <- chan significa que el canal es exclusivo para leer
// func MessageWrite(conn net.Conn) {
// 	for message := range messagesGlobales {
// 		// Se escriben mensajes que están siendo recibidos a través del canal
// 		fmt.Fprintln(conn, message)
// 	}
// }
