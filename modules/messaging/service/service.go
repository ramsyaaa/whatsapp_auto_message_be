package service

import (
	"context"
	"go_whatsapp/models"
)

type MessageService interface {
	SaveMessage(ctx context.Context, recipientNumber, messageContent string) (models.Message, error)
	GetRecentMessages(ctx context.Context, limit int) ([]models.Message, error)
	UpdateMessageStatus(ctx context.Context, id uint, status string) error
}
