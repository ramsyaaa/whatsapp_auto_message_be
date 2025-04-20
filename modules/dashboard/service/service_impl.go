package service

import (
	"context"
	"go_whatsapp/modules/dashboard/repository"
	"time"
)

type dashboardService struct {
	repo repository.DashboardRepository
}

func NewDashboardService(repo repository.DashboardRepository) DashboardService {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) GetDashboardStats(ctx context.Context) (map[string]interface{}, error) {
	// Get total messages
	totalMessages, err := s.repo.GetTotalMessages(ctx)
	if err != nil {
		return nil, err
	}
	
	// Get total broadcasts
	totalBroadcasts, err := s.repo.GetTotalBroadcasts(ctx)
	if err != nil {
		return nil, err
	}
	
	// Get total recipients
	totalRecipients, err := s.repo.GetTotalRecipients(ctx)
	if err != nil {
		return nil, err
	}
	
	// Get broadcast success rate
	successRate, err := s.repo.GetBroadcastSuccessRate(ctx)
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"total_messages":   totalMessages,
		"total_broadcasts": totalBroadcasts,
		"total_recipients": totalRecipients,
		"success_rate":     successRate,
	}, nil
}

func (s *dashboardService) GetMessageActivity(ctx context.Context, days int) (map[string]interface{}, error) {
	// Get message counts by date
	messagesByDate, err := s.repo.GetMessagesByDate(ctx, days)
	if err != nil {
		return nil, err
	}
	
	// Convert to arrays for Chart.js
	dates := make([]string, 0, len(messagesByDate))
	counts := make([]int64, 0, len(messagesByDate))
	
	// Use Asia/Jakarta timezone
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback to UTC if timezone loading fails
		loc = time.UTC
	}
	
	// Calculate the start date (n days ago)
	now := time.Now().In(loc)
	startDate := now.AddDate(0, 0, -days+1).Truncate(24 * time.Hour)
	
	// Ensure we have data for all dates in order
	for i := 0; i < days; i++ {
		date := startDate.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		dates = append(dates, date.Format("Jan 02")) // Format for display
		counts = append(counts, messagesByDate[dateStr])
	}
	
	return map[string]interface{}{
		"labels": dates,
		"data":   counts,
	}, nil
}

func (s *dashboardService) GetBroadcastStatus(ctx context.Context) (map[string]interface{}, error) {
	// Get broadcast status counts
	statusCounts, err := s.repo.GetBroadcastStatusCounts(ctx)
	if err != nil {
		return nil, err
	}
	
	// Convert to arrays for Chart.js
	labels := []string{"Success", "Pending", "Failed"}
	data := []int64{
		statusCounts["Success"],
		statusCounts["Pending"],
		statusCounts["Failed"],
	}
	
	return map[string]interface{}{
		"labels": labels,
		"data":   data,
	}, nil
}

func (s *dashboardService) GetRecentBroadcasts(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	return s.repo.GetRecentBroadcasts(ctx, limit)
}

func (s *dashboardService) GetHourlyMessageStats(ctx context.Context, date time.Time) (map[string]interface{}, error) {
	// Get message counts by hour
	messagesByHour, err := s.repo.GetMessagesByHour(ctx, date)
	if err != nil {
		return nil, err
	}
	
	// Convert to arrays for Chart.js
	hours := make([]string, 24)
	counts := make([]int64, 24)
	
	for i := 0; i < 24; i++ {
		hours[i] = formatHour(i)
		counts[i] = messagesByHour[i]
	}
	
	return map[string]interface{}{
		"labels": hours,
		"data":   counts,
		"date":   date.Format("2006-01-02"),
	}, nil
}

func (s *dashboardService) GetTopRecipients(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	return s.repo.GetTopRecipients(ctx, limit)
}

// Helper function to format hour
func formatHour(hour int) string {
	if hour == 0 {
		return "12 AM"
	} else if hour < 12 {
		return time.Date(0, 0, 0, hour, 0, 0, 0, time.UTC).Format("3 PM")
	} else if hour == 12 {
		return "12 PM"
	} else {
		return time.Date(0, 0, 0, hour-12, 0, 0, 0, time.UTC).Format("3 PM")
	}
}
