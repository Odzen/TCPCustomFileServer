package server

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func PrintVariable() {
	fmt.Printf("Host %s uses port %s and protcol %s\n", os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("PROTOCOL_TYPE"))
}
