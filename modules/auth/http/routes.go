package http

import (
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, handler *AuthHandler) {
	app.Get("/authenticate", handler.HandleQR)
	app.Post("/log-out", handler.HandleLogout)
}
