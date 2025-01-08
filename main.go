package main

import (
	"go_whatsapp/config"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using environment variables instead")
	}

	db, client := config.Connect()
	config.Route(db, client)

	if client.Store.ID == nil {
		log.Println("No WhatsApp session found, please scan QR code first")
	}
}
