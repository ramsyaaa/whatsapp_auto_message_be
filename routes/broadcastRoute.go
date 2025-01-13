package routes

import (
	"go_whatsapp/modules/broadcast/http"
	"go_whatsapp/modules/broadcast/repository"
	"go_whatsapp/modules/broadcast/service"

	"github.com/gofiber/fiber/v2"
	"go.mau.fi/whatsmeow"
	"gorm.io/gorm"
)

func BroadcastRouter(app *fiber.App, db *gorm.DB, client *whatsmeow.Client) {
	broadcastRepo := repository.NewBroadcastRepository(db)
	broadcastService := service.NewBroadcastService(broadcastRepo)
	broadcastHandler := http.NewBroadcastHandler(broadcastService, client)

	http.BroadcastRoutes(app, broadcastHandler)
}
