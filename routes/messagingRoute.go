package routes

import (
	"go_whatsapp/modules/messaging/http"

	"github.com/gofiber/fiber/v2"
	"go.mau.fi/whatsmeow"
)

func MessagingRouter(app *fiber.App, client *whatsmeow.Client) {
	authHandler := http.NewMessagingHandler(client)
	http.MessagingRoutes(app, authHandler)
}
