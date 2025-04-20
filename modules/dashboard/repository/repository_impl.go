package repository

import (
	"context"
	"go_whatsapp/models"
	"time"

	"gorm.io/gorm"
)

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetTotalMessages(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Message{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalBroadcasts(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.BroadcastJob{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalRecipients(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.BroadcastRecipient{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetMessagesByDate(ctx context.Context, days int) (map[string]int64, error) {
	// Get messages grouped by date for the last 'days' days
	result := make(map[string]int64)

	// Use Asia/Jakarta timezone
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback to UTC if timezone loading fails
		loc = time.UTC
	}

	// Calculate the start date (n days ago)
	now := time.Now().In(loc)
	startDate := now.AddDate(0, 0, -days+1).Truncate(24 * time.Hour)

	// Initialize all dates in the range with zero counts
	for i := 0; i < days; i++ {
		date := startDate.AddDate(0, 0, i)
		result[date.Format("2006-01-02")] = 0
	}

	// Query the database for message counts by date
	type DateCount struct {
		Date  string
		Count int64
	}

	var dateCounts []DateCount

	// SQL query to count messages by date
	err = r.db.WithContext(ctx).
		Model(&models.Message{}).
		Select("TO_CHAR(sent_at, 'YYYY-MM-DD') as date, COUNT(*) as count").
		Where("sent_at >= ?", startDate).
		Group("TO_CHAR(sent_at, 'YYYY-MM-DD')").
		Order("date").
		Scan(&dateCounts).Error

	if err != nil {
		return nil, err
	}

	// Update the result map with actual counts
	for _, dc := range dateCounts {
		result[dc.Date] = dc.Count
	}

	return result, nil
}

func (r *dashboardRepository) GetBroadcastStatusCounts(ctx context.Context) (map[string]int64, error) {
	// Get counts of broadcast recipients by status
	result := map[string]int64{
		"Success": 0,
		"Pending": 0,
		"Failed":  0,
	}

	type StatusCount struct {
		Status string
		Count  int64
	}

	var statusCounts []StatusCount

	err := r.db.WithContext(ctx).
		Model(&models.BroadcastRecipient{}).
		Select("broadcast_status as status, COUNT(*) as count").
		Group("broadcast_status").
		Scan(&statusCounts).Error

	if err != nil {
		return nil, err
	}

	// Update the result map with actual counts
	for _, sc := range statusCounts {
		result[sc.Status] = sc.Count
	}

	return result, nil
}

func (r *dashboardRepository) GetRecentBroadcasts(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	var broadcasts []models.BroadcastJob
	err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&broadcasts).Error
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(broadcasts))
	for i, broadcast := range broadcasts {
		// Count recipients for this broadcast
		var totalRecipients int64
		r.db.Model(&models.BroadcastRecipient{}).Where("broadcast_job_id = ?", broadcast.ID).Count(&totalRecipients)

		// Count successful recipients
		var successfulRecipients int64
		r.db.Model(&models.BroadcastRecipient{}).Where("broadcast_job_id = ? AND broadcast_status = ?", broadcast.ID, "Success").Count(&successfulRecipients)

		// Calculate success rate
		successRate := 0.0
		if totalRecipients > 0 {
			successRate = float64(successfulRecipients) / float64(totalRecipients) * 100
		}

		result[i] = map[string]interface{}{
			"id":                broadcast.ID,
			"client_info":       broadcast.ClientInfo,
			"broadcast_plan_at": broadcast.BroadcastPlanAt,
			"broadcast_status":  broadcast.BroadcastJobStatus,
			"created_at":        broadcast.CreatedAt,
			"total_recipients":  totalRecipients,
			"success_rate":      successRate,
			"broadcast_message": broadcast.BroadcastMessage,
		}
	}

	return result, nil
}

func (r *dashboardRepository) GetBroadcastSuccessRate(ctx context.Context) (float64, error) {
	// Get total number of recipients
	var totalRecipients int64
	err := r.db.WithContext(ctx).Model(&models.BroadcastRecipient{}).Count(&totalRecipients).Error
	if err != nil {
		return 0, err
	}

	if totalRecipients == 0 {
		return 0, nil
	}

	// Get number of successful recipients
	var successfulRecipients int64
	err = r.db.WithContext(ctx).Model(&models.BroadcastRecipient{}).Where("broadcast_status = ?", "Success").Count(&successfulRecipients).Error
	if err != nil {
		return 0, err
	}

	// Calculate success rate
	return float64(successfulRecipients) / float64(totalRecipients) * 100, nil
}

func (r *dashboardRepository) GetMessagesByHour(ctx context.Context, date time.Time) (map[int]int64, error) {
	// Get messages grouped by hour for a specific date
	result := make(map[int]int64)

	// Initialize all hours with zero counts
	for i := 0; i < 24; i++ {
		result[i] = 0
	}

	// Use Asia/Jakarta timezone
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback to UTC if timezone loading fails
		loc = time.UTC
	}

	// Set the date to the beginning of the day in the specified timezone
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	endDate := startDate.AddDate(0, 0, 1)

	type HourCount struct {
		Hour  int
		Count int64
	}

	var hourCounts []HourCount

	// SQL query to count messages by hour
	err = r.db.WithContext(ctx).
		Model(&models.Message{}).
		Select("EXTRACT(HOUR FROM sent_at) as hour, COUNT(*) as count").
		Where("sent_at >= ? AND sent_at < ?", startDate, endDate).
		Group("EXTRACT(HOUR FROM sent_at)").
		Order("hour").
		Scan(&hourCounts).Error

	if err != nil {
		return nil, err
	}

	// Update the result map with actual counts
	for _, hc := range hourCounts {
		result[hc.Hour] = hc.Count
	}

	return result, nil
}

func (r *dashboardRepository) GetTopRecipients(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	// Get top recipients by message count
	type RecipientCount struct {
		RecipientNumber string
		Count           int64
	}

	var recipientCounts []RecipientCount

	err := r.db.WithContext(ctx).
		Model(&models.Message{}).
		Select("recipient_number, COUNT(*) as count").
		Group("recipient_number").
		Order("count DESC").
		Limit(limit).
		Scan(&recipientCounts).Error

	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(recipientCounts))
	for i, rc := range recipientCounts {
		result[i] = map[string]interface{}{
			"recipient_number": rc.RecipientNumber,
			"message_count":    rc.Count,
		}
	}

	return result, nil
}
