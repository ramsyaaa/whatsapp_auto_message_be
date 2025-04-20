package http

import (
	"github.com/gofiber/fiber/v2"
)

func BroadcastRoutes(app *fiber.App, handler *BroadcastHandler) {
	// First, register all specific routes
	app.Post("/broadcast/create", handler.CreateBroadcast)
	app.Post("/broadcast/import-recipient", handler.ImportRecipient)
	app.Post("/broadcast/pecatu/import-recipient", handler.ImportPecatuRecipient)
	app.Get("/broadcast/detail/:broadcast_id", handler.BroadcastDetail)
	app.Get("/broadcast/list", handler.GetAllBroadcasts)
	app.Get("/broadcast/recipients/:broadcast_id", handler.GetAllRecipientByBroadcastID)
	app.Delete("/broadcast/recipient/:recipient_id", handler.DeleteRecipient) // Added for deleting recipients
	app.Delete("/broadcast/:broadcast_id", handler.DeleteBroadcast)           // Added for deleting broadcasts
	app.Post("/broadcast/update-status", handler.UpdateBroadcastStatus)       // Added for updating broadcast status
	app.Post("/broadcast/send", handler.HandleSendBroadcast)
	app.Post("/broadcast/pecatu/send", handler.HandlePecatuBroadcast)
	app.Post("/broadcast/send-single", handler.HandleSendSingleBroadcast)
	app.Post("/broadcast/pecatu/send-single", handler.HandleSendSinglePecatuBroadcast)

	// Then, register the generic routes
	app.Get("/broadcast", handler.GetAllBroadcasts) // For compatibility with reports page

	// Finally, register the wildcard route (must be last to avoid conflicts)
	app.Get("/broadcast/:broadcast_id", handler.BroadcastDetail) // For compatibility with reports page
}
