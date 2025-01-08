package routes

import (
	"go_whatsapp/modules/auth/http"

	"github.com/gofiber/fiber/v2"
	"go.mau.fi/whatsmeow"
)

func AuthRouter(app *fiber.App, client *whatsmeow.Client) {
	authHandler := http.NewAuthHandler(client)
	http.AuthRoutes(app, authHandler)
}
