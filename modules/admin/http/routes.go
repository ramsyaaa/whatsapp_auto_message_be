package http

import (
	"go_whatsapp/middleware"

	"github.com/gofiber/fiber/v2"
)

func AdminRoutes(app *fiber.App, handler *AdminHandler) {
	// Public routes (no authentication required)
	app.Post("/admin/login", handler.HandleLogin)
	
	// Protected routes (authentication required)
	adminGroup := app.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware())
	
	// Admin user management
	adminGroup.Post("/register", handler.HandleRegister) // This would typically be restricted
	adminGroup.Get("/profile", handler.HandleGetProfile)
}
