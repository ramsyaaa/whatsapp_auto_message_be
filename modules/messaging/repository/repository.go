package repository

import (
	"context"
	"go_whatsapp/models"
)

type MessageRepository interface {
	SaveMessage(ctx context.Context, message models.Message) (models.Message, error)
	GetRecentMessages(ctx context.Context, limit int) ([]models.Message, error)
	UpdateMessageStatus(ctx context.Context, id uint, status string) error
}
