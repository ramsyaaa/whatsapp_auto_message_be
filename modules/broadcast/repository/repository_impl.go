package repository

import (
	"context"
	"go_whatsapp/models"
	"time"

	"gorm.io/gorm"
)

type broadcastRepository struct {
	db *gorm.DB
}

func NewBroadcastRepository(db *gorm.DB) BroadcastRepository {
	return &broadcastRepository{db: db}
}

func (r *broadcastRepository) CreateBroadcast(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
	broadcast := models.BroadcastJob{
		ClientInfo:       data["client_info"].(string),
		BroadcastPlanAt:  parseTime(data["broadcast_plan_at"].(string)),
		CreatedAt:        parseTime(data["created_at"].(string)),
		BroadcastCode:    data["broadcast_code"].(string),
		BroadcastMessage: data["broadcast_message"].(string),
	}
	err := r.db.WithContext(ctx).Create(&broadcast).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func parseTime(value string) time.Time {
	// Use Asia/Jakarta timezone
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback to UTC if timezone loading fails
		loc = time.UTC
	}
	parsedTime, _ := time.Parse(time.RFC3339, value) // Handle error as needed
	return parsedTime.In(loc)
}

func (r *broadcastRepository) ImportRecipient(ctx context.Context, broadcastID int, data []map[string]interface{}) (map[string]interface{}, error) {
	var recipients []models.BroadcastRecipient
	for _, recipient := range data {
		broadcastRecipient := models.BroadcastRecipient{
			BroadcastJobID:  uint(broadcastID),
			WhatsappNumber:  recipient["whatsapp_number"].(string),
			BroadcastStatus: "Pending",
		}
		recipients = append(recipients, broadcastRecipient)
	}

	err := r.db.WithContext(ctx).Create(&recipients).Error
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"success": true}, nil
}

func (r *broadcastRepository) ImportPecatuRecipient(ctx context.Context, broadcastID int, data []map[string]interface{}) (map[string]interface{}, error) {
	var recipients []models.BroadcastRecipient
	for _, recipient := range data {
		recipientName, ok := recipient["name"].(string)
		if !ok {
			recipientName = ""
		}
		recipientUniqueIdentifier, ok := recipient["identifier"].(string)
		if !ok {
			recipientUniqueIdentifier = ""
		}
		whatsappNumber, ok := recipient["whatsapp_number"].(string)
		if !ok {
			whatsappNumber = ""
		}

		broadcastRecipient := models.BroadcastRecipient{
			BroadcastJobID:            uint(broadcastID),
			RecipientName:             recipientName,
			RecipientUniqueIdentifier: recipientUniqueIdentifier,
			WhatsappNumber:            whatsappNumber,
			BroadcastStatus:           "Pending",
		}
		recipients = append(recipients, broadcastRecipient)
	}

	err := r.db.WithContext(ctx).Create(&recipients).Error
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"success": true}, nil
}

func (r *broadcastRepository) BroadcastDetail(ctx context.Context, broadcastID int) (map[string]interface{}, error) {
	var broadcast models.BroadcastJob
	err := r.db.WithContext(ctx).Where("id = ?", broadcastID).First(&broadcast).Error
	if err != nil {
		return nil, err
	}

	detail := map[string]interface{}{
		"id":                broadcast.ID,
		"client_info":       broadcast.ClientInfo,
		"broadcast_plan_at": broadcast.BroadcastPlanAt,
		"created_at":        broadcast.CreatedAt,
		"broadcast_code":    broadcast.BroadcastCode,
		"broadcast_message": broadcast.BroadcastMessage,
	}

	return detail, nil
}

func (r *broadcastRepository) GetBroadcastMessage(ctx context.Context, broadcastID int) (map[string]interface{}, error) {
	var broadcast models.BroadcastJob
	err := r.db.WithContext(ctx).Where("id = ?", broadcastID).First(&broadcast).Error
	if err != nil {
		return nil, err
	}

	detail := map[string]interface{}{
		"broadcast_message": broadcast.BroadcastMessage,
	}

	return detail, nil
}

