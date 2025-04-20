package http

import (
	"context"
	"go_whatsapp/helper"
	"go_whatsapp/modules/dashboard/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	service service.DashboardService
}

func NewDashboardHandler(service service.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

func (h *DashboardHandler) GetDashboardStats(c *fiber.Ctx) error {
	ctx := context.Background()
	
	stats, err := h.service.GetDashboardStats(ctx)
	if err != nil {
		response := helper.APIResponse("Failed to get dashboard stats", http.StatusInternalServerError, "ERROR", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	
	response := helper.APIResponse("Dashboard stats retrieved successfully", http.StatusOK, "OK", stats)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *DashboardHandler) GetMessageActivity(c *fiber.Ctx) error {
	ctx := context.Background()
	
	// Get days parameter, default to 7
	days, err := strconv.Atoi(c.Query("days", "7"))
	if err != nil || days < 1 || days > 30 {
		days = 7
	}
	
	activity, err := h.service.GetMessageActivity(ctx, days)
	if err != nil {
		response := helper.APIResponse("Failed to get message activity", http.StatusInternalServerError, "ERROR", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	
	response := helper.APIResponse("Message activity retrieved successfully", http.StatusOK, "OK", activity)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *DashboardHandler) GetBroadcastStatus(c *fiber.Ctx) error {
	ctx := context.Background()
	
	status, err := h.service.GetBroadcastStatus(ctx)
	if err != nil {
		response := helper.APIResponse("Failed to get broadcast status", http.StatusInternalServerError, "ERROR", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	
	response := helper.APIResponse("Broadcast status retrieved successfully", http.StatusOK, "OK", status)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *DashboardHandler) GetRecentBroadcasts(c *fiber.Ctx) error {
	ctx := context.Background()
	
	// Get limit parameter, default to 5
	limit, err := strconv.Atoi(c.Query("limit", "5"))
	if err != nil || limit < 1 || limit > 20 {
		limit = 5
	}
	
	broadcasts, err := h.service.GetRecentBroadcasts(ctx, limit)
	if err != nil {
		response := helper.APIResponse("Failed to get recent broadcasts", http.StatusInternalServerError, "ERROR", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	
	response := helper.APIResponse("Recent broadcasts retrieved successfully", http.StatusOK, "OK", broadcasts)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *DashboardHandler) GetHourlyMessageStats(c *fiber.Ctx) error {
	ctx := context.Background()
	
	// Get date parameter, default to today
	dateStr := c.Query("date", "")
	var date time.Time
	var err error
	
	if dateStr == "" {
		// Use Asia/Jakarta timezone
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			// Fallback to UTC if timezone loading fails
			loc = time.UTC
		}
		date = time.Now().In(loc)
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			response := helper.APIResponse("Invalid date format", http.StatusBadRequest, "ERROR", nil)
			return c.Status(http.StatusBadRequest).JSON(response)
		}
	}
	
	stats, err := h.service.GetHourlyMessageStats(ctx, date)
	if err != nil {
		response := helper.APIResponse("Failed to get hourly message stats", http.StatusInternalServerError, "ERROR", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	
	response := helper.APIResponse("Hourly message stats retrieved successfully", http.StatusOK, "OK", stats)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *DashboardHandler) GetTopRecipients(c *fiber.Ctx) error {
	ctx := context.Background()
	
	// Get limit parameter, default to 5
	limit, err := strconv.Atoi(c.Query("limit", "5"))
	if err != nil || limit < 1 || limit > 20 {
		limit = 5
	}
	
	recipients, err := h.service.GetTopRecipients(ctx, limit)
	if err != nil {
		response := helper.APIResponse("Failed to get top recipients", http.StatusInternalServerError, "ERROR", nil)
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	
	response := helper.APIResponse("Top recipients retrieved successfully", http.StatusOK, "OK", recipients)
	return c.Status(http.StatusOK).JSON(response)
}
