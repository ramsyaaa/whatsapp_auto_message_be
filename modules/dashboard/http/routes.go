package http

import (
	"github.com/gofiber/fiber/v2"
)

func DashboardRoutes(app *fiber.App, handler *DashboardHandler) {
	app.Get("/dashboard/stats", handler.GetDashboardStats)
	app.Get("/dashboard/message-activity", handler.GetMessageActivity)
	app.Get("/dashboard/broadcast-status", handler.GetBroadcastStatus)
	app.Get("/dashboard/recent-broadcasts", handler.GetRecentBroadcasts)
	app.Get("/dashboard/hourly-stats", handler.GetHourlyMessageStats)
	app.Get("/dashboard/top-recipients", handler.GetTopRecipients)
}
