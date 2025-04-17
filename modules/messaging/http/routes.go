package http

import (
	"github.com/gofiber/fiber/v2"
)

func MessagingRoutes(app *fiber.App, handler *MessagingHandler) {
	app.Post("/send-message", handler.HandleSendMessage)
	app.Get("/recent-messages", handler.GetRecentMessages)
}
