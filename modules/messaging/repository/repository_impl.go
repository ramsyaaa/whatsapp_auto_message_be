package repository

import (
	"context"
	"go_whatsapp/models"

	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) SaveMessage(ctx context.Context, message models.Message) (models.Message, error) {
	err := r.db.WithContext(ctx).Create(&message).Error
	return message, err
}

func (r *messageRepository) GetRecentMessages(ctx context.Context, limit int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&messages).Error
	return messages, err
}

func (r *messageRepository) UpdateMessageStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&models.Message{}).Where("id = ?", id).Update("status", status).Error
}
