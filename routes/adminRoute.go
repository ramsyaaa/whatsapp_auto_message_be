package routes

import (
	"go_whatsapp/modules/admin/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AdminRouter(app *fiber.App, db *gorm.DB) {
	adminHandler := http.NewAdminHandler(db)
	http.AdminRoutes(app, adminHandler)
}
