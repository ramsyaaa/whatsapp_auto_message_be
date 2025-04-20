package routes

import (
	"go_whatsapp/modules/dashboard/http"
	"go_whatsapp/modules/dashboard/repository"
	"go_whatsapp/modules/dashboard/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DashboardRouter(app *fiber.App, db *gorm.DB) {
	dashboardRepo := repository.NewDashboardRepository(db)
	dashboardService := service.NewDashboardService(dashboardRepo)
	dashboardHandler := http.NewDashboardHandler(dashboardService)

	http.DashboardRoutes(app, dashboardHandler)
}