func (r *broadcastRepository) GetAllRecipientByBroadcastID(ctx context.Context, broadcastID int) (map[string]interface{}, error) {
	var recipients []models.BroadcastRecipient
	// Order by broadcast_at DESC (newest first), then by status (Success first, then Pending, then Failed)
	err := r.db.WithContext(ctx).Where("broadcast_job_id = ?", broadcastID).
		Order("CASE WHEN broadcast_status = 'Success' THEN 0 WHEN broadcast_status = 'Pending' THEN 1 ELSE 2 END").
		Order("broadcasted_at DESC NULLS LAST").
		Find(&recipients).Error
	if err != nil {
		return nil, err
	}

	recipientMaps := make([]map[string]interface{}, len(recipients))
	for i, recipient := range recipients {
		recipientMaps[i] = map[string]interface{}{
			"id":                          int(recipient.ID),
			"recipient_name":              recipient.RecipientName,
			"recipient_unique_identifier": recipient.RecipientUniqueIdentifier,
			"broadcast_status":            recipient.BroadcastStatus,
			"broadcast_at":                recipient.BroadcastedAt,
			"whatsapp_number":             recipient.WhatsappNumber,
		}
	}

	return map[string]interface{}{
		"recipients": recipientMaps,
	}, nil
}

func (r *broadcastRepository) GetPendingRecipientsByBroadcastID(ctx context.Context, broadcastID int) (map[string]interface{}, error) {
	var recipients []models.BroadcastRecipient
	// Even though we're only getting Pending status, maintain consistent ordering by broadcasted_at
	err := r.db.WithContext(ctx).Where("broadcast_job_id = ? AND broadcast_status = ?", broadcastID, "Pending").
		Order("broadcasted_at DESC NULLS LAST").
		Find(&recipients).Error
	if err != nil {
		return nil, err
	}

	recipientMaps := make([]map[string]interface{}, len(recipients))
	for i, recipient := range recipients {
		recipientMaps[i] = map[string]interface{}{
			"id":                          int(recipient.ID),
			"recipient_name":              recipient.RecipientName,
			"recipient_unique_identifier": recipient.RecipientUniqueIdentifier,
			"broadcast_status":            recipient.BroadcastStatus,
			"broadcast_at":                recipient.BroadcastedAt,
			"whatsapp_number":             recipient.WhatsappNumber,
		}
	}

	return map[string]interface{}{
		"recipients": recipientMaps,
	}, nil
}

func (r *broadcastRepository) GetPaginatedRecipientsByBroadcastID(ctx context.Context, broadcastID int, page int, limit int, search string) (map[string]interface{}, error) {
	var recipients []models.BroadcastRecipient
	var totalCount int64

	// Calculate offset
	offset := (page - 1) * limit

	// Create base query
	query := r.db.WithContext(ctx).Model(&models.BroadcastRecipient{}).Where("broadcast_job_id = ?", broadcastID)

	// Add search condition if search string is provided
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("(whatsapp_number LIKE ? OR recipient_name LIKE ?)", searchTerm, searchTerm)
	}

	// Get total count for pagination
	query.Count(&totalCount)

	// Get paginated results with ordering
	// Order by broadcast_at DESC (newest first), then by status (Success first, then Pending, then Failed)
	err := query.Order("CASE WHEN broadcast_status = 'Success' THEN 0 WHEN broadcast_status = 'Pending' THEN 1 ELSE 2 END").
		Order("broadcasted_at DESC NULLS LAST").
		Offset(offset).Limit(limit).Find(&recipients).Error
	if err != nil {
		return nil, err
	}

	recipientMaps := make([]map[string]interface{}, len(recipients))
	for i, recipient := range recipients {
		recipientMaps[i] = map[string]interface{}{
			"id":                          int(recipient.ID),
			"recipient_name":              recipient.RecipientName,
			"recipient_unique_identifier": recipient.RecipientUniqueIdentifier,
			"broadcast_status":            recipient.BroadcastStatus,
			"broadcast_at":                recipient.BroadcastedAt,
			"whatsapp_number":             recipient.WhatsappNumber,
		}
	}

	// Calculate total pages
	totalPages := (int(totalCount) + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	return map[string]interface{}{
		"recipients": recipientMaps,
		"pagination": map[string]interface{}{
			"total_count":  totalCount,
			"total_pages":  totalPages,
			"current_page": page,
			"limit":        limit,
		},
	}, nil
}

