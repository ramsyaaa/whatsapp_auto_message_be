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
	parsedTime, _ := time.Parse(time.RFC3339, value) // Menangani error sesuai kebutuhan
	return parsedTime.UTC()
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
	err := r.db.WithContext(ctx).Where("broadcast_job_id = ?", broadcastID).Find(&recipients).Error
	if err != nil {
		return nil, err
	}

	recipientMaps := make([]map[string]interface{}, len(recipients))
	for i, recipient := range recipients {
		recipientMaps[i] = map[string]interface{}{
			"id":                          int(recipient.ID),
			"recipient_name":              recipient.RecipientName,
			"recipient_unique_identifier": recipient.RecipientUniqueIdentifier,
			"whatsapp_number":             recipient.WhatsappNumber,
		}
	}

	return map[string]interface{}{
		"recipients": recipientMaps,
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

	recipient.BroadcastStatus = status
	recipient.BroadcastedAt = time.Now() // Set the broadcasted_at using current timestamp
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
