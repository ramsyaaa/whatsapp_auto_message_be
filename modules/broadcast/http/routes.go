package http

import (
	"github.com/gofiber/fiber/v2"
)

func BroadcastRoutes(app *fiber.App, handler *BroadcastHandler) {
	app.Post("/broadcast/create", handler.CreateBroadcast)
	app.Post("/broadcast/import-recipient", handler.ImportRecipient)
	app.Post("/broadcast/pecatu/import-recipient", handler.ImportPecatuRecipient)
	app.Get("/broadcast/detail/:broadcast_id", handler.BroadcastDetail)
	app.Get("/broadcast/:broadcast_id", handler.BroadcastDetail) // Added for compatibility with reports page
	app.Get("/broadcast/list", handler.GetAllBroadcasts)
	app.Get("/broadcast", handler.GetAllBroadcasts) // Added for compatibility with reports page
	app.Get("/broadcast/recipients/:broadcast_id", handler.GetAllRecipientByBroadcastID)
	app.Delete("/broadcast/recipient/:recipient_id", handler.DeleteRecipient) // Added for deleting recipients
	app.Post("/broadcast/update-status", handler.UpdateBroadcastStatus)       // Added for updating broadcast status
	app.Post("/broadcast/send", handler.HandleSendBroadcast)
	app.Post("/broadcast/pecatu/send", handler.HandlePecatuBroadcast)
	app.Post("/broadcast/send-single", handler.HandleSendSingleBroadcast)
	app.Post("/broadcast/pecatu/send-single", handler.HandleSendSinglePecatuBroadcast)
}