func (r *broadcastRepository) UpdateBroadcastStatus(ctx context.Context, broadcastID int, status string) (map[string]interface{}, error) {
	var broadcast models.BroadcastJob
	err := r.db.WithContext(ctx).Where("id = ?", broadcastID).First(&broadcast).Error
	if err != nil {
		return nil, err
	}

	broadcast.BroadcastJobStatus = status
	err = r.db.WithContext(ctx).Save(&broadcast).Error
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "Broadcast status updated successfully",
	}, nil
}

func (r *broadcastRepository) UpdateRecipientBroadcastStatus(ctx context.Context, recipientID int, broadcastID int, status string) (map[string]interface{}, error) {
	var recipient models.BroadcastRecipient
	err := r.db.WithContext(ctx).Where("id = ? AND broadcast_job_id = ?", recipientID, broadcastID).First(&recipient).Error
	if err != nil {
		return nil, err
	}

	// Use Asia/Jakarta timezone
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback to UTC if timezone loading fails
		loc = time.UTC
	}
	now := time.Now().In(loc)

	recipient.BroadcastStatus = status
	recipient.BroadcastedAt = now // Set the broadcasted_at using current timestamp with Jakarta timezone
	err = r.db.WithContext(ctx).Save(&recipient).Error
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "Recipient broadcast status updated successfully",
	}, nil
}

func (r *broadcastRepository) IsAnyRecipientInBroadcast(ctx context.Context, broadcastID int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.BroadcastRecipient{}).Where("broadcast_job_id = ?", broadcastID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *broadcastRepository) DeleteRecipient(ctx context.Context, recipientID int) (map[string]interface{}, error) {
	// Find the recipient first to make sure it exists
	var recipient models.BroadcastRecipient
	err := r.db.WithContext(ctx).Where("id = ?", recipientID).First(&recipient).Error
	if err != nil {
		return nil, err
	}

	// Delete the recipient
	err = r.db.WithContext(ctx).Delete(&recipient).Error
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "Recipient deleted successfully",
	}, nil
}

func (r *broadcastRepository) DeleteBroadcast(ctx context.Context, broadcastID int) (map[string]interface{}, error) {
	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Delete all recipients associated with the broadcast
	if err := tx.Where("broadcast_job_id = ?", broadcastID).Delete(&models.BroadcastRecipient{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Delete the broadcast
	var broadcast models.BroadcastJob
	if err := tx.Where("id = ?", broadcastID).First(&broadcast).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Delete(&broadcast).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "Broadcast and all its recipients deleted successfully",
	}, nil
}

func (r *broadcastRepository) GetAllBroadcasts(ctx context.Context) ([]map[string]interface{}, error) {
	var broadcasts []models.BroadcastJob
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&broadcasts).Error
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(broadcasts))
	for i, broadcast := range broadcasts {
		// Count recipients for this broadcast
		var recipientCount int64
		r.db.Model(&models.BroadcastRecipient{}).Where("broadcast_job_id = ?", broadcast.ID).Count(&recipientCount)

		// Check if this broadcast has any Pecatu recipients (with name and identifier)
		var pecatuCount int64
		r.db.Model(&models.BroadcastRecipient{}).Where("broadcast_job_id = ? AND recipient_name != '' AND recipient_unique_identifier != ''", broadcast.ID).Count(&pecatuCount)

		// Determine broadcast type based on recipients
		broadcastType := "regular"
		if pecatuCount > 0 {
			broadcastType = "pecatu"
		}

		result[i] = map[string]interface{}{
			"id":                   broadcast.ID,
			"client_info":          broadcast.ClientInfo,
			"broadcast_plan_at":    broadcast.BroadcastPlanAt,
			"broadcast_job_status": broadcast.BroadcastJobStatus,
			"created_at":           broadcast.CreatedAt,
			"broadcast_code":       broadcast.BroadcastCode,
			"broadcast_message":    broadcast.BroadcastMessage,
			"recipient_count":      recipientCount,
			"broadcast_type":       broadcastType,
		}
	}

	return result, nil
}
