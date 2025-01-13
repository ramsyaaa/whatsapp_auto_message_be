package http

import (
	"github.com/gofiber/fiber/v2"
)

func BroadcastRoutes(app *fiber.App, handler *BroadcastHandler) {
	app.Post("/broadcast/create", handler.CreateBroadcast)
	app.Post("/broadcast/import-recipient", handler.ImportRecipient)
	app.Post("/broadcast/pecatu/import-recipient", handler.ImportPecatuRecipient)
	app.Get("/broadcast/detail/:broadcast_id", handler.BroadcastDetail)
	app.Post("/broadcast/send", handler.HandleSendBroadcast)
	app.Post("/broadcast/pecatu/send", handler.HandlePecatuBroadcast)

}
