package service

import (
	"context"
	"go_whatsapp/models"
	"go_whatsapp/modules/messaging/repository"
	"time"
)

type messageService struct {
	repo repository.MessageRepository
}

func NewMessageService(repo repository.MessageRepository) MessageService {
	return &messageService{repo: repo}
}

func (s *messageService) SaveMessage(ctx context.Context, recipientNumber, messageContent string) (models.Message, error) {
	// Use Asia/Jakarta timezone
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback to UTC if timezone loading fails
		loc = time.UTC
	}
	now := time.Now().In(loc)
	message := models.Message{
		RecipientNumber: recipientNumber,
		MessageContent:  messageContent,
		Status:          "Sent",
		SentAt:          now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	return s.repo.SaveMessage(ctx, message)
}

func (s *messageService) GetRecentMessages(ctx context.Context, limit int) ([]models.Message, error) {
	return s.repo.GetRecentMessages(ctx, limit)
}

func (s *messageService) UpdateMessageStatus(ctx context.Context, id uint, status string) error {
	return s.repo.UpdateMessageStatus(ctx, id, status)
}
