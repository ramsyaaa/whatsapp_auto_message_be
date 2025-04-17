package routes

import (
	"go_whatsapp/modules/messaging/http"
	"go_whatsapp/modules/messaging/repository"
	"go_whatsapp/modules/messaging/service"

	"github.com/gofiber/fiber/v2"
	"go.mau.fi/whatsmeow"
	"gorm.io/gorm"
)

func MessagingRouter(app *fiber.App, client *whatsmeow.Client, db *gorm.DB) {
	messageRepo := repository.NewMessageRepository(db)
	messageService := service.NewMessageService(messageRepo)
	messagingHandler := http.NewMessagingHandler(client, messageService)
	http.MessagingRoutes(app, messagingHandler)
}
